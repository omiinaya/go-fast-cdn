package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

func main() {
	// Parse command line flags
	rollback := flag.Bool("rollback", false, "Rollback the media unification migration")
	skipBackup := flag.Bool("skip-backup", false, "Skip creating a backup before migration")
	flag.Parse()

	// Load executable path
	util.LoadExPath()

	// Create backup manager
	backupManager := database.NewBackupManager()

	// Execute migration or rollback with backup
	var err error
	if *rollback {
		fmt.Println("=== Staging Migration Rollback ===")
		err = handleMigrationRollback(backupManager, *skipBackup)
	} else {
		fmt.Println("=== Staging Migration Execution ===")
		err = handleMigrationExecution(backupManager, *skipBackup)
	}

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("=== Migration completed successfully ===")
	os.Exit(0)
}

// handleMigrationExecution handles the execution of the migration with backup
func handleMigrationExecution(backupManager *database.BackupManager, skipBackup bool) error {
	var backupPath string
	var err error

	// Step 1: Create backup
	if !skipBackup {
		fmt.Println("Step 1: Creating database backup...")
		backupPath, err = backupManager.CreateBackup()
		if err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		fmt.Printf("✓ Backup created successfully: %s\n", backupPath)
	} else {
		fmt.Println("Step 1: Skipping backup creation as requested...")
	}

	// Step 2: Run migration
	fmt.Println("Step 2: Running media unification migration...")

	// Connect to database
	database.ConnectToDB()

	// Execute migration
	err = database.RunMediaMigration()
	if err != nil {
		// Migration failed, attempt rollback
		fmt.Printf("✗ Migration failed: %v\n", err)
		if !skipBackup && backupPath != "" {
			fmt.Println("Attempting to restore from backup...")
			if restoreErr := backupManager.RestoreBackup(backupPath); restoreErr != nil {
				return fmt.Errorf("migration failed and backup restore also failed: %w (restore error: %v)", err, restoreErr)
			}
			fmt.Printf("✓ Database restored from backup: %s\n", backupPath)
			return fmt.Errorf("migration failed but database was restored from backup: %w", err)
		}
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("✓ Media unification migration completed successfully!")

	// Step 3: Verify migration
	fmt.Println("Step 3: Verifying migration...")
	if err := verifyMigration(); err != nil {
		fmt.Printf("⚠ Migration verification warning: %v\n", err)
	} else {
		fmt.Println("✓ Migration verification completed successfully!")
	}

	return nil
}

// handleMigrationRollback handles the rollback of the migration with backup
func handleMigrationRollback(backupManager *database.BackupManager, skipBackup bool) error {
	var backupPath string
	var err error

	// Step 1: Create backup before rollback
	if !skipBackup {
		fmt.Println("Step 1: Creating database backup before rollback...")
		backupPath, err = backupManager.CreateBackup()
		if err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		fmt.Printf("✓ Backup created successfully: %s\n", backupPath)
	} else {
		fmt.Println("Step 1: Skipping backup creation as requested...")
	}

	// Step 2: Run rollback
	fmt.Println("Step 2: Rolling back media unification migration...")

	// Connect to database
	database.ConnectToDB()

	// Execute rollback
	err = database.RollbackMediaMigration()
	if err != nil {
		// Rollback failed, attempt restore
		fmt.Printf("✗ Rollback failed: %v\n", err)
		if !skipBackup && backupPath != "" {
			fmt.Println("Attempting to restore from backup...")
			if restoreErr := backupManager.RestoreBackup(backupPath); restoreErr != nil {
				return fmt.Errorf("rollback failed and backup restore also failed: %w (restore error: %v)", err, restoreErr)
			}
			fmt.Printf("✓ Database restored from backup: %s\n", backupPath)
			return fmt.Errorf("rollback failed but database was restored from backup: %w", err)
		}
		return fmt.Errorf("rollback failed: %w", err)
	}

	fmt.Println("✓ Media unification migration rollback completed successfully!")

	// Step 3: Verify rollback
	fmt.Println("Step 3: Verifying rollback...")
	if err := verifyRollback(); err != nil {
		fmt.Printf("⚠ Rollback verification warning: %v\n", err)
	} else {
		fmt.Println("✓ Rollback verification completed successfully!")
	}

	return nil
}

// verifyMigration verifies that the migration was successful
func verifyMigration() error {
	db := database.DB

	// Check if media table exists and has data
	var mediaCount int64
	if err := db.Table("media").Count(&mediaCount).Error; err != nil {
		return fmt.Errorf("failed to count media records: %w", err)
	}

	// Check if images and docs tables still exist (they should)
	var imageCount int64
	var docCount int64
	if err := db.Table("images").Count(&imageCount).Error; err != nil {
		return fmt.Errorf("failed to count image records: %w", err)
	}
	if err := db.Table("docs").Count(&docCount).Error; err != nil {
		return fmt.Errorf("failed to count doc records: %w", err)
	}

	fmt.Printf("  - Media table records: %d\n", mediaCount)
	fmt.Printf("  - Images table records: %d\n", imageCount)
	fmt.Printf("  - Docs table records: %d\n", docCount)

	// Verify that the total number of records matches
	if mediaCount != imageCount+docCount {
		return fmt.Errorf("media record count (%d) does not match sum of images (%d) and docs (%d)", mediaCount, imageCount, docCount)
	}

	return nil
}

// verifyRollback verifies that the rollback was successful
func verifyRollback() error {
	db := database.DB

	// Check if media table still exists (it shouldn't)
	var mediaCount int64
	if err := db.Table("media").Count(&mediaCount).Error; err == nil {
		return fmt.Errorf("media table still exists after rollback")
	}

	// Check if images and docs tables still exist (they should)
	var imageCount int64
	var docCount int64
	if err := db.Table("images").Count(&imageCount).Error; err != nil {
		return fmt.Errorf("failed to count image records: %w", err)
	}
	if err := db.Table("docs").Count(&docCount).Error; err != nil {
		return fmt.Errorf("failed to count doc records: %w", err)
	}

	fmt.Printf("  - Images table records: %d\n", imageCount)
	fmt.Printf("  - Docs table records: %d\n", docCount)

	return nil
}
