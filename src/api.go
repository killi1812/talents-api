package main

import (
	"talents-api/app"
	"talents-api/handlers"
	"talents-api/middleware"
	"talents-api/repository"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "talents-api/docs"
)

type TalentApi struct {
	talentHandler *handlers.TalentHandler
	adminHandler  *handlers.AdminHandler
	apiKeyRepo    *repository.APIKeyRepository
}

func (a *TalentApi) NewGinApi(r *gin.Engine) {
	// Serve static files
	r.Static("/public", "./public")
	r.StaticFile("/", "./public/index.html")
	r.StaticFile("/search", "./public/search.html")

	// Version endpoint - unprotected
	r.GET("/version", app.VersionHandler)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(a.apiKeyRepo))
	{
		api.GET("/talents", a.talentHandler.GetTalents)
		api.POST("/talents", a.talentHandler.CreateTalent)
	}

	// Admin Routes
	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware(a.apiKeyRepo))
	admin.Use(middleware.AdminOnly())
	{
		admin.POST("/keys", a.adminHandler.CreateAPIKey)
	}
}
