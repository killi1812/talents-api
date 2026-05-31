package handlers

import (
	"net/http"
	"strconv"
	"talents-api/models"
	"talents-api/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
//	@Description	get talents, optional search and pagination
//	@Tags			talents
//	@Accept			json
//	@Produce		json
//	@Param			q		query		string	false	"Search query"
//	@Param			page	query		int		false	"Page number (default 1)"
//	@Param			limit	query		int		false	"Items per page (default 10)"
//	@Success		200		{object}	models.PaginatedResponse{data=[]models.Talent}
//	@Security		ApiKeyAuth
//	@Router			/api/talents [get]
func (h *TalentHandler) GetTalents(c *gin.Context) {
	query := c.Query("q")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	skip := int64((page - 1) * limit)

	var talents []models.Talent
	var total int64
	var err error

	if query != "" {
		zap.S().Debugf("Searching talents with query: %s, page: %d, limit: %d", query, page, limit)
		talents, total, err = h.repo.Search(c.Request.Context(), query, int64(limit), skip)
	} else {
		zap.S().Debugf("Fetching talents, page: %d, limit: %d", page, limit)
		talents, total, err = h.repo.GetAll(c.Request.Context(), int64(limit), skip)
	}

	if err != nil {
		zap.S().Errorf("Failed to get talents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.PaginatedResponse{
		Data:  talents,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

// CreateTalent godoc
//
//	@Summary	Create talent
//	@Tags		admin
//	@Accept		json
//	@Produce	json
//	@Param		talent	body		models.Talent	true	"Talent"
//	@Success	201		{object}	models.Talent
//	@Security	ApiKeyAuth
//	@Router		/admin/talents [post]
func (h *TalentHandler) CreateTalent(c *gin.Context) {
	var talent models.Talent
	if err := c.ShouldBindJSON(&talent); err != nil {
		zap.S().Warnf("Failed to bind talent JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	zap.S().Infof("Creating talent: %s", talent.Name)
	if err := h.repo.Create(c.Request.Context(), &talent); err != nil {
		zap.S().Errorf("Failed to create talent in DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, talent)
}

// UpdateTalent godoc
//
//	@Summary	Update talent
//	@Tags		admin
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string			true	"Talent ID"
//	@Param		talent	body		models.Talent	true	"Talent"
//	@Success	200		{object}	models.Talent
//	@Security	ApiKeyAuth
//	@Router		/admin/talents/{id} [put]
func (h *TalentHandler) UpdateTalent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.S().Warnf("Invalid UUID format for update: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	var talent models.Talent
	if err := c.ShouldBindJSON(&talent); err != nil {
		zap.S().Warnf("Failed to bind talent JSON for update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}
	talent.Id = id

	zap.S().Infof("Updating talent %s: %s", id, talent.Name)
	if err := h.repo.Update(c.Request.Context(), &talent); err != nil {
		zap.S().Errorf("Failed to update talent in DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, talent)
}

// DeleteTalent godoc
//
//	@Summary	Delete talent
//	@Tags		admin
//	@Produce	json
//	@Param		id	path		string	true	"Talent ID"
//	@Success	204	{object}	interface{}
//	@Security	ApiKeyAuth
//	@Router		/admin/talents/{id} [delete]
func (h *TalentHandler) DeleteTalent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		zap.S().Warnf("Invalid UUID format for delete: %s", idStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	zap.S().Infof("Deleting talent: %s", id)
	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		zap.S().Errorf("Failed to delete talent from DB: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
