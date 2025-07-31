package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting Media Migration Verification")
	fmt.Println("=====================================")

	// Load executable path
	util.LoadExPath()

	// Connect to database
	database.ConnectToDB()

	// Run verification
	if err := verifyMigration(); err != nil {
		log.Printf("Migration verification failed: %v", err)
		os.Exit(1)
	}

	fmt.Println("\nVerification completed successfully!")
}

func verifyMigration() error {
	// Check if media table exists
	if !database.DB.Migrator().HasTable(&models.Media{}) {
		return fmt.Errorf("media table does not exist - migration may not have been run")
	}

	// Check if original tables still exist
	if !database.DB.Migrator().HasTable(&models.Image{}) {
		return fmt.Errorf("images table does not exist")
	}
	if !database.DB.Migrator().HasTable(&models.Doc{}) {
		return fmt.Errorf("docs table does not exist")
	}

	fmt.Println("\nStep 1: Verifying table existence...")
	fmt.Println("✓ All required tables exist (media, images, docs)")

	// Count records in all tables
	var imageCount, docCount, mediaCount int64
	database.DB.Model(&models.Image{}).Count(&imageCount)
	database.DB.Model(&models.Doc{}).Count(&docCount)
	database.DB.Model(&models.Media{}).Count(&mediaCount)

	fmt.Println("\nStep 2: Verifying record counts...")
	fmt.Printf("  Images table: %d records\n", imageCount)
	fmt.Printf("  Docs table: %d records\n", docCount)
	fmt.Printf("  Media table: %d records\n", mediaCount)

	// Verify counts match
	expectedMediaCount := imageCount + docCount
	if mediaCount != expectedMediaCount {
		return fmt.Errorf("media count mismatch: expected %d, got %d", expectedMediaCount, mediaCount)
	}
	fmt.Println("✓ Record counts match (images + docs = media)")

	// Verify image data was migrated correctly
	fmt.Println("\nStep 3: Verifying image migration...")
	var images []models.Image
	database.DB.Find(&images)
	imageMigrationErrors := 0

	for _, img := range images {
		var media models.Media
		if err := database.DB.Where("file_name = ? AND type = ?", img.FileName, models.MediaTypeImage).First(&media).Error; err != nil {
			log.Printf("ERROR: Failed to find migrated image %s: %v", img.FileName, err)
			imageMigrationErrors++
			continue
		}

		// Check checksum
		if !equalByteSlices(media.Checksum, img.Checksum) {
			log.Printf("ERROR: Checksum mismatch for image %s", img.FileName)
			imageMigrationErrors++
			continue
		}

		// Check media type
		if media.Type != models.MediaTypeImage {
			log.Printf("ERROR: Incorrect media type for image %s: expected %s, got %s",
				img.FileName, models.MediaTypeImage, media.Type)
			imageMigrationErrors++
			continue
		}

		// For images, width and height should be null (not set in original model)
		if media.Width != nil || media.Height != nil {
			log.Printf("WARNING: Width/height set for migrated image %s (should be null)", img.FileName)
		}
	}

	if imageMigrationErrors == 0 {
		fmt.Printf("✓ All %d images migrated correctly\n", len(images))
	} else {
		fmt.Printf("✗ %d errors found in image migration\n", imageMigrationErrors)
	}

	// Verify document data was migrated correctly
	fmt.Println("\nStep 4: Verifying document migration...")
	var docs []models.Doc
	database.DB.Find(&docs)
	docMigrationErrors := 0

	for _, doc := range docs {
		var media models.Media
		if err := database.DB.Where("file_name = ? AND type = ?", doc.FileName, models.MediaTypeDocument).First(&media).Error; err != nil {
			log.Printf("ERROR: Failed to find migrated document %s: %v", doc.FileName, err)
			docMigrationErrors++
			continue
		}

		// Check checksum
		if !equalByteSlices(media.Checksum, doc.Checksum) {
			log.Printf("ERROR: Checksum mismatch for document %s", doc.FileName)
			docMigrationErrors++
			continue
		}

		// Check media type
		if media.Type != models.MediaTypeDocument {
			log.Printf("ERROR: Incorrect media type for document %s: expected %s, got %s",
				doc.FileName, models.MediaTypeDocument, media.Type)
			docMigrationErrors++
			continue
		}

		// For documents, width and height should be null
		if media.Width != nil || media.Height != nil {
			log.Printf("WARNING: Width/height set for migrated document %s (should be null)", doc.FileName)
		}
	}

	if docMigrationErrors == 0 {
		fmt.Printf("✓ All %d documents migrated correctly\n", len(docs))
	} else {
		fmt.Printf("✗ %d errors found in document migration\n", docMigrationErrors)
	}

	// Check for orphaned media records (media without corresponding image or doc)
	fmt.Println("\nStep 5: Checking for orphaned media records...")
	orphanedMediaCount := 0

	// Get all media records
	var allMedia []models.Media
	database.DB.Find(&allMedia)

	for _, media := range allMedia {
		if media.Type == models.MediaTypeImage {
			var img models.Image
			if err := database.DB.Where("file_name = ?", media.FileName).First(&img).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Printf("WARNING: Orphaned image media record found: %s", media.FileName)
					orphanedMediaCount++
				}
			}
		} else if media.Type == models.MediaTypeDocument {
			var doc models.Doc
			if err := database.DB.Where("file_name = ?", media.FileName).First(&doc).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					log.Printf("WARNING: Orphaned document media record found: %s", media.FileName)
					orphanedMediaCount++
				}
			}
		} else {
			log.Printf("WARNING: Media record with unknown type: %s (type: %s)", media.FileName, media.Type)
			orphanedMediaCount++
		}
	}

	if orphanedMediaCount == 0 {
		fmt.Println("✓ No orphaned media records found")
	} else {
		fmt.Printf("✗ %d orphaned media records found\n", orphanedMediaCount)
	}

	// Check migration record
	fmt.Println("\nStep 6: Verifying migration record...")
	var migrationRecord database.MigrationRecord
	if err := database.DB.Where("name = ?", "media_unification_2024").First(&migrationRecord).Error; err != nil {
		log.Printf("WARNING: Migration record not found: %v", err)
	} else {
		if !migrationRecord.Completed {
			log.Printf("WARNING: Migration record not marked as completed")
		} else {
			fmt.Println("✓ Migration record found and marked as completed")
		}
	}

	// Summary
	fmt.Println("\nVerification Summary:")
	fmt.Println("===================")
	fmt.Printf("Total images: %d\n", imageCount)
	fmt.Printf("Total documents: %d\n", docCount)
	fmt.Printf("Total media records: %d\n", mediaCount)
	fmt.Printf("Image migration errors: %d\n", imageMigrationErrors)
	fmt.Printf("Document migration errors: %d\n", docMigrationErrors)
	fmt.Printf("Orphaned media records: %d\n", orphanedMediaCount)

	totalErrors := imageMigrationErrors + docMigrationErrors + orphanedMediaCount
	if totalErrors == 0 {
		fmt.Println("\n✓ All verification checks passed! Migration was successful.")
		return nil
	} else {
		fmt.Printf("\n✗ %d total errors found. Migration may have issues.\n", totalErrors)
		return fmt.Errorf("%d verification errors found", totalErrors)
	}
}

// Helper function to compare byte slices
func equalByteSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
