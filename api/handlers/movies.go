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
	fmt.Println("PostCreatMovie executado com sucesso")
}


func PutUpdateMovie(w http.ResponseWriter, r *http.Request) {
	idParametro := r.PathValue("id")

	id, err := strconv.Atoi(idParametro)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	var filme models.Movie
	err = json.NewDecoder(r.Body).Decode(&filme)
	if err != nil {
		http.Error(w, "Estrutura inválida.", http.StatusBadRequest)
		return
	}

	atualizado := storage.Update(id, filme)
	if !atualizado {
		http.Error(w, "Filme não encontrado - Tente buscar outro ID.", http.StatusNotFound)
		return
	}

	filme.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filme)
}

func DeleteMovie(w http.ResponseWriter, r *http.Request) {
	idParametro := r.PathValue("id")

	id, err := strconv.Atoi(idParametro)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	deletado := storage.Delete(id)
	if !deletado {
		http.Error(w, "Filme não encontrado - Tente buscar outro ID.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}