package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jay-bhogayata/blogapi/config"
	"github.com/jay-bhogayata/blogapi/database"
	"github.com/jay-bhogayata/blogapi/handlers"
	"github.com/jay-bhogayata/blogapi/logger"
	openSearchClient "github.com/jay-bhogayata/blogapi/opensearch"
	"github.com/jay-bhogayata/blogapi/router"
	"github.com/jay-bhogayata/blogapi/server"

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

	if len(os.Args) <= 1 {
		logger.Log.Info("Arguments not found")
		logger.Log.Info("Starting blog API as per default choice...")
		InitApi(cfg)
	}

	choice, err := strconv.Atoi(os.Args[1:][0])
	if err != nil {
		fmt.Println(`
1. Start Blog API
2. Migration UP
3. Migration Down

Invalid Argument, Choose between 1 to 3`)
		return
	}

	for {
		switch choice {

		case 1:
			InitApi(cfg)

		case 2:
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

		case 3:
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

		case 4:
			return

		default:
			logger.Log.Warn("Invalid choice")
			return
		}
	}
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

	handlers := handlers.NewHandlers(&handlers.Handlers{
		DB:               db,
		OpenSearchClient: osc,
		Query:            query,
		Logger:           logger.Log,
		Config:           cfg,
	})

	router := router.NewRouter(cfg, handlers)

	server := server.NewServer(cfg, router)

	err := server.Start()
	if err != nil {
		log.Fatalf("error in starting the server")
	}
}
