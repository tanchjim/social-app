package database

import (
	"fmt"
	"time"

	"social-app/internal/config"
	"social-app/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init initializes database connection
func Init(cfg *config.Config) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	var err error
	DB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // Disable prepared statements for Supabase pooler
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// Get underlying sql.DB to configure connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Auto migrate (ignore "already exists" errors)
	if err := autoMigrate(); err != nil {
		// Log but don't fail on "already exists" errors
		fmt.Printf("Migration warning: %v\n", err)
	}

	return nil
}

// autoMigrate runs auto migration for all models
func autoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Content{},
		&model.Comment{},
		&model.Like{},
		&model.RefreshToken{},
	)
}

// Close closes database connection
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns database instance
func GetDB() *gorm.DB {
	return DB
}
