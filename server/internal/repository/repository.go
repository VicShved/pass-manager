package repository

import (
	"context"
	"errors"
)

// RepoInterface interface
type RepoInterface interface {
	Register(ctx context.Context, userID string, login string, hashPassword string) error
	Login(ctx context.Context, login string, hashPassword string) (string, error)
	GetFileStorage(fileName string) (FileStoragerInterface, error)
	SaveData(ctx context.Context, userID string, desc string, dataType string, fileName string, secretKey string) (uint32, error)
	GetUserData(ctx context.Context, userID string, rowID uint32) (UserData, error)
	GetUserDatas(ctx context.Context, userID string, dataType string) ([]UserData, error)
	CloseConn() error
}

// ErrLoginConflict this login exists
var ErrLoginConflict = errors.New("Login conflict")

// ErrLoginPassword login/password invalid
var ErrLoginPassword = errors.New("Login/Password error")

// ErrUserIdFromToken invlid userID from token
var ErrUserIdFromToken = errors.New("DB hasnt User from token")
