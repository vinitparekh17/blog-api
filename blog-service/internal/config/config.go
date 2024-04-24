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
	OpenSearch  struct {
		URL      string
		UserName string
		Password string
	}
}

func init() {
	godotenv.Load()
}

func LoadConfig() (*Config, error) {
	var cfg Config

	cfg.Database.DBURL = os.Getenv("DATABASE_URL")
	if cfg.Database.DBURL == "" {
		return nil, errors.New("DATABASE_URL env not found")
	}

	cfg.OpenSearch.URL = os.Getenv("OPENSEARCH_URL")
	if cfg.OpenSearch.URL == "" {
		return nil, errors.New("OPENSEARCH_URL env not found")
	}

	cfg.OpenSearch.UserName = os.Getenv("OPENSEARCH_USERNAME")
	if cfg.OpenSearch.UserName == "" {
		return nil, errors.New("OPENSEARCH_USERNAME env not found")
	}

	cfg.OpenSearch.Password = os.Getenv("OPENSEARCH_PASSWORD")
	if cfg.OpenSearch.Password == "" {
		return nil, errors.New("OPENSEARCH_PASSWORD env not found")
	}

	cfg.JWTSecret = os.Getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET env not found")
	}

	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		log.Println("No SERVER_PORT found in env file, using default port 8080")
		cfg.Server.Port = "8080"
	}
	if _, err := strconv.Atoi(cfg.Server.Port); err != nil {
		log.Println("Invalid SERVER_PORT found in env file, using default port 8080")
		cfg.Server.Port = "8080"
	}

	cfg.Env = os.Getenv("ENV")
	if cfg.Env == "" {
		log.Println("No ENV found in env file, using default development")
		cfg.Env = "development"
	}

	return &cfg, nil
}
