package handlers

import (
	"net/http"
	"talents-api/models"
	"talents-api/repository"

	"github.com/gin-gonic/gin"
)

type TalentHandler struct {
	repo *repository.TalentRepository
}

func NewTalentHandler(repo *repository.TalentRepository) *TalentHandler {
	return &TalentHandler{repo: repo}
}

// GetTalents godoc
//
//	@Summary		List talents
//	@Description	get talents, optional search
//	@Tags			talents
//	@Accept			json
//	@Produce		json
//	@Param			q	query	string	false	"Search query"
//	@Success		200	{array}	models.Talent
//	@Router			/api/talents [get]
//	@Security		ApiKeyAuth
func (h *TalentHandler) GetTalents(c *gin.Context) {
	query := c.Query("q")
	var talents []models.Talent
	var err error

	if query != "" {
		talents, err = h.repo.Search(c.Request.Context(), query)
	} else {
		talents, err = h.repo.GetAll(c.Request.Context())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, talents)
}

// CreateTalent godoc
//
//	@Summary	Create talent
//	@Tags		talents
//	@Accept		json
//	@Produce	json
//	@Param		talent	body		models.Talent	true	"Talent"
//	@Success	201		{object}	models.Talent
//	@Router		/api/talents [post]
//	@Security	ApiKeyAuth
func (h *TalentHandler) CreateTalent(c *gin.Context) {
	var talent models.Talent
	if err := c.ShouldBindJSON(&talent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Create(c.Request.Context(), &talent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, talent)
}
