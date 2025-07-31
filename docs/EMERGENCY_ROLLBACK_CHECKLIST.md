# Emergency Media Migration Rollback Checklist

This checklist provides a quick reference for rolling back the media migration in an emergency situation. For detailed instructions, see [MEDIA_MIGRATION_ROLLBACK_PLAN.md](MEDIA_MIGRATION_ROLLBACK_PLAN.md).

## Immediate Actions

- [ ] **Alert the team**: Notify development lead, DevOps, and product manager
- [ ] **Stop the application**: Prevent new data from being written during rollback
- [ ] **Assess the situation**: Determine the severity and impact of the issue

## Rollback Procedure

### Option 1: Built-in Rollback (Preferred)

- [ ] **Run rollback script**:
  - Linux/macOS: `./scripts/media_migration.sh --rollback`
  - Windows: `scripts\media_migration.bat --rollback`
  - Direct: `go run cmd/media_migration/main.go --rollback`
- [ ] **Verify rollback success**:
  - Linux/macOS: `./scripts/verify_media_migration.sh`
  - Windows: `scripts\verify_media_migration.bat`
- [ ] **Check verification results**:
  - [ ] Media table no longer exists
  - [ ] Images and docs tables are intact
  - [ ] All original data is present

### Option 2: Backup Restoration (If built-in rollback fails)

- [ ] **List available backups**:
  - Linux/macOS: `./scripts/db_backup.sh list`
  - Windows: `scripts\db_backup.bat list`
- [ ] **Identify pre-migration backup**: Select the most recent backup before migration
- [ ] **Restore from backup**:
  - Linux/macOS: `./scripts/db_backup.sh restore -backup /path/to/backup -force`
  - Windows: `scripts\db_backup.bat restore -backup \path\to\backup -force`
- [ ] **Verify restoration success**:
  - Linux/macOS: `./scripts/verify_media_migration.sh`
  - Windows: `scripts\verify_media_migration.bat`

### Option 3: Manual Restoration (If both above fail)

- [ ] **Stop the application** (if not already stopped)
- [ ] **Create emergency backup**: `./scripts/db_backup.sh create -output /path/to/emergency_backup`
- [ ] **Replace database file**: `cp /path/to/pre_migration_backup.db /path/to/current/database.db`
- [ ] **Verify database state**:
  ```bash
  sqlite3 database.db ".tables"
  sqlite3 database.db "SELECT COUNT(*) FROM images;"
  sqlite3 database.db "SELECT COUNT(*) FROM docs;"
  ```

## Post-Rollback Actions

- [ ] **Restart the application**:
  - Service: `sudo systemctl start go-fast-cdn` (Linux) or `net start go-fast-cdn` (Windows)
  - Terminal: `go run main.go`
- [ ] **Verify application functionality**:
  - [ ] Image upload and retrieval works
  - [ ] Document upload and retrieval works
  - [ ] All media-related API endpoints respond correctly
  - [ ] User interface displays media correctly
- [ ] **Monitor system**:
  - [ ] Check application logs for errors
  - [ ] Monitor system performance
  - [ ] Verify database integrity

## Communication

- [ ] **Notify stakeholders**:
  - [ ] Development team
  - [ ] Quality assurance team
  - [ ] Customer support team
  - [ ] End users (if necessary)
- [ ] **Document the incident**:
  - [ ] Record timeline of events
  - [ ] Document root cause (if determined)
  - [ ] Capture all error messages and logs
  - [ ] Note resolution steps taken

## Follow-up Tasks

- [ ] **Investigate root cause**:
  - [ ] Analyze logs and error messages
  - [ ] Examine failed migration data
  - [ ] Review migration scripts for issues
- [ ] **Plan for re-migration**:
  - [ ] Schedule new migration window
  - [ ] Update migration plan based on lessons learned
  - [ ] Increase testing and validation
- [ ] **Update documentation**:
  - [ ] Update rollback plan with lessons learned
  - [ ] Document root cause and resolution
  - [ ] Share findings with the team

## Emergency Contacts

| Role | Name | Contact | Availability |
|------|------|---------|-------------|
| Development Lead | [Name] | [Email/Phone] | 24/7 |
| DevOps Lead | [Name] | [Email/Phone] | 24/7 |
| Product Manager | [Name] | [Email/Phone] | Business Hours |
| Stakeholder | [Name] | [Email/Phone] | Business Hours |

## Quick Commands Reference

### Linux/macOS

```bash
# Stop application
sudo systemctl stop go-fast-cdn

# Built-in rollback
./scripts/media_migration.sh --rollback

# Verify rollback
./scripts/verify_media_migration.sh

# List backups
./scripts/db_backup.sh list

# Restore from backup
./scripts/db_backup.sh restore -backup /path/to/backup -force

# Start application
sudo systemctl start go-fast-cdn
```

### Windows

```cmd
REM Stop application
net stop go-fast-cdn

REM Built-in rollback
scripts\media_migration.bat --rollback

REM Verify rollback
scripts\verify_media_migration.bat

REM List backups
scripts\db_backup.bat list

REM Restore from backup
scripts\db_backup.bat restore -backup \path\to\backup -force

REM Start application
net start go-fast-cdn
```

## Emergency Scripts

- **Linux/macOS**: [`scripts/emergency_rollback.sh`](../scripts/emergency_rollback.sh)
- **Windows**: [`scripts/emergency_rollback.bat`](../scripts/emergency_rollback.bat)

These scripts provide an interactive, step-by-step guide for performing an emergency rollback.

---

**Remember**: This checklist is for emergency use only. Always refer to the detailed [MEDIA_MIGRATION_ROLLBACK_PLAN.md](MEDIA_MIGRATION_ROLLBACK_PLAN.md) for comprehensive instructions and information.