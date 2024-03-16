package middleware

import (
	"errors"
	"github.com/ellexo2456/FilmLib/internal/domain"
	"net/http"
	"time"
)

type Middleware struct {
	authUsecase domain.AuthUsecase
}

func New(au domain.AuthUsecase) *Middleware {
	return &Middleware{authUsecase: au}
}

func (m *Middleware) IsAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				domain.WriteError(w, err.Error(), http.StatusUnauthorized)
				return
			}

			domain.WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		if c.Expires.After(time.Now()) {
			domain.WriteError(w, "cookie is expired", http.StatusUnauthorized)
		}

		sessionToken := c.Value
		exists, err := m.authUsecase.IsAuth(sessionToken)
		if err != nil {
			domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
			return
		}
		if !exists {
			domain.WriteError(w, domain.ErrUnauthorized.Error(), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
