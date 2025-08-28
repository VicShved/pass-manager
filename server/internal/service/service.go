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

// GetService return prepared service
func GetService(repo repository.RepoInterface, conf *config.ServerConfigStruct) *PassManageService {
	return &PassManageService{repo: repo, conf: conf}
}

// Register register new user by login/password
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

// Login user by login/password
func (s *PassManageService) Login(ctx context.Context, login string, password string) (string, error) {
	userID, err := (*s).repo.Login(ctx, login, utils.HashSha256(password))
	logger.Log.Debug("service.Login", zap.String("userID=", userID), zap.Error(err))
	if err != nil {
		return "", err
	}
	tokenStr, _ := GetJWTTokenString(&userID, s.conf.SecretKey)
	return tokenStr, err
}

// SaveData save user data
func (s *PassManageService) SaveData(
	ctx context.Context,
	userID string,
	description string,
	dataType repository.DataType,
	buf []byte,
) (rowID uint32, fileSize uint64, err error) {

	secretKey, err := generateHexKeyString(secretKeyLength)
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

// PostCard save card data
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

// PostLogPass save user login/password data
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

// GetData retirn user data
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

// GetCard return user card data
func (s *PassManageService) GetCard(ctx context.Context, userID string, rowID uint32) (card CardStruct, err error) {
	buf, err := s.GetData(ctx, userID, rowID)
	if err != nil {
		return card, err
	}
	err = json.Unmarshal(buf, &card)
	return card, err
}

// GetLogPass return user login/password data
func (s *PassManageService) GetLogPass(ctx context.Context, userID string, rowID uint32) (logPass LogPassStruct, err error) {
	buf, err := s.GetData(ctx, userID, rowID)
	if err != nil {
		return logPass, err
	}
	err = json.Unmarshal(buf, &logPass)
	return logPass, err
}

// PostFile save user file
func (s *PassManageService) PostFile(
	ctx context.Context,
	stream grpc.ClientStreamingServer[pb.PostFileRequest, pb.PostDataResponse],
	userID string,
	dataType repository.DataType,
) (rowID uint32, length uint64, err error) {
	var secretKey string
	var fileName string
	var description string
	var fileSize uint64
	newFileName := getNewFileName("")
	// get filestorage
	fileStorage, err := s.repo.GetFileStorage(newFileName)
	if err != nil {
		logger.Log.Error("PostFile", zap.Error(err))
		return rowID, fileSize, err
	}
	defer fileStorage.Close()
	// read stream to file
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return rowID, fileSize, err
		}
		if fileName == "" {
			fileName = req.FileName
			description = req.Description
			_, err := fileStorage.OpenWrite()
			if err != nil {
				return rowID, fileSize, err
			}
		}
		n, err := fileStorage.Write(req.GetChunk())
		if err != nil {
			logger.Log.Error("PostFile fileStorage.Write", zap.Error(err))
			return rowID, fileSize, err
		}
		fileSize += uint64(n)
	}
	rowID, err = s.repo.SaveData(ctx, userID, description, string(dataType), newFileName, secretKey)

	return rowID, fileSize, nil
}

// GetFile return saved user file
func (s *PassManageService) GetFile(
	ctx context.Context,
	userID string,
	rowID uint32,
	stream grpc.ServerStreamingServer[pb.GetFileResponse],

) (fileSize uint64, err error) {
	userData, err := s.repo.GetUserData(ctx, userID, rowID)
	if err != nil {
		return fileSize, err
	}
	// get filestorage
	fileStorage, err := s.repo.GetFileStorage(userData.FileName)
	if err != nil {
		logger.Log.Error("service.GetFile", zap.Error(err))
		return fileSize, err
	}
	_, err = fileStorage.OpenRead()
	if err != nil {
		return fileSize, err
	}
	defer fileStorage.Close()
	// read stream to file
	buf := make([]byte, 4096) //  4096 вынести в настройки
	var n int
	for {
		n, err = fileStorage.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fileSize, err
		}
		message := pb.GetFileResponse{Chunk: buf[:n], FileName: userData.FileName}
		err = stream.Send(&message)
		if err != nil {
			return fileSize, err
		}
		fileSize += uint64(n)
	}
	return fileSize, nil
}

func (s *PassManageService) GetDataInfo(ctx context.Context, userID string, dataType int) ([]UserData, error) {
	convert := map[int]string{
		0: "",
		1: string(repository.DataTypeLoginPassword),
		2: string(repository.DataTypeCard),
		3: string(repository.DataTypeFile),
	}
	dataTypeStr := convert[dataType]
	userDatas, err := s.repo.GetUserDatas(ctx, userID, dataTypeStr)
	if err != nil {
		return nil, err
	}
	results := make([]UserData, len(userDatas))
	for i, ud := range userDatas {
		results[i].RowID = uint32(ud.ID)
		results[i].Desc = ud.Description
		results[i].DataType = ud.DataType
	}
	return results, err

}
