package http_test

import (
	"bytes"
	"encoding/json"
	"github.com/ellexo2456/FilmLib/internal/domain"
	"github.com/ellexo2456/FilmLib/internal/domain/mocks"
	films_http "github.com/ellexo2456/FilmLib/internal/films/delivery/http"
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

func TestAddFilm(t *testing.T) {
	tests := []struct {
		name                 string
		getBody              func() []byte
		setUCaseExpectations func(usecase *mocks.FilmsUsecase, film domain.Film)
		ctx                  context.Context
		status               int
	}{
		{
			name: "GoodCase/Common",
			getBody: func() []byte {
				return []byte(`{"title":"Film Title", "description": "Film Description", "releaseDate": "2022-01-01", "rating": 8.5, "actors": [{"id": 4}, {"id": 5}]}`)
			},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film domain.Film) {
				usecase.On("Add", film).Return(1, nil)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusOK,
		},
		{
			name: "BadCase/EmptyTitle",
			getBody: func() []byte {
				return []byte(`{"title":"", "description": "Film Description", "releaseDate": "2022-01-01", "rating": 8.5, "actors": [{"id": 4}, {"id": 5}]}`)
			},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film domain.Film) {
				usecase.On("Add", film).Return(0, domain.ErrBadRequest)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/EmptyDescription",
			getBody: func() []byte {
				return []byte(`{"title":"Film Title", "description": "", "releaseDate": "2022-01-01", "rating": 8.5, "actors": [{"id": 4}, {"id": 5}]}`)
			},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film domain.Film) {
				usecase.On("Add", film).Return(0, domain.ErrBadRequest)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/InvalidReleaseDate",
			getBody: func() []byte {
				return []byte(`{"title":"Film Title", "description": "Film Description", "releaseDate": "invalid_date", "rating": 8.5, "actors": [{"id": 4}, {"id": 5}]}`)
			},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film domain.Film) {
				usecase.On("Add", film).Return(0, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/NegativeRating",
			getBody: func() []byte {
				return []byte(`{"title":"Film Title", "description": "Film Description", "releaseDate": "2022-01-01", "rating": -1.5, "actors": [{"id": 4}, {"id": 5}]}`)
			},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film domain.Film) {
				usecase.On("Add", film).Return(0, domain.ErrBadRequest)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/NoModerRole",
			getBody: func() []byte {
				return []byte(`{"title":"Film Title", "description": "Film Description", "releaseDate": "2022-01-01", "rating": 8.5, "actors": [{"id": 4}, {"id": 5}]}`)
			},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film domain.Film) {
				usecase.On("Add", film).Return(0, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Usr}),
			status: http.StatusForbidden,
		},
		{
			name: "BadCase/NoUserContext",
			getBody: func() []byte {
				return []byte(`{"title":"Film Title", "description": "Film Description", "releaseDate": "2022-01-01", "rating": 8.5, "actors": [{"id": 4}, {"id": 5}]}`)
			},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, film domain.Film) {
				usecase.On("Add", film).Return(0, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.Background(),
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.FilmsUsecase)
			var film domain.Film
			json.Unmarshal(test.getBody(), &film)
			test.setUCaseExpectations(mockUsecase, film)

			req := httptest.NewRequest("POST", "/api/v1/films", bytes.NewReader(test.getBody()))
			req = req.WithContext(test.ctx)
			rec := httptest.NewRecorder()

			handler := &films_http.FilmsHandler{FilmsUsecase: mockUsecase}
			handler.AddFilm(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestGetFilms(t *testing.T) {
	tests := []struct {
		name                 string
		queryParams          map[string]string
		setUCaseExpectations func(usecase *mocks.FilmsUsecase, titleDir domain.SortDirection, releaseDateDir domain.SortDirection)
		status               int
	}{
		{
			name:        "GoodCase/WithSortParams",
			queryParams: map[string]string{"sortTitle": "Asc", "sortReleaseDate": "Desc"},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, titleDir domain.SortDirection, releaseDateDir domain.SortDirection) {
				usecase.On("GetAll", titleDir, releaseDateDir).Return([]domain.Film{}, nil)
			},
			status: http.StatusOK,
		},
		{
			name:        "GoodCase/WithoutSortParams",
			queryParams: map[string]string{},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, titleDir domain.SortDirection, releaseDateDir domain.SortDirection) {
				usecase.On("GetAll", titleDir, releaseDateDir).Return([]domain.Film{}, nil)
			},
			status: http.StatusOK,
		},
		{
			name:        "GoodCase/InvalidSortTitle",
			queryParams: map[string]string{"sortTitle": "Invalid", "sortReleaseDate": "Desc"},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, titleDir domain.SortDirection, releaseDateDir domain.SortDirection) {
				usecase.On("GetAll", titleDir, releaseDateDir).Return([]domain.Film{}, nil)
			},
			status: http.StatusOK,
		},
		{
			name:        "GoodCase/InvalidSortReleaseDate",
			queryParams: map[string]string{"sortTitle": "Asc", "sortReleaseDate": "Invalid"},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, titleDir domain.SortDirection, releaseDateDir domain.SortDirection) {
				usecase.On("GetAll", titleDir, releaseDateDir).Return([]domain.Film{}, nil)
			},
			status: http.StatusOK,
		},
		{
			name:        "GoodCase/InternalServerError",
			queryParams: map[string]string{},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, titleDir domain.SortDirection, releaseDateDir domain.SortDirection) {
				usecase.On("GetAll", titleDir, releaseDateDir).Return([]domain.Film{}, domain.ErrInternalServerError)
			},
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.FilmsUsecase)
			var titleDir, releaseDateDir domain.SortDirection
			if sortTitle, ok := test.queryParams[domain.TitleParam]; ok {
				titleDir = domain.SortDirection(sortTitle)
			}
			if sortReleaseDate, ok := test.queryParams[domain.ReleaseDateParam]; ok {
				releaseDateDir = domain.SortDirection(sortReleaseDate)
			}
			test.setUCaseExpectations(mockUsecase, titleDir, releaseDateDir)

			req := httptest.NewRequest("GET", "/api/v1/films", nil)
			q := req.URL.Query()
			for key, value := range test.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()

			handler := &films_http.FilmsHandler{FilmsUsecase: mockUsecase}
			handler.GetFilms(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestSearchFilms(t *testing.T) {
	tests := []struct {
		name                 string
		queryParams          map[string]string
		setUCaseExpectations func(usecase *mocks.FilmsUsecase, searchStr string)
		status               int
	}{
		{
			name:        "GoodCase/WithSearchStr",
			queryParams: map[string]string{"searchStr": "film title"},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, searchStr string) {
				usecase.On("Search", searchStr).Return([]domain.Film{}, nil)
			},
			status: http.StatusOK,
		},
		{
			name:        "BadCase/EmptySearchStr",
			queryParams: map[string]string{"searchStr": ""},
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, searchStr string) {
				usecase.On("Search", searchStr).Return(nil, domain.ErrBadRequest)
			},
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.FilmsUsecase)
			test.setUCaseExpectations(mockUsecase, test.queryParams[domain.SearchParam])

			req := httptest.NewRequest("GET", "/api/v1/films/search", nil)
			q := req.URL.Query()
			for key, value := range test.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			rec := httptest.NewRecorder()

			handler := &films_http.FilmsHandler{FilmsUsecase: mockUsecase}
			handler.Search(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestDeleteFilm(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(usecase *mocks.FilmsUsecase, id int)
		ctx                  context.Context
		id                   string
		status               int
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, id int) {
				usecase.On("Remove", id).Return(nil)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "1",
			status: http.StatusNoContent,
		},
		{
			name: "BadCase/InvalidRole",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Usr}),
			id:     "1",
			status: http.StatusForbidden,
		},
		{
			name: "BadCase/InvalidID",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "invalid_id",
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/EmptyID",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "",
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/NoUserContext",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrBadRequest).Maybe()
			},
			ctx:    context.Background(),
			id:     "1",
			status: http.StatusInternalServerError,
		},
		{
			name: "BadCase/OutOfRangeFilmId",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "1234563456789",
			status: http.StatusNotFound,
		},
		{
			name: "BadCase/NegativeFilmId",
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase, id int) {
				usecase.On("Remove", id).Return(domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			id:     "-3",
			status: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.FilmsUsecase)
			id, err := strconv.Atoi(test.id)
			if err != nil {
				id = 0
			}
			if test.setUCaseExpectations != nil {
				test.setUCaseExpectations(mockUsecase, id)
			}

			req := httptest.NewRequest("DELETE", "/films/"+test.id, nil)
			req = req.WithContext(test.ctx)
			rec := httptest.NewRecorder()

			mux := http.NewServeMux()
			films_http.NewFilmsHandler(mux, mockUsecase)

			mux.ServeHTTP(rec, req)
			assert.Equal(t, test.status, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestModifyFilm(t *testing.T) {
	tests := []struct {
		name                 string
		requestBody          io.Reader
		setUCaseExpectations func(usecase *mocks.FilmsUsecase)
		ctx                  context.Context
		status               int
	}{
		{
			name:        "GoodCase/Common",
			requestBody: strings.NewReader(`{"id":1,"title":"New Title", "description": "New Description", "releaseDate": "2023-01-01", "rating": 9.0}`),
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				var d pgtype.Date
				d.Scan("2023-01-01")

				usecase.On("Modify", domain.Film{
					ID:          1,
					Title:       "New Title",
					Description: "New Description",
					ReleaseDate: d,
					Rating:      9.0,
				}).Return(domain.Film{
					ID:          1,
					Title:       "New Title",
					Description: "New Description",
					ReleaseDate: d,
					Rating:      9.0,
				}, nil)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusOK,
		},
		{
			name:        "BadCase/EmptyTitle",
			requestBody: strings.NewReader(`{"id":1,"title":"", "description": "New Description", "releaseDate": "2023-01-01", "rating": 9.0}`),
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				var d pgtype.Date
				d.Scan("2023-01-01")

				usecase.On("Modify", domain.Film{
					ID:          1,
					Title:       "",
					Description: "New Description",
					ReleaseDate: d,
					Rating:      9.0,
				}).Return(domain.Film{}, domain.ErrBadRequest)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name:        "BadCase/InvalidReleaseDate",
			requestBody: strings.NewReader(`{"id":1,"title":"New Title", "description": "New Description", "releaseDate": "invalid_date", "rating": 9.0}`),
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				usecase.On("Modify", mock.Anything).Return(domain.Film{}, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name:        "BadCase/NegativeRating",
			requestBody: strings.NewReader(`{"id":1,"title":"New Title", "description": "New Description", "releaseDate": "2023-01-01", "rating": -1.5}`),
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				usecase.On("Modify", mock.Anything).Return(domain.Film{}, domain.ErrBadRequest)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusBadRequest,
		},
		{
			name:        "BadCase/NoModerRole",
			requestBody: strings.NewReader(`{"id":1,"title":"New Title", "description": "New Description", "releaseDate": "2023-01-01", "rating": 9.0}`),
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				usecase.On("Modify", mock.Anything).Return(domain.Film{}, domain.ErrBadRequest)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Usr}),
			status: http.StatusForbidden,
		},
		{
			name:        "BadCase/FutureReleaseDate",
			requestBody: strings.NewReader(`{"id":1,"title":"New Title", "description": "New Description", "releaseDate": "3023-01-01", "rating": 9.0}`),
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				usecase.On("Modify", mock.Anything).Return(domain.Film{}, domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusNotFound,
		},
		{
			name:        "BadCase/PastReleaseDate",
			requestBody: strings.NewReader(`{"id":1,"title":"New Title", "description": "New Description", "releaseDate": "1000-01-01", "rating": 9.0}`),
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				usecase.On("Modify", mock.Anything).Return(domain.Film{}, domain.ErrOutOfRange)
			},
			ctx:    context.WithValue(context.Background(), domain.SessionContextKey, domain.SessionContext{Role: domain.Moder}),
			status: http.StatusNotFound,
		},
		{
			name:        "BadCase/NoUserContext",
			requestBody: strings.NewReader(`{"id":1,"title":"New Title", "description": "New Description", "releaseDate": "2023-01-01", "rating": 9.0}`),
			setUCaseExpectations: func(usecase *mocks.FilmsUsecase) {
				usecase.On("Modify", mock.Anything).Return(domain.Film{}, domain.ErrBadRequest).Maybe()
			},
			ctx:    context.Background(),
			status: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.FilmsUsecase)
			test.setUCaseExpectations(mockUsecase)

			req := httptest.NewRequest("PUT", "/films", test.requestBody)
			req = req.WithContext(test.ctx)
			rec := httptest.NewRecorder()

			handler := &films_http.FilmsHandler{FilmsUsecase: mockUsecase}
			handler.ModifyFilm(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}
