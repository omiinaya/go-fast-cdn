package config

import (
	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/config"
)

// FileTypeHandler handles file type configuration requests
type FileTypeHandler struct {
	fileTypeConfig *config.FileTypeConfig
}

// NewFileTypeHandler creates a new file type handler
func NewFileTypeHandler() *FileTypeHandler {
	return &FileTypeHandler{
		fileTypeConfig: config.GetDefaultFileTypeConfig(),
	}
}

// GetFileTypeConfig returns the complete file type configuration
func (h *FileTypeHandler) GetFileTypeConfig(c *gin.Context) {
	c.JSON(200, h.fileTypeConfig)
}

// GetSupportedFileTypes returns all supported file extensions grouped by category
func (h *FileTypeHandler) GetSupportedFileTypes(c *gin.Context) {
	response := gin.H{
		"image":    h.fileTypeConfig.GetExtensionsByCategory(config.MediaTypeImage),
		"document": h.fileTypeConfig.GetExtensionsByCategory(config.MediaTypeDocument),
		"video":    h.fileTypeConfig.GetExtensionsByCategory(config.MediaTypeVideo),
		"audio":    h.fileTypeConfig.GetExtensionsByCategory(config.MediaTypeAudio),
	}
	c.JSON(200, response)
}

// GetSupportedMimeTypes returns all supported MIME types grouped by category
func (h *FileTypeHandler) GetSupportedMimeTypes(c *gin.Context) {
	response := gin.H{
		"image":    h.fileTypeConfig.GetMimeTypesByCategory(config.MediaTypeImage),
		"document": h.fileTypeConfig.GetMimeTypesByCategory(config.MediaTypeDocument),
		"video":    h.fileTypeConfig.GetMimeTypesByCategory(config.MediaTypeVideo),
		"audio":    h.fileTypeConfig.GetMimeTypesByCategory(config.MediaTypeAudio),
	}
	c.JSON(200, response)
}
