package server

import (
	"context"

	"github.com/VicShved/pass-manager/server/internal/service"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/golang-jwt/jwt/v4"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func getUserFromContext(ctx context.Context) (userID string) {
	var token *jwt.Token
	var err error
	md, exists := metadata.FromIncomingContext(ctx)
	if !exists {
		logger.Log.Warn("authUnaryInterceptor hasnt metadata")
	}
	tokens := md.Get(authorizationTokenName)
	logger.Log.Debug("AuthUnaryInterceptor", zap.Any("Tokens=", tokens))
	if (len(tokens) > 0) && (len(tokens[0]) > 0) {
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
	// Если токен не содержит ид пользователя, то не сильно ругаюсь
	if userID == "" {
		logger.Log.Warn("Empty userID")
	}
	return userID
}

func AuthStreamInterceptor(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := stream.Context()
	userID := getUserFromContext(ctx)
	newStream := grpc_middleware.WrapServerStream(stream)
	md, _ := metadata.FromIncomingContext(ctx)
	md.Set("userID", userID)
	newStream.WrappedContext = metadata.NewIncomingContext(ctx, md)
	return handler(srv, newStream)
}

func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Log.Debug("In authUnaryInterceptor")
	userID := getUserFromContext(ctx)
	md, _ := metadata.FromIncomingContext(ctx)
	md.Set("userID", userID)
	newCtx := metadata.NewIncomingContext(ctx, md)
	logger.Log.Debug("Out authUnaryInterceptor")
	return handler(newCtx, req)
}
