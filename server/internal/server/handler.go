package server

import (
	"context"

	"github.com/VicShved/pass-manager/server/internal/service"
	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
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

func (s GServer) Register(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var response pb.LoginResponse
	// Validate login / password
	err := utils.ValidateLoginPassword(in.Login, in.Password)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, err.Error())
	}
	// Register
	tokenStr, err := s.serv.Register(in.Login, in.Password)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, err.Error())
	}
	// Token str add to header
	authHeader := metadata.Pairs(authorizationTokenName, tokenStr)
	_ = grpc.SetHeader(ctx, authHeader)

	response.Token = tokenStr
	return &response, nil
}

func (s GServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var response pb.LoginResponse
	// Validate login / password
	err := utils.ValidateLoginPassword(in.Login, in.Password)
	if err != nil {
		return &response, status.Errorf(codes.InvalidArgument, err.Error())
	}
	// Login
	tokenStr, err := s.serv.Login(in.Login, in.Password)
	if err != nil {
		return &response, status.Errorf(codes.PermissionDenied, err.Error())
	}
	// Token str add to header
	authHeader := metadata.Pairs(authorizationTokenName, tokenStr)
	_ = grpc.SetHeader(ctx, authHeader)

	response.Token = tokenStr
	return &response, nil
}

func (s GServer) PostCard(ctx context.Context, in *pb.PostCardRequest) (*pb.PostCardResponse, error) {
	var response pb.PostCardResponse
	userID := getUserID(ctx)
	if userID == "" {
		return &response, status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	card := service.CardStruct{CardNumber: in.CardNumber, CardValid: in.Valid, CardCode: in.Code}
	s.serv.PostCard(&ctx, card, in.Description, userID)

	return &response, nil
}

func (s GServer) PostFile(stream grpc.ClientStreamingServer[pb.PostFileRequest, pb.PostFileResponse]) error {
	logger.Log.Info("Start PostFile")
	var fileName string
	var fileSize uint64
	ctx := stream.Context()
	userID := getUserID(ctx)
	if userID == "" {
		return status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	// for {
	// 	req, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if fileName == "" {
	// 		fileName = req.FileName
	// 		file, err := os.Create(newFileName)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		defer file.Close()
	// 	}
	// 	file, err := os.OpenFile(newFileName, os.O_APPEND|os.O_WRONLY, 0644)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	_, err = file.Write(req.GetChunk())
	// 	if err != nil {
	// 		file.Close()
	// 		return err
	// 	}
	// 	file.Close()
	// 	filelSize += uint64(len(req.GetChunk()))
	// }
	fileName, fileSize, err := s.serv.PostFile(stream, userID)
	if err != nil {
		return err
	}
	logger.Log.Info("Finish PostFile")
	return stream.SendAndClose(&pb.PostFileResponse{
		FileName: fileName,
		FileSize: fileSize,
	})

}
