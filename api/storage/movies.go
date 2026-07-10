package storage

import (
	"encoding/json"
	"os"
	"fmt"

	"sipub_teste/api/models"
)


// Variável global para guardar os filmes em memória.
var movies []models.Movie


func LoadMovies(filename string) error {
	
	fmt.Println("Entramos na loadMovies...")

	// Lê o arquivo JSON
	arquivo_de_filmes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	fmt.Printf("Tipo de dados: %T\n", arquivo_de_filmes)
	

	// Converte o JSON para o slice de Movie.
	err = json.Unmarshal(arquivo_de_filmes, &movies)
	if err != nil {
		return err
	}
	
	//  Percorrendo o slice de filmes.
	fmt.Println("===== Listando alguns filmes ====")
	for indice, filme := range movies {
		if indice >= 90 && indice <= 100 {
			fmt.Printf("-> filme_id: %d| Ano: %v | Nome: %s \n", filme.ID, filme.Year, filme.Title)
		}
	}
	fmt.Println("===== Fim da listagem ====")
	
	fmt.Printf("Foram carregados %d filmes.\n", len(movies))
	return nil
}


func GetAll() []models.Movie {
	// Método que é responsável por listar todos os filmes existentes.
	fmt.Printf("Entramos no métood GetAll - Temos %d filmes", len(movies))
	return movies
}


func GetByID(id int) (models.Movie, bool) {
	// Método responsável por retornar um filme específico.
	fmt.Println("Entramos em GetByID")

	var filme_requisitado models.Movie
	var encontrado bool

	for _, filme := range movies {
		if id == filme.ID {
			filme_requisitado = filme
			encontrado = true
			fmt.Printf("Escontramos o filme pedido: %d", filme.ID)
			return filme_requisitado, encontrado
		}
	}

	return models.Movie{}, false
}


func Create(filme models.Movie) models.Movie {
	novo_id := descobreProximoId()
	filme.ID = novo_id
	movies = append(movies, filme)  
	return filme
}


func Update(id int, filmeNovo models.Movie) bool {
	fmt.Println("\nEntramos no método PUT...")
	for indice, filme := range movies {
		if filme.ID == id {
			filmeNovo.ID = id 
			movies[indice] = filmeNovo
			return true
		}
	}
	return false
}


func descobreProximoId() int {
	// Função criada para atribuir um novo id enquanto testo em memória.

	maiorID := 0
	fmt.Println("Buscando próximo id...")

	for _, filme := range movies {
		if filme.ID > maiorID {
			maiorID = filme.ID
		}
	}

	return maiorID + 1
}


func Delete(id int) bool {
	fmt.Println("Entramos no método DELETE...")
	for indice, filme := range movies {
		if filme.ID == id {
			movies = append(movies[:indice], movies[indice+1:]...)
			return true
		}
	}
	return false
}