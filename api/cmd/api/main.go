package main

import (
	"fmt"
	"log"
	"net/http"

	httpadapter "sipub_teste/api/internal/adapter/http"
	"sipub_teste/api/internal/adapter/memory"
	"sipub_teste/api/internal/usecase"
)

// main é o ponto de entrada da aplicação. Monta as dependências (repositório
// em memória, service e handler HTTP), carrega os filmes a partir de
// "movies.json", registra as rotas da API de filmes e sobe o servidor HTTP
// na porta 8080.
func main() {

	repo := memory.NewMovieRepository()

	if erro := repo.LoadMovies("movies.json"); erro != nil {
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
