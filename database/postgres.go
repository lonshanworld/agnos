package database

import (
	"agnos_candidate_assignment/config"
	"agnos_candidate_assignment/models"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresConnection(configuration *config.Config) (*gorm.DB, error) {
	dsn := configuration.DatabaseUrl
	if !strings.Contains(dsn, "preferSimpleProtocol") {
		if strings.Contains(dsn, "?") {
			dsn += "&preferSimpleProtocol=true"
		} else {
			dsn += "?preferSimpleProtocol=true"
		}
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(1)
	sqlDB.SetMaxIdleConns(0)

	if err := db.AutoMigrate(
		&models.Hospital{},
		&models.Staff{},
		&models.Patient{},
	); err != nil {
		log.Printf("auto migrate error: %v", err)
		return nil, err
	}

	return db, nil
}

// for seeding
func NewPostgresConnectionNoMigrate(configuration *config.Config) (*gorm.DB, error) {
	dsn2 := configuration.DatabaseUrl
	if !strings.Contains(dsn2, "preferSimpleProtocol") {
		if strings.Contains(dsn2, "?") {
			dsn2 += "&preferSimpleProtocol=true"
		} else {
			dsn2 += "?preferSimpleProtocol=true"
		}
	}

	db, err := gorm.Open(postgres.Open(dsn2), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(2)
	sqlDB.SetMaxIdleConns(1)

	return db, nil
}
