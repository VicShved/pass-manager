package server

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/VicShved/pass-manager/server/internal/repository"
	"github.com/VicShved/pass-manager/server/internal/service"
	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024 * 100

var lis *bufconn.Listener

var serverAddress = ":8080"
var secretKey = "VerySecret"

func setup() (*GServer, *bufconn.Listener) {
	err := logger.InitLogger("DEBUG")
	if err != nil {
		log.Fatal(err.Error())
	}
	conf := config.GetServerConfig()
	conf.ServerAddress = serverAddress
	conf.SecretKey = secretKey
	repo, _ := repository.GetTestRepo(conf)
	serv := service.GetService(repo, conf)
	gserver, _ := GetServer(serv, conf)
	lis = bufconn.Listen(bufSize)
	go func() {
		if err := gserver.server.Serve(lis); err != nil {
			log.Fatal("Start serve error")
		}
	}()
	return gserver, lis
}

func close(gserver *GServer, lis *bufconn.Listener) {
	lis.Close()
	gserver.server.GracefulStop()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestMain(m *testing.M) {
	gserver, lis := setup()
	time.Sleep(time.Millisecond * 100)
	code := m.Run()
	close(gserver, lis)
	os.Exit(code)
}

// func doLogin(login string, password string) (codes.Code, string, error) {
// 	ctx := context.Background()
// 	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
// 	// grpc.NewClient("bufnet", grpc.WithTransportCredentials((insecure.NewCredentials())))
// 	if err != nil {
// 		logger.Log.Error("doLogin", zap.Error(err))
// 		return 999, "", err
// 	}
// 	defer conn.Close()

// 	client := pb.NewPassManagerServiceClient(conn)
// 	var header metadata.MD
// 	reqData := pb.LoginRequest{Login: login, Password: password}
// 	_, err = client.Login(ctx, &reqData, grpc.Header(&header))
// 	grpcStatus := codes.OK
// 	if status.Code(err) != codes.OK {
// 		logger.Log.Warn("doLogin", zap.String("Error", err.Error()))
// 		st, ok := status.FromError(err)
// 		if ok {
// 			grpcStatus = st.Code()
// 		}
// 		return grpcStatus, "", err
// 	}
// 	authToken := header.Get(authorizationTokenName)[0]
// 	logger.Log.Info("doLogin", zap.String("tokenStr", authToken))
// 	if len(authToken) == 0 {
// 		return 999, "", errors.New("Сервер не возвратил auth token")
// 	}
// 	testAuthToken := authToken
// 	return 0, testAuthToken, nil
// }

// func TestDoLogin(t *testing.T) {

// 	type want struct {
// 		status   codes.Code
// 		tokenStr string
// 	}
// 	var tests = []struct {
// 		name     string
// 		login    string
// 		password string
// 		want     want
// 	}{
// 		{
// 			name:     "good request",
// 			login:    "1",
// 			password: "password1",
// 			want: want{
// 				status:   codes.OK,
// 				tokenStr: "",
// 			},
// 		},
// 		{
// 			name:     "bad password",
// 			login:    "1",
// 			password: "passwordBad",
// 			want: want{
// 				status:   codes.PermissionDenied,
// 				tokenStr: "",
// 			},
// 		},
// 		{
// 			name:     "bad login",
// 			login:    "bad",
// 			password: "password1",
// 			want: want{
// 				status:   codes.PermissionDenied,
// 				tokenStr: "",
// 			},
// 		},
// 	}
// 	for _, tst := range tests {
// 		statusCode, tokenStr, _ := doLogin(tst.login, tst.password)
// 		assert.Equal(t, tst.want.status, statusCode)
// 		if statusCode == codes.OK {
// 			token, _, err := service.ParseTokenUserID(tokenStr, secretKey)
// 			assert.Nil(t, err, "Err")
// 			assert.True(t, token.Valid)
// 		}
// 	}
// }

// func doRegister(login string, password string) (codes.Code, string, error) {
// 	ctx := context.Background()
// 	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
// 	if err != nil {
// 		logger.Log.Error("doRegister", zap.Error(err))
// 		return 999, "", err
// 	}
// 	defer conn.Close()

// 	client := pb.NewPassManagerServiceClient(conn)
// 	var header metadata.MD
// 	reqData := pb.LoginRequest{Login: login, Password: password}
// 	_, err = client.Register(ctx, &reqData, grpc.Header(&header))
// 	grpcStatus := codes.OK
// 	if status.Code(err) != codes.OK {
// 		logger.Log.Warn("doRegister", zap.String("Error", err.Error()))
// 		st, ok := status.FromError(err)
// 		if ok {
// 			grpcStatus = st.Code()
// 		}
// 		return grpcStatus, "", err
// 	}
// 	authToken := header.Get(authorizationTokenName)[0]
// 	if len(authToken) == 0 {
// 		return 999, "", errors.New("Сервер не возвратил auth token")
// 	}
// 	return 0, authToken, nil
// }

// func TestDoRegister(t *testing.T) {

// 	type want struct {
// 		status   codes.Code
// 		tokenStr string
// 	}
// 	var tests = []struct {
// 		name     string
// 		login    string
// 		password string
// 		want     want
// 	}{
// 		{
// 			name:     "bad register",
// 			login:    "1",
// 			password: "password1",
// 			want: want{
// 				status:   codes.InvalidArgument,
// 				tokenStr: "",
// 			},
// 		},
// 		{
// 			name:     "good register",
// 			login:    "2",
// 			password: "password2",
// 			want: want{
// 				status:   codes.OK,
// 				tokenStr: "",
// 			},
// 		},
// 	}
// 	for _, tst := range tests {
// 		statusCode, tokenStr, _ := doRegister(tst.login, tst.password)
// 		assert.Equal(t, tst.want.status, statusCode)
// 		if statusCode == codes.OK {
// 			token, _, err := service.ParseTokenUserID(tokenStr, secretKey)
// 			assert.Nil(t, err, "Err")
// 			assert.True(t, token.Valid)
// 			// test new login exists
// 			statusCode, tokenStr, err = doLogin(tst.login, tst.password)
// 			assert.Nil(t, err)
// 			assert.Equal(t, statusCode, codes.OK)
// 			if (err == nil) && (statusCode == codes.OK) {
// 				statusCode, _, err = doPostCard(tokenStr)
// 				assert.Nil(t, err)
// 				assert.Equal(t, statusCode, codes.OK)
// 			}
// 		}
// 	}
// }

func doPostCard(tokenStr string) (codes.Code, string, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
		return 999, "", err
	}
	defer conn.Close()

	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(authorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	var headers metadata.MD
	reqData := pb.PostCardRequest{CardNumber: "0110-0220", Valid: "03/30", Code: "111"}
	_, err = client.PostCard(ctx, &reqData, grpc.Header(&headers))
	grpcStatus := codes.OK
	if status.Code(err) != codes.OK {
		log.Printf(err.Error())
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = st.Code()
		}
		return grpcStatus, "", err
	}
	return 0, "", nil
}

func TestDoPostCard(t *testing.T) {

	type want struct {
		status   codes.Code
		tokenStr string
	}
	var tests = []struct {
		name     string
		login    string
		password string
		tokenStr string
		want     want
	}{
		{
			name:     "goodpst card",
			login:    "1",
			password: "password1",
			tokenStr: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJ1c2VySUQxIn0.cTLx3NPyazsSt1Ny09-ZVJHT-ham_r1zpgGCGU4PHO9yHlQplk7YdcdDGU5rkTXp1NCdlHxhYMuDnfhigD45uw",
			want: want{
				status:   codes.OK,
				tokenStr: "",
			},
		},
		{
			name:     "bad token",
			login:    "2",
			password: "password2",
			tokenStr: "eyJhbGciO",
			want: want{
				status:   codes.PermissionDenied,
				tokenStr: "",
			},
		},
	}
	for _, tst := range tests {
		statusCode, _, err := doPostCard(tst.tokenStr)
		assert.Nil(t, err)
		assert.Equal(t, tst.want.status, statusCode)
	}
}

func doPostFile(tokenStr string) (codes.Code, string, error) {
	// create connection
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
		return 999, "", err
	}
	defer conn.Close()
	// open file
	fileName := "server_test.go"
	file, err := os.Open(fileName)
	if err != nil {
		return codes.OK, "", err
	}
	defer file.Close()
	// create client & stream
	client := pb.NewPassManagerServiceClient(conn)
	md := metadata.Pairs(authorizationTokenName, tokenStr)
	ctx = metadata.NewOutgoingContext(ctx, md)
	stream, err := client.PostFile(ctx)
	if err != nil {
		return status.Code(err), "", err
	}

	buffer := make([]byte, 32)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Log.Error("doPostFile", zap.Error(err))
			return codes.OK, "", err
		}
		err = stream.Send(
			&pb.PostFileRequest{
				FileName: fileName,
				Chunk:    buffer[:n],
			},
		)
		if err != nil {
			logger.Log.Error("doPostFile", zap.Error(err))
			return codes.OK, "", err
		}
	}
	result, err := stream.CloseAndRecv()
	if err != nil {
		logger.Log.Error("doPostFile", zap.Error(err))
	}
	logger.Log.Info("Finish doPostFile with results", zap.String("filename: ", result.FileName), zap.Int("filesize: ", int(result.FileSize)))
	return 0, "", err
}

func TestDoPostFile(t *testing.T) {

	type want struct {
		status   codes.Code
		tokenStr string
	}
	var tests = []struct {
		name     string
		login    string
		password string
		tokenStr string
		fileName string
		want     want
	}{
		{
			name:     "goodpst card",
			login:    "1",
			password: "password1",
			tokenStr: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJ1c2VySUQxIn0.cTLx3NPyazsSt1Ny09-ZVJHT-ham_r1zpgGCGU4PHO9yHlQplk7YdcdDGU5rkTXp1NCdlHxhYMuDnfhigD45uw",
			fileName: "",
			want: want{
				status:   codes.OK,
				tokenStr: "",
			},
		},
		// {
		// 	name:     "good register",
		// 	login:    "2",
		// 	password: "password2",
		// 	want: want{
		// 		status:   codes.OK,
		// 		tokenStr: "",
		// 	},
		// },
	}
	for _, tst := range tests {
		statusCode, _, err := doPostFile(tst.tokenStr, tst.fileName)
		assert.Nil(t, err)
		assert.Equal(t, tst.want.status, statusCode)
	}
}

// func post(tokenStr string, url string) (*pb.PostResponse, error) {
// 	conn, err := grpc.NewClient(baseURL, grpc.WithTransportCredentials((insecure.NewCredentials())))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()

// 	c := pb.NewShortenerServiceClient(conn)
// 	md := metadata.Pairs(middware.AuthorizationCookName, tokenStr)
// 	ctx := metadata.NewOutgoingContext(context.Background(), md)
// 	var header metadata.MD
// 	return c.Post(ctx, &pb.PostRequest{Url: url}, grpc.Header(&header))
// }
// func TestPost(t *testing.T) {
// 	tokenStr, err := getAuthToken()
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	response, err := post(tokenStr, "https://pract.org")
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	fmt.Println("response = ", response.GetResult())
// }

// func get(tokenStr string, shortUrl string) (*pb.GetResponse, error) {
// 	conn, err := grpc.NewClient(baseURL, grpc.WithTransportCredentials((insecure.NewCredentials())))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer conn.Close()
// 	c := pb.NewShortenerServiceClient(conn)
// 	md := metadata.Pairs(middware.AuthorizationCookName, tokenStr) // middware.AuthorizationCookName, authToken
// 	ctx := metadata.NewOutgoingContext(context.Background(), md)
// 	var header metadata.MD
// 	return c.Get(ctx, &pb.GetRequest{Key: shortUrl}, grpc.Header(&header))

// }

// func TestGet(t *testing.T) {
// 	url := "https://pract.org"
// 	tokenStr, _ := getAuthToken()
// 	postResponse, _ := post(tokenStr, url)
// 	respUrl := postResponse.GetResult()
// 	splits := strings.Split(respUrl, "/")
// 	response, err := get(tokenStr, splits[1])
// 	if err != nil {
// 		log.Print(err)
// 	}
// 	fmt.Println("postResponse ", postResponse.GetResult())
// 	assert.Equal(t, url, response.GetUrl())
// 	fmt.Println("response = ", response.String())
// }
