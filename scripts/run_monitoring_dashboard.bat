@echo off
REM Monitoring Dashboard Runner Script for Unified Media Repository (Windows)
REM This script starts the monitoring dashboard for visualizing post-deployment monitoring data

echo ==========================================
echo Unified Media Repository Monitoring Dashboard Runner
echo ==========================================

REM Set up variables
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..
set DASHBOARD_DIR=%PROJECT_ROOT%\cmd\monitoring_dashboard
set REPORTS_DIR=%PROJECT_ROOT%\post-deployment-monitoring-reports

REM Function to print section headers
:print_section
echo.
echo ==========================================
echo %~1
echo ==========================================
goto :eof

REM Check prerequisites
call :print_section "Checking Prerequisites"

where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go to continue.
    pause
    exit /b 1
)

echo ✅ Go is installed:
go version

REM Check if the dashboard directory exists
if not exist "%DASHBOARD_DIR%" (
    echo ❌ Dashboard directory not found: %DASHBOARD_DIR%
    pause
    exit /b 1
)

REM Create reports directory if it doesn't exist
if not exist "%REPORTS_DIR%" mkdir "%REPORTS_DIR%"

REM Check if monitoring data exists
if not exist "%REPORTS_DIR%\post-deployment-monitoring-results.json" (
    echo ⚠️  Monitoring data not found. The dashboard will show sample data.
    echo    Run the post-deployment monitoring first to generate real data:
    echo    scripts\run_post_deployment_monitoring.bat
    echo.
)

REM Change to the dashboard directory
cd /d "%DASHBOARD_DIR%"

REM Start the dashboard server
call :print_section "Starting Monitoring Dashboard"

echo Starting monitoring dashboard server...
echo Monitoring reports directory: %REPORTS_DIR%
echo.

REM Set environment variables for the dashboard
set MONITORING_REPORTS_DIR=%REPORTS_DIR%
set MONITORING_DASHBOARD_PORT=:8080

REM Start the dashboard in the background
start "Monitoring Dashboard" go run main.go

REM Wait for the dashboard to start
echo Waiting for the dashboard to start...
timeout /t 3 /nobreak >nul

REM Check if the dashboard is running
curl -s http://localhost:8080 >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Dashboard is running
    echo.
    echo 🌐 Access the dashboard at: http://localhost:8080
    echo.
    echo Close this window to keep the dashboard running in the background
    echo or press any key to stop the dashboard server...
    echo.
    
    REM Open the dashboard in the default browser
    start "" http://localhost:8080
    
    REM Wait for user input
    pause >nul
    
    REM Stop the dashboard
    echo.
    echo Stopping dashboard server...
    taskkill /FI "WINDOWTITLE eq Monitoring Dashboard*" /F >nul 2>&1
    taskkill /IM go.exe /F >nul 2>&1
) else (
    echo ❌ Failed to start the dashboard
    taskkill /FI "WINDOWTITLE eq Monitoring Dashboard*" /F >nul 2>&1
    taskkill /IM go.exe /F >nul 2>&1
    pause
    exit /b 1
)

pause