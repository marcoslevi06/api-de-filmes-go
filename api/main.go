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
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Varuável global para guardar os filmes em memória.
var movies []Movie

func main(){
	if erro := loadMovies("teste_tecnico/movies.json"); erro != nil {
		log.Fatal("Erro ao carregar filmes:", erro)
	}

	fmt.Printf("Carregando %d filmes \n", len(movies))

	log.Fatal(http.ListenAndServe(":8080"), nil)
}


type Movie struct {
    ID    int    `json:"id"`
    Title string `json:"title"`
    Year  string `json:"year"`
}

func loadMovies(filename string) error {
	file, erro := os.ReadFile(filename)
	if erro != nil {
		return erro
	}

	return json.Unmarshal(file, &movies)
}
