package repository

import (
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
)

var users = map[string][3]string{"1": {"1", "CxTVAaWURCoBxoWVQbyz6BZNGD0yk3uFGDVEL2nVyU4", "userID1"}}

type TestRepository struct {
	conf *config.ServerConfigStruct
}

// GetGormRepo(dns string)
func GetTestRepo(conf *config.ServerConfigStruct) (*TestRepository, error) {
	repo := &TestRepository{
		conf: conf,
	}
	return repo, nil
}

// CloseConn Close connection
func (r TestRepository) CloseConn() error {
	return nil
}

func (r TestRepository) Register(userID string, login string, hashPassword string) error {
	logger.Log.Debug("Register", zap.String("login", login), zap.String("hashPassword", hashPassword))
	// users[login] = make([]string{login, hashPassword, userID})
	return nil
}

func (r TestRepository) Login(login string, hashPassword string) (string, error) {
	logger.Log.Debug("Login", zap.String("login", login), zap.String("hashPassword", hashPassword))
	value, exists := users[login]
	if !exists {
		return "", ErrLoginPassword
	}
	return value[2], nil
}
