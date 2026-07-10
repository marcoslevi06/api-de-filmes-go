package handlers

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"

	"sipub_teste/api/storage"
	"sipub_teste/api/models"
)

func GetMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[GET] Entramos em Handler.GetMovies")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storage.GetAll())

}

func GetMovie(w http.ResponseWriter, r *http.Request) {

	id_parametro := r.PathValue("id")

	id, err := strconv.Atoi(id_parametro)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	filme, found := storage.GetByID(id)
	if !found {
		fmt.Println("[GetMovie] - Filme não encontrado.")
		http.Error(w, "Filme não encontrado - Tente buscar outro ID.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filme)
}


func PostCreateMovie(w http.ResponseWriter, r *http.Request) {
	var filme models.Movie

	err := json.NewDecoder(r.Body).Decode(&filme)
	if err != nil {
		http.Error(w, "Estrutura inválida.", http.StatusBadRequest)
		return
	}

	filmeCriado := storage.Create(filme)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(filmeCriado)
}