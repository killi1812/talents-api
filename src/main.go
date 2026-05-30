package main

import (
	"os"

	"talents-api/app"
	"talents-api/handlers"
	"talents-api/repository"

	"go.uber.org/zap"
)

//	@title			Talent API
//	@version		1.0
//	@description	Talent API with API Key auth
//	@host			localhost:8080
//	@BasePath		/

// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						X-API-KEY
// @description				Enter your API key (admin or generated)
func main() {
	app.Setup()
	zap.S().Infof("Starting %s version %s (build: %s, commit: %s)", app.APP_NAME, app.Version, app.Build, app.CommitHash)

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
