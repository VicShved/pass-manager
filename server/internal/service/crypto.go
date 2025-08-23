package service

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"

	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"
)

const lecretKeyLength int = 32

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	return b, err
}

func generateHexKeyString(size int) (string, error) {
	key, err := generateRandom(size)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func encryptData(source []byte, keyString string) (target []byte, err error) {
	key := make([]byte, hex.DecodedLen(len([]byte(keyString))))
	hex.Decode(key, []byte(keyString))
	aesblock, err := aes.NewCipher(key)
	if err != nil {
		logger.Log.Error("encryptData", zap.Error(err))
		return target, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		logger.Log.Error("encryptData", zap.Error(err))
		return target, err
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	t := aesgcm.Seal(nil, nonce, source, nil)

	return t, nil
}

func encryptData2Hex(source []byte, keyString string) (target string, err error) {
	// key := make([]byte, hex.DecodedLen(len([]byte(keyString))))
	// hex.Decode(key, []byte(keyString))
	// aesblock, err := aes.NewCipher(key)
	// if err != nil {
	// 	logger.Log.Error("encryptData", zap.Error(err))
	// 	return target, err
	// }

	// aesgcm, err := cipher.NewGCM(aesblock)
	// if err != nil {
	// 	logger.Log.Error("encryptData", zap.Error(err))
	// 	return target, err
	// }

	// nonce := key[len(key)-aesgcm.NonceSize():]

	// t := aesgcm.Seal(nil, nonce, source, nil)

	t, err := encryptData(source, keyString)
	if err != nil {
		return target, err
	}
	target = hex.EncodeToString(t)
	return target, nil
}

func decryptHexData(sourceString string, keyString string) (target []byte, err error) {
	key, err := hex.DecodeString(keyString)
	if err != nil {
		logger.Log.Error("decryptData", zap.Error(err))
		return target, err
	}

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		logger.Log.Error("decryptData", zap.Error(err))
		return target, err
	}
	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		logger.Log.Error("decryptData", zap.Error(err))
		return target, err
	}

	// создаём вектор инициализации
	nonce := key[len(key)-aesgcm.NonceSize():]

	source, err := hex.DecodeString(sourceString)
	if err != nil {
		logger.Log.Error("decryptData", zap.Error(err))
		return target, err
	}
	target, err = aesgcm.Open(nil, nonce, source, nil)
	if err != nil {
		logger.Log.Error("decryptData", zap.Error(err))
		return target, err
	}
	return target, nil
}
