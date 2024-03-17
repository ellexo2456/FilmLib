package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"

	"github.com/ellexo2456/FilmLib/internal/domain"
	logs "github.com/ellexo2456/FilmLib/internal/logger"
)

const insertQuery = `
	INSERT INTO actor (name, sex, birthdate)
	VALUES 
		($1, $2, $3)
	RETURNING id
`

const deleteQuery = `
	DELETE FROM actor
	WHERE id = $1
`

const updateQuery = `
	UPDATE actor
	SET name = $1, sex = $2, birthdate = $3 
	WHERE id = $5 
	RETURNING id, name, sex, birthdate
`

const selectByIdQuery = `
	SELECT id, name, sex, birthdate  
	FROM actor
`

const selectAllQuery = `
	SELECT a.id,
       a.name,
       a.sex,
       a.birthdate,
       COALESCE(f.id, 0),
       COALESCE(f.title, ""),
       COALESCE(f.description, ""),
       COALESCE(f.release_date, ""),
       COALESCE(f.rating, 0)
	FROM actor a
         LEFT JOIN film_actor fa ON a.id = fa.actor_id
         LEFT JOIN film f ON f.id = fa.film_id
    	
`

type actorPostgresqlRepository struct {
	db  domain.PgxPoolIface
	ctx context.Context
}

func NewActorPostgresqlRepository(pool domain.PgxPoolIface, ctx context.Context) domain.ActorRepository {
	return &actorPostgresqlRepository{
		db:  pool,
		ctx: ctx,
	}
}

func (r *actorPostgresqlRepository) Insert(actor domain.Actor) (int, error) {
	row := r.db.QueryRow(r.ctx, insertQuery, actor.Name, actor.Sex, actor.Birthdate)

	var id int
	err := row.Scan(
		&id,
	)

	if err != nil {
		logs.LogError(logs.Logger, "actor/postgres", "Insert", err, err.Error())
		return 0, err
	}
	return id, nil
}

func (r *actorPostgresqlRepository) Delete(id int) error {
	res, err := r.db.Exec(r.ctx, deleteQuery, id)
	if err != nil {
		logs.LogError(logs.Logger, "actor/postgres", "Delete", err, err.Error())
		return err
	}

	if res.RowsAffected() == 0 {
		logs.LogError(logs.Logger, "actor/postgres", "Delete", domain.ErrOutOfRange, domain.ErrOutOfRange.Error())
		return domain.ErrOutOfRange
	}

	return nil
}

func (r *actorPostgresqlRepository) Update(actor domain.Actor) (domain.Actor, error) {
	row := r.db.QueryRow(r.ctx, updateQuery, actor.Name, actor.Sex, actor.Birthdate)

	err := row.Scan(
		&actor.ID,
		&actor.Name,
		&actor.Sex,
		&actor.Birthdate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		logs.LogError(logs.Logger, "actor/postgres", "Update", err, err.Error())
		return domain.Actor{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "actor/postgres", "Update", err, err.Error())
		return domain.Actor{}, err
	}

	return actor, nil
}

func (r *actorPostgresqlRepository) SelectById(id int) (domain.Actor, error) {
	row := r.db.QueryRow(r.ctx, selectByIdQuery, id)

	var actor domain.Actor
	err := row.Scan(
		&actor.ID,
		&actor.Name,
		&actor.Sex,
		&actor.Birthdate,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		logs.LogError(logs.Logger, "actor/postgres", "SelectById", err, err.Error())
		return domain.Actor{}, domain.ErrNotFound
	}
	if err != nil {
		logs.LogError(logs.Logger, "actor/postgres", "SelectById", err, err.Error())
		return domain.Actor{}, err
	}

	return actor, nil
}

func (r *actorPostgresqlRepository) SelectAll() ([]domain.Actor, error) {
	rows, err := r.db.Query(r.ctx, selectAllQuery)
	if err != nil {
		logs.LogError(logs.Logger, "actor/postgres", "SelectAll", err, err.Error())
		return nil, err
	}
	defer rows.Close()

	var actors []domain.Actor
	var actor domain.Actor
	var film domain.Film
	prevActorID := 0
	if rows.Next() {
		err = rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.Sex,
			&actor.Birthdate,
			&film.ID,
			&film.Title,
			&film.Description,
			&film.ReleaseDate,
			&film.Rating,
		)

		actors = append(actors, actor)
		if film.ID != 0 {
			actors[0].Films = append(actors[0].Films, film)
		}
		prevActorID = actor.ID
	}

	for rows.Next() {
		err = rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.Sex,
			&actor.Birthdate,
			&film.ID,
			&film.Title,
			&film.Description,
			&film.ReleaseDate,
			&film.Rating,
		)
		if err != nil {
			return nil, err
		}

		if actor.ID != prevActorID {
			actors = append(actors, actor)
			prevActorID = actor.ID
		}
		if film.ID != 0 {
			actors[len(actors)-1].Films = append(actors[len(actors)-1].Films, film)
		}
	}

	if len(actors) == 0 {
		return nil, domain.ErrNotFound
	}

	return actors, nil
}
