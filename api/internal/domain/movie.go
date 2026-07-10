package domain

// Movie representa um filme do domínio, com seu identificador, título e
// ano de lançamento.
type Movie struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Year  string `json:"year"`
}
