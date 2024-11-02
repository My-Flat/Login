package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Config struct to hold the database configuration
type Config struct {
	DBHost                   string
	DBUser                   string
	DBPassword               string
	DBName                   string
	DBPort                   string
	DBInstanceConnectionName string
	DBCredentialsFile        string
}

// LoadDBConfig loads the database configuration from environment variables
func LoadDBConfig() (Config, error) {
	var missingVars []string

	err := godotenv.Load() // Add this line to load the .env file
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load the database configuration from environment variables
	cfg := Config{
		DBHost:                   os.Getenv("DB_HOST"),
		DBUser:                   os.Getenv("DB_USER"),
		DBPassword:               os.Getenv("DB_PASSWORD"),
		DBName:                   os.Getenv("DB_NAME"),
		DBPort:                   os.Getenv("DB_PORT"),
		DBInstanceConnectionName: os.Getenv("DB_INSTANCE_CONNECTION_NAME"),
		DBCredentialsFile:        os.Getenv("DB_CREDENTIALS_FILE"),
	}

	// Check for missing environment variables
	if cfg.DBHost == "" {
		missingVars = append(missingVars, "DB_HOST")
	}
	if cfg.DBUser == "" {
		missingVars = append(missingVars, "DB_USER")
	}
	if cfg.DBPassword == "" {
		missingVars = append(missingVars, "DB_PASSWORD")
	}
	if cfg.DBName == "" {
		missingVars = append(missingVars, "DB_NAME")
	}
	if cfg.DBPort == "" {
		missingVars = append(missingVars, "DB_PORT")
	}
	if cfg.DBInstanceConnectionName == "" {
		missingVars = append(missingVars, "DB_INSTANCE_CONNECTION_NAME")
	}
	if cfg.DBCredentialsFile == "" {
		missingVars = append(missingVars, "DB_CREDENTIALS_FILE")
	}

	if len(missingVars) > 0 {
		return cfg, fmt.Errorf("missing environment variables: %v", missingVars)
	}

	return cfg, nil
}

func ConnectDB(cfg Config) (*gorm.DB, error) {
	// dsn := fmt.Sprintf("%s:%s@unix(/cloudsql/%s)/%s?parseTime=true",
	// 	cfg.DBUser, cfg.DBPassword, cfg.DBInstanceConnectionName, cfg.DBName,
	// )

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
	// Connect using GORM
	sqlDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database: " + err.Error())
	}

	sql, err := sqlDB.DB()

	if err != nil {
		return nil, fmt.Errorf("failed to configure database connection pool: %w", err)
	}

	err = sql.Ping()
	if err != nil {
		log.Fatal(err)
	}

	sql.SetMaxIdleConns(10)
	sql.SetMaxOpenConns(100)
	sql.SetConnMaxLifetime(time.Hour)

	log.Println("Connected to Cloud SQL database")
	return sqlDB, nil
}
