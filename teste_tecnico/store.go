package main

import (
	"encoding/json"
	"os"
	"sync"
)

type MovieStore struct {
	mu     sync.Mutex
	movies []Movie
	nextID int
}

func NewMovieStore(path string) (*MovieStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var movies []Movie
	if err := json.Unmarshal(data, &movies); err != nil {
		return nil, err
	}

	maxID := 0
	for _, m := range movies {
		if m.ID > maxID {
			maxID = m.ID
		}
	}

	return &MovieStore{movies: movies, nextID: maxID + 1}, nil
}

func (s *MovieStore) GetAll() []Movie {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.movies
}

func (s *MovieStore) GetByID(id int) (Movie, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, m := range s.movies {
		if m.ID == id {
			return m, true
		}
	}
	return Movie{}, false
}

func (s *MovieStore) Create(m Movie) Movie {
	s.mu.Lock()
	defer s.mu.Unlock()
	m.ID = s.nextID
	s.nextID++
	s.movies = append(s.movies, m)
	return m
}

func (s *MovieStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, m := range s.movies {
		if m.ID == id {
			s.movies = append(s.movies[:i], s.movies[i+1:]...)
			return true
		}
	}
	return false
}
