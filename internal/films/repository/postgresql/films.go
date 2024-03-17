package postgres

import (
	"context"
	"errors"
	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
	"github.com/jackc/pgx/v5"
	"math"
)

const insertQuery = `
	INSERT INTO film (title, description, release_date, rating)
	VALUES 
		($1, $2, $3, $4)
	RETURNING id
`

const selectAllQuery = `
	SELECT id, title, description, release_date, rating 
	FROM film
`

const searchQuery = `
	SELECT f.id, f.title, f.description, f.release_date, f.rating
	FROM film f
         LEFT JOIN film_actor fa ON fa.film_id = f.id
         LEFT JOIN actor a ON a.id = fa.actor_id
	WHERE (f.title ILIKE '%hero%'
    OR a.name ILIKE '%hero%')
`

const deleteQuery = `
	DELETE FROM film
	WHERE id = $1
`

const updateQuery = `
	UPDATE film
	SET title = $1, description = $2, release_date = $3, rating = $4 
	WHERE id = $5 
	RETURNING id, title, description, release_date, rating
`

const selectByIdQuery = `
	SELECT id, title, description, release_date, rating
	FROM film
	WHERE id = $1
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
	tx, err := r.db.Begin(r.ctx)
	if err != nil {
		return 0, domain.ErrInternalServerError
	}
	defer tx.Rollback(r.ctx)

	row := tx.QueryRow(r.ctx, insertQuery, film.Title, film.Description, film.ReleaseDate, film.Rating)

	var id int
	err = row.Scan(
		&id,
	)
	if err != nil {
		logs.LogError(logs.Logger, "films/postgres", "Insert", err, err.Error())
		return 0, err
	}

	var rows [][]interface{}
	for _, a := range film.Actors {
		rows = append(rows, []interface{}{id, a.ID})
	}

	rowsCount, err := tx.CopyFrom(
		r.ctx,
		pgx.Identifier{"film_actor"},
		[]string{"film_id", "actor_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		logs.LogError(logs.Logger, "films/postgres", "Insert", err, err.Error())
		return 0, err
	}
	if rowsCount == 0 {
		logs.LogError(logs.Logger, "films/postgres", "Insert", domain.ErrInternalServerError, "can`t insert rows to film_actor")
		return 0, domain.ErrInternalServerError
	}

	err = tx.Commit(r.ctx)
	if err != nil {
		logs.LogError(logs.Logger, "films/postgres", "Insert", domain.ErrInternalServerError, "can`t commit changes")
		return 0, err
	}

	return id, nil
}

func (r *filmsPostgresqlRepository) SelectAll() ([]domain.Film, error) {
	rows, err := r.db.Query(r.ctx, selectAllQuery)
	if err != nil {
		logs.LogError(logs.Logger, "films/postgres", "SelectAll", err, err.Error())
		return nil, err
	}
	defer rows.Close()

	var films []domain.Film
	var film domain.Film

	for rows.Next() {
		err = rows.Scan(
			&film.ID,
			&film.Title,
			&film.Description,
			&film.ReleaseDate,
			&film.Rating,
		)

		if err != nil {
			return nil, err
		}

		film.Rating = math.Trunc(film.Rating*10) / 10
		films = append(films, film)
	}

	if len(films) == 0 {
		return []domain.Film{}, nil
	}
	return films, nil
}

func (r *filmsPostgresqlRepository) Search(searchStr string) ([]domain.Film, error) {
	searchStr = "%" + searchStr + "%"
	rows, err := r.db.Query(r.ctx, searchQuery, searchStr)
	if err != nil {
		logs.LogError(logs.Logger, "films/postgres", "Search", err, err.Error())
		return nil, err
	}
	defer rows.Close()

	var films []domain.Film
	var film domain.Film
	for rows.Next() {
		err = rows.Scan(
			&film.ID,
			&film.Title,
			&film.Description,
			&film.ReleaseDate,
			&film.Rating,
		)

		if err != nil {
			logs.LogError(logs.Logger, "films/postgres", "Search", err, err.Error())
			return nil, err
		}

		film.Rating = math.Trunc(film.Rating*10) / 10
		films = append(films, film)
	}

	if len(films) == 0 {
		return nil, domain.ErrNotFound
	}
	return films, nil
}

func (r *filmsPostgresqlRepository) Delete(id int) error {
	res, err := r.db.Exec(r.ctx, deleteQuery, id)
	if err != nil {
		logs.LogError(logs.Logger, "films/postgres", "Delete", err, err.Error())
		return err
	}

	if res.RowsAffected() == 0 {
		logs.LogError(logs.Logger, "films/postgres", "Delete", domain.ErrOutOfRange, domain.ErrOutOfRange.Error())
		return domain.ErrOutOfRange
	}

	return nil
}

func (r *filmsPostgresqlRepository) Update(film domain.Film) (domain.Film, error) {
	row := r.db.QueryRow(r.ctx, updateQuery, film.Title, film.Description, film.ReleaseDate, film.Rating, film.ID)

	err := row.Scan(
		&film.ID,
		&film.Title,
		&film.Description,
		&film.ReleaseDate,
		&film.Rating,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		logs.LogError(logs.Logger, "films/postgres", "Update", err, err.Error())
		return domain.Film{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "films/postgres", "Update", err, err.Error())
		return domain.Film{}, err
	}

	return film, nil
}

func (r *filmsPostgresqlRepository) SelectById(id int) (domain.Film, error) {
	row := r.db.QueryRow(r.ctx, selectByIdQuery, id)

	var film domain.Film
	err := row.Scan(
		&film.ID,
		&film.Title,
		&film.Description,
		&film.ReleaseDate,
		&film.Rating,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		logs.LogError(logs.Logger, "films/postgres", "SelectById", err, err.Error())
		return domain.Film{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "films/postgres", "SelectById", err, err.Error())
		return domain.Film{}, err
	}

	return film, nil
}
