package client

import (
	"context"
	"fmt"

	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DoGetCard(tokenStr string, rowID uint32) (grpcStatus codes.Code, cardStr string, err error) {
	ctx := context.Background()
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("DoGetCard", zap.Error(err))
		return grpcStatus, cardStr, err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	var headers metadata.MD
	reqData := pb.GetDataRequest{RowId: rowID}
	response, err := client.GetCard(ctx, &reqData, grpc.Header(&headers))
	grpcStatus = codes.OK
	if status.Code(err) != codes.OK {
		logger.Log.Error("DoGetCard", zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, cardStr, err
	}
	cardStr = fmt.Sprintf("Card number = %s\nValid = %s\nCode = %s", response.CardNumber, response.Valid, response.Code)

	return grpcStatus, cardStr, nil
}
