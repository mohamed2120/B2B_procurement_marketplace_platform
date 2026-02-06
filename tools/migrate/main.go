package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func main() {
	var serviceName string
	var dsn string
	var dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode string

	flag.StringVar(&serviceName, "service", "", "Service name (required)")
	flag.StringVar(&dsn, "dsn", "", "Database DSN (optional, can use individual params)")
	flag.StringVar(&dbHost, "host", getEnv("DB_HOST", "localhost"), "Database host")
	flag.StringVar(&dbPort, "port", getEnv("DB_PORT", "5432"), "Database port")
	flag.StringVar(&dbUser, "user", getEnv("DB_USER", "b2b_user"), "Database user")
	flag.StringVar(&dbPassword, "password", getEnv("DB_PASSWORD", "b2b_password"), "Database password")
	flag.StringVar(&dbName, "dbname", getEnv("DB_NAME", "b2b_platform"), "Database name")
	flag.StringVar(&dbSSLMode, "sslmode", getEnv("DB_SSLMODE", "disable"), "SSL mode")
	flag.Parse()

	if serviceName == "" {
		fmt.Fprintf(os.Stderr, "Error: --service is required\n")
		os.Exit(1)
	}

	// Build DSN if not provided
	if dsn == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)
	}

	// Connect to database
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

	// Get migration directory
	// Find repo root by looking for services/ directory
	repoRoot := findRepoRoot()
	migrationDir := filepath.Join(repoRoot, "services", serviceName, "migrations")
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		fmt.Printf("No migrations folder found for service '%s' at %s, skipping...\n", serviceName, migrationDir)
		os.Exit(0)
	}

	// Read migration files
	files, err := filepath.Glob(filepath.Join(migrationDir, "*.sql"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read migration files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Printf("No migration files found for service '%s', skipping...\n", serviceName)
		os.Exit(0)
	}

	// Sort files by name
	sort.Strings(files)

	// Create schema_migrations table if it doesn't exist
	// Use a shared schema for tracking migrations across all services
	if _, err := db.Exec(`
		CREATE SCHEMA IF NOT EXISTS migrations;
		CREATE TABLE IF NOT EXISTS migrations.schema_migrations (
			id SERIAL PRIMARY KEY,
			service_name VARCHAR(255) NOT NULL,
			migration_file VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			UNIQUE(service_name, migration_file)
		);
	`); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create schema_migrations table: %v\n", err)
		os.Exit(1)
	}

	// Run migrations
	for _, file := range files {
		fileName := filepath.Base(file)

		// Check if migration already applied
		var count int
		err := db.QueryRow(
			"SELECT COUNT(*) FROM migrations.schema_migrations WHERE service_name = $1 AND migration_file = $2",
			serviceName, fileName,
		).Scan(&count)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to check migration status: %v\n", err)
			os.Exit(1)
		}

		if count > 0 {
			fmt.Printf("Skipping already applied migration: %s\n", fileName)
			continue
		}

		// Read migration file
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read migration file %s: %v\n", file, err)
			os.Exit(1)
		}

		sqlContent := string(sqlBytes)
		if strings.TrimSpace(sqlContent) == "" {
			fmt.Printf("Skipping empty migration file: %s\n", fileName)
			continue
		}

		// Execute migration in a transaction
		tx, err := db.Begin()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to begin transaction: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Running migration: %s\n", fileName)
		if _, err := tx.Exec(sqlContent); err != nil {
			tx.Rollback()
			fmt.Fprintf(os.Stderr, "Failed to execute migration %s: %v\n", fileName, err)
			os.Exit(1)
		}

		// Record migration as applied
		if _, err := tx.Exec(
			"INSERT INTO migrations.schema_migrations (service_name, migration_file) VALUES ($1, $2)",
			serviceName, fileName,
		); err != nil {
			tx.Rollback()
			fmt.Fprintf(os.Stderr, "Failed to record migration: %v\n", err)
			os.Exit(1)
		}

		if err := tx.Commit(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to commit transaction: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✅ Applied migration: %s\n", fileName)
	}

	fmt.Printf("✅ Migrations completed successfully for service '%s'!\n", serviceName)
}

func findRepoRoot() string {
	// Start from current working directory
	dir, err := os.Getwd()
	if err != nil {
		// Fallback: try to get executable directory
		exe, err := os.Executable()
		if err == nil {
			dir = filepath.Dir(exe)
		}
	}

	// Walk up the directory tree to find repo root (has services/ directory)
	for {
		servicesPath := filepath.Join(dir, "services")
		if _, err := os.Stat(servicesPath); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			break
		}
		dir = parent
	}

	// Fallback: return current directory
	return "."
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
