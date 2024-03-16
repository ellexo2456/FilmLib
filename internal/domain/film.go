package domain

import "github.com/jackc/pgx/v5/pgtype"

type Film struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ReleaseDate pgtype.Date `json:"releaseDate"`
	Rating      float64     `json:"rating"`
}

type FilmRepository interface {
	Insert(film Film) (int, error)
}

type FilmUsecase interface {
	Add(film Film) (int, error)
}
