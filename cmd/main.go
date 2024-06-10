package main

import (
	"fmt"
	"log"
	mentor "pet"
	"pet/internal/handler"
	"pet/internal/repository"
	"pet/internal/service"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing config: %s", err.Error())
	}
	// Не забыть настроить конфиг
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     "localhost", //viper.GetString("db.host"),
		Port:     "5436",      //viper.GetString("db.post"),
		Username: "postgres",  //viper.GetString("db.username"),
		Password: "qwerty",    //viper.GetString("db.password"),
		DBname:   "postgres",  //viper.GetString("db.dbname"),
		SSLmode:  "disable",   //viper.GetString("db.sslmode"),
	})
	if err != nil {
		log.Fatalf("failed to initialized db %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	adminName := "Adminovich"
  adminUserName:="Admin"
	AdPassword:="passwordForAdmin"
	_,err=services.Authorization.CreateAdmin(adminName, adminUserName,AdPassword)
	if err != nil {
		fmt.Printf("admin created or error %s", err.Error())
	}

	server := new(mentor.Server)
	if err := server.Run("8080", handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error server init %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
