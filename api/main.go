// Manda o .json() dos filmes que você tem em mente (pode ser um exemplo com 
// 2-3 filmes só, já dá pra gente ver a estrutura) 
// e a partir dele eu te ajudo a:
// 
// Criar a struct Movie em Go correspondente
// Montar um CRUD simples:
// 
// GET /movies — listar todos
// GET /movies/{id} — buscar um
// POST /movies — criar
// PUT /movies/{id} — atualizar
// DELETE /movies/{id} — deletar
// 
// 
// Por enquanto guardamos tudo em memória (um slice), sem banco de dados — 
// assim você foca em entender o fluxo HTTP + Go puro antes de 
// complicar com persistência.

package main

import (
	"fmt"
	"log"
	"net/http"

	"sipub_teste/api/handlers"
	"sipub_teste/api/storage"
)

func main() {

	if erro := storage.LoadMovies("movies.json"); erro != nil {
		log.Fatal("Erro ao carregar filmes: ", erro)
	}

	fmt.Printf("Carregando %d filmes\n", len(storage.GetAll()))
	http.HandleFunc("GET /movies", handlers.GetMovies)

	log.Fatal(http.ListenAndServe(":8080", nil))

}