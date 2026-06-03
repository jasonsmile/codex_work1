package main

import (
	"log"
	"os"

	"drug-info/backend/config"
	"drug-info/backend/db"
	"drug-info/backend/handlers"
	"drug-info/backend/logger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := logger.Init("log"); err != nil {
		log.Fatalf("init logger failed: %v", err)
	}
	defer logger.Close()

	appConfig, err := config.Load("config.yaml")
	if err != nil {
		logger.Error("load config failed", logger.Field{Key: "error", Value: err})
		log.Fatalf("load config failed: %v", err)
	}

	database, err := db.Connect(appConfig.MySQL)
	if err != nil {
		logger.Error("connect database failed", logger.Field{Key: "error", Value: err})
		log.Fatalf("connect database failed: %v", err)
	}

	enforcer, err := db.NewRBACEnforcer(database)
	if err != nil {
		logger.Error("init casbin failed", logger.Field{Key: "error", Value: err})
		log.Fatalf("init casbin failed: %v", err)
	}

	router := gin.New()
	router.Use(logger.RecoveryMiddleware())
	router.Use(logger.AccessMiddleware())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	drugHandler := handlers.NewDrugHandler(database)
	specimenHandler := handlers.NewSpecimenHandler(database)
	userHandler := handlers.NewUserHandler(database)
	fileHandler := handlers.NewFileHandler(database, appConfig.QiniuKodo)
	traceCodeHandler := handlers.NewTraceCodeHandler(database, appConfig.BaiduOCR)

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
			protected.POST("/fileUploadAndDownload/upload", fileHandler.Upload)
			protected.GET("/fileUploadAndDownload/get", fileHandler.List)
			protected.GET("/fileUploadAndDownload/download/:id", fileHandler.Download)
			protected.POST("/files/upload", fileHandler.Upload)
			protected.GET("/files/get", fileHandler.List)
			protected.GET("/files/download/:id", fileHandler.Download)
			protected.POST("/trace_codes/recognize", traceCodeHandler.Recognize)
			protected.POST("/trace_codes/confirm", traceCodeHandler.Confirm)
			protected.GET("/trace_codes/get", traceCodeHandler.List)
			protected.POST("/users/add", userHandler.CreateUser)
			protected.GET("/users/get", userHandler.ListUsers)
			protected.POST("/users/delete", userHandler.DeleteUser)
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8888"
	}

	logger.Access("server starting", logger.Field{Key: "port", Value: port})
	if err := router.Run(":" + port); err != nil {
		logger.Error("start server failed", logger.Field{Key: "error", Value: err})
		log.Fatalf("start server failed: %v", err)
	}
}
