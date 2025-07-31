@echo off
REM File Migration Script for Windows
REM This script migrates files from separate directories to a unified media directory

set ROLLBACK=false
set CLEANUP=false

:parse_args
if "%~1"=="" goto :run_migration
if "%~1"=="--rollback" (
    set ROLLBACK=true
    shift
    goto :parse_args
)
if "%~1"=="--cleanup" (
    set CLEANUP=true
    shift
    goto :parse_args
)
echo Unknown option: %~1
echo Usage: %~nx0 [--rollback] [--cleanup]
echo   --rollback  : Rollback file migration to legacy directories
echo   --cleanup   : Clean up legacy files after successful migration
exit /b 1

:run_migration
REM Build the file migration tool
echo Building file migration tool...
go build -o bin/file_migration.exe cmd/file_migration/main.go
if %errorlevel% neq 0 (
    echo Failed to build file migration tool
    exit /b 1
)

REM Run the file migration
if "%ROLLBACK%"=="true" (
    echo Rolling back file migration...
    bin\file_migration.exe --rollback
) else if "%CLEANUP%"=="true" (
    echo Cleaning up legacy files...
    bin\file_migration.exe --cleanup
) else (
    echo Running file migration to unified media directory...
    bin\file_migration.exe
)

REM Check if migration was successful
if %errorlevel% equ 0 (
    echo File migration operation completed successfully!
) else (
    echo File migration operation failed!
    exit /b 1
)