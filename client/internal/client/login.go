package client

import (
	"context"
	"errors"

	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DoLogin(login string, password string) (grpcStatus codes.Code, tokenStr string, err error) {
	ctx := context.Background()
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("doLogin", zap.Error(err))
		return grpcStatus, tokenStr, err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	var header metadata.MD
	reqData := pb.LoginRequest{Login: login, Password: password}
	_, err = client.Login(ctx, &reqData, grpc.Header(&header))
	if status.Code(err) != codes.OK {
		logger.Log.Warn("doLogin", zap.String("Error", err.Error()))
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, tokenStr, err
	}
	authToken := header.Get(AuthorizationTokenName)[0]
	if len(authToken) == 0 {
		return grpcStatus, tokenStr, errors.New("Сервер не возвратил auth token")
	}
	return grpcStatus, authToken, nil
}
