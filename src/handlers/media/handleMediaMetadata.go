package handlers

import (
	"errors"
	"image"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

// HandleMediaMetadata handles metadata retrieval for both images and documents
func (h *MediaHandler) HandleMediaMetadata(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "File name is required",
		})
		return
	}

	mediaType := c.Query("type")
	if mediaType == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Media type is required",
		})
		return
	}

	// Get the media from the database
	media := h.repo.GetMediaByFileName(fileName)
	if len(media.Checksum) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Media not found",
		})
		return
	}

	// Verify the media type matches
	if string(media.Type) != mediaType {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Media type mismatch",
		})
		return
	}

	// Get file path
	filePath := filepath.Join(util.ExPath, "uploads", mediaType, fileName)

	// Check if file exists and get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "File does not exist",
			})
		} else {
			log.Printf("Failed to get file %s: %s\n", fileName, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	// Create base response body
	body := gin.H{
		"filename":     fileName,
		"download_url": c.Request.Host + "/api/cdn/download/" + mediaType + "/" + fileName,
		"file_size":    fileInfo.Size(),
		"type":         media.Type,
	}

	// For images, add dimensions
	if media.Type == "image" {
		if file, err := os.Open(filePath); err != nil {
			log.Printf("Failed to open the image %s: %s\n", fileName, err.Error())
		} else {
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				log.Printf("Failed to decode image %s: %s\n", fileName, err.Error())
			} else {
				body["width"] = img.Bounds().Dx()
				body["height"] = img.Bounds().Dy()
			}
		}
	}

	c.JSON(http.StatusOK, body)
}

// HandleImageMetadata provides backward compatibility for image metadata
func HandleImageMetadata(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Image name is required",
		})
		return
	}

	filePath := filepath.Join(util.ExPath, "uploads", "images", fileName)

	if fileinfo, err := os.Stat(filePath); err == nil {
		if file, err := os.Open(filePath); err != nil {
			log.Printf("Failed to open the image %s: %s\n", fileName, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			return
		} else {
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				log.Printf("Failed to decode image %s: %s\n", fileName, err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
				return
			}
			width := img.Bounds().Dx()
			height := img.Bounds().Dy()

			body := gin.H{
				"filename":     fileName,
				"download_url": c.Request.Host + "/api/cdn/download/images/" + fileName,
				"file_size":    fileinfo.Size(),
				"width":        width,
				"height":       height,
			}

			c.JSON(http.StatusOK, body)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Image does not exist",
		})
		return
	} else {
		log.Printf("Failed to get the image %s: %s\n", fileName, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		return
	}
}

// HandleDocMetadata provides backward compatibility for document metadata
func HandleDocMetadata(c *gin.Context) {
	fileName := c.Param("filename")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Doc name is required",
		})
		return
	}

	filePath := filepath.Join(util.ExPath, "uploads", "docs", fileName)
	stat, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Doc does not exist",
			})
		} else {
			log.Printf("Failed to get document %s: %s\n", fileName, err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"filename":     fileName,
		"download_url": c.Request.Host + "/api/cdn/download/docs/" + fileName,
		"file_size":    stat.Size(),
	})
}
