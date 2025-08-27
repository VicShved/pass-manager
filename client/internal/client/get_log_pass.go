package client

import (
	"context"
	"fmt"

	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// DoGetLogPass   get login/password pair from server
func (c GClient) DoGetLogPass(tokenStr string, rowID uint32) (grpcStatus codes.Code, logPassStr string, err error) {
	ctx := context.Background()
	conn, err := c.getConnection()
	if err != nil {
		logger.Log.Error("DoGetLogPass", zap.Error(err))
		return grpcStatus, logPassStr, err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	ctx = c.addToken2Context(ctx, tokenStr)
	var headers metadata.MD
	reqData := pb.GetDataRequest{RowId: rowID}
	response, err := client.GetLogPass(ctx, &reqData, grpc.Header(&headers))
	grpcStatus = codes.OK
	if status.Code(err) != codes.OK {
		logger.Log.Error("DoGetLogPass", zap.Error(err))
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, logPassStr, err
	}
	logPassStr = fmt.Sprintf("Login = %s\nPassword = %s", response.Login, response.Password)

	return grpcStatus, logPassStr, nil
}
