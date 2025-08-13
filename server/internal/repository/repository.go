package repository

import (
	"errors"
)

// RepoInterface interface
type RepoInterface interface {
	Empty() error
	CloseConn()
}

// ErrPKConflict
var ErrPKConflict = errors.New("PK conflict")
