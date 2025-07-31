@echo off
echo =====================================
echo Media Migration Verification Script
echo =====================================

REM Change to the project root directory
cd /d "%~dp0\.."

REM Run the media migration verification
echo Running media migration verification...
go run cmd/verify_media_migration/main.go

REM Check the exit code
if %ERRORLEVEL% EQU 0 (
    echo.
    echo Verification completed successfully!
    exit /b 0
) else (
    echo.
    echo Verification failed!
    exit /b 1
)