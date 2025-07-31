package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"github.com/kevinanielsen/go-fast-cdn/src/validations"
)

func (h *DocHandler) HandleDocsRename(c *gin.Context) {
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

	// Check if a file with the new name already exists
	docWithNewName := h.repo.GetDocByFileName(filteredNewName)
	if len(docWithNewName.Checksum) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "File with this name already exists",
		})
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
