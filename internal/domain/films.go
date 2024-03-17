package domain

import "github.com/jackc/pgx/v5/pgtype"

type SortDirection string

const (
	Asc  SortDirection = "Asc"
	Desc SortDirection = "Desc"
)

const (
	TitleParam       = "sortTitle"
	ReleaseDateParam = "sortReleaseDate"
	SearchParam      = "searchStr"
)

type Film struct {
	ID          int         `json:"id"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	ReleaseDate pgtype.Date `json:"releaseDate"`
	Rating      float64     `json:"rating"`
	Actors      []Actor     `json:"actors,omitempty"`
}

type FilmsRepository interface {
	Insert(film Film) (int, error)
	SelectAll() ([]Film, error)
	Search(searchStr string) ([]Film, error)
	Delete(id int) error
	Update(film Film) (Film, error)
	SelectById(id int) (Film, error)
}

type FilmsUsecase interface {
	Add(film Film) (int, error)
	GetAll(title, releaseDate SortDirection) ([]Film, error)
	Search(searchStr string) ([]Film, error)
	Remove(id int) error
	Modify(film Film) (Film, error)
}
