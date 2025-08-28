package server

import (
	"context"

	"github.com/VicShved/pass-manager/server/internal/repository"
	"github.com/VicShved/pass-manager/server/internal/service"
	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/VicShved/pass-manager/server/pkg/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func getUserID(ctx context.Context) string {
	userID := ""
	md, _ := metadata.FromIncomingContext(ctx)
	logger.Log.Debug("getUserID ", zap.Any("metadata ", md))
	users := md.Get("userID")
	if len(users) > 0 {
		userID = users[0]
	}
	return userID
}

// Register user register handler
func (s GServer) Register(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var response pb.LoginResponse
	// Validate login / password
	err := utils.ValidateLoginPassword(in.Login, in.Password)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, err.Error())
	}
	// Register
	tokenStr, err := s.serv.Register(ctx, in.Login, in.Password)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, err.Error())
	}
	// Token str add to header
	authHeader := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	_ = grpc.SetHeader(ctx, authHeader)

	response.Token = tokenStr
	return &response, nil
}

// Login user handler
func (s GServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var response pb.LoginResponse
	// Validate login / password
	err := utils.ValidateLoginPassword(in.Login, in.Password)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, err.Error())
	}
	// Login
	tokenStr, err := s.serv.Login(ctx, in.Login, in.Password)
	if err != nil {
		return &response, status.Errorf(codes.PermissionDenied, err.Error())
	}
	// Token str add to header
	authHeader := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	_ = grpc.SetHeader(ctx, authHeader)

	response.Token = tokenStr
	return &response, nil
}

// PostCard save card handler
func (s GServer) PostCard(ctx context.Context, in *pb.PostCardRequest) (*pb.PostDataResponse, error) {
	var response pb.PostDataResponse
	userID := getUserID(ctx)
	if userID == "" {
		return &response, status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	card := service.CardStruct{CardNumber: in.CardNumber, CardValid: in.Valid, CardCode: in.Code}
	rowID, length, err := s.serv.PostCard(ctx, userID, card, in.Description)
	if err != nil {
		return &response, err
	}
	response.RowId = rowID
	response.Length = length

	return &response, nil
}

// PostLogPass save login/password handler
func (s GServer) PostLogPass(ctx context.Context, in *pb.PostLogPassRequest) (*pb.PostDataResponse, error) {
	var response pb.PostDataResponse
	userID := getUserID(ctx)
	if userID == "" {
		return &response, status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	logPass := service.LogPassStruct{Login: in.Login, Password: in.Password}
	rowID, length, err := s.serv.PostLogPass(ctx, userID, logPass, in.Description)
	if err != nil {
		return &response, err
	}
	response.RowId = rowID
	response.Length = length

	return &response, nil
}

// GetCard get saved card data handler
func (s GServer) GetCard(ctx context.Context, in *pb.GetDataRequest) (*pb.GetCardResponse, error) {
	var response pb.GetCardResponse
	userID := getUserID(ctx)
	if userID == "" {
		return &response, status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	card, err := s.serv.GetCard(ctx, userID, in.RowId)
	if err != nil {
		return &response, err
	}
	response.CardNumber = card.CardNumber
	response.Valid = card.CardValid
	response.Code = card.CardCode

	return &response, nil

}

// GetLogPass get saved login/password data heandle
func (s GServer) GetLogPass(ctx context.Context, in *pb.GetDataRequest) (*pb.GetLogPassResponse, error) {
	var response pb.GetLogPassResponse
	userID := getUserID(ctx)
	if userID == "" {
		return &response, status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	logPass, err := s.serv.GetLogPass(ctx, userID, in.RowId)
	if err != nil {
		return &response, err
	}
	response.Login = logPass.Login
	response.Password = logPass.Password

	return &response, nil
}

// PostFile save user file heandler
func (s GServer) PostFile(stream grpc.ClientStreamingServer[pb.PostFileRequest, pb.PostDataResponse]) error {
	logger.Log.Info("Start PostFile")
	var rowID uint32
	var fileSize uint64
	ctx := stream.Context()
	userID := getUserID(ctx)
	if userID == "" {
		return status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	rowID, fileSize, err := s.serv.PostFile(ctx, stream, userID, repository.DataTypeFile)
	if err != nil {
		return err
	}
	logger.Log.Info("Finish PostFile")
	return stream.SendAndClose(&pb.PostDataResponse{
		RowId:  rowID,
		Length: fileSize,
	})
}

// GetFile get saved file heandler
func (s GServer) GetFile(in *pb.GetDataRequest, stream grpc.ServerStreamingServer[pb.GetFileResponse]) error {
	logger.Log.Info("Start GetFile")
	rowID := in.RowId
	ctx := stream.Context()
	userID := getUserID(ctx)
	if userID == "" {
		return status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	fileSize, err := s.serv.GetFile(ctx, userID, rowID, stream)
	if err != nil {
		return err
	}
	logger.Log.Info("Finish GetFile", zap.Uint64("FileSize", fileSize))
	return nil
}

func (s GServer) GetDataInfo(ctx context.Context, in *pb.GetDataInfoRequest) (*pb.GetDataInfoResponse, error) {
	var response pb.GetDataInfoResponse
	userID := getUserID(ctx)
	if userID == "" {
		return &response, status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	dataInfos, err := s.serv.GetDataInfo(ctx, userID, int(in.DataType))
	if err != nil {
		return &response, err
	}
	results := make([]*pb.DataInfo, len(dataInfos))
	for i, dataInfo := range dataInfos {
		results[i].RowId = dataInfo.RowID
		results[i].Desc = dataInfo.Desc
		results[i].DataType = dataInfo.DataType
	}
	response = pb.GetDataInfoResponse{UserData: results}
	return &response, nil

}
