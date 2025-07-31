package config

import (
	"encoding/json"
	"fmt"
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

// FileTypeInfo contains information about a file type
type FileTypeInfo struct {
	Extension string    `json:"extension"`
	MimeTypes []string  `json:"mimeTypes"`
	Category  MediaType `json:"category"`
}

// FileTypeConfig contains the complete file type configuration
type FileTypeConfig struct {
	FileTypes map[string]FileTypeInfo `json:"fileTypes"`
}

// GetDefaultFileTypeConfig returns the default file type configuration
func GetDefaultFileTypeConfig() *FileTypeConfig {
	return &FileTypeConfig{
		FileTypes: map[string]FileTypeInfo{
			// Images
			".jpg": {
				Extension: ".jpg",
				MimeTypes: []string{"image/jpeg", "image/jpg"},
				Category:  MediaTypeImage,
			},
			".jpeg": {
				Extension: ".jpeg",
				MimeTypes: []string{"image/jpeg", "image/jpg"},
				Category:  MediaTypeImage,
			},
			".png": {
				Extension: ".png",
				MimeTypes: []string{"image/png"},
				Category:  MediaTypeImage,
			},
			".gif": {
				Extension: ".gif",
				MimeTypes: []string{"image/gif"},
				Category:  MediaTypeImage,
			},
			".webp": {
				Extension: ".webp",
				MimeTypes: []string{"image/webp"},
				Category:  MediaTypeImage,
			},
			".bmp": {
				Extension: ".bmp",
				MimeTypes: []string{"image/bmp"},
				Category:  MediaTypeImage,
			},
			".tiff": {
				Extension: ".tiff",
				MimeTypes: []string{"image/tiff"},
				Category:  MediaTypeImage,
			},
			".svg": {
				Extension: ".svg",
				MimeTypes: []string{"image/svg+xml"},
				Category:  MediaTypeImage,
			},

			// Documents
			".txt": {
				Extension: ".txt",
				MimeTypes: []string{"text/plain", "text/plain; charset=utf-8"},
				Category:  MediaTypeDocument,
			},
			".pdf": {
				Extension: ".pdf",
				MimeTypes: []string{"application/pdf"},
				Category:  MediaTypeDocument,
			},
			".doc": {
				Extension: ".doc",
				MimeTypes: []string{"application/msword"},
				Category:  MediaTypeDocument,
			},
			".docx": {
				Extension: ".docx",
				MimeTypes: []string{"application/vnd.openxmlformats-officedocument.wordprocessingml.document"},
				Category:  MediaTypeDocument,
			},
			".xls": {
				Extension: ".xls",
				MimeTypes: []string{"application/vnd.ms-excel"},
				Category:  MediaTypeDocument,
			},
			".xlsx": {
				Extension: ".xlsx",
				MimeTypes: []string{"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"},
				Category:  MediaTypeDocument,
			},
			".ppt": {
				Extension: ".ppt",
				MimeTypes: []string{"application/vnd.ms-powerpoint"},
				Category:  MediaTypeDocument,
			},
			".pptx": {
				Extension: ".pptx",
				MimeTypes: []string{"application/vnd.openxmlformats-officedocument.presentationml.presentation"},
				Category:  MediaTypeDocument,
			},
			".rtf": {
				Extension: ".rtf",
				MimeTypes: []string{"application/rtf"},
				Category:  MediaTypeDocument,
			},
			".odt": {
				Extension: ".odt",
				MimeTypes: []string{"application/vnd.oasis.opendocument.text"},
				Category:  MediaTypeDocument,
			},
			".ods": {
				Extension: ".ods",
				MimeTypes: []string{"application/vnd.oasis.opendocument.spreadsheet"},
				Category:  MediaTypeDocument,
			},
			".odp": {
				Extension: ".odp",
				MimeTypes: []string{"application/vnd.oasis.opendocument.presentation"},
				Category:  MediaTypeDocument,
			},
			".csv": {
				Extension: ".csv",
				MimeTypes: []string{"text/csv"},
				Category:  MediaTypeDocument,
			},
			".json": {
				Extension: ".json",
				MimeTypes: []string{"application/json"},
				Category:  MediaTypeDocument,
			},
			".xml": {
				Extension: ".xml",
				MimeTypes: []string{"application/xml", "text/xml", "application/xml; charset=utf-8"},
				Category:  MediaTypeDocument,
			},
			".html": {
				Extension: ".html",
				MimeTypes: []string{"text/html", "text/html; charset=utf-8"},
				Category:  MediaTypeDocument,
			},
			".htm": {
				Extension: ".htm",
				MimeTypes: []string{"text/html", "text/html; charset=utf-8"},
				Category:  MediaTypeDocument,
			},
			".css": {
				Extension: ".css",
				MimeTypes: []string{"text/css", "text/css; charset=utf-8"},
				Category:  MediaTypeDocument,
			},
			".js": {
				Extension: ".js",
				MimeTypes: []string{"application/javascript", "text/javascript", "application/javascript; charset=utf-8", "text/javascript; charset=utf-8"},
				Category:  MediaTypeDocument,
			},
			".md": {
				Extension: ".md",
				MimeTypes: []string{"text/markdown"},
				Category:  MediaTypeDocument,
			},
			".log": {
				Extension: ".log",
				MimeTypes: []string{"text/plain"},
				Category:  MediaTypeDocument,
			},
			".ini": {
				Extension: ".ini",
				MimeTypes: []string{"text/plain"},
				Category:  MediaTypeDocument,
			},
			".yaml": {
				Extension: ".yaml",
				MimeTypes: []string{"application/x-yaml", "text/yaml", "application/yaml", "text/yaml; charset=utf-8"},
				Category:  MediaTypeDocument,
			},
			".yml": {
				Extension: ".yml",
				MimeTypes: []string{"application/x-yaml", "text/yaml", "application/yaml", "text/yaml; charset=utf-8"},
				Category:  MediaTypeDocument,
			},
			".toml": {
				Extension: ".toml",
				MimeTypes: []string{"application/toml"},
				Category:  MediaTypeDocument,
			},
			".conf": {
				Extension: ".conf",
				MimeTypes: []string{"text/plain"},
				Category:  MediaTypeDocument,
			},
			".config": {
				Extension: ".config",
				MimeTypes: []string{"text/plain"},
				Category:  MediaTypeDocument,
			},
			".zip": {
				Extension: ".zip",
				MimeTypes: []string{"application/zip"},
				Category:  MediaTypeDocument,
			},
			".rar": {
				Extension: ".rar",
				MimeTypes: []string{"application/x-rar-compressed"},
				Category:  MediaTypeDocument,
			},
			".7z": {
				Extension: ".7z",
				MimeTypes: []string{"application/x-7z-compressed"},
				Category:  MediaTypeDocument,
			},
			".tar": {
				Extension: ".tar",
				MimeTypes: []string{"application/x-tar"},
				Category:  MediaTypeDocument,
			},
			".gz": {
				Extension: ".gz",
				MimeTypes: []string{"application/gzip", "application/x-gzip"},
				Category:  MediaTypeDocument,
			},
			".bz2": {
				Extension: ".bz2",
				MimeTypes: []string{"application/x-bzip2"},
				Category:  MediaTypeDocument,
			},
			".xz": {
				Extension: ".xz",
				MimeTypes: []string{"application/x-xz"},
				Category:  MediaTypeDocument,
			},

			// Videos
			".mp4": {
				Extension: ".mp4",
				MimeTypes: []string{"video/mp4"},
				Category:  MediaTypeVideo,
			},
			".avi": {
				Extension: ".avi",
				MimeTypes: []string{"video/x-msvideo"},
				Category:  MediaTypeVideo,
			},
			".mov": {
				Extension: ".mov",
				MimeTypes: []string{"video/quicktime"},
				Category:  MediaTypeVideo,
			},
			".wmv": {
				Extension: ".wmv",
				MimeTypes: []string{"video/x-ms-wmv"},
				Category:  MediaTypeVideo,
			},
			".flv": {
				Extension: ".flv",
				MimeTypes: []string{"video/x-flv"},
				Category:  MediaTypeVideo,
			},
			".webm": {
				Extension: ".webm",
				MimeTypes: []string{"video/webm"},
				Category:  MediaTypeVideo,
			},
			".mkv": {
				Extension: ".mkv",
				MimeTypes: []string{"video/x-matroska"},
				Category:  MediaTypeVideo,
			},
			".m4v": {
				Extension: ".m4v",
				MimeTypes: []string{"video/x-m4v"},
				Category:  MediaTypeVideo,
			},
			".mpg": {
				Extension: ".mpg",
				MimeTypes: []string{"video/mpeg"},
				Category:  MediaTypeVideo,
			},
			".mpeg": {
				Extension: ".mpeg",
				MimeTypes: []string{"video/mpeg"},
				Category:  MediaTypeVideo,
			},
			".m2v": {
				Extension: ".m2v",
				MimeTypes: []string{"video/x-mpeg2"},
				Category:  MediaTypeVideo,
			},
			".3gp": {
				Extension: ".3gp",
				MimeTypes: []string{"video/3gpp"},
				Category:  MediaTypeVideo,
			},
			".ogv": {
				Extension: ".ogv",
				MimeTypes: []string{"video/ogg"},
				Category:  MediaTypeVideo,
			},
			".ts": {
				Extension: ".ts",
				MimeTypes: []string{"video/mp2t"},
				Category:  MediaTypeVideo,
			},
			".mts": {
				Extension: ".mts",
				MimeTypes: []string{"video/MP2T"},
				Category:  MediaTypeVideo,
			},
			".m2ts": {
				Extension: ".m2ts",
				MimeTypes: []string{"video/MP2T"},
				Category:  MediaTypeVideo,
			},

			// Audio
			".mp3": {
				Extension: ".mp3",
				MimeTypes: []string{"audio/mpeg"},
				Category:  MediaTypeAudio,
			},
			".wav": {
				Extension: ".wav",
				MimeTypes: []string{"audio/wav"},
				Category:  MediaTypeAudio,
			},
			".ogg": {
				Extension: ".ogg",
				MimeTypes: []string{"audio/ogg"},
				Category:  MediaTypeAudio,
			},
			".flac": {
				Extension: ".flac",
				MimeTypes: []string{"audio/flac"},
				Category:  MediaTypeAudio,
			},
			".aac": {
				Extension: ".aac",
				MimeTypes: []string{"audio/aac"},
				Category:  MediaTypeAudio,
			},
			".m4a": {
				Extension: ".m4a",
				MimeTypes: []string{"audio/mp4"},
				Category:  MediaTypeAudio,
			},
			".wma": {
				Extension: ".wma",
				MimeTypes: []string{"audio/x-ms-wma"},
				Category:  MediaTypeAudio,
			},
			".opus": {
				Extension: ".opus",
				MimeTypes: []string{"audio/opus"},
				Category:  MediaTypeAudio,
			},
			".aiff": {
				Extension: ".aiff",
				MimeTypes: []string{"audio/aiff"},
				Category:  MediaTypeAudio,
			},
			".au": {
				Extension: ".au",
				MimeTypes: []string{"audio/basic"},
				Category:  MediaTypeAudio,
			},
			".ra": {
				Extension: ".ra",
				MimeTypes: []string{"audio/x-pn-realaudio"},
				Category:  MediaTypeAudio,
			},
			".ac3": {
				Extension: ".ac3",
				MimeTypes: []string{"audio/ac3"},
				Category:  MediaTypeAudio,
			},
			".dts": {
				Extension: ".dts",
				MimeTypes: []string{"audio/vnd.dts"},
				Category:  MediaTypeAudio,
			},
			".amr": {
				Extension: ".amr",
				MimeTypes: []string{"audio/amr"},
				Category:  MediaTypeAudio,
			},
			".ape": {
				Extension: ".ape",
				MimeTypes: []string{"audio/ape"},
				Category:  MediaTypeAudio,
			},
		},
	}
}

// ToJSON converts the file type configuration to JSON
func (c *FileTypeConfig) ToJSON() (string, error) {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal file type config: %w", err)
	}
	return string(data), nil
}

// FromJSON populates the file type configuration from JSON
func (c *FileTypeConfig) FromJSON(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), c)
}

// GetFileTypeInfo returns information about a specific file type
func (c *FileTypeConfig) GetFileTypeInfo(extension string) (FileTypeInfo, error) {
	// Ensure extension starts with a dot and is lowercase
	ext := extension
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ext = strings.ToLower(ext)

	info, exists := c.FileTypes[ext]
	if !exists {
		return FileTypeInfo{}, fmt.Errorf("unsupported file extension: %s", extension)
	}
	return info, nil
}

// GetSupportedExtensions returns all supported file extensions
func (c *FileTypeConfig) GetSupportedExtensions() []string {
	extensions := make([]string, 0, len(c.FileTypes))
	for ext := range c.FileTypes {
		extensions = append(extensions, ext)
	}
	return extensions
}

// GetExtensionsByCategory returns all extensions for a specific media category
func (c *FileTypeConfig) GetExtensionsByCategory(category MediaType) []string {
	extensions := make([]string, 0)
	for ext, info := range c.FileTypes {
		if info.Category == category {
			extensions = append(extensions, ext)
		}
	}
	return extensions
}

// GetMimeTypesByCategory returns all MIME types for a specific media category
func (c *FileTypeConfig) GetMimeTypesByCategory(category MediaType) []string {
	mimeTypes := make([]string, 0)
	for _, info := range c.FileTypes {
		if info.Category == category {
			mimeTypes = append(mimeTypes, info.MimeTypes...)
		}
	}
	return mimeTypes
}

// IsSupportedExtension checks if a file extension is supported
func (c *FileTypeConfig) IsSupportedExtension(extension string) bool {
	// Ensure extension starts with a dot and is lowercase
	ext := extension
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ext = strings.ToLower(ext)

	_, exists := c.FileTypes[ext]
	return exists
}

// IsSupportedMimeType checks if a MIME type is supported
func (c *FileTypeConfig) IsSupportedMimeType(mimeType string) bool {
	for _, info := range c.FileTypes {
		for _, mt := range info.MimeTypes {
			if mt == mimeType {
				return true
			}
		}
	}
	return false
}

// GetCategoryFromExtension returns the media category for a given extension
func (c *FileTypeConfig) GetCategoryFromExtension(extension string) (MediaType, error) {
	info, err := c.GetFileTypeInfo(extension)
	if err != nil {
		return "", err
	}
	return info.Category, nil
}

// GetCategoryFromMimeType returns the media category for a given MIME type
func (c *FileTypeConfig) GetCategoryFromMimeType(mimeType string) (MediaType, error) {
	for _, info := range c.FileTypes {
		for _, mt := range info.MimeTypes {
			if mt == mimeType {
				return info.Category, nil
			}
		}
	}
	return "", fmt.Errorf("unsupported MIME type: %s", mimeType)
}
