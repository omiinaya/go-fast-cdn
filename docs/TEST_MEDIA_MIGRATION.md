# Media Migration Testing Guide

This document provides instructions for testing the media migration script in a simulated staging environment.

## Overview

The test script performs the following steps:
1. Creates a test database separate from the main application
2. Populates it with sample image and document data
3. Runs the media migration script
4. Verifies that the migration was successful
5. Tests the rollback functionality
6. Cleans up the test database

## Prerequisites

- Go installed on your system
- Access to the project repository
- Sufficient permissions to create and modify files

## Running the Test

### For Linux/macOS Users

```bash
# Make the script executable
chmod +x scripts/test_media_migration.sh

# Run the full test suite
./scripts/test_media_migration.sh

# If you need to clean up a previous test without running new tests
./scripts/test_media_migration.sh --cleanup-only
```

### For Windows Users

```cmd
# Run the full test suite
scripts\test_media_migration.bat

# If you need to clean up a previous test without running new tests
scripts\test_media_migration.bat --cleanup-only
```

### Direct Go Command

Alternatively, you can run the test directly using Go:

```bash
# Run the full test suite
go run cmd/test_media_migration/main.go

# Clean up only
go run cmd/test_media_migration/main.go --cleanup-only
```

## Test Output

The test script provides detailed output for each step:

```
Starting Media Migration Test Suite
====================================

Step 1: Setting up test database...
✓ Test database set up successfully

Step 2: Populating with sample data...
✓ Created 3 sample images and 3 sample documents

Step 3: Running media migration...
✓ Media migration completed successfully

Step 4: Verifying migration...
✓ Migration verified successfully: 3 images + 3 documents = 6 media records

Step 5: Testing rollback...
✓ Rollback verified successfully: 3 images and 3 documents preserved

Step 6: Cleaning up test database...
✓ Test database cleaned up successfully

All tests passed successfully!
```

## What the Test Verifies

### 1. Database Setup
- Creates a separate test database in `db_data_test/` directory
- Initializes the database with the required tables

### 2. Sample Data Population
- Creates 3 sample image records
- Creates 3 sample document records
- Each record has a unique filename and checksum

### 3. Migration Execution
- Runs the media migration script
- Merges images and docs tables into a unified media table
- Marks the migration as completed in the migration records

### 4. Migration Verification
- Verifies that the media table was created
- Checks that all image and document records were migrated
- Ensures checksums are preserved during migration
- Confirms that the migration record was created

### 5. Rollback Testing
- Executes the rollback functionality
- Verifies that the media table is dropped
- Confirms that the migration record is removed
- Ensures original tables and data are preserved

### 6. Cleanup
- Removes the test database directory
- Closes database connections properly

## Troubleshooting

### Test Fails During Database Setup

If the test fails during database setup:
1. Ensure you have write permissions in the project directory
2. Check if there's already a `db_data_test` directory and remove it manually
3. Verify that SQLite is properly installed

### Test Fails During Migration

If the test fails during migration:
1. Check the error message for specific details
2. Verify that the migration code in `src/database/media_migration.go` is correct
3. Ensure the database models in `src/models/` are properly defined

### Test Fails During Verification

If the test fails during verification:
1. Check that all expected records were created in the media table
2. Verify that the checksums match between original and migrated records
3. Ensure the migration record was properly created

### Test Fails During Rollback

If the test fails during rollback:
1. Check that the media table is properly dropped
2. Verify that the migration record is removed
3. Ensure original tables still contain their data

### Cleanup Fails

If cleanup fails:
1. Manually remove the `db_data_test` directory
2. Check for any locked database files
3. Ensure no other processes are using the test database

## Test Database Location

The test database is created in a separate directory to avoid conflicts with the main application:
- Location: `{project_root}/db_data_test/test.db`
- This directory is completely removed after successful test completion

## Next Steps

After successfully running the test script:
1. Review the test output to ensure all steps passed
2. If any issues were found, fix them and re-run the test
3. Once all tests pass, the migration script is ready for the staging environment
4. Consider running the test with larger datasets to ensure scalability
5. Document any performance considerations for the actual migration

## Customization

To customize the test:
1. Modify the sample data in the `populateSampleData()` function in `cmd/test_media_migration/main.go`
2. Adjust the number of test records as needed
3. Add additional verification steps if required
4. Modify the test database location by changing the constants at the top of the file