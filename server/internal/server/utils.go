package server

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/hashicorp/go-uuid"
)

// authorizationTokenName
var authorizationTokenName = "Authorization"

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

// parseTokenUserID парсит jwt из строки
func parseTokenUserID(tokenStr string, secretKey string) (*jwt.Token, string, error) {
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

func getNewUserIDToken(secretKey string) (string, string) {
	userID, _ := GetNewUUID()
	token, _ := GetJWTTokenString(&userID, secretKey)
	return userID, token
}
