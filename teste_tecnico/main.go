package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func main() {
	store, err := NewMovieStore("data/movies.json")
	if err != nil {
		log.Fatalf("failed to load movies: %v", err)
	}
	log.Printf("loaded %d movies", len(store.GetAll()))

	mux := http.NewServeMux()
	mux.HandleFunc("GET /movies", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, store.GetAll())
	})

	mux.HandleFunc("GET /movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		movie, found := store.GetByID(id)
		if !found {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		writeJSON(w, http.StatusOK, movie)
	})

	mux.HandleFunc("POST /movies", func(w http.ResponseWriter, r *http.Request) {
		var movie Movie
		if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if movie.Title == "" {
			http.Error(w, "title is required", http.StatusBadRequest)
			return
		}
		created := store.Create(movie)
		writeJSON(w, http.StatusCreated, created)
	})

	mux.HandleFunc("DELETE /movies/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, "invalid id", http.StatusBadRequest)
			return
		}
		if !store.Delete(id) {
			http.Error(w, "movie not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
