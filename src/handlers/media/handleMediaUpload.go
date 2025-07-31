package handlers

import (
	"crypto/md5"
	"image"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

// HandleMediaUpload handles the upload of both images and documents
func (h *MediaHandler) HandleMediaUpload(c *gin.Context) {
	// Get the file from the form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read file: %s", err.Error())
		return
	}

	// Get the custom filename if provided
	newName := c.PostForm("filename")

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to open file: %s", err.Error())
		return
	}
	defer file.Close()

	// Read the first 512 bytes to detect the content type
	fileBuffer := make([]byte, 512)
	_, err = file.Read(fileBuffer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %s", err.Error())
		return
	}

	// Detect the file type
	fileType := http.DetectContentType(fileBuffer)

	// Determine media type and validate
	utilMediaType, err := util.GetMediaTypeFromMIME(fileType)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid file type: %s", err.Error())
		return
	}

	// Convert util.MediaType to models.MediaType
	var mediaType models.MediaType
	switch utilMediaType {
	case util.MediaTypeImage:
		mediaType = models.MediaTypeImage
	case util.MediaTypeDocument:
		mediaType = models.MediaTypeDocument
	case util.MediaTypeVideo:
		mediaType = models.MediaType("video") // Add to models if needed
	case util.MediaTypeAudio:
		mediaType = models.MediaType("audio") // Add to models if needed
	default:
		mediaType = models.MediaTypeDocument // Default fallback
	}

	// Calculate file hash
	fileHashBuffer := md5.Sum(fileBuffer)

	// Determine the filename
	var filename string
	if newName == "" {
		filename = fileHeader.Filename
	} else {
		filename = newName + filepath.Ext(fileHeader.Filename)
	}

	// Use the more comprehensive SanitizeFilename function
	filteredFilename, err := util.SanitizeFilename(filename)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Create media object
	media := models.Media{
		FileName: filteredFilename,
		Checksum: fileHashBuffer[:],
		Type:     mediaType,
	}

	// For images, try to extract dimensions
	if mediaType == models.MediaTypeImage {
		width, height, err := h.getImageDimensions(fileHeader)
		if err == nil {
			media.Width = &width
			media.Height = &height
		}
	}

	// Check if file already exists
	mediaInDatabase := h.repo.GetMediaByCheckSum(fileHashBuffer[:])
	if len(mediaInDatabase.Checksum) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "File already exists",
		})
		return
	}

	// Save to database
	savedFilename, err := h.repo.AddMedia(media)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Use the unified media directory for all file types
	uploadPath := util.GetMediaUploadPath()

	// Ensure the upload directory exists
	if err := util.EnsureUploadDirectories(); err != nil {
		c.String(http.StatusInternalServerError, "Failed to create upload directory: %s", err.Error())
		return
	}

	// Save the file to the unified media directory
	filePath := filepath.Join(uploadPath, savedFilename)
	err = c.SaveUploadedFile(fileHeader, filePath)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to save file: %s", err.Error())
		return
	}

	// Return success response with unified URL path
	body := gin.H{
		"file_url": c.Request.Host + util.GetMediaURLPath(savedFilename),
		"type":     mediaType,
	}

	c.JSON(http.StatusOK, body)
}

// getImageDimensions extracts the dimensions of an image file
func (h *MediaHandler) getImageDimensions(fileHeader *multipart.FileHeader) (int, int, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return 0, 0, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return 0, 0, err
	}

	return img.Bounds().Dx(), img.Bounds().Dy(), nil
}

// HandleImageUpload provides backward compatibility for image uploads
func (h *MediaHandler) HandleImageUpload(c *gin.Context) {
	// Get the file from the form using the "image" field name
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read file: %s", err.Error())
		return
	}

	// Get the custom filename if provided
	newName := c.PostForm("filename")

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to open file: %s", err.Error())
		return
	}
	defer file.Close()

	// Read the first 512 bytes to detect the content type
	fileBuffer := make([]byte, 512)
	_, err = file.Read(fileBuffer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %s", err.Error())
		return
	}

	// Detect the file type
	fileType := http.DetectContentType(fileBuffer)

	// Validate that it's an image
	imageMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
		"image/bmp":  true,
	}

	if !imageMimeTypes[fileType] {
		c.String(http.StatusBadRequest, "Invalid file type")
		return
	}

	// Calculate file hash
	fileHashBuffer := md5.Sum(fileBuffer)

	// Determine the filename
	var filename string
	if newName == "" {
		filename = fileHeader.Filename
	} else {
		filename = newName + filepath.Ext(fileHeader.Filename)
	}

	// Use the more comprehensive SanitizeFilename function
	filteredFilename, err := util.SanitizeFilename(filename)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Create media object
	media := models.Media{
		FileName: filteredFilename,
		Checksum: fileHashBuffer[:],
		Type:     models.MediaTypeImage,
	}

	// For images, try to extract dimensions
	width, height, err := h.getImageDimensions(fileHeader)
	if err == nil {
		media.Width = &width
		media.Height = &height
	}

	// Check if file already exists
	mediaInDatabase := h.repo.GetImageByCheckSum(fileHashBuffer[:])
	if len(mediaInDatabase.Checksum) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "File already exists",
		})
		return
	}

	// Save to database using the backward compatibility method
	savedFilename, err := h.repo.AddImage(media)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Save the file to the legacy images directory for backward compatibility
	err = c.SaveUploadedFile(fileHeader, util.GetImagesPath()+"/"+savedFilename)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to save file: %s", err.Error())
		return
	}

	// Return success response
	body := gin.H{
		"file_url": c.Request.Host + "/download/images/" + savedFilename,
	}

	c.JSON(http.StatusOK, body)
}

// HandleDocUpload provides backward compatibility for document uploads
func (h *MediaHandler) HandleDocUpload(c *gin.Context) {
	// Get the file from the form using the "doc" field name
	fileHeader, err := c.FormFile("doc")
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read file: %s", err.Error())
		return
	}

	// Get the custom filename if provided
	newName := c.PostForm("filename")

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to open file: %s", err.Error())
		return
	}
	defer file.Close()

	// Read the first 512 bytes to detect the content type
	fileBuffer := make([]byte, 512)
	_, err = file.Read(fileBuffer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %s", err.Error())
		return
	}

	// Detect the file type
	fileType := http.DetectContentType(fileBuffer)

	// Validate that it's a document
	docMimeTypes := map[string]bool{
		"text/plain":                true,
		"text/plain; charset=utf-8": true,
		"application/msword":        true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         true,
		"application/pdf":       true,
		"application/rtf":       true,
		"application/x-freearc": true,
		"application/zip":       true,
	}

	if !docMimeTypes[fileType] {
		c.String(http.StatusBadRequest, "Invalid file type: %s", fileType)
		return
	}

	// Calculate file hash
	fileHashBuffer := md5.Sum(fileBuffer)

	// Determine the filename
	var filename string
	if newName == "" {
		filename = fileHeader.Filename
	} else {
		filename = newName + filepath.Ext(fileHeader.Filename)
	}

	// Use the more comprehensive SanitizeFilename function
	filteredFilename, err := util.SanitizeFilename(filename)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	// Create media object
	media := models.Media{
		FileName: filteredFilename,
		Checksum: fileHashBuffer[:],
		Type:     models.MediaTypeDocument,
	}

	// Check if file already exists
	mediaInDatabase := h.repo.GetDocByCheckSum(fileHashBuffer[:])
	if len(mediaInDatabase.Checksum) > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "File already exists"})
		return
	}

	// Save to database using the backward compatibility method
	savedFileName, err := h.repo.AddDoc(media)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Save the file to the legacy docs directory for backward compatibility
	err = c.SaveUploadedFile(fileHeader, util.GetDocsPath()+"/"+savedFileName)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to save file: %s", err.Error())
		return
	}

	// Return success response
	body := gin.H{
		"file_url": c.Request.Host + "/download/docs/" + savedFileName,
	}

	c.JSON(http.StatusOK, body)
}
