package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

// FileMigration handles the migration of files from separate directories to a unified media directory
type FileMigration struct {
	ImagesPath string
	DocsPath   string
	MediaPath  string
}

// NewFileMigration creates a new instance of FileMigration
func NewFileMigration() *FileMigration {
	return &FileMigration{
		ImagesPath: util.GetImagesPath(),
		DocsPath:   util.GetDocsPath(),
		MediaPath:  util.GetMediaPath(),
	}
}

// Run executes the file migration from separate directories to the unified media directory
func (fm *FileMigration) Run() error {
	log.Println("Starting file migration to unified media directory...")

	// Ensure the media directory exists
	if err := os.MkdirAll(fm.MediaPath, 0755); err != nil {
		return fmt.Errorf("failed to create media directory: %w", err)
	}

	// Migrate image files
	if err := fm.migrateFiles(fm.ImagesPath, "image"); err != nil {
		return fmt.Errorf("failed to migrate image files: %w", err)
	}

	// Migrate document files
	if err := fm.migrateFiles(fm.DocsPath, "document"); err != nil {
		return fmt.Errorf("failed to migrate document files: %w", err)
	}

	log.Println("File migration to unified media directory completed successfully!")
	return nil
}

// migrateFiles moves files from a source directory to the media directory
func (fm *FileMigration) migrateFiles(sourceDir, fileType string) error {
	log.Printf("Migrating %s files from %s...", fileType, sourceDir)

	// Check if source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		log.Printf("Source directory %s does not exist, skipping %s file migration", sourceDir, fileType)
		return nil
	}

	// Read all files in the source directory
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to read source directory %s: %w", sourceDir, err)
	}

	if len(files) == 0 {
		log.Printf("No files found in %s, skipping %s file migration", sourceDir, fileType)
		return nil
	}

	log.Printf("Found %d %s files to migrate", len(files), fileType)

	// Move each file to the media directory
	for i, file := range files {
		if file.IsDir() {
			continue // Skip subdirectories
		}

		sourcePath := filepath.Join(sourceDir, file.Name())
		destPath := filepath.Join(fm.MediaPath, file.Name())

		// Check if destination file already exists
		if _, err := os.Stat(destPath); err == nil {
			log.Printf("File %s already exists in media directory, skipping", file.Name())
			continue
		}

		// Copy the file to the new location
		if err := fm.copyFile(sourcePath, destPath); err != nil {
			log.Printf("Failed to copy file %s: %v", file.Name(), err)
			continue
		}

		// Log progress every 100 files
		if (i+1)%100 == 0 || i == len(files)-1 {
			log.Printf("Migrated %d/%d %s files", i+1, len(files), fileType)
		}
	}

	log.Printf("All %s files migrated successfully", fileType)
	return nil
}

// copyFile copies a file from source to destination
func (fm *FileMigration) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	// Preserve file permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// Rollback moves files back from the unified media directory to their original locations
func (fm *FileMigration) Rollback() error {
	log.Println("Starting file migration rollback...")

	// Ensure the legacy directories exist
	if err := os.MkdirAll(fm.ImagesPath, 0755); err != nil {
		return fmt.Errorf("failed to create images directory: %w", err)
	}

	if err := os.MkdirAll(fm.DocsPath, 0755); err != nil {
		return fmt.Errorf("failed to create docs directory: %w", err)
	}

	// Read all files in the media directory
	files, err := os.ReadDir(fm.MediaPath)
	if err != nil {
		return fmt.Errorf("failed to read media directory: %w", err)
	}

	if len(files) == 0 {
		log.Println("No files found in media directory, nothing to rollback")
		return nil
	}

	log.Printf("Found %d files to rollback", len(files))

	// Move each file back to its original location based on file type
	for i, file := range files {
		if file.IsDir() {
			continue // Skip subdirectories
		}

		sourcePath := filepath.Join(fm.MediaPath, file.Name())
		var destPath string

		// Determine destination based on file type
		mediaType, err := util.GetMediaTypeFromFilename(file.Name())
		if err != nil {
			log.Printf("Could not determine media type for file %s, skipping: %v", file.Name(), err)
			continue
		}

		switch mediaType {
		case util.MediaTypeImage:
			destPath = filepath.Join(fm.ImagesPath, file.Name())
		case util.MediaTypeDocument:
			destPath = filepath.Join(fm.DocsPath, file.Name())
		default:
			log.Printf("Unsupported media type for file %s, skipping", file.Name())
			continue
		}

		// Check if destination file already exists
		if _, err := os.Stat(destPath); err == nil {
			log.Printf("File %s already exists in legacy directory, skipping", file.Name())
			continue
		}

		// Copy the file to the legacy location
		if err := fm.copyFile(sourcePath, destPath); err != nil {
			log.Printf("Failed to copy file %s during rollback: %v", file.Name(), err)
			continue
		}

		// Log progress every 100 files
		if (i+1)%100 == 0 || i == len(files)-1 {
			log.Printf("Rolled back %d/%d files", i+1, len(files))
		}
	}

	log.Println("File migration rollback completed successfully!")
	return nil
}

// CleanupLegacyFiles removes files from the legacy directories after successful migration
func (fm *FileMigration) CleanupLegacyFiles() error {
	log.Println("Starting cleanup of legacy files...")

	// Clean up image files
	if err := fm.cleanupDirectory(fm.ImagesPath); err != nil {
		return fmt.Errorf("failed to cleanup images directory: %w", err)
	}

	// Clean up document files
	if err := fm.cleanupDirectory(fm.DocsPath); err != nil {
		return fmt.Errorf("failed to cleanup docs directory: %w", err)
	}

	log.Println("Legacy files cleanup completed successfully!")
	return nil
}

// cleanupDirectory removes all files from a directory but keeps the directory
func (fm *FileMigration) cleanupDirectory(dirPath string) error {
	// Check if directory exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		log.Printf("Directory %s does not exist, skipping cleanup", dirPath)
		return nil
	}

	// Read all files in the directory
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dirPath, err)
	}

	if len(files) == 0 {
		log.Printf("No files found in %s, nothing to cleanup", dirPath)
		return nil
	}

	log.Printf("Cleaning up %d files from %s", len(files), dirPath)

	// Remove each file
	for i, file := range files {
		if file.IsDir() {
			continue // Skip subdirectories
		}

		filePath := filepath.Join(dirPath, file.Name())
		if err := os.Remove(filePath); err != nil {
			log.Printf("Failed to remove file %s: %v", file.Name(), err)
			continue
		}

		// Log progress every 100 files
		if (i+1)%100 == 0 || i == len(files)-1 {
			log.Printf("Cleaned up %d/%d files from %s", i+1, len(files), dirPath)
		}
	}

	log.Printf("All files cleaned up from %s", dirPath)
	return nil
}
