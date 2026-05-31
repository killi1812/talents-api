package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type VersionInfo struct {
	Version    string `json:"version"`
	Build      string `json:"build"`
	Timestamp  string `json:"timestamp"`
	CommitHash string `json:"commitHash"`
}

// VersionHandler godoc
//
//	@Summary		Get the version
//	@Description	Get the version of the server
//	@Tags			version
//	@Produce		json
//	@Success		200	{object}	VersionInfo
//	@Router			/version [get]
func VersionHandler(c *gin.Context) {
	c.JSON(http.StatusOK, VersionInfo{
		Version:    Version,
		Build:      Build,
		Timestamp:  BuildTimestamp,
		CommitHash: CommitHash,
	})
}
