package client

import (
	"context"
	"io"
	"os"

	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c GClient) DoSaveFile(tokenStr string, fileName string, description string) (grpcStatus codes.Code, rowID uint32, err error) {
	// create connection
	ctx := context.Background()
	conn, err := c.getConnection()
	if err != nil {
		logger.Log.Error("DoSaveLogPass", zap.Error(err))
		return grpcStatus, rowID, err
	}
	defer conn.Close()

	// open file
	file, err := os.Open(fileName)
	if err != nil {
		logger.Log.Error("DoSaveFile", zap.Error(err))
		return grpcStatus, rowID, err
	}
	defer file.Close()
	// create client & stream
	client := pb.NewPassManagerServiceClient(conn)
	ctx = c.addToken2Context(ctx, tokenStr)
	stream, err := client.PostFile(ctx)
	if err != nil {
		return status.Code(err), rowID, err
	}

	buffer := make([]byte, 32)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Error("DoSaveFile", zap.Error(err))
			return grpcStatus, rowID, err
		}
		err = stream.Send(
			&pb.PostFileRequest{
				FileName:    fileName,
				Description: description,
				Chunk:       buffer[:n],
			},
		)
		if err != nil {
			logger.Log.Error("DoSaveFile", zap.Error(err))
			return grpcStatus, rowID, err
		}
	}
	result, err := stream.CloseAndRecv()
	if err != nil {
		logger.Log.Error("DoSaveFile", zap.Error(err))
	}
	logger.Log.Info("Finish DoSaveFile with results", zap.Uint32("rowID: ", result.RowId), zap.Int("filesize: ", int(result.Length)))
	return grpcStatus, result.RowId, err
}
