package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

// HandleMediaDelete handles the deletion of both images and documents
func (h *MediaHandler) HandleMediaDelete(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File name is required",
		})
		return
	}

	// Get the media type from the query parameter
	mediaType := c.Query("type")
	if mediaType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Media type is required",
		})
		return
	}

	// Delete the media from the database
	deletedFileName, success := h.repo.DeleteMedia(fileName)
	if !success {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Media not found",
		})
		return
	}

	// Delete the file from the unified media directory
	err := util.DeleteUnifiedMediaFile(deletedFileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete media",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Media deleted successfully",
		"fileName": deletedFileName,
	})
}

// HandleImageDelete provides backward compatibility for image deletion
func (h *MediaHandler) HandleImageDelete(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image name is required",
		})
		return
	}

	deletedFileName, success := h.repo.DeleteImage(fileName)
	if !success {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Image not found",
		})
		return
	}

	err := util.DeleteFile(deletedFileName, "images")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete image",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Image deleted successfully",
		"fileName": deletedFileName,
	})
}

// HandleDocDelete provides backward compatibility for document deletion
func (h *MediaHandler) HandleDocDelete(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Doc name is required",
		})
		return
	}

	deletedFileName, success := h.repo.DeleteDoc(fileName)
	if !success {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Document not found",
		})
		return
	}

	err := util.DeleteFile(deletedFileName, "docs")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Failed to delete document",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Document deleted successfully",
		"fileName": deletedFileName,
	})
}
