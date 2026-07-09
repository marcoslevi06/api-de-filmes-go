package handlers

import (
	"encoding/json"
	"net/http"

	"sipub_teste/api/storage"
)

func GetMovies(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(storage.GetAll())

}