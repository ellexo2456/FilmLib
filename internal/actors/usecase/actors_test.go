package usecase_test

import (
	"errors"
	"github.com/ellexo2456/FilmLib/internal/actors/usecase"
	"github.com/ellexo2456/FilmLib/internal/domain"
	"github.com/ellexo2456/FilmLib/internal/domain/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name                     string
		getActor                 func() domain.Actor
		setActorsRepoExpectation func(actorsRepo *mocks.ActorsRepository, id int, err error)
		expectedID               int
		expectedError            error
	}{
		{
			name: "GoodCase/Common",
			getActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("1988-11-16")

				return domain.Actor{
					Name:      "Keanu Reeves",
					Sex:       "Male",
					Birthdate: d,
				}
			},
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, id int, err error) {
				actorsRepo.On("Insert", mock.Anything).Return(id, err)
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name: "GoodCase/FutureBirthday",
			getActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("1988-11-16")

				return domain.Actor{
					Name:      "Keanu Reeves",
					Sex:       "Male",
					Birthdate: d,
				}
			},
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, id int, err error) {
				actorsRepo.On("Insert", mock.Anything).Return(id, err)
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "GoodCase/PastBirthdate",
			getActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("1000-01-01")

				return domain.Actor{
					Name:      "Keanu Reeves",
					Sex:       "Male",
					Birthdate: d,
				}
			},
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, id int, err error) {
				actorsRepo.On("Insert", mock.Anything).Return(id, err)
			},
			expectedID:    1,
			expectedError: nil,
		},
		{
			name: "BadCase/EmptyActor",
			getActor: func() domain.Actor {
				return domain.Actor{}
			},
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, id int, err error) {
				actorsRepo.On("Insert", mock.Anything).Return(id, err).Maybe()
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "BadCase/EmptyName",
			getActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Actor{
					Sex:       "Male",
					Birthdate: d,
				}
			},
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, id int, err error) {
				actorsRepo.On("Insert", mock.Anything).Return(id, err).Maybe()
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "BadCase/EmptySex",
			getActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Actor{
					Name:      "Keanu Reeves",
					Birthdate: d,
				}
			},
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, id int, err error) {
				actorsRepo.On("Insert", mock.Anything).Return(id, err).Maybe()
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
		{
			name: "BadCase/EmptyBirthdate",
			getActor: func() domain.Actor {
				return domain.Actor{
					Name: "Keanu Reeves",
					Sex:  "Male",
				}
			},
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, id int, err error) {
				actorsRepo.On("Insert", mock.Anything).Return(id, err).Maybe()
			},
			expectedID:    0,
			expectedError: domain.ErrBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actorsRepo := new(mocks.ActorsRepository)
			test.setActorsRepoExpectation(actorsRepo, test.expectedID, test.expectedError)

			actorsUsecase := usecase.NewActorsUsecase(actorsRepo)
			id, err := actorsUsecase.Add(test.getActor())

			assert.Equal(t, test.expectedID, id)
			assert.Equal(t, test.expectedError, err)

			actorsRepo.AssertExpectations(t)
		})
	}
}

func TestRemove(t *testing.T) {
	tests := []struct {
		name                     string
		id                       int
		setActorsRepoExpectation func(actorsRepo *mocks.ActorsRepository, err error)
		expectedError            error
	}{
		{
			name: "GoodCase/Common",
			id:   1,
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, err error) {
				actorsRepo.On("Delete", 1).Return(err)
			},
			expectedError: nil,
		},
		{
			name: "BadCase/NegativeID",
			id:   -1,
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, err error) {
				actorsRepo.On("Delete", -1).Return(err).Maybe()
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name: "BadCase/ZeroID",
			id:   0,
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, err error) {
				actorsRepo.On("Delete", 0).Return(err).Maybe()
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name: "BadCase/RepoError",
			id:   2,
			setActorsRepoExpectation: func(actorsRepo *mocks.ActorsRepository, err error) {
				actorsRepo.On("Delete", 2).Return(errors.New("repository error"))
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actorsRepo := new(mocks.ActorsRepository)
			test.setActorsRepoExpectation(actorsRepo, test.expectedError)

			actorsUsecase := usecase.NewActorsUsecase(actorsRepo)
			err := actorsUsecase.Remove(test.id)

			assert.Equal(t, test.expectedError, err)

			actorsRepo.AssertExpectations(t)
		})
	}
}

func TestModify(t *testing.T) {
	tests := []struct {
		name                      string
		getNewActor               func() domain.Actor
		setActorsRepoExpectations func(actorsRepo *mocks.ActorsRepository, oldActor domain.Actor, updatedActor domain.Actor, err error)
		getExpectedActor          func() domain.Actor
		getOldActor               func() domain.Actor
		expectedError             error
	}{
		{
			name: "GoodCase/Common",
			getNewActor: func() domain.Actor {
				return domain.Actor{
					ID:   1,
					Name: "John Doe",
					Sex:  "Male",
				}
			},
			getOldActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Actor{
					ID:        1,
					Name:      "Ever Ken",
					Sex:       "Female",
					Birthdate: d,
				}
			},
			setActorsRepoExpectations: func(actorsRepo *mocks.ActorsRepository, oldActor domain.Actor, updatedActor domain.Actor, err error) {
				actorsRepo.On("SelectById", oldActor.ID).Return(oldActor, err)
				actorsRepo.On("Update", mock.Anything).Return(updatedActor, err)
			},
			getExpectedActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Actor{
					ID:        1,
					Name:      "John Doe",
					Sex:       "Male",
					Birthdate: d,
				}
			},
			expectedError: nil,
		},
		{
			name: "BadCase/InvalidID",
			getNewActor: func() domain.Actor {
				return domain.Actor{
					ID:   -1,
					Name: "Invalid ID",
					Sex:  "Male",
				}
			},
			getOldActor: func() domain.Actor {
				return domain.Actor{}
			},
			setActorsRepoExpectations: func(actorsRepo *mocks.ActorsRepository, oldActor domain.Actor, updatedActor domain.Actor, err error) {
			},
			getExpectedActor: func() domain.Actor {
				return domain.Actor{}
			},
			expectedError: domain.ErrNotFound,
		},
		{
			name: "BadCase/FutureBirthdate",
			getNewActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("2100-01-01")

				return domain.Actor{
					ID:        2,
					Birthdate: d,
				}
			},
			getOldActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Actor{
					ID:        2,
					Name:      "Ever Ken",
					Sex:       "Female",
					Birthdate: d,
				}
			},
			setActorsRepoExpectations: func(actorsRepo *mocks.ActorsRepository, oldActor domain.Actor, updatedActor domain.Actor, err error) {
				actorsRepo.On("SelectById", oldActor.ID).Return(oldActor, nil)
				actorsRepo.On("Update", mock.Anything).Return(updatedActor, err)
			},
			getExpectedActor: func() domain.Actor {
				return domain.Actor{}
			},
			expectedError: domain.ErrOutOfRange,
		},
		{
			name: "BadCase/PastBirthdate",
			getNewActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("1000-01-01")

				return domain.Actor{
					ID:        3,
					Birthdate: d,
				}
			},
			getOldActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Actor{
					ID:        3,
					Name:      "Ever Ken",
					Sex:       "Female",
					Birthdate: d,
				}
			},
			setActorsRepoExpectations: func(actorsRepo *mocks.ActorsRepository, oldActor domain.Actor, updatedActor domain.Actor, err error) {
				actorsRepo.On("SelectById", oldActor.ID).Return(oldActor, nil)
				actorsRepo.On("Update", mock.Anything).Return(updatedActor, err)
			},
			getExpectedActor: func() domain.Actor {
				return domain.Actor{}
			},
			expectedError: domain.ErrOutOfRange,
		},
		{
			name: "BadCase/RepositoryError",
			getNewActor: func() domain.Actor {
				return domain.Actor{
					ID:   2,
					Name: "Actor",
					Sex:  "Female",
				}
			},
			getOldActor: func() domain.Actor {
				var d pgtype.Date
				d.Scan("2000-01-01")

				return domain.Actor{
					ID:        2,
					Name:      "Ever Ken",
					Sex:       "Male",
					Birthdate: d,
				}
			},
			setActorsRepoExpectations: func(actorsRepo *mocks.ActorsRepository, oldActor domain.Actor, updatedActor domain.Actor, err error) {
				actorsRepo.On("SelectById", 2).Return(updatedActor, errors.New("repository error"))
			},
			getExpectedActor: func() domain.Actor {
				return domain.Actor{}
			},
			expectedError: errors.New("repository error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actorsRepo := new(mocks.ActorsRepository)
			test.setActorsRepoExpectations(actorsRepo, test.getOldActor(), test.getExpectedActor(), test.expectedError)

			actorsUsecase := usecase.NewActorsUsecase(actorsRepo)
			updatedActor, err := actorsUsecase.Modify(test.getNewActor())

			assert.Equal(t, test.getExpectedActor(), updatedActor)
			assert.Equal(t, test.expectedError, err)

			actorsRepo.AssertExpectations(t)
		})
	}
}

func TestGetAll(t *testing.T) {
	tests := []struct {
		name                      string
		setActorsRepoExpectations func(actorsRepo *mocks.ActorsRepository, actors []domain.Actor, err error)
		getExpectedActors         func() []domain.Actor
		expectedError             error
	}{
		{
			name: "GoodCase/Common",
			setActorsRepoExpectations: func(actorsRepo *mocks.ActorsRepository, actors []domain.Actor, err error) {
				actorsRepo.On("SelectAll").Return(actors, err)
			},
			getExpectedActors: func() []domain.Actor {
				var d1, d2 pgtype.Date
				d1.Scan("2000-01-01")
				d2.Scan("2000-02-01")

				return []domain.Actor{
					{ID: 1, Name: "John Doe", Sex: "Male", Birthdate: d2},
					{ID: 2, Name: "Jane Smith", Sex: "Female", Birthdate: d1},
				}
			},
			expectedError: nil,
		},
		{
			name: "BadCase/RepoError",
			setActorsRepoExpectations: func(actorsRepo *mocks.ActorsRepository, actors []domain.Actor, err error) {
				actorsRepo.On("SelectAll").Return(actors, err)
			},
			getExpectedActors: func() []domain.Actor {
				return nil
			},
			expectedError: errors.New("some repo error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actorsRepo := new(mocks.ActorsRepository)
			test.setActorsRepoExpectations(actorsRepo, test.getExpectedActors(), test.expectedError)

			actorsUsecase := usecase.NewActorsUsecase(actorsRepo)
			actors, err := actorsUsecase.GetAll()

			assert.Equal(t, test.getExpectedActors(), actors)
			assert.Equal(t, test.expectedError, err)

			actorsRepo.AssertExpectations(t)
		})
	}
}
