package handlers

import (
	"crypto/md5"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// HandleMediaUpload handles the upload of both images and documents
func (h *MediaHandler) HandleMediaUpload(c *gin.Context) {
	fmt.Printf("[DEBUG] HandleMediaUpload called\n")

	// Get the file from the form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		fmt.Printf("[DEBUG] Failed to read file: %s\n", err.Error())
		c.String(http.StatusBadRequest, "Failed to read file: %s", err.Error())
		return
	}
	fmt.Printf("[DEBUG] File header received: Name=%s, Size=%d, Header=%+v\n", fileHeader.Filename, fileHeader.Size, fileHeader.Header)

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
	bytesRead, err := file.Read(fileBuffer)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read file: %s", err.Error())
		return
	}

	// Log detailed debugging information
	fmt.Printf("File: %s, Size: %d, Bytes read: %d\n", fileHeader.Filename, fileHeader.Size, bytesRead)
	fmt.Printf("File buffer content (first 100 bytes): %x\n", fileBuffer[:min(100, bytesRead)])
	fmt.Printf("File buffer content (as string): %s\n", string(fileBuffer[:min(100, bytesRead)]))

	// Detect the file type
	fileType := http.DetectContentType(fileBuffer)
	// Log the detected file type for debugging
	fmt.Printf("Detected file type: %s\n", fileType)

	// Also try to detect from file extension
	ext := filepath.Ext(fileHeader.Filename)
	fmt.Printf("File extension: %s\n", ext)

	// Enhanced fallback mechanism: if detected as application/octet-stream or unknown, try to determine from extension
	if fileType == "application/octet-stream" || !util.IsSupportedMIMEType(fileType) {
		fmt.Printf("File detected as %s, trying fallback detection from extension: %s\n", fileType, ext)

		// Try to get MIME type from extension using our utility function
		if mimeTypeFromExt, err := util.GetMIMETypeFromExtension(ext); err == nil {
			fileType = mimeTypeFromExt
			fmt.Printf("Fallback: setting file type to %s based on %s extension\n", fileType, ext)
		} else {
			fmt.Printf("Warning: Could not determine MIME type from extension %s: %v\n", ext, err)
		}
		fmt.Printf("Final file type after fallback: %s\n", fileType)
	}

	// Reset file position to beginning for hash calculation
	if _, err := file.Seek(0, 0); err != nil {
		fmt.Printf("Warning: Failed to reset file position: %v\n", err)
	}

	// Determine media type and validate
	utilMediaType, err := util.GetMediaTypeFromMIME(fileType)
	if err != nil {
		fmt.Printf("Failed to get media type from MIME '%s': %v\n", fileType, err)
		c.String(http.StatusBadRequest, "Invalid file type: %s", err.Error())
		return
	}
	fmt.Printf("Successfully determined media type: %s\n", utilMediaType)

	// Convert util.MediaType to models.MediaType
	var mediaType models.MediaType
	switch utilMediaType {
	case util.MediaTypeImage:
		mediaType = models.MediaTypeImage
	case util.MediaTypeDocument:
		mediaType = models.MediaTypeDocument
	case util.MediaTypeVideo:
		mediaType = models.MediaTypeVideo
	case util.MediaTypeAudio:
		mediaType = models.MediaTypeAudio
	default:
		mediaType = models.MediaTypeDocument // Default fallback
	}

	// Calculate file hash from the entire file content
	// Reset file position and read the entire file
	if _, err := file.Seek(0, 0); err != nil {
		c.String(http.StatusInternalServerError, "Failed to reset file position: %s", err.Error())
		return
	}

	// Read the entire file for hash calculation
	fullFileContent, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read full file content: %s", err.Error())
		return
	}

	fileHashBuffer := md5.Sum(fullFileContent)
	fmt.Printf("Calculated file hash from full content: %x\n", fileHashBuffer)
	fmt.Printf("Full file size: %d bytes\n", len(fullFileContent))

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

	// For images, try to extract dimensions (non-critical operation)
	if mediaType == models.MediaTypeImage {
		fmt.Printf("[DEBUG] Attempting to extract image dimensions\n")
		width, height, err := h.getImageDimensions(fileHeader)
		if err == nil {
			media.Width = &width
			media.Height = &height
			fmt.Printf("[DEBUG] Image dimensions extracted: %dx%d\n", width, height)
		} else {
			fmt.Printf("[DEBUG] Failed to extract image dimensions (non-critical): %v\n", err)
			// Don't fail the upload if dimension extraction fails
		}
	}

	// Check if file already exists by checksum (content-based duplicate detection)
	fmt.Printf("Checking if file exists in database with hash: %x\n", fileHashBuffer)
	mediaInDatabase := h.repo.GetMediaByCheckSum(fileHashBuffer[:])
	fmt.Printf("Database query result - Checksum length: %d, FileName: %s\n", len(mediaInDatabase.Checksum), mediaInDatabase.FileName)
	if len(mediaInDatabase.Checksum) > 0 {
		fmt.Printf("File with same content already exists in database: %s\n", mediaInDatabase.FileName)
		c.JSON(http.StatusConflict, gin.H{
			"error":         "File with this content already exists",
			"existing_file": mediaInDatabase.FileName,
		})
		return
	}
	fmt.Printf("File content not found in database, proceeding with upload\n")

	// Check if file already exists by filename (name-based duplicate detection)
	fmt.Printf("Checking if file with name '%s' already exists in database\n", filteredFilename)
	mediaWithSameName := h.repo.GetMediaByFileName(filteredFilename)
	if len(mediaWithSameName.Checksum) > 0 {
		fmt.Printf("File with same name already exists in database: %s\n", mediaWithSameName.FileName)
		c.JSON(http.StatusConflict, gin.H{
			"error":         "File with this name already exists",
			"existing_file": mediaWithSameName.FileName,
		})
		return
	}
	fmt.Printf("File name not found in database, proceeding with upload\n")

	// Save to database
	fmt.Printf("[DEBUG] Saving media to database: %+v\n", media)
	savedFilename, err := h.repo.AddMedia(media)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to save media to database: %v\n", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	fmt.Printf("[DEBUG] Successfully saved media to database with filename: %s\n", savedFilename)

	// Use the unified media directory for all file types
	uploadPath := util.GetMediaUploadPath()

	// Ensure the upload directory exists
	if err := util.EnsureUploadDirectories(); err != nil {
		c.String(http.StatusInternalServerError, "Failed to create upload directory: %s", err.Error())
		return
	}

	// Save the file to the unified media directory
	filePath := filepath.Join(uploadPath, savedFilename)
	fmt.Printf("[DEBUG] Saving file to path: %s\n", filePath)
	err = c.SaveUploadedFile(fileHeader, filePath)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to save file: %s\n", err.Error())
		c.String(http.StatusInternalServerError, "Failed to save file: %s", err.Error())
		return
	}
	fmt.Printf("[DEBUG] File saved successfully to: %s\n", filePath)

	// Return success response with unified URL path
	body := gin.H{
		"file_url": c.Request.Host + util.GetMediaURLPath(savedFilename),
		"type":     mediaType,
		"filename": savedFilename,
		"checksum": fmt.Sprintf("%x", fileHashBuffer),
		"note":     "Duplicate detection is based on file content (checksum), not filename",
	}

	c.JSON(http.StatusOK, body)
}

// getImageDimensions extracts the dimensions of an image file
func (h *MediaHandler) getImageDimensions(fileHeader *multipart.FileHeader) (int, int, error) {
	fmt.Printf("[DEBUG] getImageDimensions called for file: %s\n", fileHeader.Filename)

	file, err := fileHeader.Open()
	if err != nil {
		fmt.Printf("[DEBUG] Failed to open file for dimension extraction: %v\n", err)
		return 0, 0, err
	}
	defer file.Close()

	fmt.Printf("[DEBUG] Attempting to decode image\n")
	img, format, err := image.Decode(file)
	if err != nil {
		fmt.Printf("[DEBUG] Failed to decode image: %v\n", err)
		return 0, 0, err
	}

	fmt.Printf("[DEBUG] Successfully decoded image format: %s, dimensions: %dx%d\n", format, img.Bounds().Dx(), img.Bounds().Dy())
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

	// Calculate file hash from the entire file content
	// Reset file position and read the entire file
	if _, err := file.Seek(0, 0); err != nil {
		c.String(http.StatusInternalServerError, "Failed to reset file position: %s", err.Error())
		return
	}

	// Read the entire file for hash calculation
	fullFileContent, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read full file content: %s", err.Error())
		return
	}

	fileHashBuffer := md5.Sum(fullFileContent)
	fmt.Printf("Calculated file hash from full content: %x\n", fileHashBuffer)
	fmt.Printf("Full file size: %d bytes\n", len(fullFileContent))

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

	// Check if file already exists by checksum (content-based duplicate detection)
	mediaInDatabase := h.repo.GetImageByCheckSum(fileHashBuffer[:])
	if len(mediaInDatabase.Checksum) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":         "File with this content already exists",
			"existing_file": mediaInDatabase.FileName,
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
		"type":     models.MediaTypeImage,
		"filename": savedFilename,
		"checksum": fmt.Sprintf("%x", fileHashBuffer),
		"note":     "Duplicate detection is based on file content (checksum), not filename",
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

	// Calculate file hash from the entire file content
	// Reset file position and read the entire file
	if _, err := file.Seek(0, 0); err != nil {
		c.String(http.StatusInternalServerError, "Failed to reset file position: %s", err.Error())
		return
	}

	// Read the entire file for hash calculation
	fullFileContent, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read full file content: %s", err.Error())
		return
	}

	fileHashBuffer := md5.Sum(fullFileContent)
	fmt.Printf("Calculated file hash from full content: %x\n", fileHashBuffer)
	fmt.Printf("Full file size: %d bytes\n", len(fullFileContent))

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

	// Check if file already exists by filename (name-based duplicate detection)
	docWithSameName := h.repo.GetDocByFileName(filteredFilename)
	if len(docWithSameName.Checksum) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":         "File with this name already exists",
			"existing_file": docWithSameName.FileName,
		})
		return
	}

	// Check if file already exists by checksum (content-based duplicate detection)
	mediaInDatabase := h.repo.GetDocByCheckSum(fileHashBuffer[:])
	if len(mediaInDatabase.Checksum) > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":         "File with this content already exists",
			"existing_file": mediaInDatabase.FileName,
		})
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
		"type":     models.MediaTypeDocument,
		"filename": savedFileName,
		"checksum": fmt.Sprintf("%x", fileHashBuffer),
		"note":     "Duplicate detection is based on file content (checksum), not filename",
	}

	c.JSON(http.StatusOK, body)
}
