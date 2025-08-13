package service

import (
	"github.com/VicShved/pass-manager/server/internal/repository"
	"github.com/VicShved/pass-manager/server/pkg/config"
)

// Service struct
type Service struct {
	repo repository.RepoInterface
	conf *config.ServerConfigStruct
}

// GetService(repo repository.RepoInterface, baseurl string)
func GetService(repo repository.RepoInterface, conf *config.ServerConfigStruct) *Service {
	return &Service{repo: repo, conf: conf}
}
