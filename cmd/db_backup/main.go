package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kevinanielsen/go-fast-cdn/src/database"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

func main() {
	// Load executable path
	util.LoadExPath()

	// Create backup manager
	backupManager := database.NewBackupManager()

	// Define command-line flags
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	restoreCmd := flag.NewFlagSet("restore", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// Create command flags
	createOutputPath := createCmd.String("output", "", "Custom output path for backup (optional)")

	// Restore command flags
	restoreBackupPath := restoreCmd.String("backup", "", "Path to backup file to restore from")
	restoreForce := restoreCmd.Bool("force", false, "Skip confirmation prompt for restore")

	// Delete command flags
	deleteBackupPath := deleteCmd.String("backup", "", "Path to backup file to delete")
	deleteForce := deleteCmd.Bool("force", false, "Skip confirmation prompt for delete")

	// Check if at least one command is provided
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// Parse the command
	switch os.Args[1] {
	case "create":
		err := createCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("Error parsing create command: %v", err)
		}
		handleCreateCommand(backupManager, *createOutputPath)
	case "restore":
		err := restoreCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("Error parsing restore command: %v", err)
		}
		handleRestoreCommand(backupManager, *restoreBackupPath, *restoreForce)
	case "list":
		err := listCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("Error parsing list command: %v", err)
		}
		handleListCommand(backupManager)
	case "delete":
		err := deleteCmd.Parse(os.Args[2:])
		if err != nil {
			log.Fatalf("Error parsing delete command: %v", err)
		}
		handleDeleteCommand(backupManager, *deleteBackupPath, *deleteForce)
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Database Backup Tool")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  db_backup <command> [options]")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  create     Create a new database backup")
	fmt.Println("  restore    Restore database from a backup")
	fmt.Println("  list       List all available backups")
	fmt.Println("  delete     Delete a specific backup")
	fmt.Println("")
	fmt.Println("Create Command Options:")
	fmt.Println("  -output    Custom output path for backup (optional)")
	fmt.Println("             Example: db_backup create -output /path/to/custom/location")
	fmt.Println("")
	fmt.Println("Restore Command Options:")
	fmt.Println("  -backup    Path to backup file to restore from (required)")
	fmt.Println("             Example: db_backup restore -backup /path/to/backup/file")
	fmt.Println("")
	fmt.Println("Delete Command Options:")
	fmt.Println("  -backup    Path to backup file to delete (required)")
	fmt.Println("             Example: db_backup delete -backup /path/to/backup/file")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Create a backup in the default location")
	fmt.Println("  db_backup create")
	fmt.Println("")
	fmt.Println("  # Create a backup in a custom location")
	fmt.Println("  db_backup create -output /my/backups")
	fmt.Println("")
	fmt.Println("  # List all available backups")
	fmt.Println("  db_backup list")
	fmt.Println("")
	fmt.Println("  # Restore from a specific backup")
	fmt.Println("  db_backup restore -backup /path/to/db_backup_20230101-120000.db")
	fmt.Println("")
	fmt.Println("  # Delete a specific backup")
	fmt.Println("  db_backup delete -backup /path/to/db_backup_20230101-120000.db")
}

func handleCreateCommand(backupManager *database.BackupManager, outputPath string) {
	fmt.Println("Creating database backup...")

	var backupPath string
	var err error

	if outputPath != "" {
		// If custom output path is provided, we need to handle it differently
		// For now, we'll just use the default backup manager and then move the file
		backupPath, err = backupManager.CreateBackup()
		if err != nil {
			log.Fatalf("Failed to create backup: %v", err)
		}

		// Move the backup file to the custom location
		customPath := filepath.Join(outputPath, filepath.Base(backupPath))
		if err := os.Rename(backupPath, customPath); err != nil {
			log.Fatalf("Failed to move backup to custom location: %v", err)
		}
		backupPath = customPath
	} else {
		backupPath, err = backupManager.CreateBackup()
		if err != nil {
			log.Fatalf("Failed to create backup: %v", err)
		}
	}

	fmt.Printf("✓ Backup created successfully: %s\n", backupPath)
}

func handleRestoreCommand(backupManager *database.BackupManager, backupPath string, force bool) {
	if backupPath == "" {
		fmt.Println("Error: Backup path is required for restore command")
		fmt.Println("Use: db_backup restore -backup /path/to/backup/file")
		os.Exit(1)
	}

	// Check if the backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		log.Fatalf("Backup file does not exist: %s", backupPath)
	}

	fmt.Printf("Restoring database from backup: %s\n", backupPath)
	fmt.Println("Warning: This will overwrite the current database.")

	if !force {
		fmt.Print("Are you sure you want to continue? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Restore operation cancelled.")
			return
		}
	}

	err := backupManager.RestoreBackup(backupPath)
	if err != nil {
		log.Fatalf("Failed to restore database: %v", err)
	}

	fmt.Printf("✓ Database restored successfully from: %s\n", backupPath)
}

func handleListCommand(backupManager *database.BackupManager) {
	fmt.Println("Available backups:")

	backups, err := backupManager.ListBackups()
	if err != nil {
		log.Fatalf("Failed to list backups: %v", err)
	}

	if len(backups) == 0 {
		fmt.Println("No backups found.")
		return
	}

	for i, backup := range backups {
		fmt.Printf("%d. %s\n", i+1, backup)
	}
}

func handleDeleteCommand(backupManager *database.BackupManager, backupPath string, force bool) {
	if backupPath == "" {
		fmt.Println("Error: Backup path is required for delete command")
		fmt.Println("Use: db_backup delete -backup /path/to/backup/file")
		os.Exit(1)
	}

	// Check if the backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		log.Fatalf("Backup file does not exist: %s", backupPath)
	}

	fmt.Printf("Deleting backup: %s\n", backupPath)

	if !force {
		fmt.Print("Are you sure you want to continue? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Delete operation cancelled.")
			return
		}
	}

	err := backupManager.DeleteBackup(backupPath)
	if err != nil {
		log.Fatalf("Failed to delete backup: %v", err)
	}

	fmt.Printf("✓ Backup deleted successfully: %s\n", backupPath)
}
