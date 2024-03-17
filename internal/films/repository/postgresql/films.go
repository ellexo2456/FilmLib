package postgres

import (
	"context"

	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
)

const insertQuery = `
	INSERT INTO film (title, description, release_date, rating)
	VALUES 
		($1, $2, $3, $4)
	RETURNING id
`

type filmsPostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewFilmsPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.FilmsRepository {
	return &filmsPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *filmsPostgresqlRepository) Insert(film domain.Film) (int, error) {
	row := r.db.QueryRow(r.ctx, insertQuery, film.Title, film.Description, film.ReleaseDate, film.Rating)

	var id int
	err := row.Scan(
		&id,
	)

	if err != nil {
		logs.LogError(logs.Logger, "films/postgres", "Insert", err, err.Error())
		return 0, err
	}
	return id, nil
}
