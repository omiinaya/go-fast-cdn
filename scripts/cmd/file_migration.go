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
	rollback := flag.Bool("rollback", false, "Rollback the file migration to legacy directories")
	cleanup := flag.Bool("cleanup", false, "Clean up legacy files after successful migration")
	flag.Parse()

	// Load executable path
	util.LoadExPath()

	// Create file migration instance
	fileMigration := database.NewFileMigration()

	// Execute migration, rollback, or cleanup
	var err error
	switch {
	case *rollback:
		fmt.Println("Rolling back file migration...")
		err = fileMigration.Rollback()
		if err != nil {
			log.Fatalf("Failed to rollback file migration: %v", err)
		}
		fmt.Println("File migration rolled back successfully!")
	case *cleanup:
		fmt.Println("Cleaning up legacy files...")
		err = fileMigration.CleanupLegacyFiles()
		if err != nil {
			log.Fatalf("Failed to cleanup legacy files: %v", err)
		}
		fmt.Println("Legacy files cleanup completed successfully!")
	default:
		fmt.Println("Running file migration to unified media directory...")
		err = fileMigration.Run()
		if err != nil {
			log.Fatalf("Failed to run file migration: %v", err)
		}
		fmt.Println("File migration completed successfully!")
		fmt.Println("You can now run with --cleanup to remove files from legacy directories")
	}

	os.Exit(0)
}
