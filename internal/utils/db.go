package utils

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func ConnectDB(config *Config, dialector ...gorm.Dialector) (*gorm.DB, error) {
	// Allow providing custom gorm.Dialector for mocks
	var dialectorVar gorm.Dialector
	if len(dialector) == 0 {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
			config.DBHost,
			config.DBUserName,
			config.DBUserPassword,
			config.DBName,
			config.DBPort,
		)
		dialectorVar = postgres.Open(dsn)
	} else {
		dialectorVar = dialector[0]
	}

	db, err := gorm.Open(dialectorVar, &gorm.Config{
		Logger: GormLogger{},
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

type AnyString struct{}

func (a AnyString) Match(v driver.Value) bool {
	_, ok := v.(string)
	return ok
}
