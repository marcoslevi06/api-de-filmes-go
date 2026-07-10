package mongo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"sipub_teste/api/internal/domain"
)

const operationTimeout = 5 * time.Second

// MovieRepository é o adapter que guarda os filmes numa coleção MongoDB.
// Implementa a porta domain.MovieRepository.
type MovieRepository struct {
	collection *mongo.Collection
}

// movieDocument é a representação em BSON de um filme, isolada do
// domain.Movie para que o domínio não precise conhecer detalhes de
// persistência do MongoDB. O id inteiro do domínio é usado diretamente
// como _id do documento, mantendo o contrato da API (IDs inteiros)
// idêntico ao do adapter em memória.
type movieDocument struct {
	ID    int    `bson:"_id"`
	Title string `bson:"title"`
	Year  string `bson:"year"`
}

// fromDomain converte um domain.Movie para o formato de documento
// persistido no MongoDB.
func fromDomain(filme domain.Movie) movieDocument {
	return movieDocument{ID: filme.ID, Title: filme.Title, Year: filme.Year}
}

// toDomain converte um documento lido do MongoDB de volta para o tipo
// domain.Movie usado pelo restante da aplicação.
func (d movieDocument) toDomain() domain.Movie {
	return domain.Movie{ID: d.ID, Title: d.Title, Year: d.Year}
}

// NewMovieRepository conecta ao MongoDB usando a uri informada, seleciona
// o banco/coleção indicados e retorna o MovieRepository pronto para uso.
// Faz um Ping para garantir que o servidor está acessível antes de
// retornar; devolve erro caso a conexão ou o ping falhem.
func NewMovieRepository(ctx context.Context, uri, dbName, collectionName string) (*MovieRepository, error) {
	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar no MongoDB: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()
	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("erro ao pingar o MongoDB: %w", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MovieRepository{collection: collection}, nil
}

// LoadMovies lê o arquivo JSON indicado por filename e, apenas se a
// coleção ainda estiver vazia, insere todos os filmes no MongoDB. Isso
// evita duplicar os dados a cada reinicialização da aplicação. Retorna
// erro se o arquivo não puder ser lido/decodificado ou se a contagem/
// inserção no banco falhar.
func (r *MovieRepository) LoadMovies(ctx context.Context, filename string) error {
	countCtx, cancel := context.WithTimeout(ctx, operationTimeout)
	defer cancel()

	total, err := r.collection.CountDocuments(countCtx, bson.D{})
	if err != nil {
		return fmt.Errorf("erro ao contar filmes existentes: %w", err)
	}
	if total > 0 {
		fmt.Printf("Coleção já possui %d filmes, pulando carga inicial.\n", total)
		return nil
	}

	arquivoDeFilmes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var filmes []domain.Movie
	if err := json.Unmarshal(arquivoDeFilmes, &filmes); err != nil {
		return err
	}

	documentos := make([]any, len(filmes))
	for i, filme := range filmes {
		documentos[i] = fromDomain(filme)
	}

	insertCtx, cancelInsert := context.WithTimeout(ctx, operationTimeout*4)
	defer cancelInsert()
	if _, err := r.collection.InsertMany(insertCtx, documentos); err != nil {
		return fmt.Errorf("erro ao inserir filmes: %w", err)
	}

	fmt.Printf("Foram carregados %d filmes no MongoDB.\n", len(filmes))
	return nil
}

// GetAll retorna todos os filmes armazenados na coleção MongoDB. Em caso
// de erro de comunicação com o banco, loga o erro e retorna um slice
// vazio.
func (r *MovieRepository) GetAll() []domain.Movie {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println("[GetAll] erro ao buscar filmes:", err)
		return []domain.Movie{}
	}

	var documentos []movieDocument
	if err := cursor.All(ctx, &documentos); err != nil {
		fmt.Println("[GetAll] erro ao decodificar filmes:", err)
		return []domain.Movie{}
	}

	filmes := make([]domain.Movie, len(documentos))
	for i, doc := range documentos {
		filmes[i] = doc.toDomain()
	}
	return filmes
}

// GetByID busca, na coleção MongoDB, o filme cujo _id seja igual ao
// informado. Retorna o filme encontrado e true, ou um domain.Movie
// zerado e false caso o filme não exista ou ocorra erro de comunicação
// com o banco.
func (r *MovieRepository) GetByID(id int) (domain.Movie, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	var documento movieDocument
	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: id}}).Decode(&documento)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			fmt.Println("[GetByID] erro ao buscar filme:", err)
		}
		return domain.Movie{}, false
	}

	return documento.toDomain(), true
}

// Create atribui ao filme informado o próximo ID disponível (calculado
// por proximoID) e o insere na coleção MongoDB. Retorna o filme já com o
// ID definido, mesmo que a inserção falhe (o erro é apenas logado, para
// manter a mesma assinatura do adapter em memória).
func (r *MovieRepository) Create(filme domain.Movie) domain.Movie {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	filme.ID = r.proximoID(ctx)

	if _, err := r.collection.InsertOne(ctx, fromDomain(filme)); err != nil {
		fmt.Println("[Create] erro ao inserir filme:", err)
	}

	return filme
}

// proximoID calcula o próximo ID a ser usado em um novo filme, buscando
// o maior _id já existente na coleção e somando 1. Retorna 1 se a
// coleção estiver vazia ou se ocorrer erro de comunicação com o banco.
func (r *MovieRepository) proximoID(ctx context.Context) int {
	opts := options.FindOne().SetSort(bson.D{{Key: "_id", Value: -1}})

	var ultimo movieDocument
	err := r.collection.FindOne(ctx, bson.D{}, opts).Decode(&ultimo)
	if err != nil {
		return 1
	}

	return ultimo.ID + 1
}

// Update procura o filme com o ID informado e o substitui por filmeNovo
// (preservando o ID original). Retorna true se um filme com esse ID foi
// encontrado e atualizado, ou false caso contrário ou em caso de erro de
// comunicação com o banco.
func (r *MovieRepository) Update(id int, filmeNovo domain.Movie) bool {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	filmeNovo.ID = id
	resultado, err := r.collection.ReplaceOne(ctx, bson.D{{Key: "_id", Value: id}}, fromDomain(filmeNovo))
	if err != nil {
		fmt.Println("[Update] erro ao atualizar filme:", err)
		return false
	}

	return resultado.MatchedCount > 0
}

// Delete procura o filme com o ID informado e o remove da coleção
// MongoDB. Retorna true se o filme foi encontrado e removido, ou false
// caso contrário ou em caso de erro de comunicação com o banco.
func (r *MovieRepository) Delete(id int) bool {
	ctx, cancel := context.WithTimeout(context.Background(), operationTimeout)
	defer cancel()

	resultado, err := r.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		fmt.Println("[Delete] erro ao remover filme:", err)
		return false
	}

	return resultado.DeletedCount > 0
}
