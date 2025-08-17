package repository

import (
	"time"
)

type LifeTimeModel struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type User struct {
	LifeTimeModel
	UserID       string `gorm:"size:36,unique"`
	Login        string `gorm:"size:256,unique"`
	HashPassword string `gorm:"type:bytes"`
	UserDatas    []UserData
}

type UserData struct {
	LifeTimeModel
	UserID      uint
	DataType    string `gorm:"size:16"`
	Description string `gorm:"type:text"`
	FileName    string `gorm:"size:512"`
	SecretKey   string `gorm:"size:1024"`
}
