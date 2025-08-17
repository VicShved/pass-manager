package service

import (
	"github.com/VicShved/pass-manager/server/internal/repository"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/VicShved/pass-manager/server/pkg/utils"
	"go.uber.org/zap"
)

// PassManageService struct
type PassManageService struct {
	repo repository.RepoInterface
	conf *config.ServerConfigStruct
}

// GetService(repo repository.RepoInterface, config)
func GetService(repo repository.RepoInterface, conf *config.ServerConfigStruct) *PassManageService {
	return &PassManageService{repo: repo, conf: conf}
}

func (s *PassManageService) Register(login string, password string) (string, error) {
	userID, _ := GetNewUUID()
	err := (*s).repo.Register(userID, login, utils.HashSha256(password))
	if err != nil {
		return "", err
	}
	tokenStr, _ := GetJWTTokenString(&userID, s.conf.SecretKey)
	logger.Log.Debug("", zap.String("userID", userID), zap.Any("err", err))
	return tokenStr, err
}

func (s *PassManageService) Login(login string, password string) (string, error) {
	userID, err := (*s).repo.Login(login, utils.HashSha256(password))
	if err != nil {
		return "", err
	}
	tokenStr, _ := GetJWTTokenString(&userID, s.conf.SecretKey)

	return tokenStr, err
}
