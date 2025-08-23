package repository

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LifeTimeModel struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type User struct {
	ID uint `gorm:"primarykey"`
	LifeTimeModel
	UserID       string `gorm:"size:36;unique"`
	Login        string `gorm:"size:256;unique"`
	HashPassword string `gorm:"type:bytes"`
	UserDatas    []UserData
}

type UserData struct {
	ID uint `gorm:"primarykey"`
	LifeTimeModel
	UserID      uint
	DataType    string `gorm:"type:varchar(16)"`
	Description string `gorm:"type:text"`
	FileName    string `gorm:"size:512"`
	SecretKey   string `gorm:"size:1024"`
}

type DataType string

const (
	DataTypeCard          DataType = "card"
	DataTypeLoginPassword DataType = "logpass"
	DataTypeFile          DataType = "file"
)

func (d *DataType) Scan(value interface{}) error {
	if value == nil {
		*d = ""
		return nil
	}
	dt, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("unsupported Scan value type: %T", value)
	}
	*d = DataType(dt)
	return nil
}

func (d DataType) Value() (driver.Value, error) {
	if d == "" {
		return nil, nil
	}
	return string(d), nil
}
