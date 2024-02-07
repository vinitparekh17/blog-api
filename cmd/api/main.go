package main

import (
	"context"
	"log"

	"github.com/jay-bhogayata/blogapi/config"
	"github.com/jay-bhogayata/blogapi/database"
	"github.com/jay-bhogayata/blogapi/handlers"
	"github.com/jay-bhogayata/blogapi/logger"
	"github.com/jay-bhogayata/blogapi/router"
	"github.com/jay-bhogayata/blogapi/server"
)

func main() {

	logger.Init()

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf(err.Error())
	}

	db, err := database.Init(context.Background(), cfg.Database.DBURL)
	if err != nil {
		log.Fatalf(err.Error())
	}

	query := database.New(db)

	handlers := handlers.NewHandlers(db, query, logger.Log)

	router := router.NewRouter(handlers)

	server := server.NewServer(cfg, router)

	err = server.Start()
	if err != nil {
		log.Fatalf("error in starting the server")
	}
}
