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

	enforcer, err := db.NewRBACEnforcer(database)
	if err != nil {
		log.Fatalf("init casbin failed: %v", err)
	}

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	drugHandler := handlers.NewDrugHandler(database)
	specimenHandler := handlers.NewSpecimenHandler(database)
	userHandler := handlers.NewUserHandler(database)

	api := router.Group("/api")
	{
		api.POST("/users/login", userHandler.Login)

		protected := api.Group("")
		protected.Use(handlers.RBACMiddleware(enforcer))
		{
			protected.POST("/drugs/add", drugHandler.CreateDrug)
			protected.GET("/drugs/get", drugHandler.ListDrugs)
			protected.POST("/specimens/add", specimenHandler.CreateApplication)
			protected.POST("/specimens/import/preview", specimenHandler.PreviewImportApplications)
			protected.POST("/specimens/import", specimenHandler.ImportApplications)
			protected.GET("/specimens/get", specimenHandler.ListApplications)
			protected.POST("/users/add", userHandler.CreateUser)
			protected.GET("/users/get", userHandler.ListUsers)
			protected.POST("/users/delete", userHandler.DeleteUser)
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8888"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("start server failed: %v", err)
	}
}
