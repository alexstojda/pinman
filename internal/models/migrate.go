package models

import (
	"fmt"
	"gorm.io/gorm"
)

func Migrate(gormDb *gorm.DB) error {
	err := gormDb.AutoMigrate(
		&User{},
		&League{},
		&Location{},
	)
	if err != nil {
		return fmt.Errorf("migrating models: %w", err)
	}
	return nil
}
