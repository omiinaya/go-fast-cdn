package util

import (
	"os"
	"path/filepath"
)

var ExPath string

// LoadExPath loads the executable path and stores it in the ExPath variable.
// It uses os.Executable to get the path of the current executable.
// The filepath.Dir function is used to get the directory containing the executable.
// If there is an error getting the executable path, it will panic.
func LoadExPath() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	ExPath = exPath
}

// GetUploadsPath returns the path to the uploads directory
func GetUploadsPath() string {
	return filepath.Join(ExPath, "uploads")
}

// GetMediaPath returns the path to the unified media directory
func GetMediaPath() string {
	return filepath.Join(ExPath, "uploads", "media")
}

// GetImagesPath returns the path to the images directory (for backward compatibility)
func GetImagesPath() string {
	return filepath.Join(ExPath, "uploads", "images")
}

// GetDocsPath returns the path to the documents directory (for backward compatibility)
func GetDocsPath() string {
	return filepath.Join(ExPath, "uploads", "docs")
}

// EnsureUploadDirectories ensures that all upload directories exist
func EnsureUploadDirectories() error {
	// Create the main uploads directory if it doesn't exist
	if err := os.MkdirAll(GetUploadsPath(), 0755); err != nil {
		return err
	}

	// Create the unified media directory if it doesn't exist
	if err := os.MkdirAll(GetMediaPath(), 0755); err != nil {
		return err
	}

	// Create legacy directories for backward compatibility
	if err := os.MkdirAll(GetImagesPath(), 0755); err != nil {
		return err
	}

	if err := os.MkdirAll(GetDocsPath(), 0755); err != nil {
		return err
	}

	return nil
}
