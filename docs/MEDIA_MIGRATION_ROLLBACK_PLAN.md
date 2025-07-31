# Media Migration Rollback Plan

## Overview

This document outlines the rollback plan for the media unification migration that merges the separate `images` and `docs` tables into a unified `media` table. This plan is designed to ensure a quick and safe revert to the original state if any issues arise during or after the migration.

## Scenarios Requiring Rollback

### 1. Data Corruption
- **Scenario**: Data integrity issues detected during or after migration
- **Indicators**:
  - Checksum mismatches between original and migrated data
  - Missing records in the media table
  - Orphaned records in the media table
  - Data validation failures

### 2. Application Errors
- **Scenario**: Application fails to function correctly after migration
- **Indicators**:
  - Application crashes or errors when accessing media
  - API endpoints returning unexpected errors
  - User interface issues related to media display
  - Authentication or authorization problems with media access

### 3. Performance Issues
- **Scenario**: Significant performance degradation after migration
- **Indicators**:
  - Slower response times for media operations
  - Increased database load
  - Memory usage spikes
  - Timeout errors during media operations

### 4. Migration Process Failures
- **Scenario**: Migration process fails to complete
- **Indicators**:
  - Migration script errors or crashes
  - Incomplete migration (partial data transfer)
  - Migration stuck in progress
  - Transaction rollback during migration

### 5. Verification Failures
- **Scenario**: Post-migration verification detects issues
- **Indicators**:
  - Verification script reports errors
  - Record count mismatches
  - Data type conversion errors
  - Migration record not properly created

## Rollback Timeline

The rollback should be initiated as soon as possible after detecting issues:

- **Critical Issues** (data corruption, application crashes): Immediate rollback (within 15 minutes)
- **Major Issues** (significant functionality loss): Within 1 hour
- **Minor Issues** (performance degradation, non-critical bugs): Within 4 hours or during scheduled maintenance

## Rollback Procedures

### Option 1: Using the Built-in Rollback Functionality (Preferred)

This is the first and preferred method for rolling back the migration. The migration script includes built-in rollback functionality that reverts the database changes.

#### Step-by-Step Instructions:

1. **Stop the Application**
   ```bash
   # Stop any running instances of the application
   # This prevents new data from being written during rollback
   ```

2. **Run the Rollback Script**
   
   **On Linux/macOS:**
   ```bash
   ./scripts/media_migration.sh --rollback
   ```
   
   **On Windows:**
   ```cmd
   scripts\media_migration.bat --rollback
   ```
   
   **Using Go directly:**
   ```bash
   go run cmd/media_migration/main.go --rollback
   ```

3. **Verify Rollback Success**
   ```bash
   # Run the verification script to confirm rollback
   ./scripts/verify_media_migration.sh
   ```
   
   The verification should confirm:
   - The `media` table no longer exists
   - The `images` and `docs` tables are intact with all original data
   - No migration completion marker exists

4. **Restart the Application**
   ```bash
   # Restart the application with the original database schema
   ```

### Option 2: Database Backup Restoration

If the built-in rollback functionality fails or encounters issues, restore from a pre-migration backup.

#### Step-by-Step Instructions:

1. **Stop the Application**
   ```bash
   # Stop any running instances of the application
   ```

2. **List Available Backups**
   ```bash
   # List all available backups
   ./scripts/db_backup.sh list
   ```

3. **Identify the Pre-Migration Backup**
   - Look for the most recent backup created before the migration
   - Note the full path to this backup file

4. **Restore from Backup**
   ```bash
   # Restore from the identified backup
   ./scripts/db_backup.sh restore -backup /path/to/pre_migration_backup.db -force
   ```
   
   The `-force` flag skips the confirmation prompt, which is useful for automation.

5. **Verify Restoration Success**
   ```bash
   # Run the verification script to confirm the database state
   ./scripts/verify_media_migration.sh
   ```
   
   The verification should confirm:
   - The `media` table does not exist
   - The `images` and `docs` tables contain all original data
   - The database is in a consistent state

6. **Restart the Application**
   ```bash
   # Restart the application with the restored database
   ```

### Option 3: Manual Database Restoration

If both the built-in rollback and backup restoration fail, manual intervention may be required.

#### Step-by-Step Instructions:

1. **Stop the Application**
   ```bash
   # Stop any running instances of the application
   ```

2. **Create a Current State Backup**
   ```bash
   # Create a backup of the current state for analysis
   ./scripts/db_backup.sh create -output /path/to/emergency_backup
   ```

3. **Replace Database File**
   ```bash
   # Replace the current database file with a known good backup
   cp /path/to/pre_migration_backup.db /path/to/current/database.db
   ```

4. **Verify Restoration**
   ```bash
   # Manually verify the database state
   sqlite3 /path/to/current/database.db ".tables"
   sqlite3 /path/to/current/database.db "SELECT COUNT(*) FROM images;"
   sqlite3 /path/to/current/database.db "SELECT COUNT(*) FROM docs;"
   ```

5. **Restart the Application**
   ```bash
   # Restart the application with the restored database
   ```

## Verification Steps After Rollback

After performing any rollback procedure, follow these verification steps:

### 1. Database Schema Verification
```bash
# Check that the media table no longer exists
sqlite3 database.db ".tables"
# The output should NOT include "media"

# Check that the images and docs tables exist
sqlite3 database.db ".schema images"
sqlite3 database.db ".schema docs"
```

### 2. Data Integrity Verification
```bash
# Verify record counts match pre-migration state
sqlite3 database.db "SELECT COUNT(*) FROM images;"
sqlite3 database.db "SELECT COUNT(*) FROM docs;"

# Compare these counts with the pre-migration backup counts
sqlite3 /path/to/pre_migration_backup.db "SELECT COUNT(*) FROM images;"
sqlite3 /path/to/pre_migration_backup.db "SELECT COUNT(*) FROM docs;"
```

### 3. Application Functionality Verification
- Start the application
- Test image upload and retrieval
- Test document upload and retrieval
- Verify all media-related API endpoints
- Check user interface functionality
- Verify authentication and authorization for media access

### 4. Run Verification Script
```bash
# Run the automated verification script
./scripts/verify_media_migration.sh
```

The script should report:
- No media table exists
- Images and docs tables exist with expected record counts
- No migration completion marker

## Notification Protocol

### Who to Notify

1. **Immediate Notification** (within 15 minutes of rollback initiation):
   - Development Team Lead
   - DevOps/Infrastructure Team
   - Product Manager
   - Stakeholders

2. **Post-Rollback Notification** (within 1 hour of rollback completion):
   - All Development Team members
   - Quality Assurance Team
   - Customer Support Team
   - End Users (via appropriate communication channels)

### Notification Template

**Subject**: URGENT: Media Migration Rollback Initiated

**Body**:
```
Team,

A rollback of the media migration has been initiated due to [brief description of issue].

Rollback Details:
- Time Initiated: [timestamp]
- Rollback Method: [built-in script/backup restoration/manual]
- Expected Duration: [estimated time]

Impact:
- Application may be temporarily unavailable
- Recent media uploads may be lost
- System will be restored to pre-migration state

Next Steps:
1. Complete rollback process
2. Verify system functionality
3. Investigate root cause
4. Plan for re-migration

We will provide updates as the situation progresses.

Regards,
[Your Name/Team]
```

## Post-Rollback Steps

### 1. System Stabilization
- Monitor system performance and logs
- Verify all application functionality
- Check for any residual issues
- Ensure all services are running normally

### 2. Data Analysis
- Analyze the cause of the rollback
- Review logs and error messages
- Examine the failed migration data
- Document findings for future reference

### 3. Issue Resolution
- Fix the identified issues
- Update migration scripts if necessary
- Enhance error handling and validation
- Improve backup and restoration procedures

### 4. Re-migration Planning
- Schedule a new migration window
- Update migration plan based on lessons learned
- Increase testing and validation
- Prepare additional rollback options

### 5. Documentation Update
- Update this rollback plan with any lessons learned
- Document the root cause and resolution
- Share findings with the team
- Update migration procedures

## Rollback Scripts and Tools

### 1. Built-in Rollback Script
- **Location**: [`cmd/media_migration/main.go`](cmd/media_migration/main.go)
- **Usage**: `go run cmd/media_migration/main.go --rollback`
- **Functionality**: 
  - Drops the media table
  - Removes migration completion marker
  - Preserves original images and docs tables

### 2. Shell Script Wrappers
- **Linux/macOS**: [`scripts/media_migration.sh`](scripts/media_migration.sh)
- **Windows**: [`scripts/media_migration.bat`](scripts/media_migration.bat)
- **Usage**: `./scripts/media_migration.sh --rollback`

### 3. Database Backup Tool
- **Location**: [`cmd/db_backup/main.go`](cmd/db_backup/main.go)
- **Shell Wrapper**: [`scripts/db_backup.sh`](scripts/db_backup.sh)
- **Functions**:
  - Create backups: `./scripts/db_backup.sh create`
  - List backups: `./scripts/db_backup.sh list`
  - Restore backups: `./scripts/db_backup.sh restore -backup /path/to/backup`
  - Delete backups: `./scripts/db_backup.sh delete -backup /path/to/backup`

### 4. Verification Script
- **Location**: [`cmd/verify_media_migration/main.go`](cmd/verify_media_migration/main.go)
- **Shell Wrapper**: [`scripts/verify_media_migration.sh`](scripts/verify_media_migration.sh)
- **Functionality**:
  - Verifies database schema
  - Checks data integrity
  - Validates record counts
  - Reports any issues

## Testing the Rollback Plan

Before executing the actual migration, test the rollback plan:

1. **Create a Test Environment**
   ```bash
   # Run the migration test script
   ./scripts/test_media_migration.sh
   ```

2. **Test Rollback Functionality**
   ```bash
   # After running the test migration, test the rollback
   cd db_data_test
   sqlite3 test.db ".tables"
   # Verify media table exists
   
   # Run rollback
   go run cmd/media_migration/main.go --rollback
   
   # Verify rollback success
   sqlite3 test.db ".tables"
   # Verify media table no longer exists
   ```

3. **Test Backup Restoration**
   ```bash
   # Create a backup
   ./scripts/db_backup.sh create
   
   # List backups and note the path
   ./scripts/db_backup.sh list
   
   # Restore from backup
   ./scripts/db_backup.sh restore -backup /path/to/backup -force
   ```

## Emergency Contacts

| Role | Name | Contact | Availability |
|------|------|---------|-------------|
| Development Lead | [Name] | [Email/Phone] | 24/7 |
| DevOps Lead | [Name] | [Email/Phone] | 24/7 |
| Product Manager | [Name] | [Email/Phone] | Business Hours |
| Stakeholder | [Name] | [Email/Phone] | Business Hours |

## Lessons Learned and Continuous Improvement

After any rollback event:

1. **Document the Incident**
   - Record the timeline of events
   - Document the root cause
   - Capture all error messages and logs
   - Note the resolution steps taken

2. **Analyze the Process**
   - Evaluate the effectiveness of the rollback
   - Identify any gaps in the rollback plan
   - Assess the time taken for each step
   - Review communication effectiveness

3. **Update Procedures**
   - Revise this rollback plan based on findings
   - Improve migration scripts to prevent recurrence
   - Enhance monitoring and alerting
   - Update training materials

4. **Share Knowledge**
   - Conduct a post-mortem meeting
   - Share findings with the team
   - Update documentation
   - Implement process improvements

## Conclusion

This rollback plan provides a comprehensive approach to reverting the media migration if issues arise. By following these procedures, the team can quickly and safely restore the system to its original state while minimizing downtime and data loss.

Regular testing and updating of this plan will ensure its effectiveness when needed. All team members should be familiar with these procedures and their roles in the rollback process.