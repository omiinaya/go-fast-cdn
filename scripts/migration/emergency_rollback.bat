@echo off
setlocal enabledelayedexpansion

REM Emergency Media Migration Rollback Script
REM This script provides a quick reference for rolling back the media migration in an emergency situation.

echo ================================================================
echo EMERGENCY MEDIA MIGRATION ROLLBACK SCRIPT
echo ================================================================
echo.

REM Get the directory where this script is located
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR:~0,-1%
cd /d "%PROJECT_ROOT%"

echo This script will guide you through the emergency rollback process.
echo Please follow the instructions carefully.
echo.

REM Step 1: Stop the application
echo STEP 1: STOP THE APPLICATION
echo ============================
echo Before proceeding, you must stop any running instances of the application.
echo.
echo Common ways to stop the application:
echo - If running as a service: net stop go-fast-cdn
echo - If running in a terminal: Press Ctrl+C
echo - If running as a background process: taskkill /f /im "go.exe"
echo.
set /p "STOPPED=Have you stopped the application? (y/N): "
if /i not "%STOPPED%"=="y" (
    echo Please stop the application before continuing.
    pause
    exit /b 1
)

REM Step 2: Choose rollback method
echo.
echo STEP 2: CHOOSE ROLLBACK METHOD
echo ==============================
echo Select the rollback method to use:
echo.
echo 1. Built-in rollback functionality (RECOMMENDED)
echo    - Uses the migration script's built-in rollback
echo    - Fast and safe if the migration script is working
echo.
echo 2. Restore from backup
echo    - Restores the database from a pre-migration backup
echo    - Use if built-in rollback fails
echo.
echo 3. Exit this script
echo.
set /p "CHOICE=Enter your choice (1-3): "

if "%CHOICE%"=="1" (
    echo You selected: Built-in rollback functionality
    echo.
    echo Running built-in rollback...
    echo.
    
    REM Run the rollback script
    call scripts\media_migration.bat --rollback
    if !errorlevel! neq 0 (
        echo.
        echo ✗ Rollback failed!
        echo.
        echo Please try restoring from backup instead.
        pause
        exit /b 1
    ) else (
        echo.
        echo ✓ Rollback completed successfully!
    )
) else if "%CHOICE%"=="2" (
    echo You selected: Restore from backup
    echo.
    
    REM List available backups
    echo Available backups:
    echo ==================
    call scripts\db_backup.bat list
    if !errorlevel! neq 0 (
        echo Failed to list backups. Please check the backup system.
        pause
        exit /b 1
    )
    
    echo.
    set /p "BACKUP_PATH=Enter the full path to the pre-migration backup file: "
    
    if "%BACKUP_PATH%"=="" (
        echo No backup path provided. Exiting.
        pause
        exit /b 1
    )
    
    if not exist "%BACKUP_PATH%" (
        echo Backup file not found: %BACKUP_PATH%
        pause
        exit /b 1
    )
    
    echo.
    echo WARNING: This will overwrite the current database with the backup.
    echo All data added after the backup was created will be lost.
    echo.
    set /p "CONFIRM=Continue with restore? (y/N): "
    
    if /i "%CONFIRM%"=="y" (
        echo Restoring from backup...
        echo.
        
        call scripts\db_backup.bat restore -backup "%BACKUP_PATH%" -force
        if !errorlevel! neq 0 (
            echo.
            echo ✗ Restore failed!
            pause
            exit /b 1
        ) else (
            echo.
            echo ✓ Restore completed successfully!
        )
    ) else (
        echo Restore cancelled.
        pause
        exit /b 0
    )
) else if "%CHOICE%"=="3" (
    echo Exiting script.
    pause
    exit /b 0
) else (
    echo Invalid choice. Please run the script again.
    pause
    exit /b 1
)

REM Step 3: Verify rollback
echo.
echo STEP 3: VERIFY ROLLBACK
echo ========================
echo Verifying that the rollback was successful...
echo.

call scripts\verify_media_migration.bat
if !errorlevel! neq 0 (
    echo.
    echo ✗ Verification failed!
    echo.
    echo The rollback may not have been completed successfully.
    echo Please check the verification output and take appropriate action.
    pause
    exit /b 1
) else (
    echo.
    echo ✓ Verification completed successfully!
    echo The rollback was successful and the system is in a consistent state.
)

REM Step 4: Restart the application
echo.
echo STEP 4: RESTART THE APPLICATION
echo ===============================
echo The rollback has been completed and verified.
echo You can now restart the application.
echo.
echo Common ways to start the application:
echo - As a service: net start go-fast-cdn
echo - In a terminal: go run main.go
echo - As a background process: start /b go run main.go
echo.
set /p "START=Do you want to start the application now? (y/N): "

if /i "%START%"=="y" (
    echo Starting the application...
    echo.
    
    REM Try to start the application
    start "Go Fast CDN" go run main.go
    if !errorlevel! neq 0 (
        echo.
        echo ✗ Failed to start the application!
        echo Please start it manually using the appropriate command for your setup.
    ) else (
        echo.
        echo ✓ Application started successfully!
        echo Check the application logs to ensure it's running correctly.
    )
) else (
    echo Please start the application manually when ready.
)

REM Completion message
echo.
echo ================================================================
echo EMERGENCY ROLLBACK COMPLETED
echo ================================================================
echo.
echo The media migration rollback has been completed.
echo The system has been restored to its pre-migration state.
echo.
echo Next steps:
echo 1. Monitor the application for any issues
echo 2. Investigate the cause of the rollback
echo 3. Notify stakeholders of the rollback
echo 4. Plan for re-migration after fixing the issues
echo.
echo For more detailed information, see the rollback plan documentation:
echo docs\MEDIA_MIGRATION_ROLLBACK_PLAN.md
echo.

pause