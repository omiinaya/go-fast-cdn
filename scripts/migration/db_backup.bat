@echo off
REM Database Backup Script for Windows
REM This script provides a convenient way to run the database backup tool

REM Get the directory where this script is located
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..

REM Build the backup tool if it doesn't exist or if source files are newer
if not exist "%PROJECT_ROOT%\bin\db_backup.exe" (
    echo Building database backup tool...
    cd /d "%PROJECT_ROOT%"
    go build -o bin/db_backup.exe cmd/db_backup/main.go
    if %errorlevel% neq 0 (
        echo Error: Failed to build database backup tool
        exit /b 1
    )
    echo Build completed successfully.
) else (
    REM Check if source is newer than the executable
    for /f %%i in ('dir /b /o:d "%PROJECT_ROOT%\cmd\db_backup\main.go" "%PROJECT_ROOT%\bin\db_backup.exe" ^| findstr /c:"main.go"') do (
        echo Source file is newer than executable, rebuilding...
        cd /d "%PROJECT_ROOT%"
        go build -o bin/db_backup.exe cmd/db_backup/main.go
        if %errorlevel% neq 0 (
            echo Error: Failed to build database backup tool
            exit /b 1
        )
        echo Build completed successfully.
    )
)

REM Run the backup tool with all provided arguments
echo Running database backup tool...
"%PROJECT_ROOT%\bin\db_backup.exe" %*