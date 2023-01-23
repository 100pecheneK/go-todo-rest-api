package main

import (
	"log"
	"os"

	"github.com/100pecheneK/go-todo-rest-api.git/internal/server"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/handler"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("psql.host"),
		Port:     viper.GetString("psql.port"),
		Username: viper.GetString("psql.username"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   viper.GetString("psql.DBName"),
		SSLMode:  viper.GetString("psql.SSLMode"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	} else {
		log.Println("db connect successfull")
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	port := viper.GetString("port")
	log.Println("starting server on port:", port)
	if err := srv.Run(port, handlers.InitRoutes()); err != nil {
		log.Fatalf("error while running http server: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
