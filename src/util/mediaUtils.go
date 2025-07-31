package util

import (
	"errors"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/kevinanielsen/go-fast-cdn/src/config"
)

// MediaType defines the type of media stored
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeDocument MediaType = "document"
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
)

// Global file type configuration instance
var fileTypeConfig = config.GetDefaultFileTypeConfig()

// MediaTypeInfo contains information about a media type
type MediaTypeInfo struct {
	Type        MediaType
	DisplayName string
	Extensions  []string
	MimeTypes   []string
}

// GetMediaTypes returns a map of all supported media types with their information
func GetMediaTypes() map[MediaType]MediaTypeInfo {
	return map[MediaType]MediaTypeInfo{
		MediaTypeImage: {
			Type:        MediaTypeImage,
			DisplayName: "Image",
			Extensions:  fileTypeConfig.GetExtensionsByCategory(config.MediaTypeImage),
			MimeTypes:   fileTypeConfig.GetMimeTypesByCategory(config.MediaTypeImage),
		},
		MediaTypeDocument: {
			Type:        MediaTypeDocument,
			DisplayName: "Document",
			Extensions:  fileTypeConfig.GetExtensionsByCategory(config.MediaTypeDocument),
			MimeTypes:   fileTypeConfig.GetMimeTypesByCategory(config.MediaTypeDocument),
		},
		MediaTypeVideo: {
			Type:        MediaTypeVideo,
			DisplayName: "Video",
			Extensions:  fileTypeConfig.GetExtensionsByCategory(config.MediaTypeVideo),
			MimeTypes:   fileTypeConfig.GetMimeTypesByCategory(config.MediaTypeVideo),
		},
		MediaTypeAudio: {
			Type:        MediaTypeAudio,
			DisplayName: "Audio",
			Extensions:  fileTypeConfig.GetExtensionsByCategory(config.MediaTypeAudio),
			MimeTypes:   fileTypeConfig.GetMimeTypesByCategory(config.MediaTypeAudio),
		},
	}
}

// GetMediaTypeFromExtension determines the media type from a file extension
func GetMediaTypeFromExtension(extension string) (MediaType, error) {
	if extension == "" {
		return "", errors.New("empty extension")
	}

	// Ensure extension starts with a dot and is lowercase
	ext := strings.ToLower(extension)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	mediaTypes := GetMediaTypes()
	for mediaType, info := range mediaTypes {
		for _, validExt := range info.Extensions {
			if ext == validExt {
				return mediaType, nil
			}
		}
	}

	return "", errors.New("unsupported file extension: " + extension)
}

// GetMediaTypeFromMIME determines the media type from a MIME type
func GetMediaTypeFromMIME(mimeType string) (MediaType, error) {
	if mimeType == "" {
		return "", errors.New("empty MIME type")
	}

	fmt.Printf("[DEBUG] GetMediaTypeFromMIME called with MIME type: %s\n", mimeType)

	mediaTypes := GetMediaTypes()
	for mediaType, info := range mediaTypes {
		for _, validMime := range info.MimeTypes {
			if mimeType == validMime {
				fmt.Printf("[DEBUG] Found matching media type: %s for MIME: %s\n", mediaType, mimeType)
				return mediaType, nil
			}
		}
	}

	fmt.Printf("[DEBUG] No matching media type found for MIME: %s\n", mimeType)
	return "", errors.New("unsupported MIME type: " + mimeType)
}

// GetMediaTypeFromFilename determines the media type from a filename
func GetMediaTypeFromFilename(filename string) (MediaType, error) {
	if filename == "" {
		return "", errors.New("empty filename")
	}

	extension := filepath.Ext(filename)
	return GetMediaTypeFromExtension(extension)
}

// DetectMediaType detects the media type from file content (first 512 bytes)
func DetectMediaType(fileBuffer []byte) (MediaType, error) {
	if len(fileBuffer) == 0 {
		return "", errors.New("empty file buffer")
	}

	mimeType := http.DetectContentType(fileBuffer)
	return GetMediaTypeFromMIME(mimeType)
}

// IsSupportedExtension checks if a file extension is supported
func IsSupportedExtension(extension string) bool {
	_, err := GetMediaTypeFromExtension(extension)
	return err == nil
}

// IsSupportedMIMEType checks if a MIME type is supported
func IsSupportedMIMEType(mimeType string) bool {
	_, err := GetMediaTypeFromMIME(mimeType)
	return err == nil
}

// GetMediaUploadPath returns the upload path for a media file based on the unified approach
// Deprecated: Use GetMediaPath from getPath.go instead
func GetMediaUploadPath() string {
	return filepath.Join(ExPath, "uploads", "media")
}

// GetLegacyUploadPath returns the upload path for a media file based on the legacy approach
func GetLegacyUploadPath(mediaType MediaType) string {
	return filepath.Join(ExPath, "uploads", string(mediaType))
}

// GetMediaURLPath returns the URL path for accessing a media file
func GetMediaURLPath(filename string) string {
	return "/uploads/media/" + filename
}

// GetLegacyURLPath returns the URL path for accessing a media file using the legacy approach
func GetLegacyURLPath(filename string, mediaType MediaType) string {
	return "/uploads/" + string(mediaType) + "/" + filename
}

// GetMIMETypeFromExtension returns the MIME type for a given file extension
func GetMIMETypeFromExtension(extension string) (string, error) {
	if extension == "" {
		return "", errors.New("empty extension")
	}

	// Ensure extension starts with a dot and is lowercase
	ext := strings.ToLower(extension)
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	fmt.Printf("[DEBUG] GetMIMETypeFromExtension called with extension: %s (normalized: %s)\n", extension, ext)

	// First, check our custom mappings for more accurate MIME type detection
	fileInfo, err := fileTypeConfig.GetFileTypeInfo(ext)
	if err == nil {
		// Return the first (primary) MIME type for this extension
		if len(fileInfo.MimeTypes) > 0 {
			mimeType := fileInfo.MimeTypes[0]
			fmt.Printf("[DEBUG] Custom mapping found: %s -> %s\n", ext, mimeType)
			return mimeType, nil
		}
	}

	// If no custom mapping found, try to get MIME type from the system
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		fmt.Printf("[DEBUG] System detected MIME type: %s for extension: %s\n", mimeType, ext)
		return mimeType, nil
	}

	fmt.Printf("[DEBUG] No MIME type found for extension: %s\n", ext)
	return "", errors.New("unsupported file extension: " + extension)
}

// ValidateMediaFile validates a media file based on its extension and MIME type
func ValidateMediaFile(filename string, mimeType string) (MediaType, error) {
	// Get media type from filename
	mediaTypeFromExt, err := GetMediaTypeFromFilename(filename)
	if err != nil {
		return "", err
	}

	// Get media type from MIME type
	mediaTypeFromMIME, err := GetMediaTypeFromMIME(mimeType)
	if err != nil {
		return "", err
	}

	// Ensure both methods return the same media type
	if mediaTypeFromExt != mediaTypeFromMIME {
		return "", errors.New("file extension and MIME type do not match")
	}

	return mediaTypeFromExt, nil
}

// GetMediaTypeInfo returns information about a specific media type
func GetMediaTypeInfo(mediaType MediaType) (MediaTypeInfo, error) {
	mediaTypes := GetMediaTypes()
	info, exists := mediaTypes[mediaType]
	if !exists {
		return MediaTypeInfo{}, errors.New("unsupported media type: " + string(mediaType))
	}
	return info, nil
}
