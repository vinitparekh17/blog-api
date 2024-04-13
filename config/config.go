package config

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		DBURL string
	}
	Env         string
	EmailSender string
	JWTSecret   string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func LoadConfig() (*Config, error) {
	var cfg Config

	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		log.Println("No SERVER_PORT found in env file, using default port 8080")
		cfg.Server.Port = "8080"
	}
	if _, err := strconv.Atoi(cfg.Server.Port); err != nil {
		log.Println("Invalid SERVER_PORT found in env file, using default port 8080")
		cfg.Server.Port = "8080"
	}

	cfg.Database.DBURL = os.Getenv("DATABASE_URL")
	if cfg.Database.DBURL == "" {
		return nil, errors.New("DATABASE_URL env not found")
	}

	cfg.Env = os.Getenv("ENV")
	if cfg.Env == "" {
		log.Println("No ENV found in env file, using default DEV")
		cfg.Env = "DEV"
	}

	cfg.EmailSender = os.Getenv("MAILER_SENDER")
	if cfg.EmailSender == "" {
		return nil, errors.New("MAILER_SENDER env not found")
	}

	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET env not found")
	}

	return &cfg, nil
}
