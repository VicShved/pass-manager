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
)

// GormRepository struct
type GormRepository struct {
	DB   *gorm.DB
	conf *config.ServerConfigStruct
}

// GetGormDB(dns string)
func GetGormDB(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{TranslateError: true})
	return db, err
}

// GetGormRepo(dns string)
func GetGormRepo(conf *config.ServerConfigStruct) (*GormRepository, error) {
	db, err := GetGormDB(conf.DBDSN)
	if err != nil {
		return nil, err
	}
	repo := &GormRepository{
		DB:   db,
		conf: conf,
	}
	err = repo.Migrate()
	if err != nil {
		return nil, err
	}
	return repo, err
}

// Migrate()
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

func (r GormRepository) Register(userID string, login string, hashPassword string) error {
	logger.Log.Debug("", zap.String("login", login), zap.String("hashPassword", hashPassword))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
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

func (r GormRepository) Login(login string, hashPassword string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
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
