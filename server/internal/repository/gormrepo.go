package repository

import (
	"context"
	"errors"
	"time"

	"github.com/VicShved/pass-manager/server/pkg/config"
	"github.com/VicShved/pass-manager/server/pkg/logger"
	"go.uber.org/zap"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GormRepository struct
type GormRepository struct {
	DB       *gorm.DB
	conf     *config.ServerConfigStruct
	fileRepo FileStoragerRepoInterface
}

// GetGormDB(dns string)
func GetGormDB(dns string, schemaName string) (*gorm.DB, error) {
	config := gorm.Config{TranslateError: true}
	// add schema name
	if schemaName != "" {
		config.NamingStrategy = schema.NamingStrategy{
			TablePrefix: schemaName + ".",
		}
	}
	db, err := gorm.Open(
		postgres.Open(dns),
		&config,
	)
	return db, err
}

// GetGormRepo(dns string)
func GetGormRepo(conf *config.ServerConfigStruct, fileStorageRepo FileStoragerRepoInterface) (*GormRepository, error) {
	db, err := GetGormDB(conf.DBDSN, conf.SchemaName)
	if err != nil {
		return nil, err
	}
	repo := &GormRepository{
		DB:       db,
		conf:     conf,
		fileRepo: fileStorageRepo,
	}
	err = repo.Migrate()
	if err != nil {
		return nil, err
	}
	return repo, err
}

// Migrate data struct
func (r *GormRepository) Migrate() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err := r.DB.WithContext(ctx).AutoMigrate(&User{}, &UserData{})
	return err
}

// CloseConn Close connection
func (r GormRepository) CloseConn() error {
	sqlDB, _ := r.DB.DB()
	return sqlDB.Close()
}

// GetUserByUserID return User by userID
func (r GormRepository) GetUserByUserID(ctx context.Context, userID string) (User, error) {
	user := User{}
	result := r.DB.WithContext(ctx).Where(&User{UserID: userID}).First(&user)
	if result.Error != nil {
		logger.Log.Error("SaveData", zap.Error(result.Error))
		return user, ErrUserIdFromToken
	}
	return user, nil
}

// Register new user login/password
func (r GormRepository) Register(ctx context.Context, userID string, login string, hashPassword string) error {
	logger.Log.Debug("", zap.String("login", login), zap.String("hashPassword", hashPassword))
	user := User{UserID: userID, Login: login, HashPassword: hashPassword}
	result := r.DB.WithContext(ctx).Create(&user)
	if result.Error != nil {
		// проверяем на ошибка дублирования логина
		if errors.Is(result.Error, gorm.ErrCheckConstraintViolated) {
			logger.Log.Debug("login exists", zap.String("login", login))
			return ErrLoginConflict
		}
	}
	return result.Error
}

// Login user by login/password
func (r GormRepository) Login(ctx context.Context, login string, hashPassword string) (string, error) {
	user := User{}
	result := r.DB.WithContext(ctx).Where(&User{Login: login, HashPassword: hashPassword}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Log.Debug("login|Password not found", zap.String("login", login), zap.String("hashPassword", hashPassword))
			return "", ErrLoginPassword
		}
		return "", result.Error
	}
	logger.Log.Debug("", zap.Any("User", user))
	return user.UserID, result.Error
}

// GetFileStorage return filestorage by filename
func (r GormRepository) GetFileStorage(fileName string) (FileStoragerInterface, error) {
	return r.fileRepo.GetFileStorage(fileName)
}

// SaveData save user data
func (r GormRepository) SaveData(ctx context.Context, userID string, desc string, dataType string, fileName string, secretKey string) (rowID uint32, err error) {
	user, err := r.GetUserByUserID(ctx, userID)
	if err != nil {
		return rowID, err
	}
	userData := UserData{
		UserID:      user.ID,
		Description: desc,
		DataType:    dataType,
		FileName:    fileName,
		SecretKey:   secretKey,
	}
	result := r.DB.WithContext(ctx).Create(&userData)
	if result.Error != nil {
		logger.Log.Error("SaveData", zap.Error(result.Error))
		return rowID, result.Error
	}
	logger.Log.Debug("SaveData", zap.Any("Create UserData", userData))
	return uint32(userData.ID), err
}

// GetUserData return user data
func (r GormRepository) GetUserData(ctx context.Context, userID string, rowID uint32) (UserData, error) {
	userData := UserData{}
	user, err := r.GetUserByUserID(ctx, userID)
	if err != nil {
		return userData, err
	}
	logger.Log.Debug("GetUserData", zap.String("userID", userID), zap.Uint32("rowID", rowID))
	result := r.DB.WithContext(ctx).Where(&UserData{ID: uint(rowID), UserID: user.ID}).First(&userData)
	return userData, result.Error
}

// GetUserDatas return array user data
func (r GormRepository) GetUserDatas(ctx context.Context, userID string, dataType string) ([]UserData, error) {
	user, err := r.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	var userDatas []UserData
	filter := UserData{UserID: user.ID}
	if dataType != "" {
		filter.DataType = dataType
	}
	result := r.DB.WithContext(ctx).Where(&filter).Find(&userDatas)
	return userDatas, result.Error
}
