# Media Unification Migration

This document describes the migration process to merge the separate `images` and `docs` tables into a single unified `media` table as part of Phase 2 of the unification project.

## Overview

The migration script performs the following operations:

1. Creates a new `media` table based on the unified media model
2. Migrates all data from the existing `images` table to the `media` table
3. Migrates all data from the existing `docs` table to the `media` table
4. Sets the appropriate media type ("image" or "document") for each record
5. Handles image-specific fields (width, height) for image records
6. Preserves all existing data including checksums and filenames
7. Includes proper error handling and transaction management to ensure data integrity
8. Provides clear output about the migration progress and status

## Prerequisites

Before running the migration, ensure that:

1. You have a backup of your database
2. The application is not running
3. You have sufficient disk space for the migration

## Running the Migration

### Using the Shell Scripts

#### On Linux/macOS:

```bash
# Run the migration
./scripts/media_migration.sh

# Rollback the migration
./scripts/media_migration.sh --rollback
```

#### On Windows:

```cmd
REM Run the migration
scripts\media_migration.bat

REM Rollback the migration
scripts\media_migration.bat --rollback
```

### Using the Go Command Directly

```bash
# Run the migration
go run cmd/media_migration/main.go

# Rollback the migration
go run cmd/media_migration/main.go --rollback
```

### Using the Database Package

If you want to run the migration programmatically:

```go
import "github.com/kevinanielsen/go-fast-cdn/src/database"

// Run the migration
err := database.RunMediaMigration()
if err != nil {
    // Handle error
}

// Rollback the migration
err = database.RollbackMediaMigration()
if err != nil {
    // Handle error
}
```

## Migration Process Details

### Step 1: Create the Media Table

The migration creates a new `media` table with the following schema:

```go
type Media struct {
    gorm.Model
    
    FileName string    `json:"file_name" gorm:"uniqueIndex"`
    Checksum []byte    `json:"checksum"`
    Type     MediaType `json:"type" gorm:"type:varchar(20);not null;default:'document'"`
    
    // Image-specific fields (will be empty/null for non-image media)
    Width  *int `json:"width,omitempty" gorm:"default:null"`
    Height *int `json:"height,omitempty" gorm:"default:null"`
}
```

### Step 2: Migrate Images Data

The migration:

1. Fetches all records from the `images` table
2. Converts each image record to a media record using `MediaFromImage`
3. Sets the media type to "image"
4. Saves the media record to the `media` table
5. Reports progress every 100 records

### Step 3: Migrate Docs Data

The migration:

1. Fetches all records from the `docs` table
2. Converts each doc record to a media record using `MediaFromDoc`
3. Sets the media type to "document"
4. Saves the media record to the `media` table
5. Reports progress every 100 records

### Step 4: Mark Migration as Completed

The migration creates a record in the `migration_records` table to track that the migration has been completed. This prevents the migration from running multiple times.

## Rollback Process

The rollback process:

1. Drops the `media` table
2. Removes the migration completion marker
3. Restores the original state (the `images` and `docs` tables remain unchanged)

For a comprehensive rollback plan including emergency procedures, see [MEDIA_MIGRATION_ROLLBACK_PLAN.md](MEDIA_MIGRATION_ROLLBACK_PLAN.md).

### Emergency Rollback

In case of emergency, use the interactive rollback scripts:

- **Linux/macOS**: `./scripts/emergency_rollback.sh`
- **Windows**: `scripts\emergency_rollback.bat`

For a quick reference checklist, see [EMERGENCY_ROLLBACK_CHECKLIST.md](EMERGENCY_ROLLBACK_CHECKLIST.md).

## Error Handling

The migration includes comprehensive error handling:

1. All operations are wrapped in a database transaction
2. If any step fails, the transaction is rolled back
3. Detailed error messages are logged
4. The migration can be safely rerun after fixing any issues

## Post-Migration Steps

After successfully running the migration:

1. Update your application code to use the new `Media` model instead of the separate `Image` and `Doc` models
2. Update the database migration to include the `Media` model:

```go
// Instead of:
database.Migrate()

// Use:
database.MigrateWithMedia()
```

3. Test your application thoroughly to ensure all functionality works with the unified media model

## Troubleshooting

### Migration Fails

If the migration fails:

1. Check the error message for details
2. Ensure you have sufficient disk space
3. Verify that the database is not being used by another process
4. Try running the rollback and then running the migration again

### Rollback Fails

If the rollback fails:

1. Check the error message for details
2. Ensure the database is not being used by another process
3. Restore from your backup if necessary

## Data Integrity

The migration is designed to preserve all existing data:

1. All filenames are preserved exactly
2. All checksums are preserved exactly
3. Creation and update timestamps are preserved
4. No data is lost during the migration process

## Performance Considerations

For large databases:

1. The migration processes records in batches
2. Progress is reported regularly
3. The migration may take some time for very large databases
4. Ensure you have sufficient time to complete the migration before starting