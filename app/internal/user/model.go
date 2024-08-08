package user

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserID string
type Token string

type User struct {
	ID       UserID `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (u *User) HashPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) GenerateToken() Token {
	return Token(uuid.New().String())
}
