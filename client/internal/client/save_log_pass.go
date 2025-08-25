package client

import (
	"context"

	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// DoSaveLogPass save login and password to server
func DoSaveLogPass(tokenStr string, login string, password string, description string) (gprcStatus codes.Code, rowID uint32, err error) {
	ctx := context.Background()
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("DoSaveLogPass", zap.Error(err))
		return gprcStatus, rowID, err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	var headers metadata.MD
	reqData := pb.PostLogPassRequest{Login: login, Password: password, Description: description}
	response, err := client.PostLogPass(ctx, &reqData, grpc.Header(&headers))
	grpcStatus := codes.OK
	if status.Code(err) != codes.OK {
		logger.Log.Error(err.Error())
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, 0, err
	}
	return grpcStatus, response.RowId, nil
}
