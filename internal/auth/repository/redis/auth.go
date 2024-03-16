package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ellexo2456/FilmLib/internal/domain"
)

type sessionRedisRepository struct {
	client *redis.Client
}

func NewSessionRedisRepository(client *redis.Client) domain.SessionRepository {
	return &sessionRedisRepository{client}
}

func (s *sessionRedisRepository) Add(session domain.Session) error {
	if session.Token == "" {
		return domain.ErrInvalidToken
	}

	jsonData, err := json.Marshal(domain.SessionContext{
		UserID: session.UserID,
		Role:   session.Role,
	})
	if err != nil {
		return domain.ErrInvalidToken
	}

	duration := session.ExpiresAt.Sub(time.Now())
	err = s.client.Set(context.TODO(), session.Token, jsonData, duration).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *sessionRedisRepository) DeleteByToken(token string) error {
	if token == "" {
		return domain.ErrInvalidToken
	}

	err := s.client.Del(context.Background(), token).Err()
	if err != nil {
		return err
	}

	return nil
}

func (s *sessionRedisRepository) GetSessionContext(token string) (domain.SessionContext, error) {
	if token == "" {
		return domain.SessionContext{}, domain.ErrInvalidToken
	}

	r, err := s.client.Get(context.Background(), token).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return domain.SessionContext{}, domain.ErrNotFound
		}
		return domain.SessionContext{}, err
	}

	var sc domain.SessionContext
	err = json.Unmarshal([]byte(r), &sc)
	if err != nil {
		return domain.SessionContext{}, err
	}

	return sc, nil
}
