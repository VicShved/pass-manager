package server

import (
	"context"

	"github.com/VicShved/pass-manager/server/internal/service"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Log.Debug("In authUnaryInterceptor")
	var userID string
	var token *jwt.Token
	// var tokenStr string
	var err error
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		logger.Log.Warn("authUnaryInterceptor hasnt metadata")
	}
	tokens := md.Get(authorizationTokenName)
	if len(tokens) > 0 {
		token, userID, err = service.ParseTokenUserID(tokens[0], config.ServerConfig.SecretKey)
		if err != nil {
			userID = ""
		}
		// Если токен не валидный
		if !token.Valid {
			logger.Log.Warn("Not valid token")
			userID = ""
		}
	}
	// Если кука не содержит ид пользователя, то не сильно ругаюсь
	if userID == "" {
		logger.Log.Warn("Empty userID")
	}
	md.Set("userID", userID)
	newCtx := metadata.NewIncomingContext(ctx, md)

	logger.Log.Debug("Out authUnaryInterceptor")
	return handler(newCtx, req)
}
