package repository

import (
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
)

type userStruct struct {
	login        string
	hashPassword string
	userID       string
}

var users = map[string]userStruct{
	"1": {
		login:        "1",
		hashPassword: "CxTVAaWURCoBxoWVQbyz6BZNGD0yk3uFGDVEL2nVyU4",
		userID:       "userID1"},
}

type TestRepository struct {
	conf     *config.ServerConfigStruct
	fileRepo FileStoragerRepoInterface
}

// GetGormRepo(dns string)
func GetTestRepo(conf *config.ServerConfigStruct, fileRepo FileStoragerRepoInterface) (*TestRepository, error) {
	repo := &TestRepository{
		conf:     conf,
		fileRepo: fileRepo,
	}
	return repo, nil
}

// CloseConn Close connection
func (r TestRepository) CloseConn() error {
	return nil
}

func (r TestRepository) Register(userID string, login string, hashPassword string) error {
	logger.Log.Debug("Register user = ", zap.String("userID", userID), zap.String("login", login), zap.String("hashPassword", hashPassword))
	_, exists := users[login]
	if exists {
		return ErrLoginConflict
	}
	users[login] = userStruct{login: login, hashPassword: hashPassword, userID: userID}
	return nil
}

func (r TestRepository) Login(login string, hashPassword string) (string, error) {
	logger.Log.Debug("Login", zap.String("login", login), zap.String("hashPassword", hashPassword))
	value, exists := users[login]
	if !exists {
		return "", ErrLoginPassword
	}
	if value.hashPassword != hashPassword {
		return "", ErrLoginPassword
	}
	return value.userID, nil
}
