package repository

import (
	"context"
	"time"

	"github.com/VicShved/pass-manager/server/pkg/config"

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
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
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
	err := r.DB.WithContext(ctx).AutoMigrate()
	return err
}

// Empty
func (r GormRepository) Empty() error {
	return nil
}

// CloseConn Close connection
func (r GormRepository) CloseConn() {
	sqlDB, _ := r.DB.DB()
	sqlDB.Close()
}
