package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jay-bhogayata/blogapi/config"
	"github.com/jay-bhogayata/blogapi/database"
	"github.com/jay-bhogayata/blogapi/handlers"
	"github.com/jay-bhogayata/blogapi/logger"
	"github.com/jay-bhogayata/blogapi/router"
	"github.com/jay-bhogayata/blogapi/server"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	logger.Init()

	cfg, err := config.LoadConfig()

	if err != nil {
		logger.Log.Info("error in loading the config")
	}

	var choice int16

	for {

		fmt.Println(`
1. Start Blog API
2. Migration UP
3. Migration Down
4. Exit

Enter the number to choose:
( Default 1 )`)
		fmt.Scanln(&choice)

		switch choice {

		case 0:
			fmt.Println("Starting blog API as per default choice...")
			InitApi(cfg)

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

		case 4:
			return

		default:
			fmt.Println("Invalid choice")
		}
	}
}

func InitApi(cfg *config.Config) {

	db := InitDatabase(cfg.Database.DBURL)

	query := database.New(db)

	handlers := handlers.NewHandlers(cfg, db, query, logger.Log)

	router := router.NewRouter(cfg, handlers)

	server := server.NewServer(cfg, router)

	err := server.Start()
	if err != nil {
		log.Fatalf("error in starting the server")
	}
}

func InitDatabase(url string) *pgxpool.Pool {
	db, err := database.Init(context.Background(), url)
	if err != nil {
		log.Fatalf(err.Error())
	}

	return db
}
