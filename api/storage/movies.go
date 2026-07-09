package storage

import (
	"encoding/json"
	"os"

	"sipub_teste/api/models"
)


// Variável global para guardar os filmes em memória.
var movies []models.Movie

func LoadMovies(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, &movies)
}

func GetAll() []models.Movie {
	return movies
}
