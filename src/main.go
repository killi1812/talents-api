package main

import (
	"os"
	"talents-api/app"
	"talents-api/handlers"
	"talents-api/repository"

	"go.uber.org/zap"
)

func main() {
	app.Setup()

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	client, err := repository.ConnectMongo(mongoURI)
	if err != nil {
		zap.S().Fatalf("Failed to connect to MongoDB: %v", err)
	}
	db := client.Database("talentdb")

	talentRepo := repository.NewTalentRepository(db)
	apiKeyRepo := repository.NewAPIKeyRepository(db)

	talentHandler := handlers.NewTalentHandler(talentRepo)
	adminHandler := handlers.NewAdminHandler(apiKeyRepo)

	api := &TalentApi{
		talentHandler: talentHandler,
		adminHandler:  adminHandler,
		apiKeyRepo:    apiKeyRepo,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.Start(api, ":"+port)
}
