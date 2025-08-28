package server

import (
	"context"

	"github.com/VicShved/pass-manager/server/internal/service"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/golang-jwt/jwt/v4"

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
	tokens := md.Get(config.AuthorizationTokenName)
	logger.Log.Debug("getUserFromContext", zap.Any("Tokens=", tokens))
	if (len(tokens) > 0) && (len(tokens[0]) > 0) {
		token, userID, err = service.ParseTokenUserID(tokens[0], config.ServerConfig.SecretKey)
		if err != nil {
			userID = ""
		} else {
			// Если токен не валидный
			if !token.Valid {
				logger.Log.Warn("Not valid token")
				userID = ""
			}
		}
	}
	// Если токен не содержит ид пользователя, то не сильно ругаюсь
	if userID == "" {
		logger.Log.Warn("Empty userID")
	}
	return userID
}

type newStreamStruct struct {
	grpc.ServerStream
	ctx context.Context
}

// Context return stream context
func (s newStreamStruct) Context() context.Context {
	return s.ctx
}

// AuthStreamInterceptor get token from stream context, get userID from token and set userID to header (metaData)
func AuthStreamInterceptor(srv any, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	logger.Log.Debug("In authStreamInterceptor")
	ctx := stream.Context()
	userID := getUserFromContext(ctx)
	md, _ := metadata.FromIncomingContext(ctx)
	md.Set("userID", userID)
	newCtx := metadata.NewIncomingContext(ctx, md)
	newStream := grpc.ServerStream(newStreamStruct{ServerStream: stream, ctx: newCtx})
	logger.Log.Debug("Out authStreamInterceptor")
	return handler(srv, newStream)
}

// AuthUnaryInterceptor get token from context, get userID from token and set userID to header (metaData)
func AuthUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logger.Log.Debug("In authUnaryInterceptor")
	userID := getUserFromContext(ctx)
	md, _ := metadata.FromIncomingContext(ctx)
	md.Set("userID", userID)
	newCtx := metadata.NewIncomingContext(ctx, md)
	logger.Log.Debug("Out authUnaryInterceptor")
	return handler(newCtx, req)
}
