package usecase

import "sipub_teste/api/internal/domain"

// MovieService orquestra as regras de negócio dos filmes. Depende só da
// interface domain.MovieRepository, nunca de uma implementação concreta.
type MovieService struct {
	repo domain.MovieRepository
}

// NewMovieService cria um MovieService que delega toda a persistência ao
// domain.MovieRepository informado.
func NewMovieService(repo domain.MovieRepository) *MovieService {
	return &MovieService{repo: repo}
}

// ListAll retorna todos os filmes cadastrados, delegando a busca ao
// repositório configurado.
func (s *MovieService) ListAll() []domain.Movie {
	return s.repo.GetAll()
}

// GetByID busca um filme pelo seu ID, delegando ao repositório
// configurado. Retorna o filme e true se encontrado, ou um domain.Movie
// zerado e false caso contrário.
func (s *MovieService) GetByID(id int) (domain.Movie, bool) {
	return s.repo.GetByID(id)
}

// Create registra um novo filme, delegando a criação (e atribuição de ID)
// ao repositório configurado, e retorna o filme já criado.
func (s *MovieService) Create(filme domain.Movie) domain.Movie {
	return s.repo.Create(filme)
}

// Update substitui os dados do filme de ID informado pelos dados de
// filmeNovo, delegando ao repositório configurado. Retorna true se o
// filme existia e foi atualizado, ou false caso contrário.
func (s *MovieService) Update(id int, filmeNovo domain.Movie) bool {
	return s.repo.Update(id, filmeNovo)
}

// Delete remove o filme de ID informado, delegando ao repositório
// configurado. Retorna true se o filme existia e foi removido, ou false
// caso contrário.
func (s *MovieService) Delete(id int) bool {
	return s.repo.Delete(id)
}
