package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/100pecheneK/go-todo-rest-api.git/internal/server"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/handler"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/repository"
	"github.com/100pecheneK/go-todo-rest-api.git/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Todo App API
// @version 1.0
// @description API server for todolist application

// @host localhost:3000
// @basePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	logrus.Println("initializing configs")
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}
	logrus.Println("configs initialized")
	logrus.Println("loading env variables")
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}
	logrus.Println("env variables loaded")
	logrus.Println("initializing db connect")
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("psql.host"),
		Port:     viper.GetString("psql.port"),
		Username: viper.GetString("psql.username"),
		Password: os.Getenv("PSQL_PASSWORD"),
		DBName:   viper.GetString("psql.DBName"),
		SSLMode:  viper.GetString("psql.SSLMode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
	} else {
		logrus.Println("db connect successfull")
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		port := viper.GetString("port")
		logrus.Println("starting server on port:", port)
		if err := srv.Run(port, handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error while running http server: %s", err.Error())
		}
	}()

	logrus.Print("server started!")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	logrus.Print("server shutting down")
	if err := srv.Stop(context.Background()); err != nil {
		logrus.Errorf("error server: %s", err.Error())
	}
	logrus.Print("db connection closing")
	if err := db.Close(); err != nil {
		logrus.Errorf("error on db connection close: %s", err.Error())
	}
	logrus.Print("server stoped")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
