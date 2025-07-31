package database

import (
	"log"

	"github.com/kevinanielsen/go-fast-cdn/src/models"
)

// Migrate runs database migrations for all model structs using
// the global DB instance. This would typically be called on app startup.
func Migrate() {
	DB.AutoMigrate(&models.Image{}, &models.Doc{}, &models.User{}, &models.UserSession{}, &models.PasswordReset{})
}

// MigrateWithMedia runs database migrations including the new media table.
// This should be used after running the media unification migration.
func MigrateWithMedia() {
	DB.AutoMigrate(&models.Media{}, &models.Image{}, &models.Doc{}, &models.User{}, &models.UserSession{}, &models.PasswordReset{})
}

// RunMediaMigration runs the media unification migration to merge images and docs tables into media table.
// This is a convenience wrapper around the media migration package.
func RunMediaMigration() error {
	log.Println("Starting media unification migration...")

	migration := NewMediaMigration(DB)
	if err := migration.Run(); err != nil {
		log.Printf("Media unification migration failed: %v", err)
		return err
	}

	log.Println("Media unification migration completed successfully!")
	return nil
}

// RollbackMediaMigration rolls back the media unification migration.
// This is a convenience wrapper around the media migration package.
func RollbackMediaMigration() error {
	log.Println("Starting media unification migration rollback...")

	migration := NewMediaMigration(DB)
	if err := migration.Rollback(); err != nil {
		log.Printf("Media unification migration rollback failed: %v", err)
		return err
	}

	log.Println("Media unification migration rollback completed successfully!")
	return nil
}

// RunFileMigration runs the file migration to move files from separate directories to the unified media directory.
// This is a convenience wrapper around the file migration package.
func RunFileMigration() error {
	log.Println("Starting file migration to unified media directory...")

	fileMigration := NewFileMigration()
	if err := fileMigration.Run(); err != nil {
		log.Printf("File migration failed: %v", err)
		return err
	}

	log.Println("File migration to unified media directory completed successfully!")
	return nil
}

// RollbackFileMigration rolls back the file migration to legacy directories.
// This is a convenience wrapper around the file migration package.
func RollbackFileMigration() error {
	log.Println("Starting file migration rollback...")

	fileMigration := NewFileMigration()
	if err := fileMigration.Rollback(); err != nil {
		log.Printf("File migration rollback failed: %v", err)
		return err
	}

	log.Println("File migration rollback completed successfully!")
	return nil
}

// CleanupLegacyFiles removes files from the legacy directories after successful migration.
// This is a convenience wrapper around the file migration package.
func CleanupLegacyFiles() error {
	log.Println("Starting cleanup of legacy files...")

	fileMigration := NewFileMigration()
	if err := fileMigration.CleanupLegacyFiles(); err != nil {
		log.Printf("Legacy files cleanup failed: %v", err)
		return err
	}

	log.Println("Legacy files cleanup completed successfully!")
	return nil
}
