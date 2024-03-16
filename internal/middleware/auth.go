package middleware

import (
	"errors"
	"github.com/ellexo2456/FilmLib/internal/domain"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

type AuthMiddleware struct {
	authUsecase domain.AuthUsecase
}

func NewAuth(au domain.AuthUsecase) *AuthMiddleware {
	return &AuthMiddleware{authUsecase: au}
}

func (m *AuthMiddleware) IsAuth(next http.Handler) http.Handler {
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
		sc, err := m.authUsecase.RetrieveSessionContext(sessionToken)
		if err != nil {
			domain.WriteError(w, err.Error(), domain.GetStatusCode(err))
			return
		}
		if sc.UserID == 0 {
			domain.WriteError(w, "You`re unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), domain.SessionContextKey, sc)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
