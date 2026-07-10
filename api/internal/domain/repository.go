package domain

// MovieRepository é a porta (interface) que o domínio expõe para quem
// quiser guardar/buscar filmes. Quem implementa essa interface é um
// adapter (memória, Mongo, etc) — o domínio nunca sabe qual.
type MovieRepository interface {
	GetAll() []Movie
	GetByID(id int) (Movie, bool)
	Create(filme Movie) Movie
	Update(id int, filmeNovo Movie) bool
	Delete(id int) bool
}
