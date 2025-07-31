@echo off
REM Media Migration Script for Windows
REM This script runs the media unification migration to merge images and docs tables into a single media table

echo Starting media unification migration...

REM Check if rollback flag is provided
if "%1"=="--rollback" (
    echo Rolling back media unification migration...
    go run cmd/media_migration/main.go --rollback
) else (
    echo Running media unification migration...
    go run cmd/media_migration/main.go
)

REM Check the exit status
if %errorlevel% equ 0 (
    echo Media migration completed successfully!
) else (
    echo Media migration failed!
    exit /b 1
)