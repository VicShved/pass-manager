package client

import (
	"context"

	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type userData struct {
	rowID    int32
	desc     string
	dataType string
}

// DoGetCard Получение данных карты
func (c GClient) DoGetUserData(tokenStr string, dataType int32) (grpcStatus codes.Code, results []userData, err error) {
	ctx := context.Background()
	conn, err := c.getConnection()
	if err != nil {
		logger.Log.Error("DoGetUserData", zap.Error(err))
		return grpcStatus, results, err
	}
	defer conn.Close()

	ctx = c.addToken2Context(ctx, tokenStr)
	var headers metadata.MD
	reqData := pb.GetDataInfoRequest{DataType: pb.DataType(dataType)}
	client := pb.NewPassManagerServiceClient(conn)
	response, err := client.GetDataInfo(ctx, &reqData, grpc.Header(&headers))
	grpcStatus = codes.OK
	if status.Code(err) != codes.OK {
		logger.Log.Error("DoGetCard", zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, cardStr, err
	}

	return grpcStatus, cardStr, nil
}
