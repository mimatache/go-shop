package authentication

import (
	"fmt"

	"github.com/mimatache/go-shop/internal/store"
)

//go:generate mockgen -source ./authentication.go -destination mocks/authentication.go

// UserRegistry abstracts aways the storage from the logic
type UserRegistry interface {
	GetPasswordFor(email string) (string, error)
}

type invalidCredentials struct {
	msg string
}

func (i invalidCredentials) Error() string {
	return i.msg
}

// IsInvalidCredentialsError verifies if a given error refers to incalid credentials
func IsInvalidCredentialsError(err error) bool {
	switch err.(type) {
	case invalidCredentials:
		return true
	default:
		return false
	}
}

// NewInvalidCredentials created a new Invalid Credentials error
func NewInvalidCredentials(username string) error {
	return invalidCredentials{msg: fmt.Sprintf("invalid credentials for %s", username)}
}

//New creates a new User instance
func New(storage UserRegistry) *User {
	return &User{storage}
}

// User manages any user related actions
type User struct {
	storage UserRegistry
}

// IsValid checks if a username
func (u *User) IsValid(username, password string) error {
	passwd, err := u.storage.GetPasswordFor(username)
	if err != nil {
		if store.IsNotFoundError(err) {
			return NewInvalidCredentials(username)
		}
		return err
	}
	if passwd != password {
		return NewInvalidCredentials(username)
	}
	return nil
}

// Deprecated: do not use this
func (u *User) GetEmailForUser(ID uint) (string, error) {
	return "", fmt.Errorf("not implemented")
}
