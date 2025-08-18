package server

import (
	"context"

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
		return nil, status.Errorf(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь")
	}
	return &response, nil
}
