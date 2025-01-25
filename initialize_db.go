package main

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func connectDB() (*gorm.DB, error) {
	dbHost := "localhost"
	dbPort := "5433"
	dbName := "go_scraping_dev"
	dbUser := "go-scraping-user"
	dbPassword := "postgrespassword"

	// Postgres用のDSN形式を修正
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", dbHost, dbUser, dbPassword, dbName, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("db connection failed: %v", err)
	}
	return db, nil
}

func migrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(&ItemMaster{}, &LatestItem{}); err != nil {
		return fmt.Errorf("db migration failed: %w", err)
	}
	return nil
}
