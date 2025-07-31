package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"gorm.io/gorm"
)

const (
	TestDbFolder = "db_data_test"
	TestDbName   = "test.db"
)

var testDB *gorm.DB

func main() {
	// Parse command line flags
	cleanupOnly := flag.Bool("cleanup-only", false, "Only clean up test database without running tests")
	flag.Parse()

	// Load executable path
	util.LoadExPath()

	// Set up test database path
	testDbPath := filepath.Join(util.ExPath, TestDbFolder)
	testDbFile := filepath.Join(testDbPath, TestDbName)

	// Handle cleanup-only mode
	if *cleanupOnly {
		fmt.Println("Cleaning up test database...")
		if err := cleanupTestDatabase(testDbPath); err != nil {
			log.Fatalf("Failed to clean up test database: %v", err)
		}
		fmt.Println("Test database cleaned up successfully!")
		return
	}

	fmt.Println("Starting Media Migration Test Suite")
	fmt.Println("====================================")

	// Step 1: Create test database
	if err := setupTestDatabase(testDbPath, testDbFile); err != nil {
		log.Fatalf("Failed to set up test database: %v", err)
	}

	// Step 2: Populate with sample data
	if err := populateSampleData(); err != nil {
		log.Fatalf("Failed to populate sample data: %v", err)
	}

	// Step 3: Run migration
	if err := runMigration(); err != nil {
		log.Printf("Migration test failed: %v", err)
		cleanupTestDatabase(testDbPath)
		os.Exit(1)
	}

	// Step 4: Verify migration
	if err := verifyMigration(); err != nil {
		log.Printf("Migration verification failed: %v", err)
		cleanupTestDatabase(testDbPath)
		os.Exit(1)
	}

	// Step 5: Test rollback
	if err := testRollback(); err != nil {
		log.Printf("Rollback test failed: %v", err)
		cleanupTestDatabase(testDbPath)
		os.Exit(1)
	}

	// Step 6: Clean up
	if err := cleanupTestDatabase(testDbPath); err != nil {
		log.Printf("Cleanup failed: %v", err)
		os.Exit(1)
	}

	fmt.Println("\nAll tests passed successfully!")
}

func setupTestDatabase(testDbPath, testDbFile string) error {
	fmt.Println("\nStep 1: Setting up test database...")

	// Create test database directory if it doesn't exist
	if _, err := os.Stat(testDbPath); os.IsNotExist(err) {
		if err := os.Mkdir(testDbPath, 0755); err != nil {
			return fmt.Errorf("failed to create test database directory: %w", err)
		}
	}

	// Create test database file if it doesn't exist
	if _, err := os.Stat(testDbFile); os.IsNotExist(err) {
		file, err := os.Create(testDbFile)
		if err != nil {
			return fmt.Errorf("failed to create test database file: %w", err)
		}
		file.Close()
	}

	// Connect to test database
	var err error
	testDB, err = gorm.Open(sqlite.Open(testDbFile), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to test database: %w", err)
	}

	// Set the global DB variable to our test database
	database.DB = testDB

	// Run initial migrations
	if err := testDB.AutoMigrate(&models.Image{}, &models.Doc{}, &models.Config{}); err != nil {
		return fmt.Errorf("failed to run initial migrations: %w", err)
	}

	fmt.Println("✓ Test database set up successfully")
	return nil
}

func populateSampleData() error {
	fmt.Println("\nStep 2: Populating with sample data...")

	// Create sample images
	sampleImages := []models.Image{
		{FileName: "image1.jpg", Checksum: generateChecksum("image1.jpg")},
		{FileName: "image2.png", Checksum: generateChecksum("image2.png")},
		{FileName: "image3.gif", Checksum: generateChecksum("image3.gif")},
	}

	for _, img := range sampleImages {
		if err := testDB.Create(&img).Error; err != nil {
			return fmt.Errorf("failed to create sample image %s: %w", img.FileName, err)
		}
	}

	// Create sample documents
	sampleDocs := []models.Doc{
		{FileName: "document1.pdf", Checksum: generateChecksum("document1.pdf")},
		{FileName: "document2.docx", Checksum: generateChecksum("document2.docx")},
		{FileName: "document3.txt", Checksum: generateChecksum("document3.txt")},
	}

	for _, doc := range sampleDocs {
		if err := testDB.Create(&doc).Error; err != nil {
			return fmt.Errorf("failed to create sample document %s: %w", doc.FileName, err)
		}
	}

	// Verify data was created
	var imageCount, docCount int64
	testDB.Model(&models.Image{}).Count(&imageCount)
	testDB.Model(&models.Doc{}).Count(&docCount)

	fmt.Printf("✓ Created %d sample images and %d sample documents\n", imageCount, docCount)
	return nil
}

func runMigration() error {
	fmt.Println("\nStep 3: Running media migration...")

	// Run the migration
	if err := database.RunMediaMigration(); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("✓ Media migration completed successfully")
	return nil
}

func verifyMigration() error {
	fmt.Println("\nStep 4: Verifying migration...")

	// Check if media table exists
	if !testDB.Migrator().HasTable(&models.Media{}) {
		return fmt.Errorf("media table was not created")
	}

	// Count records in all tables
	var imageCount, docCount, mediaCount int64
	testDB.Model(&models.Image{}).Count(&imageCount)
	testDB.Model(&models.Doc{}).Count(&docCount)
	testDB.Model(&models.Media{}).Count(&mediaCount)

	// Verify counts match
	expectedMediaCount := imageCount + docCount
	if mediaCount != expectedMediaCount {
		return fmt.Errorf("media count mismatch: expected %d, got %d", expectedMediaCount, mediaCount)
	}

	// Verify image data was migrated correctly
	var images []models.Image
	testDB.Find(&images)
	for _, img := range images {
		var media models.Media
		if err := testDB.Where("file_name = ? AND type = ?", img.FileName, models.MediaTypeImage).First(&media).Error; err != nil {
			return fmt.Errorf("failed to find migrated image %s: %w", img.FileName, err)
		}
		if string(media.Checksum) != string(img.Checksum) {
			return fmt.Errorf("checksum mismatch for image %s", img.FileName)
		}
	}

	// Verify document data was migrated correctly
	var docs []models.Doc
	testDB.Find(&docs)
	for _, doc := range docs {
		var media models.Media
		if err := testDB.Where("file_name = ? AND type = ?", doc.FileName, models.MediaTypeDocument).First(&media).Error; err != nil {
			return fmt.Errorf("failed to find migrated document %s: %w", doc.FileName, err)
		}
		if string(media.Checksum) != string(doc.Checksum) {
			return fmt.Errorf("checksum mismatch for document %s", doc.FileName)
		}
	}

	// Check migration record
	var migrationRecord database.MigrationRecord
	if err := testDB.Where("name = ?", "media_unification_2024").First(&migrationRecord).Error; err != nil {
		return fmt.Errorf("migration record not found: %w", err)
	}
	if !migrationRecord.Completed {
		return fmt.Errorf("migration record not marked as completed")
	}

	fmt.Printf("✓ Migration verified successfully: %d images + %d documents = %d media records\n",
		imageCount, docCount, mediaCount)
	return nil
}

func testRollback() error {
	fmt.Println("\nStep 5: Testing rollback...")

	// Get counts before rollback
	var mediaCountBefore int64
	testDB.Model(&models.Media{}).Count(&mediaCountBefore)

	// Run rollback
	if err := database.RollbackMediaMigration(); err != nil {
		return fmt.Errorf("rollback failed: %w", err)
	}

	// Verify media table was dropped
	if testDB.Migrator().HasTable(&models.Media{}) {
		return fmt.Errorf("media table was not dropped during rollback")
	}

	// Verify migration record was removed
	var migrationRecord database.MigrationRecord
	if err := testDB.Where("name = ?", "media_unification_2024").First(&migrationRecord).Error; err == nil {
		return fmt.Errorf("migration record was not removed during rollback")
	}

	// Verify original tables still exist and have data
	var imageCount, docCount int64
	testDB.Model(&models.Image{}).Count(&imageCount)
	testDB.Model(&models.Doc{}).Count(&docCount)

	if imageCount == 0 {
		return fmt.Errorf("no images found after rollback")
	}
	if docCount == 0 {
		return fmt.Errorf("no documents found after rollback")
	}

	fmt.Printf("✓ Rollback verified successfully: %d images and %d documents preserved\n", imageCount, docCount)
	return nil
}

func cleanupTestDatabase(testDbPath string) error {
	fmt.Println("\nStep 6: Cleaning up test database...")

	// Close database connection if open
	if testDB != nil {
		sqlDB, err := testDB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	// Remove test database directory and all contents
	if err := os.RemoveAll(testDbPath); err != nil {
		return fmt.Errorf("failed to remove test database directory: %w", err)
	}

	fmt.Println("✓ Test database cleaned up successfully")
	return nil
}

func generateChecksum(filename string) []byte {
	hash := sha256.Sum256([]byte(filename + time.Now().String()))
	return hash[:]
}
