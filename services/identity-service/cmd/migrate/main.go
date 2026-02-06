package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "b2b_user")
	dbPassword := getEnv("DB_PASSWORD", "b2b_password")
	dbName := getEnv("DB_NAME", "b2b_platform")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	// Read and execute migration files
	migrationDir := filepath.Join("migrations")
	files, err := filepath.Glob(filepath.Join(migrationDir, "*.sql"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read migration files: %v\n", err)
		os.Exit(1)
	}

	for _, file := range files {
		fmt.Printf("Running migration: %s\n", file)
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read migration file %s: %v\n", file, err)
			os.Exit(1)
		}

		if _, err := db.Exec(string(sqlBytes)); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to execute migration %s: %v\n", file, err)
			os.Exit(1)
		}
	}

	fmt.Println("Migrations completed successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
