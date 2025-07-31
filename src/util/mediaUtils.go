package util

import (
	"errors"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// MediaType defines the type of media stored
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeDocument MediaType = "document"
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
)

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
			Extensions:  []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tiff", ".svg"},
			MimeTypes: []string{
				"image/jpeg",
				"image/jpg",
				"image/png",
				"image/gif",
				"image/webp",
				"image/bmp",
				"image/tiff",
				"image/svg+xml",
			},
		},
		MediaTypeDocument: {
			Type:        MediaTypeDocument,
			DisplayName: "Document",
			Extensions:  []string{".txt", ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".rtf", ".odt", ".ods", ".odp", ".zip", ".rar", ".7z"},
			MimeTypes: []string{
				"text/plain",
				"text/plain; charset=utf-8",
				"application/pdf",
				"application/msword",
				"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				"application/vnd.ms-excel",
				"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				"application/vnd.ms-powerpoint",
				"application/vnd.openxmlformats-officedocument.presentationml.presentation",
				"application/rtf",
				"application/vnd.oasis.opendocument.text",
				"application/vnd.oasis.opendocument.spreadsheet",
				"application/vnd.oasis.opendocument.presentation",
				"application/zip",
				"application/x-rar-compressed",
				"application/x-7z-compressed",
			},
		},
		MediaTypeVideo: {
			Type:        MediaTypeVideo,
			DisplayName: "Video",
			Extensions:  []string{".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm", ".mkv", ".m4v"},
			MimeTypes: []string{
				"video/mp4",
				"video/x-msvideo",
				"video/quicktime",
				"video/x-ms-wmv",
				"video/x-flv",
				"video/webm",
				"video/x-matroska",
				"video/x-m4v",
			},
		},
		MediaTypeAudio: {
			Type:        MediaTypeAudio,
			DisplayName: "Audio",
			Extensions:  []string{".mp3", ".wav", ".ogg", ".flac", ".aac", ".m4a", ".wma"},
			MimeTypes: []string{
				"audio/mpeg",
				"audio/wav",
				"audio/ogg",
				"audio/flac",
				"audio/aac",
				"audio/mp4",
				"audio/x-ms-wma",
			},
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

	mediaTypes := GetMediaTypes()
	for mediaType, info := range mediaTypes {
		for _, validMime := range info.MimeTypes {
			if mimeType == validMime {
				return mediaType, nil
			}
		}
	}

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

	// Try to get MIME type from the system
	mimeType := mime.TypeByExtension(ext)
	if mimeType != "" {
		return mimeType, nil
	}

	// If system doesn't recognize it, check our custom mappings
	mediaTypes := GetMediaTypes()
	for _, info := range mediaTypes {
		for i, validExt := range info.Extensions {
			if ext == validExt && i < len(info.MimeTypes) {
				return info.MimeTypes[i], nil
			}
		}
	}

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
