package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"sipub_teste/api/internal/domain"
	"sipub_teste/api/internal/usecase"
)

// MovieHandler é o adapter que recebe requisições HTTP e chama o
// MovieService. Não conhece como os filmes são armazenados.
type MovieHandler struct {
	service *usecase.MovieService
}

// NewMovieHandler cria um MovieHandler associado ao MovieService informado,
// que será usado para atender todas as requisições HTTP de filmes.
func NewMovieHandler(service *usecase.MovieService) *MovieHandler {
	return &MovieHandler{service: service}
}

// GetMovies trata GET /movies. Busca todos os filmes através do
// MovieService e responde com um JSON contendo a lista completa.
func (h *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entramos na rota GetMovies...")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.service.ListAll())

}

// GetMovie trata GET /movies/{id}. Lê o ID da URL, valida que é um número
// inteiro e busca o filme correspondente através do MovieService.
// Responde 400 se o ID for inválido, 404 se o filme não existir, ou o
// filme em JSON em caso de sucesso.
func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	id_parametro := r.PathValue("id")
	fmt.Println("Entramos na rota GetByID... ID requisitado:", id_parametro)

	id, err := strconv.Atoi(id_parametro)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	filme, found := h.service.GetByID(id)
	if !found {
		fmt.Println("[GetMovie] - Filme não encontrado.")
		http.Error(w, "Filme não encontrado - Tente buscar outro ID.", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filme)
}

// PostCreateMovie trata POST /movies. Decodifica o corpo da requisição
// como um domain.Movie, delega a criação ao MovieService (que atribui o
// ID) e responde 201 com o filme criado em JSON. Responde 400 se o corpo
// não puder ser decodificado.
func (h *MovieHandler) PostCreateMovie(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entramos na rota PostCreateMovie...")

	var filme domain.Movie

	err := json.NewDecoder(r.Body).Decode(&filme)
	if err != nil {
		http.Error(w, "Estrutura inválida.", http.StatusBadRequest)
		return
	}

	filmeCriado := h.service.Create(filme)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(filmeCriado)
	fmt.Println("PostCreatMovie executado com sucesso")
}

// PutUpdateMovie trata PUT /movies/{id}. Lê o ID da URL e o novo filme do
// corpo da requisição, então delega a atualização ao MovieService.
// Responde 400 se o ID ou o corpo forem inválidos, 404 se o filme não
// existir, ou o filme atualizado em JSON em caso de sucesso.
func (h *MovieHandler) PutUpdateMovie(w http.ResponseWriter, r *http.Request) {
	idParametro := r.PathValue("id")
	fmt.Println("Entramos na rota PutUpdateMovie... ID requisitado:", idParametro)

	id, err := strconv.Atoi(idParametro)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	var filme domain.Movie
	err = json.NewDecoder(r.Body).Decode(&filme)
	if err != nil {
		http.Error(w, "Estrutura inválida.", http.StatusBadRequest)
		return
	}

	atualizado := h.service.Update(id, filme)
	if !atualizado {
		http.Error(w, "Filme não encontrado - Tente buscar outro ID.", http.StatusNotFound)
		return
	}

	filme.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filme)
}

// DeleteMovie trata DELETE /movies/{id}. Lê o ID da URL e delega a remoção
// ao MovieService. Responde 400 se o ID for inválido, 404 se o filme não
// existir, ou 204 (sem corpo) em caso de sucesso.
func (h *MovieHandler) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	idParametro := r.PathValue("id")
	fmt.Println("Entramos na rota DeleteMovie... ID requisitado:", idParametro)

	id, err := strconv.Atoi(idParametro)
	if err != nil {
		http.Error(w, "ID inválido.", http.StatusBadRequest)
		return
	}

	deletado := h.service.Delete(id)
	if !deletado {
		http.Error(w, "Filme não encontrado - Tente buscar outro ID.", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
