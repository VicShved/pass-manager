package server

import (
	"context"

	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Log.Debug("In authUnaryInterceptor")
	var userID string
	var token *jwt.Token
	var tokenStr string
	var err error
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		logger.Log.Warn("authUnaryInterceptor has`nt metadata")
	}
	tokens := md.Get(authorizationTokenName)
	if len(tokens) == 0 {
		userID, tokenStr = getNewUserIDToken(config.ServerConfig.SecretKey)
	}
	if len(tokens) > 0 {
		token, userID, err = parseTokenUserID(tokens[0], config.ServerConfig.SecretKey)
		if err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "Доступ запрещен")
		}
		// Если токен не валидный,  то создаю нвый userID
		if !token.Valid {
			logger.Log.Warn("Not valid token")
			userID, tokenStr = getNewUserIDToken(config.ServerConfig.SecretKey)
		}
	}
	// Если кука не содержит ид пользователя, то возвращаю 401
	if userID == "" {
		logger.Log.Warn("Empty userID")
		return nil, status.Errorf(codes.PermissionDenied, "Доступ запрещен")
	}
	md.Set("userID", userID)
	md.Set(authorizationTokenName, tokenStr)
	newCtx := metadata.NewIncomingContext(ctx, md)

	logger.Log.Debug("Exit from authUnaryInterceptor", zap.String("token:", tokenStr))
	authHeader := metadata.Pairs(authorizationTokenName, tokenStr)
	_ = grpc.SetHeader(newCtx, authHeader)
	return handler(newCtx, req)
}
