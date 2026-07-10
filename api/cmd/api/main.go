package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	httpadapter "sipub_teste/api/internal/adapter/http"
	mongoadapter "sipub_teste/api/internal/adapter/mongo"
	"sipub_teste/api/internal/usecase"
)

// getEnv retorna o valor da variável de ambiente chave, ou padrao caso ela
// não esteja definida.
func getEnv(chave, padrao string) string {
	if valor := os.Getenv(chave); valor != "" {
		return valor
	}
	return padrao
}

// main é o ponto de entrada da aplicação. Monta as dependências
// (repositório MongoDB, service e handler HTTP), carrega os filmes a
// partir de "movies.json" (apenas na primeira execução), registra as
// rotas da API de filmes e sobe o servidor HTTP na porta 8080.
func main() {

	ctx := context.Background()

	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "sipub")
	mongoCollection := getEnv("MONGO_COLLECTION", "movies")

	repo, erro := mongoadapter.NewMovieRepository(ctx, mongoURI, mongoDB, mongoCollection)
	if erro != nil {
		log.Fatal("Erro ao conectar no MongoDB: ", erro)
	}

	if erro := repo.LoadMovies(ctx, "movies.json"); erro != nil {
		log.Fatal("Erro ao carregar filmes: ", erro)
	}

	fmt.Printf("Carregando %d filmes\n", len(repo.GetAll()))

	service := usecase.NewMovieService(repo)
	handler := httpadapter.NewMovieHandler(service)

	http.HandleFunc("GET /movies", handler.GetMovies)
	http.HandleFunc("GET /movies/{id}", handler.GetMovie)
	http.HandleFunc("POST /movies", handler.PostCreateMovie)
	http.HandleFunc("PUT /movies/{id}", handler.PutUpdateMovie)
	http.HandleFunc("DELETE /movies/{id}", handler.DeleteMovie)

	log.Fatal(http.ListenAndServe(":8080", nil))

}
