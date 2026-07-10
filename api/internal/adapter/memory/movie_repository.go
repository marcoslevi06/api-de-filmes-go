package memory

import (
	"encoding/json"
	"fmt"
	"os"

	"sipub_teste/api/internal/domain"
)

// MovieRepository é o adapter que guarda os filmes em memória (um slice).
// Implementa a porta domain.MovieRepository.
type MovieRepository struct {
	movies []domain.Movie
}

// NewMovieRepository cria um MovieRepository vazio, pronto para ter os
// filmes carregados (via LoadMovies) ou inseridos manualmente.
func NewMovieRepository() *MovieRepository {
	return &MovieRepository{}
}

// LoadMovies lê o arquivo JSON indicado por filename, faz o unmarshal do
// conteúdo para o slice de domain.Movie do repositório e imprime um
// resumo do carregamento no console. Retorna erro se o arquivo não puder
// ser lido ou se o conteúdo não for um JSON válido.
func (r *MovieRepository) LoadMovies(filename string) error {

	fmt.Println("Entramos na loadMovies...")

	// Lê o arquivo JSON
	arquivo_de_filmes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	fmt.Printf("Tipo de dados: %T\n", arquivo_de_filmes)

	// Converte o JSON para o slice de Movie.
	err = json.Unmarshal(arquivo_de_filmes, &r.movies)
	if err != nil {
		return err
	}

	//  Percorrendo o slice de filmes.
	fmt.Println("===== Listando alguns filmes ====")
	for indice, filme := range r.movies {
		if indice >= 90 && indice <= 100 {
			fmt.Printf("-> filme_id: %d| Ano: %v | Nome: %s \n", filme.ID, filme.Year, filme.Title)
		}
	}
	fmt.Println("===== Fim da listagem ====")

	fmt.Printf("Foram carregados %d filmes.\n", len(r.movies))
	return nil
}

// GetAll retorna todos os filmes atualmente armazenados em memória.
func (r *MovieRepository) GetAll() []domain.Movie {
	// Método que é responsável por listar todos os filmes existentes.
	fmt.Printf("Entramos no métood GetAll - Temos %d filmes", len(r.movies))
	return r.movies
}

// GetByID procura, no slice em memória, o filme cujo ID seja igual ao
// informado. Retorna o filme encontrado e true, ou um domain.Movie zerado
// e false caso nenhum filme com esse ID exista.
func (r *MovieRepository) GetByID(id int) (domain.Movie, bool) {
	// Método responsável por retornar um filme específico.
	fmt.Println("Entramos em GetByID")

	var filme_requisitado domain.Movie
	var encontrado bool

	for _, filme := range r.movies {
		if id == filme.ID {
			filme_requisitado = filme
			encontrado = true
			fmt.Printf("Escontramos o filme pedido: %d", filme.ID)
			return filme_requisitado, encontrado
		}
	}

	return domain.Movie{}, false
}

// Create atribui ao filme informado o próximo ID disponível (calculado
// por descobreProximoId), adiciona-o ao slice em memória e retorna o
// filme já com o ID definido.
func (r *MovieRepository) Create(filme domain.Movie) domain.Movie {
	novo_id := r.descobreProximoId()
	filme.ID = novo_id
	r.movies = append(r.movies, filme)
	return filme
}

// Update procura o filme com o ID informado e o substitui por filmeNovo
// (preservando o ID original). Retorna true se um filme com esse ID foi
// encontrado e atualizado, ou false caso contrário.
func (r *MovieRepository) Update(id int, filmeNovo domain.Movie) bool {
	fmt.Println("\nEntramos no método PUT...")
	for indice, filme := range r.movies {
		if filme.ID == id {
			filmeNovo.ID = id
			r.movies[indice] = filmeNovo
			return true
		}
	}
	return false
}

// descobreProximoId calcula o próximo ID a ser usado em um novo filme,
// retornando o maior ID já existente no slice em memória mais 1 (ou 1 se
// não houver nenhum filme cadastrado).
func (r *MovieRepository) descobreProximoId() int {
	// Função criada para atribuir um novo id enquanto testo em memória.

	maiorID := 0
	fmt.Println("Buscando próximo id...")

	for _, filme := range r.movies {
		if filme.ID > maiorID {
			maiorID = filme.ID
		}
	}

	return maiorID + 1
}

// Delete procura o filme com o ID informado e o remove do slice em
// memória. Retorna true se o filme foi encontrado e removido, ou false
// caso nenhum filme com esse ID exista.
func (r *MovieRepository) Delete(id int) bool {
	fmt.Println("Entramos no método DELETE...")
	for indice, filme := range r.movies {
		if filme.ID == id {
			r.movies = append(r.movies[:indice], r.movies[indice+1:]...)
			return true
		}
	}
	return false
}
