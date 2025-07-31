package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
)

// HandleAllMedia handles retrieval of all media files (both images and documents)
func (h *MediaHandler) HandleAllMedia(c *gin.Context) {
	// Get optional media type filter from query parameter
	mediaType := c.Query("type")

	var entries []models.Media

	if mediaType != "" {
		// Filter by media type
		entries = h.repo.GetMediaByType(models.MediaType(mediaType))
	} else {
		// Get all media
		entries = h.repo.GetAllMedia()
	}

	c.JSON(http.StatusOK, entries)
}

// HandleAllImages provides backward compatibility for retrieving all images
func (h *MediaHandler) HandleAllImages(c *gin.Context) {
	entries := h.repo.GetAllImages()
	c.JSON(http.StatusOK, entries)
}

// HandleAllDocs provides backward compatibility for retrieving all documents
func (h *MediaHandler) HandleAllDocs(c *gin.Context) {
	entries := h.repo.GetAllDocs()
	c.JSON(http.StatusOK, entries)
}
