# Staging Migration Script - Quick Start Guide

This guide provides a quick overview of how to use the staging migration script to safely execute the media unification migration.

## What is this?

The staging migration script is a safe way to execute the media unification migration on your staging environment. It:

1. **Creates a backup** of your current database
2. **Runs the migration** to unify images and docs tables
3. **Verifies the migration** was successful
4. **Automatically rolls back** if anything goes wrong

## Quick Start

### For Linux/macOS:

```bash
# Run the migration with backup (recommended)
./scripts/staging_migration.sh
```

### For Windows:

```cmd
# Run the migration with backup (recommended)
scripts\staging_migration.bat
```

## Other Options

### Skip Backup (Not Recommended)

```bash
# Linux/macOS
./scripts/staging_migration.sh --skip-backup

# Windows
scripts\staging_migration.bat --skip-backup
```

### Rollback Migration

```bash
# Linux/macOS
./scripts/staging_migration.sh --rollback

# Windows
scripts\staging_migration.bat --rollback
```

### Rollback Without Backup (Not Recommended)

```bash
# Linux/macOS
./scripts/staging_migration.sh --rollback-skip-backup

# Windows
scripts\staging_migration.bat --rollback-skip-backup
```

## What Happens During Migration?

1. **Backup Creation**: A timestamped backup is created in the `backups/` directory
2. **Migration Execution**: The script runs the media unification migration
3. **Verification**: It checks that the migration was successful
4. **Completion**: You'll see a success message if everything worked

## What If Something Goes Wrong?

If the migration fails:

1. The script will automatically detect the failure
2. It will restore your database from the backup
3. You'll see detailed error messages explaining what went wrong
4. Your database will be in the same state as before the migration

## Before You Run

1. **Stop your application** - Make sure no other processes are using the database
2. **Check disk space** - Ensure you have enough space for the backup (same size as your database)
3. **Test first** - Consider running in a test environment first

## After Migration

1. **Check the output** - Look for any warnings or errors
2. **Test your application** - Make sure everything works as expected
3. **Keep the backup** - Don't delete the backup until you're confident everything is working

## Need Help?

- Check the full documentation: [`docs/STAGING_MIGRATION.md`](docs/STAGING_MIGRATION.md)
- Review error messages in the script output
- Ensure you have the latest version of the script

## Example Output

```
=== Staging Migration Script ===
This script will:
1. Create a backup of the current database
2. Run the media unification migration
3. Verify the migration was successful
4. Automatically rollback if the migration fails

Mode: Migration

=== Staging Migration Execution ===
Step 1: Creating database backup...
✓ Backup created successfully: /path/to/project/backups/db_backup_20240101-120000.db
Step 2: Running media unification migration...
2024/01/01 12:00:00 Starting media unification migration...
2024/01/01 12:00:00 Step 1: Creating media table...
2024/01/01 12:00:00 Media table created successfully
2024/01/01 12:00:00 Step 2: Migrating images data...
2024/01/01 12:00:00 Found 150 images to migrate
2024/01/01 12:00:01 All images migrated successfully
2024/01/01 12:00:01 Step 3: Migrating docs data...
2024/01/01 12:00:01 Found 75 documents to migrate
2024/01/01 12:00:01 All documents migrated successfully
2024/01/01 12:00:01 Step 4: Marking migration as completed...
2024/01/01 12:00:01 Media unification migration completed successfully!
✓ Media unification migration completed successfully!
Step 3: Verifying migration...
  - Media table records: 225
  - Images table records: 150
  - Docs table records: 75
✓ Migration verification completed successfully!

=== Migration completed successfully ===

=== Staging migration completed successfully! ===