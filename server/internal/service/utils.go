package service

import (
	// "encoding/base64"
	"errors"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-uuid"
)

// ErrParseToken Error parse token
var ErrParseToken = errors.New("Error parse token")

const tokenPrefix string = "Bearer"

// SigningMethod set as jwt.SigningMethodHS512
var SigningMethod = jwt.SigningMethodHS512

// GetNewUUID return UUID
func GetNewUUID() (string, error) {
	id, err := uuid.GenerateUUID()
	return string(id), err
}

type UserData struct {
	RowID    uint32
	Desc     string
	DataType string
}

// custClaims struct
type custClaims struct {
	jwt.RegisteredClaims
	UserID string
}

// ParseTokenUserID парсит jwt из строки
func ParseTokenUserID(tokenStr string, secretKey string) (*jwt.Token, string, error) {
	splits := strings.Split(tokenStr, " ")
	if len(splits) != 2 {
		return nil, "", ErrParseToken
	}
	if splits[0] != tokenPrefix {
		return nil, "", ErrParseToken
	}
	claims := &custClaims{}
	token, err := jwt.ParseWithClaims(splits[1], claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return token, (*claims).UserID, err
}

// GetJWTTokenString (userID *string)
func GetJWTTokenString(userID *string, secretKey string) (string, error) {
	claim := custClaims{
		UserID: *userID,
	}
	token := jwt.NewWithClaims(SigningMethod, claim)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	tokenStr = tokenPrefix + " " + tokenStr
	return tokenStr, err
}

func getNewFileName(fileName string) string {
	uuid, _ := GetNewUUID()
	splits := strings.Split(fileName, ".")
	ext := "unk"
	if len(splits) > 1 {
		ext = splits[len(splits)-1]
	}
	return uuid + "." + ext
}
