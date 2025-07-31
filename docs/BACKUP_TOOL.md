# Database Backup Tool

This document describes the database backup tool for the go-fast-cdn project.

## Overview

The database backup tool is designed to create and restore backups of the SQLite database used by the go-fast-cdn application. It provides a command-line interface for performing backup and restore operations.

## Features

- Create complete database backups with timestamps
- Restore database from backup files
- List all available backups
- Delete specific backup files
- Backup file integrity verification
- Error handling and status reporting

## Usage

### Building the Tool

The backup tool can be built using the provided scripts:

**Linux/macOS:**
```bash
chmod +x scripts/db_backup.sh
./scripts/db_backup.sh
```

**Windows:**
```cmd
scripts\db_backup.bat
```

Alternatively, you can build it manually:
```bash
go build -o bin/db_backup cmd/db_backup/main.go
```

### Commands

#### Create a Backup

Create a new database backup in the default location:
```bash
./scripts/db_backup.sh create
```

Create a backup in a custom location:
```bash
./scripts/db_backup.sh create -output /path/to/custom/location
```

#### Restore from Backup

Restore the database from a specific backup file:
```bash
./scripts/db_backup.sh restore -backup /path/to/backup/file
```

#### List Backups

List all available backup files:
```bash
./scripts/db_backup.sh list
```

#### Delete Backup

Delete a specific backup file:
```bash
./scripts/db_backup.sh delete -backup /path/to/backup/file
```

### Command Reference

```
db_backup <command> [options]

Commands:
  create     Create a new database backup
  restore    Restore database from a backup
  list       List all available backups
  delete     Delete a specific backup

Create Command Options:
  -output    Custom output path for backup (optional)
             Example: db_backup create -output /path/to/custom/location

Restore Command Options:
  -backup    Path to backup file to restore from (required)
             Example: db_backup restore -backup /path/to/backup/file

Delete Command Options:
  -backup    Path to backup file to delete (required)
             Example: db_backup delete -backup /path/to/backup/file
```

## Backup Location

By default, backups are stored in the `backups` directory relative to the application's executable path. The backup files are named with a timestamp pattern: `db_backup_YYYYMMDD-HHMMSS.db`.

## Backup Process

1. The tool creates a backup directory if it doesn't exist
2. It generates a timestamp for the backup filename
3. It copies the current database file to the backup location
4. It verifies the integrity of the backup file
5. It reports the success or failure of the operation

## Restore Process

1. The tool verifies that the backup file exists
2. It checks the integrity of the backup file
3. It creates a backup of the current database before restoring (for safety)
4. It copies the backup file to the database location
5. It reports the success or failure of the operation

## Error Handling

The tool includes comprehensive error handling:

- Backup directory creation failures
- File copy failures
- Backup file integrity verification failures
- Missing backup files
- Permission issues

All errors are reported with descriptive messages to help with troubleshooting.

## Safety Features

- **Backup verification**: Each backup is verified to ensure it's a valid SQLite database
- **Pre-restore backup**: Before restoring, the current database is backed up for safety
- **Confirmation prompts**: Destructive operations (restore, delete) require confirmation
- **Partial cleanup**: If a backup operation fails partway through, partial files are cleaned up

## Integration with Phase 1 Unification Project

This backup tool is a critical component of Phase 1 of the unification project. It ensures that:

1. The current database state is preserved before migration
2. A restore point is available if the migration encounters issues
3. Both images and documents tables are included in the backup
4. The backup can be restored independently of the main application

## Troubleshooting

### Permission Errors

If you encounter permission errors:
- Ensure the application has write access to the backup directory
- Check file permissions on the database file
- Run the tool with appropriate permissions

### Database Lock Errors

If the database is locked:
- Stop the main application before creating a backup
- Ensure no other processes are accessing the database file

### Backup File Integrity Issues

If backup verification fails:
- Check disk space
- Verify the source database file is not corrupted
- Try creating the backup again