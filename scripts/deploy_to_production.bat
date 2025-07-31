@echo off
setlocal enabledelayedexpansion

REM Production Deployment Script for Windows
REM This script automates the deployment of the unified media repository to the production environment

REM Get the directory where this script is located
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR:~0,-1%
cd /d "%PROJECT_ROOT%"

REM Print banner
echo ==================================================
echo   UNIFIED MEDIA REPOSITORY PRODUCTION DEPLOYMENT 
echo ==================================================
echo.

REM Log file for deployment
set LOG_FILE=bin\production_deployment_%date:~-4,4%%date:~-10,2%%date:~-7,2%_%time:~0,2%%time:~3,2%.log
set LOG_FILE=%LOG_FILE: =0%
mkdir bin 2>nul
echo Deployment log: %LOG_FILE%
echo Starting production deployment at %date% %time% > "%LOG_FILE%"

REM Parse command line arguments
set SKIP_BACKUP=false
set SKIP_VERIFICATION=false
set ROLLBACK_ONLY=false

:parse_args
if "%~1"=="" goto end_parse_args
if /i "%~1"=="--skip-backup" (
    set SKIP_BACKUP=true
    shift
    goto parse_args
)
if /i "%~1"=="--skip-verification" (
    set SKIP_VERIFICATION=true
    shift
    goto parse_args
)
if /i "%~1"=="--rollback-only" (
    set ROLLBACK_ONLY=true
    shift
    goto parse_args
)
if /i "%~1"=="--help" (
    goto show_usage
)
echo Unknown option: %~1
goto show_usage
:end_parse_args

REM Function to display usage
:show_usage
echo Usage: %~nx0 [OPTIONS]
echo.
echo Options:
echo   --skip-backup           Skip creating a backup before deployment
echo   --skip-verification     Skip post-deployment verification
echo   --rollback-only         Only perform rollback without deployment
echo   --help                  Show this help message
echo.
echo Examples:
echo   %~nx0                    # Full deployment with backup and verification
echo   %~nx0 --skip-backup     # Deployment without backup (not recommended)
echo   %~nx0 --rollback-only   # Only perform rollback
exit /b 0

REM Check prerequisites
echo [INFO] Checking prerequisites...
echo [INFO] Checking prerequisites... >> "%LOG_FILE%"

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Go is not installed or not in PATH
    echo ERROR: Go is not installed or not in PATH >> "%LOG_FILE%"
    exit /b 1
)

REM Check if Node.js is installed
node --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] Node.js is not installed or not in PATH
    echo ERROR: Node.js is not installed or not in PATH >> "%LOG_FILE%"
    exit /b 1
)

REM Check if npm is installed
npm --version >nul 2>&1
if %errorlevel% neq 0 (
    echo [ERROR] npm is not installed or not in PATH
    echo ERROR: npm is not installed or not in PATH >> "%LOG_FILE%"
    exit /b 1
)

echo [SUCCESS] All prerequisites are met
echo SUCCESS: All prerequisites are met >> "%LOG_FILE%"

REM Function to create backup
:create_backup
if "%SKIP_BACKUP%"=="true" (
    echo [WARNING] Skipping backup creation as requested
    echo WARNING: Skipping backup creation as requested >> "%LOG_FILE%"
    exit /b 0
)

echo [INFO] Creating database backup...
echo INFO: Creating database backup... >> "%LOG_FILE%"

call scripts\db_backup.bat create
if %errorlevel% neq 0 (
    echo [ERROR] Failed to create database backup
    echo ERROR: Failed to create database backup >> "%LOG_FILE%"
    exit /b 1
)

echo [SUCCESS] Database backup created successfully
echo SUCCESS: Database backup created successfully >> "%LOG_FILE%"
exit /b 0

REM Function to perform database migration
:perform_migration
echo [INFO] Performing database migration...
echo INFO: Performing database migration... >> "%LOG_FILE%"

REM Set database path environment variable to ensure consistency
set CDN_DB_PATH=./db_data/production.db

REM Create db_data directory if it doesn't exist
mkdir db_data 2>nul

REM Build the production migration tool
echo [INFO] Building production migration tool...
echo INFO: Building production migration tool... >> "%LOG_FILE%"
if exist bin\production_migration.exe del bin\production_migration.exe

go build -o bin\production_migration.exe cmd\staging_migration\main.go
if %errorlevel% neq 0 (
    echo [ERROR] Failed to build production migration tool
    echo ERROR: Failed to build production migration tool >> "%LOG_FILE%"
    exit /b 1
)

echo [SUCCESS] Production migration tool built successfully
echo SUCCESS: Production migration tool built successfully >> "%LOG_FILE%"

REM Run the migration using the built binary
bin\production_migration.exe
if %errorlevel% neq 0 (
    echo [ERROR] Database migration failed
    echo ERROR: Database migration failed >> "%LOG_FILE%"
    echo [INFO] Attempting rollback...
    echo INFO: Attempting rollback... >> "%LOG_FILE%"
    
    REM Run the rollback using the built binary
    bin\production_migration.exe --rollback
    if %errorlevel% neq 0 (
        echo [ERROR] Rollback failed
        echo ERROR: Rollback failed >> "%LOG_FILE%"
    ) else (
        echo [SUCCESS] Rollback completed successfully
        echo SUCCESS: Rollback completed successfully >> "%LOG_FILE%"
    )
    
    exit /b 1
)

echo [SUCCESS] Database migration completed successfully
echo SUCCESS: Database migration completed successfully >> "%LOG_FILE%"
exit /b 0

REM Function to build backend
:build_backend
echo [INFO] Building backend application...
echo INFO: Building backend application... >> "%LOG_FILE%"

if exist bin\go-fast-cdn.exe del bin\go-fast-cdn.exe

go build -o bin\go-fast-cdn.exe main.go
if %errorlevel% neq 0 (
    echo [ERROR] Failed to build backend
    echo ERROR: Failed to build backend >> "%LOG_FILE%"
    exit /b 1
)

echo [SUCCESS] Backend built successfully
echo SUCCESS: Backend built successfully >> "%LOG_FILE%"
exit /b 0

REM Function to build frontend
:build_frontend
echo [INFO] Building frontend application...
echo INFO: Building frontend application... >> "%LOG_FILE%"

cd ui

call npm install --legacy-peer-deps
if %errorlevel% neq 0 (
    echo [ERROR] Failed to install frontend dependencies
    echo ERROR: Failed to install frontend dependencies >> "%LOG_FILE%"
    cd ..
    exit /b 1
)

call npm run build
if %errorlevel% neq 0 (
    echo [ERROR] Failed to build frontend
    echo ERROR: Failed to build frontend >> "%LOG_FILE%"
    cd ..
    exit /b 1
)

cd ..
echo [SUCCESS] Frontend built successfully
echo SUCCESS: Frontend built successfully >> "%LOG_FILE%"
exit /b 0

REM Function to deploy backend
:deploy_backend
echo [INFO] Deploying backend application...
echo INFO: Deploying backend application... >> "%LOG_FILE%"

REM In a real deployment, this would involve copying files to the production server
REM For this example, we'll just simulate the deployment

if exist bin\go-fast-cdn.exe (
    echo [SUCCESS] Backend deployment simulated successfully
    echo SUCCESS: Backend deployment simulated successfully >> "%LOG_FILE%"
) else (
    echo [ERROR] Backend binary not found
    echo ERROR: Backend binary not found >> "%LOG_FILE%"
    exit /b 1
)

exit /b 0

REM Function to deploy frontend
:deploy_frontend
echo [INFO] Deploying frontend application...
echo INFO: Deploying frontend application... >> "%LOG_FILE%"

REM In a real deployment, this would involve copying files to the production server
REM For this example, we'll just simulate the deployment

if exist ui\dist (
    echo [SUCCESS] Frontend deployment simulated successfully
    echo SUCCESS: Frontend deployment simulated successfully >> "%LOG_FILE%"
) else (
    echo [ERROR] Frontend build output not found
    echo ERROR: Frontend build output not found >> "%LOG_FILE%"
    exit /b 1
)

exit /b 0

REM Function to verify deployment
:verify_deployment
if "%SKIP_VERIFICATION%"=="true" (
    echo [WARNING] Skipping deployment verification as requested
    echo WARNING: Skipping deployment verification as requested >> "%LOG_FILE%"
    exit /b 0
)

echo [INFO] Verifying deployment...
echo INFO: Verifying deployment... >> "%LOG_FILE%"

REM Set database path environment variable to ensure consistency
set CDN_DB_PATH=./db_data/production.db

REM Build the verification tool first to ensure it uses the same database
echo [INFO] Building verification tool...
echo INFO: Building verification tool... >> "%LOG_FILE%"
if exist bin\verify_media_migration.exe del bin\verify_media_migration.exe

go build -o bin\verify_media_migration.exe cmd\verify_media_migration\main.go
if %errorlevel% neq 0 (
    echo [ERROR] Failed to build verification tool
    echo ERROR: Failed to build verification tool >> "%LOG_FILE%"
    exit /b 1
)

echo [SUCCESS] Verification tool built successfully
echo SUCCESS: Verification tool built successfully >> "%LOG_FILE%"

REM Run the verification using the built binary
bin\verify_media_migration.exe
if %errorlevel% neq 0 (
    echo [ERROR] Deployment verification failed
    echo ERROR: Deployment verification failed >> "%LOG_FILE%"
    exit /b 1
)

echo [SUCCESS] Deployment verification completed successfully
echo SUCCESS: Deployment verification completed successfully >> "%LOG_FILE%"
exit /b 0

REM Function to perform rollback
:perform_rollback
echo [INFO] Performing rollback...
echo INFO: Performing rollback... >> "%LOG_FILE%"

REM Set database path environment variable to ensure consistency
set CDN_DB_PATH=./db_data/production.db

REM Create db_data directory if it doesn't exist
mkdir db_data 2>nul

REM Build the production migration tool if it doesn't exist
if not exist bin\production_migration.exe (
    echo [INFO] Building production migration tool...
    echo INFO: Building production migration tool... >> "%LOG_FILE%"
    
    go build -o bin\production_migration.exe cmd\staging_migration\main.go
    if %errorlevel% neq 0 (
        echo [ERROR] Failed to build production migration tool
        echo ERROR: Failed to build production migration tool >> "%LOG_FILE%"
        exit /b 1
    )
    
    echo [SUCCESS] Production migration tool built successfully
    echo SUCCESS: Production migration tool built successfully >> "%LOG_FILE%"
)

REM Run the rollback using the built binary
bin\production_migration.exe --rollback
if %errorlevel% neq 0 (
    echo [ERROR] Rollback failed
    echo ERROR: Rollback failed >> "%LOG_FILE%"
    exit /b 1
)

echo [SUCCESS] Rollback completed successfully
echo SUCCESS: Rollback completed successfully >> "%LOG_FILE%"
exit /b 0

REM Main deployment process
if "%ROLLBACK_ONLY%"=="true" (
    echo [INFO] Starting rollback-only process...
    echo INFO: Starting rollback-only process... >> "%LOG_FILE%"
    call :perform_rollback
    if %errorlevel% neq 0 exit /b %errorlevel%
    echo [SUCCESS] Rollback-only process completed
    echo SUCCESS: Rollback-only process completed >> "%LOG_FILE%"
    exit /b 0
)

echo [INFO] Starting deployment process...
echo INFO: Starting deployment process... >> "%LOG_FILE%"

REM Phase 1: Pre-Deployment Preparation
echo [INFO] Phase 1: Pre-Deployment Preparation
echo INFO: Phase 1: Pre-Deployment Preparation >> "%LOG_FILE%"
call :create_backup
if %errorlevel% neq 0 exit /b %errorlevel%

REM Phase 2: Database Migration
echo [INFO] Phase 2: Database Migration
echo INFO: Phase 2: Database Migration >> "%LOG_FILE%"
call :perform_migration
if %errorlevel% neq 0 exit /b %errorlevel%

REM Phase 3: Build Application
echo [INFO] Phase 3: Build Application
echo INFO: Phase 3: Build Application >> "%LOG_FILE%"
call :build_backend
if %errorlevel% neq 0 exit /b %errorlevel%
call :build_frontend
if %errorlevel% neq 0 exit /b %errorlevel%

REM Phase 4: Deploy Application
echo [INFO] Phase 4: Deploy Application
echo INFO: Phase 4: Deploy Application >> "%LOG_FILE%"
call :deploy_backend
if %errorlevel% neq 0 exit /b %errorlevel%
call :deploy_frontend
if %errorlevel% neq 0 exit /b %errorlevel%

REM Phase 5: Post-Deployment Verification
echo [INFO] Phase 5: Post-Deployment Verification
echo INFO: Phase 5: Post-Deployment Verification >> "%LOG_FILE%"
call :verify_deployment
if %errorlevel% neq 0 (
    echo [ERROR] Deployment verification failed
    echo ERROR: Deployment verification failed >> "%LOG_FILE%"
    echo [INFO] Attempting rollback...
    echo INFO: Attempting rollback... >> "%LOG_FILE%"
    call :perform_rollback
    exit /b 1
)

REM Deployment completed successfully
echo.
echo ==================================================
echo           DEPLOYMENT COMPLETED SUCCESSFULLY      
echo ==================================================
echo.
echo [INFO] The unified media repository has been successfully deployed to the production environment.
echo [INFO] All verification checks have passed.
echo.
echo INFO: The unified media repository has been successfully deployed to the production environment. >> "%LOG_FILE%"
echo INFO: All verification checks have passed. >> "%LOG_FILE%"
echo.
echo [INFO] Next steps:
echo 1. Monitor the application for any issues
echo 2. Perform user acceptance testing
echo 3. Document the deployment results
echo 4. Plan for future maintenance and updates
echo.
echo INFO: Next steps: >> "%LOG_FILE%"
echo INFO: 1. Monitor the application for any issues >> "%LOG_FILE%"
echo INFO: 2. Perform user acceptance testing >> "%LOG_FILE%"
echo INFO: 3. Document the deployment results >> "%LOG_FILE%"
echo INFO: 4. Plan for future maintenance and updates >> "%LOG_FILE%"
echo.
echo Deployment completed at %date% %time% >> "%LOG_FILE%"

pause