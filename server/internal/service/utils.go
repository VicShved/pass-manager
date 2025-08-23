package service

import (
	// "encoding/base64"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-uuid"
)

// SigningMethod
var SigningMethod = jwt.SigningMethodHS512

// func GetNewUUID
func GetNewUUID() (string, error) {
	userID, err := uuid.GenerateUUID()
	return string(userID), err
}

// custClaims struct
type custClaims struct {
	jwt.RegisteredClaims
	UserID string
}

// ParseTokenUserID парсит jwt из строки
func ParseTokenUserID(tokenStr string, secretKey string) (*jwt.Token, string, error) {
	claims := &custClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	return token, (*claims).UserID, err
}

// GetJWTTokenString(userID *string)
func GetJWTTokenString(userID *string, secretKey string) (string, error) {
	claim := custClaims{
		UserID: *userID,
	}
	token := jwt.NewWithClaims(SigningMethod, claim)
	tokenStr, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", nil
	}
	return tokenStr, err
}

func GetNewUserIDToken(secretKey string) (string, string) {
	userID, _ := GetNewUUID()
	token, _ := GetJWTTokenString(&userID, secretKey)
	return userID, token
}

// func generateSecretKey() string {
// 	key := make([]byte, lengthSecretKey)
// 	return base64.URLEncoding.EncodeToString(key)
// }

func getNewFileName(fileName string) string {
	uuid, _ := GetNewUUID()
	splits := strings.Split(fileName, ".")
	ext := "unk"
	if len(splits) > 1 {
		ext = splits[len(splits)-1]
	}
	return uuid + "." + ext
}
