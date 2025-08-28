package server

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/VicShved/pass-manager/server/internal/repository"
	"github.com/VicShved/pass-manager/server/internal/service"
	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/gorm"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

var serverAddress = ""
var serverPort = "8080"
var secretKey = "VerySecret"

func setup() (*GServer, *bufconn.Listener, *gorm.DB) {
	err := logger.InitLogger("DEBUG")
	if err != nil {
		log.Fatal(err.Error())
	}
	conf := config.GetServerConfig()
	conf.ServerAddress = serverAddress
	conf.ServerPort = serverPort
	conf.EnableTLS = false
	conf.SecretKey = secretKey
	conf.SchemaName = "test_passmanager"
	conf.DBDSN = "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	fileRepo, err := repository.GetFileStorageRepo("")
	if err != nil {
		log.Fatal(err.Error())
	}
	repo, err := repository.GetGormRepo(conf, fileRepo)
	if err != nil {
		log.Fatal(err.Error())
	}
	serv := service.GetService(repo, conf)
	gserver, _ := GetServer(serv, conf)
	lis = bufconn.Listen(bufSize)
	go func() {
		if err := gserver.server.Serve(lis); err != nil {
			log.Fatal("Start serve error")
		}
	}()
	repo.DB = repo.DB.Begin(&sql.TxOptions{Isolation: sql.LevelReadUncommitted})
	return gserver, lis, repo.DB
}

func close(gserver *GServer, lis *bufconn.Listener) {
	lis.Close()
	gserver.server.GracefulStop()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func doRegister(login string, password string) (codes.Code, string, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		logger.Log.Error("doRegister", zap.Error(err))
		return 999, "", err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	var header metadata.MD
	reqData := pb.LoginRequest{Login: login, Password: password}
	_, err = client.Register(ctx, &reqData, grpc.Header(&header))
	grpcStatus := codes.OK
	if status.Code(err) != codes.OK {
		logger.Log.Warn("doRegister", zap.String("Error", err.Error()))
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, "", err
	}
	authToken := header.Get(config.AuthorizationTokenName)[0]
	if len(authToken) == 0 {
		return 999, "", errors.New("Сервер не возвратил auth token")
	}
	return 0, authToken, nil
}

func doLogin(login string, password string) (codes.Code, string, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	// grpc.NewClient("bufnet", grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		logger.Log.Error("doLogin", zap.Error(err))
		return 999, "", err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	var header metadata.MD
	reqData := pb.LoginRequest{Login: login, Password: password}
	_, err = client.Login(ctx, &reqData, grpc.Header(&header))
	grpcStatus := codes.OK
	if status.Code(err) != codes.OK {
		logger.Log.Warn("doLogin", zap.String("Error", err.Error()))
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, "", err
	}
	authToken := header.Get(config.AuthorizationTokenName)[0]
	logger.Log.Info("doLogin", zap.String("tokenStr", authToken))
	if len(authToken) == 0 {
		return 999, "", errors.New("Сервер не возвратил auth token")
	}
	testAuthToken := authToken
	return 0, testAuthToken, nil
}

func doPostCard(tokenStr string, cardNumber string, cardValid string, cardCode string, description string) (codes.Code, uint32, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
		return 0, 0, err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	var headers metadata.MD
	reqData := pb.PostCardRequest{CardNumber: cardNumber, Valid: cardValid, Code: cardCode, Description: description}
	response, err := client.PostCard(ctx, &reqData, grpc.Header(&headers))
	grpcStatus := codes.OK
	if status.Code(err) != codes.OK {
		log.Printf(err.Error())
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, 0, err
	}
	return grpcStatus, response.RowId, nil
}

func doPostLogPass(tokenStr string, login string, password string, description string) (codes.Code, uint32, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
		return 0, 0, err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	var headers metadata.MD
	reqData := pb.PostLogPassRequest{Login: login, Password: password, Description: description}
	response, err := client.PostLogPass(ctx, &reqData, grpc.Header(&headers))
	grpcStatus := codes.OK
	if status.Code(err) != codes.OK {
		log.Printf(err.Error())
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, 0, err
	}
	return grpcStatus, response.RowId, nil
}

func doGetCard(tokenStr string, rowID uint32) (grpcCode codes.Code, card service.CardStruct, err error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
		return grpcCode, card, err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	var headers metadata.MD
	reqData := pb.GetDataRequest{RowId: rowID}
	response, err := client.GetCard(ctx, &reqData, grpc.Header(&headers))
	grpcStatus := codes.OK
	if status.Code(err) != codes.OK {
		log.Printf(err.Error())
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, card, err
	}
	card.CardNumber = response.CardNumber
	card.CardValid = response.Valid
	card.CardCode = response.Code
	return grpcStatus, card, nil
}

func doGetLogPass(tokenStr string, rowID uint32) (grpcCode codes.Code, logPass service.LogPassStruct, err error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
		return grpcCode, logPass, err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	var headers metadata.MD
	reqData := pb.GetDataRequest{RowId: rowID}
	response, err := client.GetLogPass(ctx, &reqData, grpc.Header(&headers))
	grpcStatus := codes.OK
	if status.Code(err) != codes.OK {
		log.Printf(err.Error())
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, logPass, err
	}
	logPass.Login = response.Login
	logPass.Password = response.Password

	return grpcStatus, logPass, nil
}

func doPostFile(tokenStr string, fileName string) (grpcCode codes.Code, rowID uint32, err error) {
	// create connection
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		logger.Log.Error("doPostFile", zap.Error(err))
		return grpcCode, rowID, err
	}
	defer conn.Close()
	// open file
	file, err := os.Open(fileName)
	if err != nil {
		logger.Log.Error("doPostFile", zap.Error(err))
		return grpcCode, rowID, err
	}
	defer file.Close()
	// create client & stream
	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
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
			logger.Log.Error("doPostFile", zap.Error(err))
			return grpcCode, rowID, err
		}
		err = stream.Send(
			&pb.PostFileRequest{
				FileName:    fileName,
				Description: "Description " + fileName,
				Chunk:       buffer[:n],
			},
		)
		if err != nil {
			logger.Log.Error("doPostFile", zap.Error(err))
			return grpcCode, rowID, err
		}
	}
	result, err := stream.CloseAndRecv()
	if err != nil {
		logger.Log.Error("doPostFile", zap.Error(err))
	}
	logger.Log.Info("Finish doPostFile with results", zap.Uint32("rowID: ", result.RowId), zap.Int("filesize: ", int(result.Length)))
	return grpcCode, result.RowId, err
}

func doGetFile(tokenStr string, rowID uint32) (grpcCode codes.Code, fileName string, err error) {
	var file *os.File
	var fileSize int
	// create connection
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		logger.Log.Error("doGetFile", zap.Error(err))
		return grpcCode, fileName, err
	}
	defer conn.Close()
	// create client & stream
	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(config.AuthorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	inData := pb.GetDataRequest{RowId: rowID}
	stream, err := client.GetFile(ctx, &inData)
	if err != nil {
		logger.Log.Error("doGetFile", zap.Error(err))
		return status.Code(err), fileName, err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Error("doGetFile", zap.Error(err))
			return grpcCode, fileName, err
		}
		if fileName == "" {
			fileName = resp.FileName
			names := strings.Split(fileName, ".")
			fileName = names[0] + strconv.FormatInt(time.Now().Unix(), 10) + "." + names[1]
			file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0662)
			if err != nil {
				logger.Log.Error("doGetFile", zap.Error(err))
				return grpcCode, fileName, err
			}
			defer file.Close()

		}
		n, err := file.Write(resp.GetChunk())
		fileSize += n
		if err != nil {
			logger.Log.Error("doGetFile", zap.Error(err))
			return grpcCode, fileName, err
		}

	}
	if err != nil {
		logger.Log.Error("doGetFile", zap.Error(err))
	}
	logger.Log.Info("Finish doGetFile with results", zap.String("fileName: ", fileName), zap.Int("filesize: ", fileSize))
	return grpcCode, fileName, err
}
