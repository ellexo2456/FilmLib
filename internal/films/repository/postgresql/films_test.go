package postgres_test

import (
	"context"
	"errors"
	"github.com/ellexo2456/FilmLib/internal/domain"
	postgres "github.com/ellexo2456/FilmLib/internal/films/repository/postgresql"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
	"testing"
)

const insertQuery = `
	INSERT INTO film
`

const selectAllQuery = `
	SELECT id, title, description, release_date, rating 
	FROM film
`

const deleteQuery = `
	DELETE FROM film
	WHERE id = \$1
`

func TestInsertIntoFilm(t *testing.T) {
	tests := []struct {
		name         string
		getFilm      func() domain.Film
		getInsertErr func() error
		getCopyErr   func() error
	}{
		{
			name: "GoodCase/Common",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Film{
					ID:          1,
					Title:       "Matrix Reloaded",
					Description: "Description",
					Rating:      8.5,
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: 1}, {ID: 2}},
				}
			},
		},
		{
			name: "BadCase/FutureDate",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2100-01-01")

				return domain.Film{
					Title:       "Matrix Reloaded",
					Description: "Description",
					Rating:      8.5,
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: 1}, {ID: 2}},
				}
			},
			getInsertErr: func() error {
				var err pgconn.PgError
				err.Code = domain.DateOutOfRangeErrCode
				return &err
			},
		},
		{
			name: "BadCase/EmptyActors",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Film{
					Title:       "Matrix Reloaded",
					Description: "Description",
					Rating:      8.5,
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: -1}, {ID: -2}},
				}
			},
			getCopyErr: func() error {
				return nil
			},
		},
		{
			name: "BadCase/ActorsWithNegativeIds",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Film{
					Title:       "Matrix Reloaded",
					Description: "Description",
					Rating:      8.5,
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: 1}, {ID: 2}},
				}
			},
			getCopyErr: func() error {
				return errors.New("some db err")
			},
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	r := postgres.NewFilmsPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			film := test.getFilm()

			mockDB.ExpectBegin()

			row := mockDB.NewRows([]string{"id"}).
				AddRow(film.ID)
			eq := mockDB.ExpectQuery(insertQuery).
				WithArgs(film.Title, film.Description, film.ReleaseDate, film.Rating)

			if test.getInsertErr == nil {
				eq.WillReturnRows(row)
			} else {
				eq.WillReturnError(test.getInsertErr())
			}

			cp := mockDB.ExpectCopyFrom(
				pgx.Identifier{"film_actor"},
				[]string{"film_id", "actor_id"})

			if test.getCopyErr == nil || len(film.Actors) == 0 {
				cp.WillReturnResult(int64(len(film.Actors)))
			} else {
				cp.WillReturnError(test.getCopyErr())
			}

			mockDB.ExpectCommit()

			id, err := r.Insert(film)
			if test.getCopyErr == nil && test.getInsertErr == nil {
				require.Nil(t, err)
				require.Equal(t, id, film.ID)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
		})
	}
}

func TestSelectAll(t *testing.T) {
	tests := []struct {
		name     string
		getFilms func() []domain.Film
		good     bool
	}{
		{
			name: "GoodCase/Common",
			getFilms: func() []domain.Film {
				var d1, d2 pgtype.Date
				d1.Scan("2000-01-01")
				d2.Scan("2002-02-01")

				return []domain.Film{
					{
						ID:          1,
						Title:       "some t1",
						Description: "desc",
						ReleaseDate: d1,
						Rating:      9.5,
					},
					{
						ID:          1,
						Title:       "some t2",
						Description: "desc",
						ReleaseDate: d2,
						Rating:      6.4,
					},
				}
			},
			good: true,
		},
		{
			name: "BadCase/EmptyFilms",
			getFilms: func() []domain.Film {
				return []domain.Film{}
			},
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	r := postgres.NewFilmsPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			expectedFilms := test.getFilms()
			rows := mockDB.NewRows([]string{"id", "title", "description", "release_date",
				"rating"})

			if len(expectedFilms) != 0 {
				rows.AddRow(expectedFilms[0].ID, expectedFilms[0].Title, expectedFilms[0].Description,
					expectedFilms[0].ReleaseDate, expectedFilms[0].Rating).
					AddRow(expectedFilms[1].ID, expectedFilms[1].Title, expectedFilms[1].Description,
						expectedFilms[1].ReleaseDate, expectedFilms[1].Rating)
			}
			eq := mockDB.ExpectQuery(selectAllQuery)
			if test.good {
				eq.WillReturnRows(rows)
			}
			if expectedFilms[0].ID == 0 {
				eq.WillReturnError(errors.New("some db err"))
			}

			films, err := r.SelectAll()
			if test.good {
				require.Equal(t, expectedFilms, films)
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}

func TestDelete(t *testing.T) {
	tests := []struct {
		name   string
		filmID int
		result pgconn.CommandTag
		good   bool
	}{
		{
			name:   "GoodCase/Common",
			filmID: 1,
			result: pgxmock.NewResult("", 1),
			good:   true,
		},
		{
			name:   "BadCase/NegativeVideoID",
			result: pgxmock.NewResult("", 0),
			filmID: -1,
		},
		{
			name:   "BadCase/OutOfRangeVideoID",
			result: pgxmock.NewResult("", 0),
			filmID: 123456789,
		},
		{
			name:   "BadCase/ZeroVideoID",
			result: pgxmock.NewResult("", 0),
			filmID: 0,
		},
	}

	mockDB, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	defer mockDB.Close()

	r := postgres.NewFilmsPostgresqlRepository(mockDB, context.Background())

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockDB.ExpectExec(deleteQuery).
				WithArgs(test.filmID).
				WillReturnResult(test.result)

			err = r.Delete(test.filmID)
			if test.good {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
			}

			err = mockDB.ExpectationsWereMet()
			require.Nil(t, err)
		})
	}
}
