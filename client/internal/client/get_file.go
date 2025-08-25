package client

import (
	"context"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func DoGetFile(tokenStr string, rowID uint32) (grpcStatus codes.Code, fileName string, err error) {
	var file *os.File
	var fileSize int
	ctx := context.Background()
	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Log.Error("DoGetFile", zap.Error(err))
		return grpcStatus, fileName, err
	}
	defer conn.Close()

	// create client & stream
	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	inData := pb.GetDataRequest{RowId: rowID}
	stream, err := client.GetFile(ctx, &inData)
	if err != nil {
		logger.Log.Error("DoGetFile", zap.Error(err))
		return status.Code(err), fileName, err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Error("DoGetFile", zap.Error(err))
			return grpcStatus, fileName, err
		}
		if fileName == "" {
			fileName = resp.FileName
			names := strings.Split(fileName, ".")
			fileName = names[0] + strconv.FormatInt(time.Now().Unix(), 10) + "." + names[1]
			file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0662)
			if err != nil {
				logger.Log.Error("DoGetFile", zap.Error(err))
				return grpcStatus, fileName, err
			}
			defer file.Close()

		}
		n, err := file.Write(resp.GetChunk())
		fileSize += n
		if err != nil {
			logger.Log.Error("doGetFile", zap.Error(err))
			return grpcStatus, fileName, err
		}

	}
	if err != nil {
		logger.Log.Error("doGetFile", zap.Error(err))
	}
	logger.Log.Info("Finish doGetFile with results", zap.String("fileName: ", fileName), zap.Int("filesize: ", fileSize))

	return grpcStatus, fileName, nil
}
