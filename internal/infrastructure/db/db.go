package db

import (
	"fmt"
	"log"

	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn :=
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
			cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode,
		)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	} else {
		log.Println("Database connection successful")
	}

	return db, nil
}
