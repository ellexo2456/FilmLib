package http_test

import (
	"bytes"
	actor_http "github.com/ellexo2456/FilmLib/internal/actors/delivery/http"
	"github.com/ellexo2456/FilmLib/internal/domain"
	"github.com/ellexo2456/FilmLib/internal/domain/mocks"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestAddActor(t *testing.T) {
	tests := []struct {
		name                 string
		getBody              func() []byte
		setUCaseExpectations func(usecase *mocks.ActorsUsecase)
		ctx                  context.Context
		status               int
	}{
		{
			name: "GoodCase/Common",
			getBody: func() []byte {
				return []byte(`{ "name":"john", "sex": "M", "birthdate": "2000-01-01" }`)
			},
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Add", mock.Anything).Return(1, nil)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusOK,
		},
		{
			name: "BadCase/EmptyName",
			getBody: func() []byte {
				return []byte(`{"name":"", "sex": "M", "birthdate": "2000-01-01"}`)
			},
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Add", mock.Anything).Return(0, domain.ErrBadRequest)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/InvalidJson",
			getBody: func() []byte {
				return []byte(`{"name":john, "sex": "M", "birthdate": "2003-01-01"}`)
			},
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Add", mock.Anything).Return(0, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/EmptyJson",
			getBody: func() []byte {
				return []byte(`{"name":"john", "sex": "M", "birthdate": "2000-01-01"}`)
			},
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Add", mock.Anything).Return(0, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/InvalidSex",
			getBody: func() []byte {
				return []byte(`{"name":"john", sex: "MAT", birthdate: "2000-01-01"}`)
			},
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Add", mock.Anything).Return(0, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/FutureDate",
			getBody: func() []byte {
				return []byte(`{"name":"john", "sex": "M", "birthdate": "2033-01-01"}`)
			},
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Add", mock.Anything).Return(0, domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/PastDate",
			getBody: func() []byte {
				return []byte(`{"name":"john", "sex": "M", "birthdate": "1000-01-01"}`)
			},
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Add", mock.Anything).Return(0, domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/NoUserContext",
			getBody: func() []byte {
				return []byte(`{"name":"john", "sex": "M", "birthdate": "2000-01-01"}`)
			},
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Add", mock.Anything).Return(0, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.Background(),
			status: http.StatusInternalServerError,
		},
		{
			name: "BadCase/InvalidRole",
			getBody: func() []byte {
				return []byte(`{"name":"john", "sex": "M", "birthdate": "2000-01-01"}`)
			},
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Add", mock.Anything).Return(0, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Usr}),
			status: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.ActorsUsecase)
			test.setUCaseExpectations(mockUsecase)

			req := httptest.NewRequest("POST", "/actors", bytes.NewReader(test.getBody()))
			req = req.WithContext(test.ctx)
			rec := httptest.NewRecorder()

			handler := &actor_http.ActorsHandler{ActorsUsecase: mockUsecase}
			handler.AddActor(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestDeleteActor(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(usecase *mocks.ActorsUsecase, id int)
		ctx                  context.Context
		id                   string
		status               int
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase, id int) {
				usecase.On("Remove", id).Return(nil)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "1",
			status: http.StatusNoContent,
		},
		{
			name: "BadCase/InvalidRole",
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Usr}),
			id:     "1",
			status: http.StatusForbidden,
		},
		{
			name: "BadCase/InvalidID",
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "invalid_id",
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/EmptyID",
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "",
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/NoUserContext",
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrBadRequest).Maybe()
			},
			ctx:    context.Background(),
			id:     "1",
			status: http.StatusInternalServerError,
		},

		{
			name: "BadCase/OutOfRangeVideoId",
			setUCaseExpectations: func(fvu *mocks.ActorsUsecase, id int) {
				fvu.On("Remove", id).Return(domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "1234563456789",
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/NegativeVideoId",
			setUCaseExpectations: func(fvu *mocks.ActorsUsecase, id int) {
				fvu.On("Remove", id).Return(domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "-3",
			status: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.ActorsUsecase)
			id, err := strconv.Atoi(test.id)
			if err != nil {
				id = 0
			}
			if test.setUCaseExpectations != nil {
				test.setUCaseExpectations(mockUsecase, id)
			}

			req := httptest.NewRequest("DELETE", "/actors/"+test.id, nil)
			req = req.WithContext(test.ctx)
			rec := httptest.NewRecorder()

			mux := http.NewServeMux()
			actor_http.NewActorsHandler(mux, mockUsecase)

			mux.ServeHTTP(rec, req)
			assert.Equal(t, test.status, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestModifyActor(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          io.Reader
		setUCaseExpectations func(usecase *mocks.ActorsUsecase)
		ctx                  context.Context
		status               int
	}{
		{
			name:        "GoodCase/Common",
			requestBody: strings.NewReader(`{"name":"John","sex":"M","birthdate":"2000-01-01"}`),
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				var d pgtype.Date
				d.Scan("2000-01-01")

				usecase.On("Modify", domain.Actor{
					Name:      "John",
					Sex:       "M",
					Birthdate: d,
				}).Return(domain.Actor{
					ID:        1,
					Name:      "John",
					Sex:       "M",
					Birthdate: d,
				}, nil)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusOK,
		},
		{
			name:        "BadCase/EmptyName",
			requestBody: strings.NewReader(`{"name":"","sex":"M","birthdate":"2000-01-01"}`),
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				var d pgtype.Date
				d.Scan("2000-01-01")

				usecase.On("Modify", domain.Actor{
					Name:      "",
					Sex:       "M",
					Birthdate: d,
				}).Return(domain.Actor{}, domain.ErrBadRequest)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name:        "BadCase/InvalidSex",
			requestBody: strings.NewReader(`{"name":"John","sex":"X","birthdate":"2000-01-01"}`),
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				var d pgtype.Date
				d.Scan("2000-01-01")

				usecase.On("Modify", domain.Actor{
					Name:      "John",
					Sex:       "X",
					Birthdate: d,
				}).Return(domain.Actor{}, domain.ErrBadRequest)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name:        "BadCase/FutureDate",
			requestBody: strings.NewReader(`{"name":"John","sex":"M","birthdate":"2030-01-01"}`),
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				var d pgtype.Date
				d.Scan("2030-01-01")

				usecase.On("Modify", domain.Actor{
					Name:      "John",
					Sex:       "M",
					Birthdate: d,
				}).Return(domain.Actor{}, domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusNotFound,
		},
		{
			name:        "BadCase/PastDate",
			requestBody: strings.NewReader(`{"name":"John","sex":"M","birthdate":"1500-01-01"}`),
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				var d pgtype.Date
				d.Scan("1500-01-01")

				usecase.On("Modify", domain.Actor{
					Name:      "John",
					Sex:       "M",
					Birthdate: d,
				}).Return(domain.Actor{}, domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusNotFound,
		},
		{
			name:        "BadCase/InvalidJSON",
			requestBody: strings.NewReader(`invalid json`),
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Modify", mock.Anything).Return(domain.Actor{}, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name:        "BadCase/InvalidRole",
			requestBody: strings.NewReader(`{"name":"John","sex":"M","birthdate":"2000-01-01"}`),
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Modify", mock.Anything).Return(domain.Actor{}, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Usr}),
			status: http.StatusForbidden,
		},
		{
			name:        "BadCase/NoUserContext",
			requestBody: strings.NewReader(`{"name":"John","sex":"M","birthdate":"2000-01-01"}`),
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("Modify", mock.Anything).Return(domain.Actor{}, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.Background(),
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.ActorsUsecase)
			test.setUCaseExpectations(mockUsecase)

			req := httptest.NewRequest("PUT", "/actors", test.requestBody)
			req = req.WithContext(test.ctx)
			rec := httptest.NewRecorder()

			handler := &actor_http.ActorsHandler{ActorsUsecase: mockUsecase}
			handler.ModifyActor(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetActors(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(usecase *mocks.ActorsUsecase)
		status               int
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				var d pgtype.Date
				d.Scan("2000-01-01")

				actors := []domain.Actor{
					{ID: 1, Name: "John", Sex: "M", Birthdate: d},
					{ID: 2, Name: "Jane", Sex: "F", Birthdate: d},
				}
				usecase.On("GetAll").Return(actors, nil)
			},
			status: http.StatusOK,
		},
		{
			name: "BadCase/InternalServerError",
			setUCaseExpectations: func(usecase *mocks.ActorsUsecase) {
				usecase.On("GetAll").Return(nil, domain.ErrInternalServerError)
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.ActorsUsecase)
			test.setUCaseExpectations(mockUsecase)

			req := httptest.NewRequest("GET", "/actors", nil)
			rec := httptest.NewRecorder()

			handler := &actor_http.ActorsHandler{ActorsUsecase: mockUsecase}
			handler.GetActors(rec, req)

			assert.Equal(t, test.status, rec.Code)
		})
	}
}
