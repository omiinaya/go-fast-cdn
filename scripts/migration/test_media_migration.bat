@echo off
REM Media Migration Test Script (Windows)
REM This script creates a test database, populates it with sample data,
REM runs the media migration, verifies it, tests rollback, and cleans up.

echo Media Migration Test Suite
echo ==========================

REM Check if cleanup-only flag is provided
if "%1"=="--cleanup-only" (
    echo Running cleanup only...
    go run cmd/test_media_migration/main.go --cleanup-only
    exit /b %errorlevel%
)

REM Run the test suite
go run cmd/test_media_migration/main.go

REM Check the exit status
if %errorlevel% equ 0 (
    echo.
    echo All tests passed successfully!
    echo The migration script is ready for staging environment.
) else (
    echo.
    echo Tests failed! Please check the output above for details.
    exit /b 1
)