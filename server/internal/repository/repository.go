package repository

import (
	"errors"
)

// RepoInterface interface
type RepoInterface interface {
	Register(userID string, login string, hashPassword string) error
	Login(login string, hashPassword string) (string, error)
	CloseConn() error
}

var ErrLoginConflict = errors.New("Login conflict")
var ErrLoginPassword = errors.New("Login/Password error")
