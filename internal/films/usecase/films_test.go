package usecase_test

import (
	"errors"
	"github.com/ellexo2456/FilmLib/internal/domain"
	"github.com/ellexo2456/FilmLib/internal/domain/mocks"
	"github.com/ellexo2456/FilmLib/internal/films/usecase"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name                     string
		getFilm                  func() domain.Film
		setFilmsRepoExpectations func(filmsRepo *mocks.FilmsRepository, id int, err error)
		expectedID               int
		expectedError            error
	}{
		{
			name: "GoodCase/Common",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2023-01-01")
				return domain.Film{
					Title:       "The Matrix",
					Description: "A computer hacker learns about the true nature of reality.",
					Rating:      8.7,
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: 1}, {ID: 2}},
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, id int, err error) {
				filmsRepo.On("Insert", mock.Anything).Return(id, err)
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name: "BadCase/EmptyFilm",
			getFilm: func() domain.Film {
				return domain.Film{}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, id int, err error) {
				filmsRepo.On("Insert", mock.Anything).Return(id, err).Maybe()
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "BadCase/EmptyTitle",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2023-01-01")
				return domain.Film{
					Description: "A computer hacker learns about the true nature of reality.",
					Rating:      8.7,
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: 1}, {ID: 2}},
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, id int, err error) {
				filmsRepo.On("Insert", mock.Anything).Return(id, err).Maybe()
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "BadCase/EmptyDescription",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2023-01-01")
				return domain.Film{
					Title:       "The Matrix",
					Rating:      8.7,
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: 1}, {ID: 2}},
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, id int, err error) {
				filmsRepo.On("Insert", mock.Anything).Return(id, err).Maybe()
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "BadCase/EmptyRating",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2023-01-01")
				return domain.Film{
					Title:       "The Matrix",
					Description: "A computer hacker learns about the true nature of reality.",
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: 1}, {ID: 2}},
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, id int, err error) {
				filmsRepo.On("Insert", mock.Anything).Return(id, err).Maybe()
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "BadCase/EmptyActors",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2023-01-01")
				return domain.Film{
					Title:       "The Matrix",
					Description: "A computer hacker learns about the true nature of reality.",
					Rating:      8.7,
					ReleaseDate: d,
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, id int, err error) {
				filmsRepo.On("Insert", mock.Anything).Return(id, err).Maybe()
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "BadCase/ReleaseDateGreaterThanCurrent",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("3000-01-01")
				return domain.Film{
					Title:       "The Matrix",
					Description: "A computer hacker learns about the true nature of reality.",
					Rating:      8.7,
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: 1}, {ID: 2}},
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, id int, err error) {
				filmsRepo.On("Insert", mock.Anything).Return(id, err)
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "BadCase/ReleaseDate1000Year",
			getFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("1000-01-01")
				return domain.Film{
					Title:       "The Matrix",
					Description: "A computer hacker learns about the true nature of reality.",
					Rating:      8.7,
					ReleaseDate: d,
					Actors:      []domain.Actor{{ID: 1}, {ID: 2}},
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, id int, err error) {
				filmsRepo.On("Insert", mock.Anything).Return(id, err)
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filmsRepo := new(mocks.FilmsRepository)
			test.setFilmsRepoExpectations(filmsRepo, test.expectedID, test.expectedError)

			filmsUsecase := usecase.NewFilmsUsecase(filmsRepo)
			id, err := filmsUsecase.Add(test.getFilm())

			assert.Equal(t, test.expectedID, id)
			assert.Equal(t, test.expectedError, err)

			filmsRepo.AssertExpectations(t)
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name                     string
		titleDir                 domain.SortDirection
		releaseDateDir           domain.SortDirection
		setFilmsRepoExpectations func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error)
		getFilms                 func() []domain.Film
		expectedError            error
	}{
		{
			name:     "GoodCase/SortByTitleAsc",
			titleDir: domain.Asc,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error) {
				filmsRepo.On("SelectAll").Return(films, err)
			},
			getFilms: func() []domain.Film {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return []domain.Film{
					{ID: 1, Title: "Inception", Description: "Description2", Rating: 8, ReleaseDate: d},
					{ID: 2, Title: "The Matrix", Description: "Description", Rating: 8.9, ReleaseDate: d},
				}
			},
			expectedError: nil,
		},
		{
			name:     "GoodCase/SortByTitleDesc",
			titleDir: domain.Desc,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error) {
				filmsRepo.On("SelectAll").Return(films, err)
			},
			getFilms: func() []domain.Film {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return []domain.Film{
					{ID: 2, Title: "The Matrix", Description: "Description", Rating: 8.9, ReleaseDate: d},
					{ID: 1, Title: "Inception", Description: "Description2", Rating: 8, ReleaseDate: d},
				}
			},
			expectedError: nil,
		},
		{
			name: "GoodCase/SortByRatingDesc",
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error) {
				filmsRepo.On("SelectAll").Return(films, err)
			},
			getFilms: func() []domain.Film {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return []domain.Film{
					{ID: 2, Title: "The Matrix", Description: "Description", Rating: 8.9, ReleaseDate: d},
					{ID: 1, Title: "Inception", Description: "Description2", Rating: 8, ReleaseDate: d},
				}
			},
			expectedError: nil,
		},
		{
			name:           "GoodCase/SortByReleaseDateAsc",
			releaseDateDir: domain.Asc,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error) {
				filmsRepo.On("SelectAll").Return(films, err)
			},
			getFilms: func() []domain.Film {
				var d1, d2 pgtype.Date
				d1.Scan("2000-01-01")
				d2.Scan("2000-02-01")

				return []domain.Film{
					{ID: 1, Title: "Inception", Description: "Description2", Rating: 8, ReleaseDate: d1},
					{ID: 2, Title: "The Matrix", Description: "Description", Rating: 8.9, ReleaseDate: d2},
				}
			},
			expectedError: nil,
		},
		{
			name:           "GoodCase/SortByReleaseDateDesc",
			releaseDateDir: domain.Desc,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error) {
				filmsRepo.On("SelectAll").Return(films, err)
			},
			getFilms: func() []domain.Film {
				var d1, d2 pgtype.Date
				d1.Scan("2000-01-01")
				d2.Scan("2000-02-01")

				return []domain.Film{
					{ID: 2, Title: "The Matrix", Description: "Description", Rating: 8.9, ReleaseDate: d2},
					{ID: 1, Title: "Inception", Description: "Description2", Rating: 8, ReleaseDate: d1},
				}
			},
			expectedError: nil,
		},
		{
			name:     "BadCase/RepoError",
			titleDir: domain.Asc,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error) {
				filmsRepo.On("SelectAll").Return(films, err)
			},
			getFilms: func() []domain.Film {
				return nil
			},
			expectedError: errors.New("some repo error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filmsRepo := new(mocks.FilmsRepository)
			test.setFilmsRepoExpectations(filmsRepo, test.getFilms(), test.expectedError)

			filmsUsecase := usecase.NewFilmsUsecase(filmsRepo)
			films, err := filmsUsecase.GetAll(test.titleDir, test.releaseDateDir)

			assert.Equal(t, test.getFilms(), films)
			assert.Equal(t, test.expectedError, err)

			filmsRepo.AssertExpectations(t)
		})
	}
}

func TestSearch(t *testing.T) {
	tests := []struct {
		name                     string
		searchStr                string
		setFilmsRepoExpectations func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error)
		getFilms                 func() []domain.Film
		expectedError            error
	}{
		{
			name:      "GoodCase/Common",
			searchStr: "matrix",
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error) {
				filmsRepo.On("Search", "matrix").Return(films, err)
			},
			getFilms: func() []domain.Film {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return []domain.Film{
					{ID: 1, Title: "The Matrix", Description: "Description", Rating: 8.9, ReleaseDate: d},
					{ID: 2, Title: "Matrix Reloaded", Description: "Description2", Rating: 8.5, ReleaseDate: d},
				}
			},
			expectedError: nil,
		},
		{
			name:      "BadCase/EmptySearchStr",
			searchStr: "",
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error) {
				filmsRepo.On("Search", "matrix").Return(films, err).Maybe()
			},
			getFilms: func() []domain.Film {
				return nil
			},
			expectedError: domain.ErrBadRequest,
		},
		{
			name:      "BadCase/RepoError",
			searchStr: "matrix",
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, films []domain.Film, err error) {
				filmsRepo.On("Search", mock.Anything).Return(nil, err)
			},
			getFilms: func() []domain.Film {
				return nil
			},
			expectedError: domain.ErrBadRequest,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filmsRepo := new(mocks.FilmsRepository)
			test.setFilmsRepoExpectations(filmsRepo, test.getFilms(), test.expectedError)

			filmsUsecase := usecase.NewFilmsUsecase(filmsRepo)
			films, err := filmsUsecase.Search(test.searchStr)

			assert.Equal(t, test.getFilms(), films)
			assert.Equal(t, test.expectedError, err)

			filmsRepo.AssertExpectations(t)
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name                     string
		id                       int
		setFilmsRepoExpectations func(filmsRepo *mocks.FilmsRepository, err error)
		expectedError            error
	}{
		{
			name: "GoodCase/Common",
			id:   1,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, err error) {
				filmsRepo.On("Delete", 1).Return(err)
			},
			expectedError: nil,
		},
		{
			name: "BadCase/NegativeID",
			id:   -1,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, err error) {
				filmsRepo.On("Delete", -1).Return(err).Maybe()
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name: "BadCase/ZeroID",
			id:   0,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, err error) {
				filmsRepo.On("Delete", 0).Return(err).Maybe()
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name: "BadCase/RepoError",
			id:   2,
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, err error) {
				filmsRepo.On("Delete", 2).Return(errors.New("repository error"))
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filmsRepo := new(mocks.FilmsRepository)
			test.setFilmsRepoExpectations(filmsRepo, test.expectedError)

			filmsUsecase := usecase.NewFilmsUsecase(filmsRepo)
			err := filmsUsecase.Remove(test.id)

			assert.Equal(t, test.expectedError, err)

			filmsRepo.AssertExpectations(t)
		})
	}
}

func TestModify(t *testing.T) {
	tests := []struct {
		name                     string
		getNewFilm               func() domain.Film
		setFilmsRepoExpectations func(filmsRepo *mocks.FilmsRepository, oldFilm domain.Film, updatedFilm domain.Film, err error)
		getExpectedFilm          func() domain.Film
		expectedError            error
	}{
		{
			name: "GoodCase/Common",
			getNewFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Film{
					ID:          1,
					Title:       "The Matrix Reloaded",
					Description: "New description",
					Rating:      9.0,
					ReleaseDate: d,
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, oldFilm domain.Film, updatedFilm domain.Film, err error) {
				filmsRepo.On("SelectById", 1).Return(oldFilm, err)
				filmsRepo.On("Update", mock.Anything).Return(updatedFilm, err)
			},
			getExpectedFilm: func() domain.Film {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Film{
					ID:          1,
					Title:       "The Matrix Reloaded",
					Description: "New description",
					Rating:      9.0,
					ReleaseDate: d,
				}
			},
			expectedError: nil,
		},
		{
			name: "BadCase/InvalidID",
			getNewFilm: func() domain.Film {
				return domain.Film{
					ID:          -1,
					Title:       "Invalid ID",
					Description: "Description",
					Rating:      8.0,
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, oldFilm domain.Film, updatedFilm domain.Film, err error) {},
			getExpectedFilm: func() domain.Film {
				return domain.Film{}
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name: "BadCase/RepositoryError",
			getNewFilm: func() domain.Film {
				return domain.Film{
					ID:          2,
					Title:       "Film",
					Description: "Description",
					Rating:      7.5,
				}
			},
			setFilmsRepoExpectations: func(filmsRepo *mocks.FilmsRepository, oldFilm domain.Film, updatedFilm domain.Film, err error) {
				filmsRepo.On("SelectById", 2).Return(domain.Film{}, errors.New("repository error"))
			},
			getExpectedFilm: func() domain.Film {
				return domain.Film{}
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filmsRepo := new(mocks.FilmsRepository)
			test.setFilmsRepoExpectations(filmsRepo, test.getExpectedFilm(), test.getNewFilm(), test.expectedError)

			filmsUsecase := usecase.NewFilmsUsecase(filmsRepo)
			updatedFilm, err := filmsUsecase.Modify(test.getNewFilm())

			assert.Equal(t, test.getExpectedFilm(), updatedFilm)
			assert.Equal(t, test.expectedError, err)

			filmsRepo.AssertExpectations(t)
		})
	}
}
