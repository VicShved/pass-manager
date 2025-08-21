package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	// "os"
	"strings"

	"github.com/VicShved/pass-manager/server/internal/repository"
	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/VicShved/pass-manager/server/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
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
	logger.Log.Debug("service.Login", zap.String("userID=", userID), zap.Error(err))
	if err != nil {
		return "", err
	}
	tokenStr, _ := GetJWTTokenString(&userID, s.conf.SecretKey)
	return tokenStr, err
}

func getNewFileName(fileName string) string {
	uuid, _ := GetNewUUID()
	splits := strings.Split(fileName, ".")
	ext := "unk"
	if len(splits) > 1 {
		ext = splits[len(splits)-1]
	}
	return uuid + "." + ext
}

type CardStruct struct {
	CardNumber string `json:"card_number"`
	CardValid  string `json:"card_valid"`
	CardCode   string `json:"card_code"`
}

func (s *PassManageService) PostCard(ctx *context.Context, card CardStruct, description string, userID string) error {
	var secretKey string
	buf, err := json.Marshal(card)
	if err != nil {
		return err
	}
	// secretKey = generateSecretKey() TODO add encrypt
	// buf := encrypt(buf, secretKey)
	iobuf := bytes.NewReader([]byte(buf))
	newFileName := getNewFileName("")
	s.repo.SaveData(iobuf)

}

func (s *PassManageService) PostFile(stream grpc.ClientStreamingServer[pb.PostFileRequest, pb.PostFileResponse], userID string) (string, uint64, error) {
	var fileName string
	var fileSize uint64
	newFileName := getNewFileName("")
	fileRepo, err := repository.GetFileStorageRepo("")
	if err != nil {
		logger.Log.Panic("PostFile", zap.Error(err))
	}
	fileStorage, err := fileRepo.GetFileStorage(newFileName)
	if err != nil {
		logger.Log.Panic("PostFile", zap.Error(err))
	}
	defer fileStorage.Close()
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return newFileName, fileSize, err
		}
		if fileName == "" {
			fileName = req.FileName
			_, err := fileStorage.OpenWrite()
			if err != nil {
				return newFileName, fileSize, err
			}
			// defer fileStorage.Close()
		}
		// file, err := os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, 0644)
		// if err != nil {
		// 	return newFileName, fileSize, err
		// }
		n, err := fileStorage.Write(req.GetChunk())
		if err != nil {
			logger.Log.Error("PostFile fileStorage.Write", zap.Error(err))

			return newFileName, fileSize, err
		}
		fileSize += uint64(n)
	}
	return newFileName, fileSize, nil
}
