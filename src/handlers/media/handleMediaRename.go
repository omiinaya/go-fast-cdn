package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/kevinanielsen/go-fast-cdn/src/validations"
)

// HandleMediaRename handles the renaming of both images and documents
func (h *MediaHandler) HandleMediaRename(c *gin.Context) {
	oldName := c.PostForm("filename")
	newName := c.PostForm("newname")
	mediaType := c.PostForm("type")

	if mediaType == "" {
		c.String(http.StatusBadRequest, "Media type is required")
		return
	}

	err := validations.ValidateRenameInput(oldName, newName)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	filteredNewName, err := util.SanitizeFilename(newName)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = util.RenameUnifiedMediaFile(oldName, filteredNewName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to rename file: %s", err.Error())
		return
	}

	err = h.repo.RenameMedia(oldName, filteredNewName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to rename file in database: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "File renamed successfully"})
}

// HandleImageRename provides backward compatibility for image renaming
func (h *MediaHandler) HandleImageRename(c *gin.Context) {
	oldName := c.PostForm("filename")
	newName := c.PostForm("newname")

	err := validations.ValidateRenameInput(oldName, newName)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	filteredNewName, err := util.SanitizeFilename(newName)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = util.RenameFile(oldName, filteredNewName, "images")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to rename file: %s", err.Error())
		return
	}

	err = h.repo.RenameImage(oldName, filteredNewName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to rename file: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "File renamed successfully"})
}

// HandleDocsRename provides backward compatibility for document renaming
func (h *MediaHandler) HandleDocsRename(c *gin.Context) {
	oldName := c.PostForm("filename")
	newName := c.PostForm("newname")

	err := validations.ValidateRenameInput(oldName, newName)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	filteredNewName, err := util.FilterFilename(newName)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	err = util.RenameFile(oldName, filteredNewName, "docs")
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to rename file: %s", err.Error())
		return
	}

	err = h.repo.RenameDoc(oldName, newName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to rename file: %s", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "File renamed successfully"})
}
