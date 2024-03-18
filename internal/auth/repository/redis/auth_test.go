package redis_test

import (
	"encoding/json"
	"github.com/ellexo2456/FilmLib/internal/auth/repository/redis"
	"github.com/ellexo2456/FilmLib/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
)

func TestAdd(t *testing.T) {
	tests := []struct {
		name    string
		session domain.Session
		good    bool
		err     error
	}{
		{
			name: "GoodCase/Common",
			session: domain.Session{
				Token:     "123",
				ExpiresAt: time.Now().Add(24 * time.Hour),
				UserID:    1,
				Role:      domain.Usr,
			},
			good: true,
			err:  nil,
		},
		{
			name: "GoodCase/SameToken",
			session: domain.Session{
				Token:     "123",
				ExpiresAt: time.Now().Add(24 * time.Hour),
				UserID:    1,
				Role:      domain.Usr,
			},
			good: true,
		},
		{
			name:    "BadCase/EmptyToken",
			session: domain.Session{},
			err:     domain.ErrInvalidToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			r := redis.NewSessionRedisRepository(db)

			if test.good {
				jsonData, _ := json.Marshal(domain.SessionContext{
					UserID: test.session.UserID,
					Role:   test.session.Role,
				})
				mock.ExpectSet(test.session.Token, jsonData, test.session.ExpiresAt.Sub(time.Now())).SetVal("")
			}

			err := r.Add(test.session)

			if test.good {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, test.err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDeleteByToken(t *testing.T) {
	tests := []struct {
		name  string
		token string
		good  bool
		err   error
	}{
		{
			name:  "GoodCase/Common",
			token: "12312dcdscsad",
			good:  true,
			err:   nil,
		},
		{
			name:  "BadCase/EmptyToken",
			token: "",
			err:   domain.ErrInvalidToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			db, mock := redismock.NewClientMock()
			defer db.Close()

			r := redis.NewSessionRedisRepository(db)

			if test.good {
				mock.ExpectDel(test.token).SetVal(1)
			}

			err := r.DeleteByToken(test.token)

			if test.good {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, test.err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
