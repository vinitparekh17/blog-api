package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jay-bhogayata/blogapi/internal/config"
	"github.com/jay-bhogayata/blogapi/internal/database"
	"github.com/jay-bhogayata/blogapi/internal/logger"
	openSearchClient "github.com/jay-bhogayata/blogapi/internal/opensearch"
	"github.com/jay-bhogayata/blogapi/internal/server"
	"github.com/jay-bhogayata/blogapi/internal/store"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig()
	logger.Init()

	if err != nil {
		logger.Log.Error("Error in loading the config file", "error", err)
		return
	}

	if len(os.Args) == 1 {
		InitApi(cfg)
		return
	}

	choice, err := strconv.Atoi(os.Args[1:][0])
	if err != nil {
		MigrationOption()
		return
	}

	for {
		switch choice {

		case 1:
			m, err := migrate.New("file://migrations", cfg.Database.DBURL)
			if err != nil {
				logger.Log.Error(err.Error())
			}
			if err := m.Up(); err != nil {
				logger.Log.Error(err.Error())
			} else {
				fmt.Println("Migration up successfull")
			}
			return

		case 2:
			m, err := migrate.New("file://migrations", cfg.Database.DBURL)
			if err != nil {
				logger.Log.Error(err.Error())
			}
			if err := m.Down(); err != nil {
				logger.Log.Error(err.Error())
			} else {
				fmt.Println("Migration down successfull")
			}
			return

		default:
			MigrationOption()
			return
		}
	}
}

func MigrationOption() {
	logger.Log.Info("____________________________________________________\n")
	logger.Log.Info("Please provide valid input\n")
	logger.Log.Info("1. Migrate Up")
	logger.Log.Info("2. Migrate Down")
	logger.Log.Info("____________________________________________________")
}

func InitApi(cfg *config.Config) {

	db, dbErr := database.Init(context.Background(), cfg.Database.DBURL)
	if dbErr != nil {
		log.Fatalf(dbErr.Error())
	}

	osc, ose := openSearchClient.NewOpenSearchClient(&openSearchClient.OpenSearchConfig{
		URL:      cfg.OpenSearch.URL,
		UserName: cfg.OpenSearch.UserName,
		Password: cfg.OpenSearch.Password,
	})

	if ose != nil {
		log.Fatalf(ose.Error())
	}

	query := database.New(db)

	store.UseBlogStore(&store.BlogStoreType{
		DB:               db,
		OpenSearchClient: osc,
		Query:            query,
		Logger:           logger.Log,
		Config:           cfg,
	})

	server.UseGRPCServer()
}
