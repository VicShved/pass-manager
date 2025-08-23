package server

import (
	"os"
	"testing"
	"time"

	"github.com/VicShved/pass-manager/server/internal/service"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
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
	}
	for _, tst := range tests {
		statusCode, _, _ := doRegister(tst.login, tst.password)
		assert.Equal(t, tst.want.status, statusCode)
	}
	logger.Log.Info("Exit from TestDoRegister")
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
				status: codes.OK,
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
	logger.Log.Info("Exit from TestDoLogin")
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
		login       string
		password    string
		want        want
	}{
		{
			name:        "goodpst card",
			cardNumber:  "1001-2002-3003-4004",
			cardValid:   "01/10",
			cardCode:    "111",
			description: "test card #1",
			login:       "1",
			password:    "password1",
			want: want{
				status: codes.OK,
				err:    nil,
			},
		},
		// {
		// 	name:       "bad token",
		// 	cardNumber: "2002-3003-4004-5005",
		// 	cardValid:  "02/20",
		// 	cardCode:   "222",
		// 	login:      "bad",
		// 	password:   "password1",
		// 	want: want{
		// 		status: codes.PermissionDenied,
		// 		err:    status.Error(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь"),
		// 	},
		// },
	}
	for _, tst := range tests {
		statusCode, tokenStr, err := doLogin(tst.login, tst.password)
		assert.Equal(t, tst.want.status, statusCode)
		statusCode, rowID, err := doPostCard(tokenStr, tst.cardNumber, tst.cardValid, tst.cardCode, tst.description)
		assert.Equal(t, tst.want.err, err)
		assert.Equal(t, tst.want.status, statusCode)
		if statusCode == codes.OK {
			assert.Greaterf(t, rowID, uint32(0), "")
			logger.Log.Debug("TestDoPostCard", zap.Int32("rowID", int32(rowID)))
		}
	}
	logger.Log.Info("Exit from TestDoPostCard")
}

func TestDoPostLogPass(t *testing.T) {

	type want struct {
		status codes.Code
		err    error
	}
	var tests = []struct {
		name        string
		extLogin    string
		extPassword string
		description string
		login       string
		password    string
		want        want
	}{
		{
			name:        "goodpst log|pass",
			extLogin:    "login1",
			extPassword: "password1",
			description: "test logpass #1",
			login:       "1",
			password:    "password1",
			want: want{
				status: codes.OK,
				err:    nil,
			},
		},
		// {
		// 	name:        "goodpst log|pass",
		// 	extLogin:    "login1",
		// 	extPassword: "password2",
		// 	description: "test logpass #1",
		// 	login:       "2",
		// 	password:    "password2",
		// 	want: want{
		// 		status: codes.OK,
		// 		err:    nil,
		// 	},
		// },
		// {
		// 	name:        "bad token",
		// 	extLogin:    "login1",
		// 	extPassword: "password2",
		// 	description: "test logpass #2",
		// 	login:       "bad",
		// 	password:    "password1",
		// 	want: want{
		// 		status: codes.PermissionDenied,
		// 		err:    status.Error(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь"),
		// 	},
		// },
	}
	for _, tst := range tests {
		statusCode, tokenStr, err := doLogin(tst.login, tst.password)
		assert.Equal(t, tst.want.status, statusCode)
		statusCode, rowID, err := doPostLogPass(tokenStr, tst.extLogin, tst.extPassword, tst.description)
		assert.Equal(t, tst.want.err, err)
		assert.Equal(t, tst.want.status, statusCode)
		if statusCode == codes.OK {
			assert.Greaterf(t, rowID, uint32(0), "")
			logger.Log.Debug("TestDoPostLogPass", zap.Int32("rowID", int32(rowID)))
		}
	}
	logger.Log.Info("Exit from TestDoPostLogPass")
}

func TestDoGetCard(t *testing.T) {

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
		login       string
		password    string
		want        want
	}{
		{
			name:        "goodpst card",
			cardNumber:  "1001-2002-3003-4004",
			cardValid:   "01/10",
			cardCode:    "111",
			description: "test card #1",
			login:       "1",
			password:    "password1",
			want: want{
				status: codes.OK,
				err:    nil,
			},
		},
		// {
		// 	name:       "bad token",
		// 	cardNumber: "2002-3003-4004-5005",
		// 	cardValid:  "02/20",
		// 	cardCode:   "222",
		// 	login:      "2",
		// 	password:   "password2",
		// 	want: want{
		// 		status: codes.OK,
		// 		err:    nil,
		// 	},
		// },
	}
	for _, tst := range tests {
		statusCode, tokenStr, err := doLogin(tst.login, tst.password)
		assert.Equal(t, tst.want.status, statusCode)
		statusCode, rowID, err := doPostCard(tokenStr, tst.cardNumber, tst.cardValid, tst.cardCode, tst.description)
		assert.Equal(t, tst.want.err, err)
		assert.Equal(t, tst.want.status, statusCode)
		if statusCode == codes.OK {
			statusCode, card, err := doGetCard(tokenStr, rowID)
			assert.Nil(t, err)
			if statusCode == codes.OK {
				assert.Equal(t, tst.cardNumber, card.CardNumber)
				assert.Equal(t, tst.cardValid, card.CardValid)
				assert.Equal(t, tst.cardCode, card.CardCode)
			}
		}
	}
	logger.Log.Info("Exit from TestDoGetCard")
}

func TestDoGetLogPass(t *testing.T) {

	type want struct {
		status codes.Code
		err    error
	}
	var tests = []struct {
		name        string
		extLogin    string
		extPassword string
		description string
		login       string
		password    string
		want        want
	}{
		{
			name:        "goodpst log|pass",
			extLogin:    "login1",
			extPassword: "password1",
			description: "test logpass #1",
			login:       "1",
			password:    "password1",
			want: want{
				status: codes.OK,
				err:    nil,
			},
		},
		// {
		// 	name:        "goodpst log|pass",
		// 	extLogin:    "login1",
		// 	extPassword: "password2",
		// 	description: "test logpass #1",
		// 	login:       "2",
		// 	password:    "password2",
		// 	want: want{
		// 		status: codes.OK,
		// 		err:    nil,
		// 	},
		// },
		// {
		// 	name:        "bad token",
		// 	extLogin:    "login1",
		// 	extPassword: "password2",
		// 	description: "test logpass #2",
		// 	login:       "bad",
		// 	password:    "password1",
		// 	want: want{
		// 		status: codes.PermissionDenied,
		// 		err:    status.Error(codes.PermissionDenied, "Отсутствует токен/Не определен пользователь"),
		// 	},
		// },
	}
	for _, tst := range tests {
		statusCode, tokenStr, err := doLogin(tst.login, tst.password)
		assert.Equal(t, tst.want.status, statusCode)
		statusCode, rowID, err := doPostLogPass(tokenStr, tst.extLogin, tst.extPassword, tst.description)
		assert.Equal(t, tst.want.err, err)
		assert.Equal(t, tst.want.status, statusCode)
		if statusCode == codes.OK {
			statusCode, logPass, err := doGetLogPass(tokenStr, rowID)
			assert.Nil(t, err)
			if statusCode == codes.OK {
				assert.Equal(t, tst.extLogin, logPass.Login)
				assert.Equal(t, tst.extPassword, logPass.Password)
			}
		}
	}
	logger.Log.Info("Exit from TestDoGetLogPass")
}

func TestDoPostFile(t *testing.T) {

	type want struct {
		status codes.Code
	}
	var tests = []struct {
		name     string
		login    string
		password string
		fileName string
		want     want
	}{
		{
			name:     "goodpst card",
			login:    "1",
			password: "password1",
			fileName: "server_test.test",
			want: want{
				status: codes.OK,
			},
		},
		// {
		// 	name:     "goodpst card",
		// 	login:    "1",
		// 	password: "password1",
		// 	fileName: "staticcheck_linux_386.tar.gz",
		// 	want: want{
		// 		status: codes.OK,
		// 	},
		// },
	}
	for _, tst := range tests {
		statusCode, tokenStr, err := doLogin(tst.login, tst.password)
		assert.Equal(t, tst.want.status, statusCode)
		statusCode, _, err = doPostFile(tokenStr, tst.fileName)
		assert.Nil(t, err)
		assert.Equal(t, tst.want.status, statusCode)
	}
}

func TestDoGetFile(t *testing.T) {

	type want struct {
		status codes.Code
	}
	var tests = []struct {
		name     string
		login    string
		password string
		fileName string
		want     want
	}{
		{
			name:     "goodpst card",
			login:    "1",
			password: "password1",
			fileName: "server_test.test",
			want: want{
				status: codes.OK,
			},
		},
		{
			name:     "goodpst card",
			login:    "1",
			password: "password1",
			fileName: "staticcheck_linux_386.tar.gz",
			want: want{
				status: codes.OK,
			},
		},
	}
	for _, tst := range tests {
		statusCode, tokenStr, err := doLogin(tst.login, tst.password)
		assert.Equal(t, tst.want.status, statusCode)
		statusCode, rowID, err := doPostFile(tokenStr, tst.fileName)
		logger.Log.Debug("TestDoGetFile", zap.Uint32("rowID = ", rowID))
		assert.Nil(t, err)
		assert.Equal(t, tst.want.status, statusCode)
		statusCode, fileName, err := doGetFile(tokenStr, rowID)
		logger.Log.Debug("TestDoGetFile", zap.String("Save  to file ", fileName))
		assert.Nil(t, err)
	}
}
