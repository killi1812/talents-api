package main

import (
	"context"
	"encoding/json"
	"os"

	"talents-api/app"
	"talents-api/handlers"
	"talents-api/models"
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

	// Seeding
	if os.Getenv("SEED_DATA") == "true" {
		filePath := "./data/talents.json"
		if _, err := os.Stat(filePath); err == nil {
			zap.S().Infof("Seeding data from %s", filePath)
			fileData, err := os.ReadFile(filePath)
			if err == nil {
				var talents []models.Talent
				if err := json.Unmarshal(fileData, &talents); err == nil {
					count, err := talentRepo.Seed(context.Background(), talents)
					if err != nil {
						zap.S().Errorf("Failed to seed talents: %v", err)
					} else {
						zap.S().Infof("Successfully seeded %d talents", count)
					}
				} else {
					zap.S().DPanicf("Failed to Unmarshal seed data: %v", err)
				}
			}
		} else {
			zap.S().Warnf("Seed file not found at %s", filePath)
		}
	}

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
