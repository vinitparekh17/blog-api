package config

import (
	"errors"
	"os"
	"strconv"

	"github.com/jay-bhogayata/blogapi/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port string
	}
	Database struct {
		DBURL string
	}
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		logger.Log.Error("error in loading .env file")
		return nil, err
	}

	var cfg Config

	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		logger.Log.Warn("no SERVER_PORT env variable provided defaulting to port 8080")
		cfg.Server.Port = "8080"
	}
	if _, err := strconv.Atoi(cfg.Server.Port); err != nil {
		logger.Log.Warn("invalid SERVER_PORT, using default port 8080")

		cfg.Server.Port = "8080"
	}

	cfg.Database.DBURL = os.Getenv("DATABASE_URL")
	if cfg.Database.DBURL == "" {
		logger.Log.Error("no SERVER_PORT env variable provided defaulting to port 8080")
		return nil, errors.New("DATABASE_URL env not found")
	}
	return &cfg, nil
}
