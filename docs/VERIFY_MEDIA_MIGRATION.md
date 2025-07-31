# Media Migration Verification Script

This document describes the verification script used to validate the data integrity after migrating from separate images and docs tables to a unified media table.

## Overview

The verification script (`cmd/verify_media_migration/main.go`) performs comprehensive checks to ensure that the media unification migration was successful and no data was lost or corrupted during the process.

## What the Verification Script Checks

The script performs the following verification steps:

### 1. Table Existence Verification
- Checks that the media table exists
- Verifies that the original images and docs tables still exist

### 2. Record Count Verification
- Counts records in the images table
- Counts records in the docs table
- Counts records in the media table
- Verifies that the total number of records in the media table equals the sum of records in the images and docs tables

### 3. Image Migration Verification
- For each image in the images table:
  - Verifies that a corresponding record exists in the media table
  - Checks that the media type is set to "image"
  - Validates that the checksum matches the original
  - Ensures width and height are null (as they weren't in the original model)

### 4. Document Migration Verification
- For each document in the docs table:
  - Verifies that a corresponding record exists in the media table
  - Checks that the media type is set to "document"
  - Validates that the checksum matches the original
  - Ensures width and height are null

### 5. Orphaned Media Record Detection
- Checks for media records that don't have corresponding records in the original tables
- Identifies media records with unknown types

### 6. Migration Record Verification
- Checks that the migration record exists and is marked as completed

## Running the Verification Script

### Prerequisites

Before running the verification script, ensure that:
1. The media unification migration has been run
2. Go is installed and properly configured
3. The project dependencies are installed (`go mod download`)

### Using the Shell Scripts

#### For Linux/macOS:
```bash
# Make the script executable
chmod +x scripts/verify_media_migration.sh

# Run the verification script
./scripts/verify_media_migration.sh
```

#### For Windows:
```cmd
# Run the verification script
scripts\verify_media_migration.bat
```

### Running Directly with Go

You can also run the verification script directly using Go:

```bash
go run cmd/verify_media_migration/main.go
```

## Understanding the Output

The verification script provides detailed output for each check:

### Successful Verification Example
```
Starting Media Migration Verification
=====================================

Step 1: Verifying table existence...
✓ All required tables exist (media, images, docs)

Step 2: Verifying record counts...
  Images table: 150 records
  Docs table: 75 records
  Media table: 225 records
✓ Record counts match (images + docs = media)

Step 3: Verifying image migration...
✓ All 150 images migrated correctly

Step 4: Verifying document migration...
✓ All 75 documents migrated correctly

Step 5: Checking for orphaned media records...
✓ No orphaned media records found

Step 6: Verifying migration record...
✓ Migration record found and marked as completed

Verification Summary:
===================
Total images: 150
Total documents: 75
Total media records: 225
Image migration errors: 0
Document migration errors: 0
Orphaned media records: 0

✓ All verification checks passed! Migration was successful.
```

### Verification with Issues Example
```
Starting Media Migration Verification
=====================================

Step 1: Verifying table existence...
✓ All required tables exist (media, images, docs)

Step 2: Verifying record counts...
  Images table: 150 records
  Docs table: 75 records
  Media table: 224 records
✗ Record count mismatch: expected 225, got 224

Step 3: Verifying image migration...
ERROR: Failed to find migrated image missing_image.jpg: record not found
✗ 1 errors found in image migration

Step 4: Verifying document migration...
✓ All 75 documents migrated correctly

Step 5: Checking for orphaned media records...
WARNING: Orphaned document media record found: unknown_file.pdf
✗ 1 orphaned media records found

Step 6: Verifying migration record...
✓ Migration record found and marked as completed

Verification Summary:
===================
Total images: 150
Total documents: 75
Total media records: 224
Image migration errors: 1
Document migration errors: 0
Orphaned media records: 1

✗ 2 total errors found. Migration may have issues.
```

## Troubleshooting

### Common Issues

1. **"media table does not exist"**
   - This error indicates that the media unification migration hasn't been run yet.
   - Solution: Run the media migration first using `./scripts/media_migration.sh`

2. **"Failed to connect to database"**
   - This indicates a database connection issue.
   - Solution: Ensure the database is properly configured and accessible.

3. **"Record count mismatch"**
   - This means the total number of records in the media table doesn't match the sum of records in the images and docs tables.
   - Solution: Check for missing or duplicate records during migration.

4. **"Checksum mismatch"**
   - This indicates that the checksum of a migrated record doesn't match the original.
   - Solution: Check for data corruption during migration.

### Next Steps if Verification Fails

If the verification script reports errors:

1. Review the error messages to understand what went wrong
2. Check the migration logs for any issues during the migration process
3. Consider rolling back the migration using `./scripts/media_migration.sh --rollback`
4. Fix any issues identified
5. Re-run the migration
6. Run the verification script again

## Best Practices

1. **Always run the verification script after migration** to ensure data integrity
2. **Take a backup before migration** to have a restore point if needed
3. **Review the verification output carefully** even if it reports success
4. **Keep the verification logs** for audit purposes
5. **Run the verification in a staging environment** before production

## Integration with CI/CD

The verification script can be integrated into CI/CD pipelines to automatically validate migrations:

```yaml
# Example GitHub Actions workflow
- name: Run Media Migration
  run: ./scripts/media_migration.sh

- name: Verify Media Migration
  run: ./scripts/verify_media_migration.sh
```

The script will exit with code 0 on success and code 1 on failure, making it easy to integrate with automated workflows.