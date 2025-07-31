package database

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/kevinanielsen/go-fast-cdn/src/util"
	"gorm.io/gorm"
)

// BackupManager handles database backup and restore operations
type BackupManager struct {
	dbPath     string
	backupPath string
}

// NewBackupManager creates a new BackupManager instance
func NewBackupManager() *BackupManager {
	// Try to find the database file in common locations
	var dbPath string

	// Check if database exists in the project directory
	projectDbPath := fmt.Sprintf("%v/%s/%s", util.ExPath, DbFolder, DbName)
	if _, err := os.Stat(projectDbPath); err == nil {
		dbPath = projectDbPath
	} else {
		// If not found, try to find it in the current working directory
		cwdDbPath := fmt.Sprintf("%s/%s/%s", ".", DbFolder, DbName)
		if _, err := os.Stat(cwdDbPath); err == nil {
			dbPath = cwdDbPath
		} else {
			// If still not found, use the default path
			dbPath = projectDbPath
		}
	}

	backupPath := fmt.Sprintf("%v/backups", util.ExPath)

	return &BackupManager{
		dbPath:     dbPath,
		backupPath: backupPath,
	}
}

// CreateBackup creates a complete backup of the database
func (bm *BackupManager) CreateBackup() (string, error) {
	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(bm.backupPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create backup directory: %v", err)
	}

	// Generate timestamp for backup filename
	timestamp := time.Now().Format("20060102-150405")
	backupFileName := fmt.Sprintf("db_backup_%s.db", timestamp)
	backupFilePath := filepath.Join(bm.backupPath, backupFileName)

	// Open source database file
	sourceFile, err := os.Open(bm.dbPath)
	if err != nil {
		return "", fmt.Errorf("failed to open source database: %v", err)
	}
	defer sourceFile.Close()

	// Create backup file
	backupFile, err := os.Create(backupFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %v", err)
	}
	defer backupFile.Close()

	// Copy database file to backup location
	bytesWritten, err := io.Copy(backupFile, sourceFile)
	if err != nil {
		// Clean up partial backup file
		os.Remove(backupFilePath)
		return "", fmt.Errorf("failed to copy database file: %v", err)
	}

	// Verify backup file integrity
	if err := bm.verifyBackup(backupFilePath); err != nil {
		os.Remove(backupFilePath)
		return "", fmt.Errorf("backup verification failed: %v", err)
	}

	log.Printf("Backup created successfully: %s (%d bytes)", backupFilePath, bytesWritten)
	return backupFilePath, nil
}

// verifyBackup verifies the integrity of a backup file
func (bm *BackupManager) verifyBackup(backupPath string) error {
	// Try to open the backup file with GORM to verify it's a valid SQLite database
	testDB, err := gorm.Open(sqlite.Open(backupPath), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to open backup file for verification: %v", err)
	}
	defer func() {
		sqlDB, _ := testDB.DB()
		sqlDB.Close()
	}()

	// Try to query the database to ensure it's accessible
	var result int
	if err := testDB.Raw("SELECT 1").Scan(&result).Error; err != nil {
		return fmt.Errorf("backup file integrity check failed: %v", err)
	}

	return nil
}

// RestoreBackup restores the database from a backup file
func (bm *BackupManager) RestoreBackup(backupFilePath string) error {
	// Verify backup file exists
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %s", backupFilePath)
	}

	// Verify backup file integrity
	if err := bm.verifyBackup(backupFilePath); err != nil {
		return fmt.Errorf("backup file verification failed: %v", err)
	}

	// Create a backup of the current database before restoring
	currentBackupPath, err := bm.CreateBackup()
	if err != nil {
		log.Printf("Warning: Failed to create pre-restore backup: %v", err)
	} else {
		log.Printf("Pre-restore backup created: %s", currentBackupPath)
	}

	// Open backup file
	backupFile, err := os.Open(backupFilePath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %v", err)
	}
	defer backupFile.Close()

	// Ensure database directory exists
	dbDir := filepath.Dir(bm.dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("failed to create database directory: %v", err)
	}

	// Create/overwrite database file
	dbFile, err := os.Create(bm.dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database file: %v", err)
	}
	defer dbFile.Close()

	// Copy backup file to database location
	bytesWritten, err := io.Copy(dbFile, backupFile)
	if err != nil {
		return fmt.Errorf("failed to restore database file: %v", err)
	}

	log.Printf("Database restored successfully from %s (%d bytes)", backupFilePath, bytesWritten)
	return nil
}

// ListBackups returns a list of all available backup files
func (bm *BackupManager) ListBackups() ([]string, error) {
	var backups []string

	files, err := os.ReadDir(bm.backupPath)
	if err != nil {
		if os.IsNotExist(err) {
			return backups, nil // No backup directory yet
		}
		return nil, fmt.Errorf("failed to read backup directory: %v", err)
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".db" {
			backups = append(backups, filepath.Join(bm.backupPath, file.Name()))
		}
	}

	return backups, nil
}

// DeleteBackup deletes a specific backup file
func (bm *BackupManager) DeleteBackup(backupFilePath string) error {
	if err := os.Remove(backupFilePath); err != nil {
		return fmt.Errorf("failed to delete backup file: %v", err)
	}

	log.Printf("Backup deleted: %s", backupFilePath)
	return nil
}
