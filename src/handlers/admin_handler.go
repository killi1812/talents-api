package handlers

import (
	"net/http"
	"talents-api/models"
	"talents-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminHandler struct {
	apiKeyRepo *repository.APIKeyRepository
}

func NewAdminHandler(apiKeyRepo *repository.APIKeyRepository) *AdminHandler {
	return &AdminHandler{apiKeyRepo: apiKeyRepo}
}

// CreateAPIKey godoc
//
//	@Summary	Create API key
//	@Tags		admin
//	@Accept		json
//	@Produce	json
//	@Param		apikey	body		models.APIKey	false	"API Key description"
//	@Success	201		{object}	models.APIKey
//	@Router		/admin/keys [post]
//	@Security	ApiKeyAuth
func (h *AdminHandler) CreateAPIKey(c *gin.Context) {
	var key models.APIKey
	c.ShouldBindJSON(&key)

	if key.Key == "" {
		key.Key = uuid.New().String()
	}

	if err := h.apiKeyRepo.Create(c.Request.Context(), &key); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, key)
}
