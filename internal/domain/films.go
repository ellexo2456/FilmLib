package domain

import "github.com/jackc/pgx/v5/pgtype"

type Film struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ReleaseDate pgtype.Date `json:"releaseDate"`
	Rating      float64     `json:"rating"`
}

type FilmsRepository interface {
	Insert(film Film) (int, error)
}

type FilmsUsecase interface {
	Add(film Film) (int, error)
}
