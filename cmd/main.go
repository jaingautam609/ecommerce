package main

import (
	"ecommerce/database"
	"ecommerce/database/routes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	_ "net/http"
	"os"
)

func main() {
	if err := database.ConnectAndMigrate(
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("databaseName"),
		os.Getenv("user"),
		os.Getenv("password"),
		database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	router := gin.Default()
	routes.ServerRoutes(router)
	if err := router.Run(":8081"); err != nil {
		logrus.Panicf("Failed to start server with error: %+v", err)
	}
}
