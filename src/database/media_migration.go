package database

import (
	"fmt"
	"log"

	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"gorm.io/gorm"
)

// MediaMigration handles the migration from separate images and docs tables to a unified media table
type MediaMigration struct {
	DB *gorm.DB
}

// NewMediaMigration creates a new instance of MediaMigration
func NewMediaMigration(db *gorm.DB) *MediaMigration {
	return &MediaMigration{DB: db}
}

// Run executes the migration from images and docs tables to media table
func (m *MediaMigration) Run() error {
	log.Println("Starting media unification migration...")

	// Begin transaction for data integrity
	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Migration failed, rolled back: %v", r)
		}
	}()

	// Check if migration has already been run
	if m.hasMigrationRun(tx) {
		log.Println("Media unification migration has already been completed. Skipping.")
		return nil
	}

	// Step 1: Create the media table
	log.Println("Step 1: Creating media table...")
	if err := m.createMediaTable(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to create media table: %w", err)
	}

	// Step 2: Migrate images data
	log.Println("Step 2: Migrating images data...")
	if err := m.migrateImagesData(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to migrate images data: %w", err)
	}

	// Step 3: Migrate docs data
	log.Println("Step 3: Migrating docs data...")
	if err := m.migrateDocsData(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to migrate docs data: %w", err)
	}

	// Step 4: Mark migration as completed
	log.Println("Step 4: Marking migration as completed...")
	if err := m.markMigrationCompleted(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to mark migration as completed: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit migration transaction: %w", err)
	}

	// Step 5: Migrate files to unified directory (outside of database transaction)
	log.Println("Step 5: Migrating files to unified directory...")
	if err := m.migrateFiles(); err != nil {
		log.Printf("Warning: File migration failed: %v", err)
		log.Println("Database migration was successful, but you may need to run file migration manually")
	}

	log.Println("Media unification migration completed successfully!")
	return nil
}

// Rollback reverts the migration, dropping the media table and restoring the original tables
func (m *MediaMigration) Rollback() error {
	log.Println("Starting media unification migration rollback...")

	// Begin transaction for data integrity
	tx := m.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Rollback failed: %v", r)
		}
	}()

	// Check if migration has been run
	if !m.hasMigrationRun(tx) {
		log.Println("Media unification migration has not been run. Nothing to rollback.")
		return nil
	}

	// Step 1: Drop the media table
	log.Println("Step 1: Dropping media table...")
	if err := tx.Migrator().DropTable(&models.Media{}); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to drop media table: %w", err)
	}

	// Step 2: Remove migration completion marker
	log.Println("Step 2: Removing migration completion marker...")
	if err := m.removeMigrationMarker(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration marker: %w", err)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback transaction: %w", err)
	}

	// Step 3: Rollback files to legacy directories (outside of database transaction)
	log.Println("Step 3: Rolling back files to legacy directories...")
	if err := m.rollbackFiles(); err != nil {
		log.Printf("Warning: File rollback failed: %v", err)
		log.Println("Database rollback was successful, but you may need to run file rollback manually")
	}

	log.Println("Media unification migration rollback completed successfully!")
	return nil
}

// createMediaTable creates the media table with the unified schema
func (m *MediaMigration) createMediaTable(tx *gorm.DB) error {
	// Create the media table using AutoMigrate
	if err := tx.AutoMigrate(&models.Media{}); err != nil {
		return fmt.Errorf("failed to create media table: %w", err)
	}
	log.Println("Media table created successfully")
	return nil
}

// migrateImagesData migrates all data from the images table to the media table
func (m *MediaMigration) migrateImagesData(tx *gorm.DB) error {
	// Get all images
	var images []models.Image
	if err := tx.Find(&images).Error; err != nil {
		return fmt.Errorf("failed to fetch images: %w", err)
	}

	log.Printf("Found %d images to migrate", len(images))

	// Convert each image to media and save
	for i, image := range images {
		media := models.MediaFromImage(image)

		// For images, we don't have width/height in the original model, so they'll be null
		// In a real implementation, you might want to extract this information from the image files

		if err := tx.Create(&media).Error; err != nil {
			return fmt.Errorf("failed to migrate image %s: %w", image.FileName, err)
		}

		// Log progress every 100 images
		if (i+1)%100 == 0 || i == len(images)-1 {
			log.Printf("Migrated %d/%d images", i+1, len(images))
		}
	}

	log.Println("All images migrated successfully")
	return nil
}

// migrateDocsData migrates all data from the docs table to the media table
func (m *MediaMigration) migrateDocsData(tx *gorm.DB) error {
	// Get all docs
	var docs []models.Doc
	if err := tx.Find(&docs).Error; err != nil {
		return fmt.Errorf("failed to fetch docs: %w", err)
	}

	log.Printf("Found %d documents to migrate", len(docs))

	// Convert each doc to media and save
	for i, doc := range docs {
		media := models.MediaFromDoc(doc)

		if err := tx.Create(&media).Error; err != nil {
			return fmt.Errorf("failed to migrate document %s: %w", doc.FileName, err)
		}

		// Log progress every 100 docs
		if (i+1)%100 == 0 || i == len(docs)-1 {
			log.Printf("Migrated %d/%d documents", i+1, len(docs))
		}
	}

	log.Println("All documents migrated successfully")
	return nil
}

// MigrationRecord tracks which migrations have been run
type MigrationRecord struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"uniqueIndex;not null"`
	Completed bool   `gorm:"not null;default:false"`
	gorm.Model
}

// hasMigrationRun checks if the media unification migration has already been completed
func (m *MediaMigration) hasMigrationRun(tx *gorm.DB) bool {
	// Ensure the migration records table exists
	if err := tx.AutoMigrate(&MigrationRecord{}); err != nil {
		log.Printf("Warning: Failed to create migration records table: %v", err)
		return false
	}

	var record MigrationRecord
	result := tx.Where("name = ?", "media_unification_2024").First(&record)
	return result.Error == nil && record.Completed
}

// markMigrationCompleted marks the media unification migration as completed
func (m *MediaMigration) markMigrationCompleted(tx *gorm.DB) error {
	record := MigrationRecord{
		Name:      "media_unification_2024",
		Completed: true,
	}

	// Use OnConflict to update if record already exists
	return tx.Where("name = ?", record.Name).
		Assign(record).
		FirstOrCreate(&record).Error
}

// removeMigrationMarker removes the migration completion marker
func (m *MediaMigration) removeMigrationMarker(tx *gorm.DB) error {
	return tx.Where("name = ?", "media_unification_2024").Delete(&MigrationRecord{}).Error
}

// migrateFiles handles the migration of files from separate directories to the unified media directory
func (m *MediaMigration) migrateFiles() error {
	fileMigration := NewFileMigration()
	return fileMigration.Run()
}

// rollbackFiles handles the rollback of files from the unified media directory to legacy directories
func (m *MediaMigration) rollbackFiles() error {
	fileMigration := NewFileMigration()
	return fileMigration.Rollback()
}
