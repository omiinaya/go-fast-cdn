@echo off
REM Staging Migration Script for Windows
REM This script runs the media unification migration with backup and error handling for staging environments

echo === Staging Migration Script ===
echo This script will:
echo 1. Create a backup of the current database
echo 2. Run the media unification migration
echo 3. Verify the migration was successful
echo 4. Automatically rollback if the migration fails
echo.

REM Check if rollback flag is provided
if "%1"=="--rollback" (
    echo Mode: Rollback
    echo This will:
    echo 1. Create a backup of the current database
    echo 2. Rollback the media unification migration
    echo 3. Verify the rollback was successful
    echo 4. Automatically restore if the rollback fails
    echo.
    go run cmd/staging_migration/main.go --rollback
) else if "%1"=="--skip-backup" (
    echo Warning: Skipping backup creation as requested
    echo This is not recommended for production environments
    echo.
    go run cmd/staging_migration/main.go --skip-backup
) else if "%1"=="--rollback-skip-backup" (
    echo Warning: Skipping backup creation before rollback as requested
    echo This is not recommended for production environments
    echo.
    go run cmd/staging_migration/main.go --rollback --skip-backup
) else (
    echo Mode: Migration
    echo.
    go run cmd/staging_migration/main.go
)

REM Check the exit status
if %errorlevel% equ 0 (
    echo.
    echo === Staging migration completed successfully! ===
) else (
    echo.
    echo === Staging migration failed! ===
    echo Please check the error messages above for details.
    exit /b 1
)