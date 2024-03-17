package domain

import (
	"time"
)

type Role int

const (
	Usr Role = iota
	Moder
)

type Key string

const SessionContextKey Key = "SessionContextKey"

type Credentials struct {
	Password []byte `json:"password"`
	Email    string `json:"email"`
}

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Password  []byte `json:"password"`
	Email     string `json:"email"`
	ImagePath string `json:"imagePath"`
	ImageData []byte `json:"imageData"`
	Role      Role
}

type Session struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    int       `json:"-"`
	Role      Role      `json:"-"`
}

type SessionContext struct {
	UserID int
	Role   Role
}

type AuthUsecase interface {
	Login(credentials Credentials) (Session, int, error)
	Logout(token string) error
	Register(user User) (int, error)
	RetrieveSessionContext(token string) (SessionContext, error)
}

type AuthRepository interface {
	GetByEmail(email string) (User, error)
	AddUser(user User) (int, error)
	UserExists(email string) (bool, error)
}

type SessionRepository interface {
	Add(session Session) error
	DeleteByToken(token string) error
	GetSessionContext(token string) (SessionContext, error)
}
