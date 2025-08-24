package client

import (
	"context"
	"errors"

	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const serverAddress string = "localhost:7777"

func doRegister(login string, password string) (grpcCode codes.Code, tokenStr string, err error) {
	ctx := context.Background()
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("doRegister", zap.Error(err))
		return grpcCode, tokenStr, err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	var header metadata.MD
	reqData := pb.LoginRequest{Login: login, Password: password}
	_, err = client.Register(ctx, &reqData, grpc.Header(&header))
	if status.Code(err) != codes.OK {
		logger.Log.Warn("doRegister", zap.String("Error", err.Error()))
		st, ok := status.FromError(err)
		if ok {
			grpcCode = st.Code()
		}
		return grpcCode, tokenStr, err
	}
	authToken := header.Get(config.AuthorizationTokenName)[0]
	if len(authToken) == 0 {
		return 999, "", errors.New("Сервер не возвратил auth token")
	}
	return grpcCode, authToken, nil
}
