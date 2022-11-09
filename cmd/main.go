package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bambi2/todo-app/pkg/config"
	"github.com/bambi2/todo-app/pkg/delivery"
	"github.com/bambi2/todo-app/pkg/repository"
	"github.com/bambi2/todo-app/pkg/server"
	"github.com/bambi2/todo-app/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error initializing configs: %s\n", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s\n", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("error occured while connecting to database: %s\n", err.Error())
	}

	repos := repository.NewRepository(db)
	service := service.NewService(repos)
	handler := delivery.NewHandler(service)
	srv := new(server.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			log.Fatalf("error occured while running http server: %s\n", err.Error())
		}
	}()

	fmt.Println("To Do App Started...")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("To Do App Stopped...")

	if err := srv.ShutDown(context.Background()); err != nil {
		log.Fatalf("error occured when shutting down the server: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		log.Fatalf("error occured when closing database connection: %s", err.Error())
	}
}
