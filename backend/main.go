package main

import (
	"log"
	"os"

	"drug-info/backend/db"
	"drug-info/backend/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	database, err := db.Connect()
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept"},
	}))

	drugHandler := handlers.NewDrugHandler(database)
	specimenHandler := handlers.NewSpecimenHandler(database)

	api := router.Group("/api")
	{
		api.POST("/drugs/add", drugHandler.CreateDrug)
		api.GET("/drugs/get", drugHandler.ListDrugs)
		api.POST("/specimens/add", specimenHandler.CreateApplication)
		api.GET("/specimens/get", specimenHandler.ListApplications)
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8888"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
