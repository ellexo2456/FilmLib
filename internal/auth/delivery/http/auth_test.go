package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/bxcodec/faker"
	auth_http "github.com/ellexo2456/FilmLib/internal/auth/delivery/http"
	"github.com/ellexo2456/FilmLib/internal/domain"
	"github.com/ellexo2456/FilmLib/internal/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name                 string
		getBody              func() []byte
		setUCaseExpectations func(session *domain.Session, uCase *mocks.AuthUsecase)
		status               int
		wantCookie           bool
		setAuth              func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session)
	}{
		{
			name: "GoodCase/Common",
			getBody: func() []byte {
				var creds domain.Credentials
				faker.FakeData(&creds.Password)
				creds.Email = "ferfg@fsf.ru"
				jsonBody, _ := json.Marshal(creds)
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				err := faker.FakeData(session)
				assert.NoError(t, err)
				session.ExpiresAt = time.Now().Add(24 * time.Hour)

				uCase.On("Login", mock.Anything).Return(*session, 1, nil)
			},
			status:     http.StatusOK,
			wantCookie: true,
		},
		{
			name: "BadCase/EmptyCredentials",
			getBody: func() []byte {
				jsonBody, _ := json.Marshal(domain.Credentials{})
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/WrongCredentials",
			getBody: func() []byte {
				var creds domain.Credentials
				faker.FakeData(&creds.Password)
				creds.Email = "fsf.ru"
				jsonBody, _ := json.Marshal(creds)
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/InvalidJson",
			getBody: func() []byte {
				return []byte(`{ "password":"3490rjuv", email: rszdxtfcyguhj@sgf.ru }`)
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/EmptyJson",
			getBody: func() []byte {
				return []byte(`{}`)
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/NoBody",
			getBody: func() []byte {
				return []byte{}
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}
				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/AlreadyAuthorized",
			getBody: func() []byte {
				var creds domain.Credentials
				faker.FakeData(&creds.Password)
				creds.Email = "ferfg@fsf.ru"
				jsonBody, _ := json.Marshal(creds)
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				err := faker.FakeData(session)
				assert.NoError(t, err)

				session.UserID = 1
				session.ExpiresAt = time.Now().Add(24 * time.Hour)
				session.Role = domain.Usr

				uCase.On("Login", mock.Anything).Return(*session, 0, nil).Maybe()
			},
			status: http.StatusConflict,
			setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
				r.AddCookie(&http.Cookie{
					Name:     "session_token",
					Value:    session.Token,
					Expires:  session.ExpiresAt,
					Path:     "/",
					HttpOnly: true,
				})

				uCase.On("RetrieveSessionContext", session.Token).
					Return(domain.SessionContext{
						UserID: session.UserID,
						Role:   session.Role,
					}, nil)
			},
		},
		{
			name: "GoodCase/AlreadyAuthorizedExpiredCookie",
			getBody: func() []byte {
				var creds domain.Credentials
				faker.FakeData(&creds.Password)
				creds.Email = "ferfg@fsf.ru"
				jsonBody, _ := json.Marshal(creds)
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				err := faker.FakeData(session)
				assert.NoError(t, err)

				session.UserID = 1
				session.ExpiresAt = time.Now()
				session.Role = domain.Usr

				uCase.On("Login", mock.Anything).Return(*session, 1, nil)
			},
			status:     http.StatusOK,
			wantCookie: true,
			setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
				r.AddCookie(&http.Cookie{
					Name:     "session_token",
					Value:    session.Token,
					Expires:  session.ExpiresAt,
					Path:     "/",
					HttpOnly: true,
				})
				uCase.On("RetrieveSessionContext", session.Token).
					Return(domain.SessionContext{
						UserID: session.UserID,
						Role:   session.Role,
					}, domain.ErrUnauthorized).Maybe()
			},
		},
		{
			name: "GoodCase/AlreadyAuthorizedWrongCookie",
			getBody: func() []byte {
				var creds domain.Credentials
				faker.FakeData(&creds.Password)
				creds.Email = "ferfg@fsf.ru"
				jsonBody, _ := json.Marshal(creds)
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				err := faker.FakeData(session)
				assert.NoError(t, err)

				session.ExpiresAt = time.Now().Add(24 * time.Hour)
				session.Role = domain.Usr

				uCase.On("Login", mock.Anything).Return(*session, 1, nil).Maybe()
			},
			status:     http.StatusOK,
			wantCookie: true,
			setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
				r.AddCookie(&http.Cookie{
					Name:     "fevk",
					Value:    session.Token,
					Expires:  session.ExpiresAt,
					Path:     "/",
					HttpOnly: true,
				})
				uCase.On("RetrieveSessionContext", session.Token).
					Return(domain.SessionContext{
						UserID: session.UserID,
						Role:   session.Role,
					}, domain.ErrBadRequest)
			},
		},
		{
			name: "BadCase/UserNotFound",
			getBody: func() []byte {
				var creds domain.Credentials
				faker.FakeData(&creds.Password)
				creds.Email = "ferfg@fsf.ru"
				jsonBody, _ := json.Marshal(creds)
				return jsonBody
			},
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}

				uCase.On("Login", mock.Anything).Return(*session, 0, domain.ErrNotFound).Maybe()
			},
			status: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.AuthUsecase)
			var session domain.Session

			test.setUCaseExpectations(&session, mockUsecase)

			req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewReader(test.getBody()))
			req = req.WithContext(context.Background())
			rec := httptest.NewRecorder()

			if test.setAuth != nil {
				test.setAuth(req, mockUsecase, &session)
			}

			handler := &auth_http.AuthHandler{AuthUsecase: mockUsecase}
			handler.Login(rec, req)

			assert.Equal(t, test.status, rec.Code)

			if test.wantCookie {
				cookies := rec.Result().Cookies()
				assert.NotEmpty(t, cookies)
				assert.Equal(t, "session_token", cookies[0].Name)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name                 string
		setUCaseExpectations func(session *domain.Session, uCase *mocks.AuthUsecase)
		status               int
	}{
		{
			name: "GoodCase/Common",
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{
					Token:     "session_token",
					ExpiresAt: time.Now().Add(24 * time.Hour),
					UserID:    1,
					Role:      domain.Usr,
				}

				uCase.On("RetrieveSessionContext", mock.Anything).Return(domain.SessionContext{UserID: 1, Role: domain.Usr}, nil)
				uCase.On("Logout", mock.Anything).Return(nil)
			},
			status: http.StatusNoContent,
		},
		{
			name: "BadCase/NotAuthorized",
			setUCaseExpectations: func(session *domain.Session, uCase *mocks.AuthUsecase) {
				*session = domain.Session{}

				uCase.On("RetrieveSessionContext", mock.Anything).
					Return(domain.SessionContext{
						UserID: 0,
						Role:   domain.Usr,
					}, domain.ErrUnauthorized)
				uCase.On("Logout", "session_token").Return(nil).Maybe()

			},
			status: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockUsecase := new(mocks.AuthUsecase)
			session := &domain.Session{}
			test.setUCaseExpectations(session, mockUsecase)

			req := httptest.NewRequest("POST", "/logout", nil)
			req.AddCookie(&http.Cookie{Name: "session_token", Value: session.Token})
			rec := httptest.NewRecorder()

			mux := http.NewServeMux()
			auth_http.NewAuthHandler(mux, mockUsecase)

			mux.ServeHTTP(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name                 string
		getBody              func() []byte
		setUCaseExpectations func(uCase *mocks.AuthUsecase, session *domain.Session)
		status               int
		auth                 bool
		setAuth              func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session)
	}{
		{
			name: "GoodCase/Common",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				user.Email = "chgvj@mail.ru"
				user.Password = []byte("password")
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(1, nil)

				err := faker.FakeData(session)
				assert.NoError(t, err)
				session.ExpiresAt = time.Now().Add(24 * time.Hour)
				session.UserID = 1
				uCase.On("Login", mock.Anything).Return(*session, 1, nil)
			},
			status: http.StatusOK,
		},
		{
			name: "BadCase/AlreadyRegistered",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				user.Email = "chgvj@mail.ru"
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(0, domain.ErrAlreadyExists)

				err := faker.FakeData(session)
				assert.NoError(t, err)
				uCase.On("Login", mock.Anything).Return(*session, 1, nil).Maybe()
			},
			status: http.StatusConflict,
		},
		{
			name: "BadCase/EmptyJson",
			getBody: func() []byte {
				return []byte("{}")
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(0, nil).Maybe()
				uCase.On("Login", mock.Anything).Return(*session, 1, domain.ErrWrongCredentials).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/EmptyBody",
			getBody: func() []byte {
				return []byte("")
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(0, nil).Maybe()
				uCase.On("Login", mock.Anything).Return(*session, 1, domain.ErrBadRequest).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/InvalidJson",
			getBody: func() []byte {
				return []byte("{043895uith,redfsvdf;vfdv4er")
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(0, nil).Maybe()
				uCase.On("Login", mock.Anything).Return(*session, 1, domain.ErrBadRequest).Maybe()
			},
			status: http.StatusBadRequest,
		},
		{
			name: "BadCase/AlreadyAuthorized",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				user.Email = "chgvj@mail.ru"
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(0, errors.New("some")).Maybe()
				err := faker.FakeData(session)

				session.ExpiresAt = time.Now().Add(24 * time.Hour)
				session.Role = domain.Usr

				assert.NoError(t, err)
				uCase.On("Login", mock.Anything).Return(*session, 0, nil).Maybe()
			},
			status: http.StatusConflict,
			auth:   true,
			setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
				r.AddCookie(&http.Cookie{
					Name:     "session_token",
					Value:    session.Token,
					Expires:  session.ExpiresAt,
					Path:     "/",
					HttpOnly: true,
				})
				uCase.On("RetrieveSessionContext", mock.Anything).
					Return(domain.SessionContext{
						UserID: session.UserID,
						Role:   session.Role}, nil)
			},
		},
		{
			name: "GoodCase/AlreadyAuthorizedExpiredCookie",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				user.Email = "chgvj@mail.ru"
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(1, nil)
				err := faker.FakeData(session)

				session.ExpiresAt = time.Now().Add(24 * time.Hour)
				session.Role = domain.Usr

				assert.NoError(t, err)
				uCase.On("Login", mock.Anything).Return(*session, 1, nil)
			},
			status: http.StatusOK,
			auth:   true,
			setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
				r.AddCookie(&http.Cookie{
					Name:     "session_token",
					Value:    session.Token,
					Expires:  time.Now(),
					Path:     "/",
					HttpOnly: true,
				})
				uCase.On("RetrieveSessionContext", mock.Anything).
					Return(domain.SessionContext{
						UserID: session.UserID,
						Role:   session.Role}, domain.ErrUnauthorized).Maybe()
			},
		},
		{
			name: "GoodCase/AlreadyAuthorizedWrongCookie",
			getBody: func() []byte {
				var user domain.User
				faker.FakeData(&user)
				user.Email = "chgvj@mail.ru"
				jsonBody, _ := json.Marshal(user)
				return jsonBody
			},
			setUCaseExpectations: func(uCase *mocks.AuthUsecase, session *domain.Session) {
				uCase.On("Register", mock.Anything).Return(1, nil)

				err := faker.FakeData(session)
				session.ExpiresAt = time.Now().Add(24 * time.Hour)
				assert.NoError(t, err)
				uCase.On("Login", mock.Anything).Return(*session, 1, nil)
			},
			status: http.StatusOK,
			auth:   true,
			setAuth: func(r *http.Request, uCase *mocks.AuthUsecase, session *domain.Session) {
				r.AddCookie(&http.Cookie{
					Name:     "fevk",
					Value:    session.Token,
					Expires:  session.ExpiresAt,
					Path:     "/",
					HttpOnly: true,
				})
				uCase.On("RetrieveSessionContext", mock.Anything).Return(nil, errors.New("some error")).Maybe()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewReader(test.getBody()))
			assert.NoError(t, err)

			mockUCase := new(mocks.AuthUsecase)
			var mockSession domain.Session
			test.setUCaseExpectations(mockUCase, &mockSession)
			if test.auth {
				test.setAuth(req, mockUCase, &mockSession)
			}

			rec := httptest.NewRecorder()
			handler := &auth_http.AuthHandler{
				AuthUsecase: mockUCase,
			}

			handler.Register(rec, req)

			assert.Equal(t, test.status, rec.Code)
			mockUCase.AssertExpectations(t)

			var result *domain.Response
			err = json.NewDecoder(rec.Result().Body).Decode(&result)
			assert.NoError(t, err)

			if test.status < 300 {
				assert.NotEmpty(t, result.Body)
				assert.Empty(t, result.Err)
			} else {
				assert.Empty(t, result.Body)
				assert.NotEmpty(t, result.Err)
			}
		})
	}
}
