package main

import (
	"log"

	"github.com/100pecheneK/go-todo-rest-api.git/internal/server"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/handler"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/service"
)

func main() {
	repos := repository.NewRepository()
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	if err := srv.Run("3000", handlers.InitRoutes()); err != nil {
		log.Fatalf("error while running http server: %s", err.Error())
	}
}
