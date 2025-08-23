package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io"

	"github.com/VicShved/pass-manager/server/internal/repository"
	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/VicShved/pass-manager/server/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// CardStruct
type CardStruct struct {
	CardNumber string `json:"card_number"`
	CardValid  string `json:"card_valid"`
	CardCode   string `json:"card_code"`
}

// LoginPassword Struct
type LogPassStruct struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// PassManageService struct
type PassManageService struct {
	repo repository.RepoInterface
	conf *config.ServerConfigStruct
}

// GetService(repo repository.RepoInterface, config)
func GetService(repo repository.RepoInterface, conf *config.ServerConfigStruct) *PassManageService {
	return &PassManageService{repo: repo, conf: conf}
}

func (s *PassManageService) Register(ctx context.Context, login string, password string) (string, error) {
	userID, _ := GetNewUUID()
	err := (*s).repo.Register(ctx, userID, login, utils.HashSha256(password))
	if err != nil {
		return "", err
	}
	tokenStr, _ := GetJWTTokenString(&userID, s.conf.SecretKey)
	logger.Log.Debug("", zap.String("userID", userID), zap.Any("err", err))
	return tokenStr, err
}

func (s *PassManageService) Login(ctx context.Context, login string, password string) (string, error) {
	userID, err := (*s).repo.Login(ctx, login, utils.HashSha256(password))
	logger.Log.Debug("service.Login", zap.String("userID=", userID), zap.Error(err))
	if err != nil {
		return "", err
	}
	tokenStr, _ := GetJWTTokenString(&userID, s.conf.SecretKey)
	return tokenStr, err
}

func (s *PassManageService) SaveData(ctx context.Context, userID string, description string, dataType repository.DataType, buf []byte) (rowID uint32, fileSize uint64, err error) {

	secretKey, err := generateHexKeyString(lecretKeyLength)
	if err != nil {
		return rowID, fileSize, err
	}

	bufString, err := encryptData2Hex(buf, secretKey)
	iobuf := bytes.NewReader([]byte(bufString))
	newFileName := getNewFileName("")

	buf = make([]byte, 1024)
	fileStorage, err := s.repo.GetFileStorage(newFileName)
	if err != nil {
		logger.Log.Panic("PostFile", zap.Error(err))
	}
	defer fileStorage.Close()
	fileStorage.OpenWrite()
	for {
		n, err := iobuf.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return rowID, fileSize, err
		}
		_, err = fileStorage.Write(buf[:n])
		if err != nil {
			logger.Log.Error("SaveData", zap.Error(err))

			return rowID, fileSize, err
		}
		fileSize += uint64(n)
	}
	rowID, err = s.repo.SaveData(ctx, userID, description, string(dataType), newFileName, secretKey)
	return rowID, fileSize, err
}

func (s *PassManageService) PostCard(ctx context.Context, userID string, card CardStruct, description string) (uint32, uint64, error) {
	var rowID uint32
	var fileSize uint64
	buf, err := json.Marshal(card)
	if err != nil {
		return rowID, fileSize, err
	}

	rowID, fileSize, err = s.SaveData(ctx, userID, description, repository.DataTypeCard, buf)
	return rowID, fileSize, err
}

func (s *PassManageService) PostLogPass(ctx context.Context, userID string, logPass LogPassStruct, description string) (uint32, uint64, error) {
	var rowID uint32
	var fileSize uint64
	buf, err := json.Marshal(logPass)
	if err != nil {
		return rowID, fileSize, err
	}
	rowID, fileSize, err = s.SaveData(ctx, userID, description, repository.DataTypeLoginPassword, buf)
	return rowID, fileSize, err
}

func (s *PassManageService) GetData(ctx context.Context, userID string, rowID uint32) (buf []byte, err error) {
	userData, err := s.repo.GetUserData(ctx, userID, rowID)
	if err != nil {
		return buf, err
	}
	fileStorage, err := s.repo.GetFileStorage(userData.FileName)
	if err != nil {
		return buf, err
	}
	file, err := fileStorage.OpenRead()
	if err != nil {
		return buf, err
	}
	defer fileStorage.Close()
	encbuf, err := io.ReadAll(file)
	if err != nil {
		return buf, err
	}
	buf, err = decryptHexData(string(encbuf), userData.SecretKey)
	return buf, err

}

func (s *PassManageService) GetCard(ctx context.Context, userID string, rowID uint32) (card CardStruct, err error) {
	buf, err := s.GetData(ctx, userID, rowID)
	if err != nil {
		return card, err
	}
	err = json.Unmarshal(buf, &card)
	return card, err
}

func (s *PassManageService) GetLogPass(ctx context.Context, userID string, rowID uint32) (logPass LogPassStruct, err error) {
	buf, err := s.GetData(ctx, userID, rowID)
	if err != nil {
		return logPass, err
	}
	err = json.Unmarshal(buf, &logPass)
	return logPass, err
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
