package domain

import "github.com/jackc/pgx/v5/pgtype"

type Sex string

const (
	M Sex = "M"
	F     = "F"
)

type Actor struct {
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Sex       Sex         `json:"sex"`
	Birthdate pgtype.Date `json:"birthdate"`
	Films     []Film      `json:"films,omitempty"`
}

type ActorRepository interface {
	Insert(actor Actor) (int, error)
	Delete(id int) error
	Update(actor Actor) (Actor, error)
	SelectById(id int) (Actor, error)
	SelectAll() ([]Actor, error)
}

type ActorUsecase interface {
	Add(actor Actor) (int, error)
	Remove(id int) error
	Modify(actor Actor) (Actor, error)
	GetAll() ([]Actor, error)
}
