package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

// HandleMediaResize handles the resizing of images with type-specific functionality
func (h *MediaHandler) HandleMediaResize(c *gin.Context) {
	body := struct {
		Filename string `json:"filename" binding:"required"`
		Width    int    `json:"width" binding:"required"`
		Height   int    `json:"height" binding:"required"`
	}{}

	if e := c.BindJSON(&body); e != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": e.Error(),
		})
		return
	}

	// Validate dimensions
	if body.Width <= 0 || body.Height <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Width and height must be positive values",
		})
		return
	}

	filename, err := util.FilterFilename(body.Filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get the media from the database to verify it exists and get its type
	media := h.repo.GetMediaByFileName(filename)
	if len(media.Checksum) == 0 {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": "Media not found",
		})
		return
	}

	// Verify it's an image - this is the type-specific check
	if media.Type != models.MediaTypeImage {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Cannot resize media of type '%s'. Only images can be resized", media.Type),
		})
		return
	}

	// Extract file extension
	fileExt := filepath.Ext(filename)
	if len(fileExt) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file extension",
		})
		return
	}

	// Remove the dot from the extension
	imgType := fileExt[1:]
	if len(imgType) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid file extension",
		})
		return
	}

	// Construct the file path using the media type
	filepath := filepath.Join(util.ExPath, "uploads", string(media.Type), filename)

	// Open the image file
	img, err := imgio.Open(filepath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to open image file: %s", err.Error()),
		})
		return
	}

	// Resize the image
	resizedImg := transform.Resize(img, body.Width, body.Height, transform.Linear)

	// Determine the encoder based on image type
	encoder, err := h.getImageEncoder(imgType)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Save the resized image
	if err := imgio.Save(filepath, resizedImg, encoder); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Failed to save resized image: %s", err.Error()),
		})
		return
	}

	// Update the image dimensions in the database
	if err := h.repo.UpdateImageDimensions(filename, body.Width, body.Height); err != nil {
		// Log the error but don't fail the operation since the image was successfully resized
		// In a production environment, you might want to use a proper logging package
		fmt.Printf("Warning: Failed to update image dimensions in database: %s\n", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "File resized successfully",
		"width":   body.Width,
		"height":  body.Height,
		"type":    media.Type,
		"message": fmt.Sprintf("Image '%s' has been resized to %dx%d pixels", filename, body.Width, body.Height),
	})
}

// getImageEncoder returns the appropriate image encoder based on the file type
func (h *MediaHandler) getImageEncoder(imgType string) (imgio.Encoder, error) {
	switch strings.ToLower(imgType) {
	case "png":
		return imgio.PNGEncoder(), nil
	case "jpg", "jpeg":
		// 75 is the default quality encoding parameter
		return imgio.JPEGEncoder(75), nil
	case "bmp":
		return imgio.BMPEncoder(), nil
	case "gif":
		// Note: bild doesn't have a GIF encoder, so we'll return an error
		return nil, fmt.Errorf("GIF encoding is not supported for resizing")
	default:
		return nil, fmt.Errorf("Image of type '%s' is not supported for resizing", imgType)
	}
}

// HandleImageResize provides backward compatibility for image resizing
func HandleImageResize(c *gin.Context) {
	body := struct {
		Filename string `json:"filename" binding:"required"`
		Width    int    `json:"width" binding:"required"`
		Height   int    `json:"height" binding:"required"`
	}{}
	if e := c.BindJSON(&body); e != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": e.Error(),
		})
		return
	}

	filename, err := util.FilterFilename(body.Filename)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	imgType := strings.Split(filename, ".")[1]

	filepath := filepath.Join(util.ExPath, "uploads", "images", filename)

	img, err := imgio.Open(filepath)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	img = transform.Resize(img, body.Width, body.Height, transform.Linear)

	// TODO: a shared accepted image type data could be added to be shared between upload and resize api
	var encoder imgio.Encoder
	switch imgType {
	case "png":
		encoder = imgio.PNGEncoder()
	case "jpg", "jpeg":
		// 75 is the default quality encoding parameter
		encoder = imgio.JPEGEncoder(75)
	case "bmp":
		encoder = imgio.BMPEncoder()
	default:
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Image of type %s is not supported", imgType),
		})
		return
	}

	if err := imgio.Save(filepath, img, encoder); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "File resized successfully",
	})
}
