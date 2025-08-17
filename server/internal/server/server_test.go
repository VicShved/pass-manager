package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"testing"

	"github.com/VicShved/pass-manager/server/internal/repository"
	"github.com/VicShved/pass-manager/server/internal/service"
	pb "github.com/VicShved/pass-manager/server/pkg/api/proto"
	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024 * 10

var lis *bufconn.Listener

var serverAddress = ":8080"
var secretKey = "VerySecret"

func setup() *GServer {
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
			log.Fatalf("Start serve error")
		}
	}()
	return gserver
}

func TestMain(m *testing.M) {
	setup()
	// time.Sleep(3 * time.Second)
	code := m.Run()
	// gserver.GracefulStop()
	os.Exit(code)
}

func doLogin(login string, password string) (uint32, string, error) {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials((insecure.NewCredentials())))
	// grpc.NewClient("bufnet", grpc.WithTransportCredentials((insecure.NewCredentials())))
	if err != nil {
		log.Fatal(err)
		return 999, "", err
	}
	defer conn.Close()

	c := pb.NewPassManagerServiceClient(conn)
	// md := metadata.Pairs(middware.AuthorizationCookName, "")
	var header metadata.MD
	reqData := pb.LoginRequest{Login: login, Password: password}
	resp, err := c.Login(ctx, &reqData)
	grpcStatus := uint32(0)
	if status.Code(err) != codes.OK {
		log.Printf(err.Error())
		st, ok := status.FromError(err)
		if ok {
			grpcStatus = uint32(st.Code())
		}
		return grpcStatus, "", err
	}
	fmt.Println("resp.Token=", resp.Jwt)
	hAuthToken := header.Get(authorizationTokenName)
	print("hAuthToken=", hAuthToken)
	authToken := resp.Jwt
	if len(authToken) == 0 {
		return 999, "", errors.New("Сервер не возвратил auth token")
	}
	fmt.Println("authToken ", authToken)
	testAuthToken := authToken
	return 0, testAuthToken, nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
func TestDoLogin(t *testing.T) {

	type wont struct {
		status   uint32
		tokenStr string
	}
	var tests = []struct {
		name     string
		login    string
		password string
		wont     wont
	}{
		{
			name:     "good request",
			login:    "1",
			password: "password1",
			wont: wont{
				status:   uint32(0),
				tokenStr: "",
			},
		},
		// {
		// 	name:     "bad request",
		// 	login:    "999",
		// 	password: "password999",
		// 	wont: wont{
		// 		status:   3,
		// 		tokenStr: "",
		// 	},
		// },
	}
	for _, tst := range tests {
		statusCode, tokenStr, _ := doLogin(tst.login, tst.password)
		assert.Equal(t, tst.wont.status, statusCode)
		if statusCode == 0 {
			fmt.Println("tokenStr=", tokenStr)
			token, _, err := service.ParseTokenUserID(tokenStr, secretKey)
			assert.Nil(t, err, "Err")
			assert.True(t, token.Valid)
		}
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
