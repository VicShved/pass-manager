package repository

import (
	"database/sql/driver"
	"fmt"
	"time"
)

// LifeTimeModel for known of rows create & update
type LifeTimeModel struct {
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

// User model for manage users of application
type User struct {
	ID uint `gorm:"primarykey"`
	LifeTimeModel
	UserID       string `gorm:"size:36;unique"`
	Login        string `gorm:"size:256;unique"`
	HashPassword string `gorm:"type:bytes"`
	UserDatas    []UserData
}

// UserData for manage user data
type UserData struct {
	ID uint `gorm:"primarykey"`
	LifeTimeModel
	UserID      uint
	DataType    string `gorm:"type:varchar(16)"`
	Description string `gorm:"type:text"`
	FileName    string `gorm:"size:512"`
	SecretKey   string `gorm:"size:1024"`
}

// DataType for user data
type DataType string

const (
	DataTypeCard          DataType = "card"
	DataTypeLoginPassword DataType = "logpass"
	DataTypeFile          DataType = "file"
)

// Scan for DataType
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

// Valuefor DataType
func (d DataType) Value() (driver.Value, error) {
	if d == "" {
		return nil, nil
	}
	return string(d), nil
}
