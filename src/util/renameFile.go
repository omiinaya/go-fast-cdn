package util

import (
	"os"
	"path/filepath"
)

func RenameFile(oldName, newName, fileType string) error {
	prefix := filepath.Join(ExPath, "uploads", fileType)

	err := os.Rename(
		filepath.Join(prefix, oldName),
		filepath.Join(prefix, newName),
	)
	if err != nil {
		return err
	}

	return nil
}

// RenameMediaFile is a wrapper for RenameFile with clearer naming for media operations
func RenameMediaFile(oldName, newName, mediaType string) error {
	return RenameFile(oldName, newName, mediaType)
}

// RenameUnifiedMediaFile handles renaming of media files in the unified media directory
func RenameUnifiedMediaFile(oldName, newName string) error {
	// Use the unified media directory instead of separate directories
	prefix := filepath.Join(ExPath, "uploads", "media")

	err := os.Rename(
		filepath.Join(prefix, oldName),
		filepath.Join(prefix, newName),
	)
	if err != nil {
		return err
	}

	return nil
}

// MoveFileToUnifiedDirectory moves a file from a legacy directory to the unified media directory
func MoveFileToUnifiedDirectory(fileName, mediaType string) error {
	oldPath := filepath.Join(ExPath, "uploads", mediaType, fileName)
	newPath := filepath.Join(ExPath, "uploads", "media", fileName)

	// Ensure the target directory exists
	if err := os.MkdirAll(filepath.Join(ExPath, "uploads", "media"), 0755); err != nil {
		return err
	}

	err := os.Rename(oldPath, newPath)
	if err != nil {
		return err
	}

	return nil
}

// CopyFileToUnifiedDirectory copies a file from a legacy directory to the unified media directory
func CopyFileToUnifiedDirectory(fileName, mediaType string) error {
	oldPath := filepath.Join(ExPath, "uploads", mediaType, fileName)
	newPath := filepath.Join(ExPath, "uploads", "media", fileName)

	// Ensure the target directory exists
	if err := os.MkdirAll(filepath.Join(ExPath, "uploads", "media"), 0755); err != nil {
		return err
	}

	// Read the source file
	data, err := os.ReadFile(oldPath)
	if err != nil {
		return err
	}

	// Write to the destination file
	err = os.WriteFile(newPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
