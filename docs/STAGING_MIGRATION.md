# Staging Migration Script

This document describes the staging migration script that safely executes the media unification migration with proper backup and error handling.

## Overview

The staging migration script (`cmd/staging_migration/main.go`) is designed to safely execute the media unification migration on staging environments. It includes:

1. **Database Backup**: Creates a backup of the current database before migration
2. **Migration Execution**: Runs the media unification migration
3. **Verification**: Verifies that the migration was successful
4. **Automatic Rollback**: If the migration fails, it automatically restores from the backup
5. **Clear Output**: Provides detailed progress and status information

## Features

### Safety Features

- **Automatic Backup**: Creates a backup before making any changes
- **Automatic Rollback**: If migration fails, restores from backup
- **Verification**: Verifies migration success after completion
- **Error Handling**: Comprehensive error handling with clear messages

### Flexibility

- **Skip Backup**: Option to skip backup creation (not recommended for production)
- **Rollback Support**: Can rollback previous migrations with the same safety features
- **Cross-Platform**: Works on both Linux/macOS and Windows

## Usage

### Linux/macOS

```bash
# Run migration with backup
./scripts/staging_migration.sh

# Run migration without backup (not recommended)
./scripts/staging_migration.sh --skip-backup

# Rollback migration with backup
./scripts/staging_migration.sh --rollback

# Rollback migration without backup (not recommended)
./scripts/staging_migration.sh --rollback-skip-backup
```

### Windows

```cmd
REM Run migration with backup
scripts\staging_migration.bat

REM Run migration without backup (not recommended)
scripts\staging_migration.bat --skip-backup

REM Rollback migration with backup
scripts\staging_migration.bat --rollback

REM Rollback migration without backup (not recommended)
scripts\staging_migration.bat --rollback-skip-backup
```

### Direct Execution

You can also run the Go program directly:

```bash
# Run migration with backup
go run cmd/staging_migration/main.go

# Run migration without backup
go run cmd/staging_migration/main.go --skip-backup

# Rollback migration with backup
go run cmd/staging_migration/main.go --rollback

# Rollback migration without backup
go run cmd/staging_migration/main.go --rollback --skip-backup
```

## Process Flow

### Migration Execution

1. **Step 1: Create Backup**
   - Creates a timestamped backup of the current database
   - Verifies backup integrity
   - Stores backup path for potential rollback

2. **Step 2: Run Migration**
   - Connects to the database
   - Executes the media unification migration
   - If migration fails:
     - Logs the error
     - Automatically restores from backup
     - Returns error with details

3. **Step 3: Verify Migration**
   - Checks that media table exists and has data
   - Verifies that images and docs tables still exist
   - Ensures record counts match expected values
   - Reports any discrepancies as warnings

### Migration Rollback

1. **Step 1: Create Backup**
   - Creates a backup of the current database before rollback
   - Verifies backup integrity
   - Stores backup path for potential restore

2. **Step 2: Run Rollback**
   - Connects to the database
   - Executes the media unification rollback
   - If rollback fails:
     - Logs the error
     - Automatically restores from backup
     - Returns error with details

3. **Step 3: Verify Rollback**
   - Checks that media table no longer exists
   - Verifies that images and docs tables still exist
   - Reports any discrepancies as warnings

## Error Handling

The script includes comprehensive error handling:

- **Backup Creation Failure**: Aborts migration if backup cannot be created
- **Migration Execution Failure**: Automatically restores from backup
- **Rollback Execution Failure**: Automatically restores from backup
- **Verification Warnings**: Reports issues but doesn't fail the operation
- **Database Connection Issues**: Reports connection errors clearly

## Backup Location

Backups are created in the `backups` directory in the project root. Backup files are named with timestamps:

```
backups/db_backup_20240101-120000.db
```

## Requirements

- Go 1.16 or higher
- Access to the database file
- Write permissions for the backup directory
- Sufficient disk space for backup (approximately same size as database)

## Best Practices

1. **Always Test First**: Run the script in a test environment before staging
2. **Monitor Disk Space**: Ensure sufficient disk space for backups
3. **Review Logs**: Check the output for any warnings or errors
4. **Keep Backups**: Retain backups until you're confident the migration was successful
5. **Document Execution**: Record when migrations are run for audit purposes

## Troubleshooting

### Common Issues

1. **Permission Denied**
   - Ensure the script has execute permissions
   - Check write permissions for the backup directory

2. **Disk Space Full**
   - Free up disk space before running the script
   - Consider cleaning up old backups

3. **Database Locked**
   - Ensure no other processes are using the database
   - Stop the application before running migration

4. **Migration Already Run**
   - The script will detect if migration has already been run
   - Use rollback option if you need to re-run the migration

### Getting Help

If you encounter issues:

1. Check the error messages in the script output
2. Verify database file exists and is accessible
3. Ensure you have the latest version of the script
4. Review the migration documentation

## Integration with CI/CD

The script can be integrated into CI/CD pipelines:

```yaml
# Example GitHub Actions workflow
- name: Run Staging Migration
  run: ./scripts/staging_migration.sh
```

```yaml
# Example rollback in CI/CD
- name: Rollback Staging Migration
  run: ./scripts/staging_migration.sh --rollback
  if: failure()