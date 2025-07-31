@echo off
REM Post-Deployment Monitoring Runner Script for Unified Media Repository (Windows)
REM This script runs the comprehensive post-deployment monitoring for the unified media repository

echo ==========================================
echo Unified Media Repository Post-Deployment Monitoring Runner
echo ==========================================

REM Set up variables
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..
set MONITORING_DIR=%PROJECT_ROOT%\cmd\post_deployment_monitoring
set REPORTS_DIR=%PROJECT_ROOT%\post-deployment-monitoring-reports

REM Create reports directory if it doesn't exist
if not exist "%REPORTS_DIR%" mkdir "%REPORTS_DIR%"

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

REM Check if the monitoring directory exists
if not exist "%MONITORING_DIR%" (
    echo ❌ Monitoring directory not found: %MONITORING_DIR%
    pause
    exit /b 1
)

REM Change to the monitoring directory
cd /d "%MONITORING_DIR%"

REM Run the monitoring tests
call :print_section "Running Post-Deployment Monitoring"

echo Running post-deployment monitoring...
go run main.go

REM Check if the monitoring ran successfully
if %errorlevel% neq 0 (
    echo ❌ Post-deployment monitoring failed
    pause
    exit /b 1
)

echo ✅ Post-deployment monitoring completed successfully

REM Move the generated reports to the reports directory
call :print_section "Organizing Reports"

if exist "post-deployment-monitoring-report.md" (
    move "post-deployment-monitoring-report.md" "%REPORTS_DIR%\" >nul
    echo ✅ Monitoring report moved to: %REPORTS_DIR%\post-deployment-monitoring-report.md
)

if exist "post-deployment-monitoring-results.json" (
    move "post-deployment-monitoring-results.json" "%REPORTS_DIR%\" >nul
    echo ✅ Monitoring results moved to: %REPORTS_DIR%\post-deployment-monitoring-results.json
)

if exist "monitoring-cpu.prof" (
    move "monitoring-cpu.prof" "%REPORTS_DIR%\" >nul
    echo ✅ CPU profile moved to: %REPORTS_DIR%\monitoring-cpu.prof
)

REM Display final status
call :print_section "Post-Deployment Monitoring Complete"

echo 🎉 Post-deployment monitoring completed successfully!
echo 📊 Reports are available in: %REPORTS_DIR%\
echo 📄 View the monitoring report: %REPORTS_DIR%\post-deployment-monitoring-report.md
echo 📈 View the detailed results: %REPORTS_DIR%\post-deployment-monitoring-results.json

REM Optional: Open the report
if exist "%REPORTS_DIR%\post-deployment-monitoring-report.md" (
    echo.
    set /p OPEN_REPORT="📖 Do you want to open the monitoring report? (y/n): "
    if /i "%OPEN_REPORT%"=="y" (
        start "" "%REPORTS_DIR%\post-deployment-monitoring-report.md"
    )
)

pause