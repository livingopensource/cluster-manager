package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// Base contains common columns for all tables.
type Base struct {
	ID        string     `gorm:"primaryKey" json:"id,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (base *Base) BeforeCreate(tx *gorm.DB) error {
	uuID := uuid.New()
	tx.Statement.SetColumn("ID", uuID.String())
	return nil
}

func DBConfig() (*gorm.DB, error) {
	db := database{
		DBType:     viper.GetString("database.type"),
		DBHost:     viper.GetString("database.host"),
		DBName:     viper.GetString("database.name"),
		DBPassword: viper.GetString("database.password"),
		DBPort:     viper.GetString("database.port"),
		DBUser:     viper.GetString("database.user"),
	}

	dbInstance, err := db.dbConfig()
	if err != nil {
		return dbInstance, err
	}
	sqlDB, err := dbInstance.DB()
	if err != nil {
		return dbInstance, err
	}

	sqlDB.SetMaxIdleConns(300)
	sqlDB.SetMaxOpenConns(380)
	sqlDB.SetConnMaxIdleTime(time.Minute * 30)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return dbInstance, err
}
