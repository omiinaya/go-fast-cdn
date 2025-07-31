package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func DeleteFile(deletedFileName string, fileType string) error {
	filePath := fmt.Sprintf("%v/uploads/%v/%v", ExPath, fileType, deletedFileName)

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMediaFile is a wrapper for DeleteFile with clearer naming for media operations
func DeleteMediaFile(deletedFileName string, mediaType string) error {
	return DeleteFile(deletedFileName, mediaType)
}

// DeleteUnifiedMediaFile handles deletion of media files from the unified media directory
func DeleteUnifiedMediaFile(deletedFileName string) error {
	// Use the unified media directory instead of separate directories
	filePath := filepath.Join(ExPath, "uploads", "media", deletedFileName)

	err := os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil
}

// GetMediaFilePath returns the file path for a media file based on the unified approach
func GetMediaFilePath(fileName string) string {
	return filepath.Join(ExPath, "uploads", "media", fileName)
}

// GetLegacyMediaFilePath returns the file path for a media file based on the legacy approach
func GetLegacyMediaFilePath(fileName string, mediaType string) string {
	return filepath.Join(ExPath, "uploads", mediaType, fileName)
}

// FileExists checks if a file exists at the given path
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// MediaFileExists checks if a media file exists in the unified media directory
func MediaFileExists(fileName string) bool {
	filePath := GetMediaFilePath(fileName)
	return FileExists(filePath)
}

// LegacyMediaFileExists checks if a media file exists in the legacy media directories
func LegacyMediaFileExists(fileName string, mediaType string) bool {
	filePath := GetLegacyMediaFilePath(fileName, mediaType)
	return FileExists(filePath)
}

// GetFileExtension returns the file extension of a filename
func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}
