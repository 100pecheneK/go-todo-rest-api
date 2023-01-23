package main

import (
	"log"

	"github.com/100pecheneK/go-todo-rest-api.git/internal/server"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/handler"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	repos := repository.NewRepository()
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
