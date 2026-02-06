package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetDBConnection returns a GORM database connection
func GetDBConnection() (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "b2b_user"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "b2b_password"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "b2b_platform"
	}

	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// CreateSchema creates a PostgreSQL schema if it doesn't exist
func CreateSchema(db *gorm.DB, schemaName string) error {
	sql := fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName)
	return db.Exec(sql).Error
}
