package server

import (
	"os"
	"testing"
	"time"

	"github.com/VicShved/pass-manager/server/internal/service"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMain(m *testing.M) {
	gserver, lis, txPoint := setup()
	defer txPoint.Rollback()
	time.Sleep(time.Millisecond * 100)
	code := m.Run()
	close(gserver, lis)
	os.Exit(code)
}

func TestDoRegister(t *testing.T) {

	type want struct {
		status codes.Code
	}
	var tests = []struct {
		name     string
		login    string
		password string
		want     want
	}{
		{
			name:     "good register",
			login:    "1",
			password: "password1",
			want: want{
				status: codes.OK,
			},
		},
		{
			name:     "bad register",
			login:    "",
			password: "password1",
			want: want{
				status: codes.InvalidArgument,
			},
		},
		{
			name:     "bad register",
			login:    "777",
			password: "",
			want: want{
				status: codes.InvalidArgument,
			},
		},
		{
			name:     "good register",
			login:    "2",
			password: "password2",
			want: want{
				status: codes.OK,
			},
		},
		{
			name:     "bad register",
			login:    "2",
			password: "password2",
			want: want{
				status: codes.InvalidArgument,
			},
		},
	}
	for _, tst := range tests {
		statusCode, _, _ := doRegister(tst.login, tst.password)
		assert.Equal(t, tst.want.status, statusCode)
		// if statusCode == codes.OK {
		// 	token, _, err := service.ParseTokenUserID(tokenStr, secretKey)
		// 	assert.Nil(t, err, "Err")
		// 	assert.True(t, token.Valid)
		// test new login exists
		// statusCode, tokenStr, err = doLogin(tst.login, tst.password)
		// assert.Nil(t, err)
		// assert.Equal(t, statusCode, codes.OK)
		// if (err == nil) && (statusCode == codes.OK) {
		// 	statusCode, _, err = doPostCard(tokenStr)
		// 	assert.Nil(t, err)
		// 	assert.Equal(t, statusCode, codes.OK)
		// }
		// }
	}
}

func TestDoLogin(t *testing.T) {

	type want struct {
		status codes.Code
	}
	var tests = []struct {
		name     string
		login    string
		password string
		want     want
	}{
		{
			name:     "good request",
			login:    "1",
			password: "password1",
			want: want{
				status: codes.PermissionDenied,
			},
		},
		{
			name:     "bad password",
			login:    "1",
			password: "passwordBad",
			want: want{
				status: codes.PermissionDenied,
			},
		},
		{
			name:     "bad login",
			login:    "bad",
			password: "password1",
			want: want{
				status: codes.PermissionDenied,
			},
		},
	}
	for _, tst := range tests {
		statusCode, tokenStr, _ := doLogin(tst.login, tst.password)
		assert.Equal(t, tst.want.status, statusCode)
		if statusCode == codes.OK {
			token, userID, err := service.ParseTokenUserID(tokenStr, secretKey)
			assert.Nil(t, err, "Err")
			assert.True(t, token.Valid)
			assert.NotEmpty(t, userID)
		}
	}
}

func TestDoPostCard(t *testing.T) {

	type want struct {
		status codes.Code
		err    error
	}
	var tests = []struct {
		name        string
		cardNumber  string
		cardValid   string
		cardCode    string
		description string
		tokenStr    string
		want        want
	}{
		{
			name:       "goodpst card",
			cardNumber: "1001-2002-3003-4004",
			cardValid:  "01/10",
			cardCode:   "111",
			tokenStr:   "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJ1c2VySUQxIn0.cTLx3NPyazsSt1Ny09-ZVJHT-ham_r1zpgGCGU4PHO9yHlQplk7YdcdDGU5rkTXp1NCdlHxhYMuDnfhigD45uw",
			want: want{
				status: codes.OK,
				err:    nil,
			},
		},
		{
			name:       "bad token",
			cardNumber: "2002-3003-4004-5005",
			cardValid:  "02/20",
			cardCode:   "222",
			tokenStr:   "eyJhbGciO",
			want: want{
				status: codes.PermissionDenied,
				err:    status.Error(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь"),
			},
		},
	}
	for _, tst := range tests {
		statusCode, _, err := doPostCard(tst.tokenStr, tst.cardNumber, tst.cardValid, tst.cardCode)
		assert.Equal(t, tst.want.err, err)
		assert.Equal(t, tst.want.status, statusCode)
	}
}

// func TestDoPostFile(t *testing.T) {

// 	type want struct {
// 		status   codes.Code
// 		tokenStr string
// 	}
// 	var tests = []struct {
// 		name     string
// 		login    string
// 		password string
// 		tokenStr string
// 		fileName string
// 		want     want
// 	}{
// 		{
// 			name:     "goodpst card",
// 			login:    "1",
// 			password: "password1",
// 			tokenStr: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOiJ1c2VySUQxIn0.cTLx3NPyazsSt1Ny09-ZVJHT-ham_r1zpgGCGU4PHO9yHlQplk7YdcdDGU5rkTXp1NCdlHxhYMuDnfhigD45uw",
// 			fileName: "server_test.tst",
// 			want: want{
// 				status:   codes.OK,
// 				tokenStr: "",
// 			},
// 		},
// 		// {
// 		// 	name:     "good register",
// 		// 	login:    "2",
// 		// 	password: "password2",
// 		// 	want: want{
// 		// 		status:   codes.OK,
// 		// 		tokenStr: "",
// 		// 	},
// 		// },
// 	}
// 	for _, tst := range tests {
// 		statusCode, _, err := doPostFile(tst.tokenStr, tst.fileName)
// 		assert.Nil(t, err)
// 		assert.Equal(t, tst.want.status, statusCode)
// 	}
// }

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
