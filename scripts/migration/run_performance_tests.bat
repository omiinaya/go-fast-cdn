@echo off
REM Performance Test Runner Script for Unified Media Repository (Windows)
REM This script runs the comprehensive performance tests for the unified media repository

echo ==========================================
echo Unified Media Repository Performance Test Runner
echo ==========================================

REM Set up variables
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..
set PERFORMANCE_DIR=%PROJECT_ROOT%\cmd\performance_tests
set REPORTS_DIR=%PROJECT_ROOT%\performance-reports

REM Create reports directory if it doesn't exist
if not exist "%REPORTS_DIR%" mkdir "%REPORTS_DIR%"

REM Function to print section headers
:print_section
echo.
echo ==========================================
echo %~1
echo ==========================================
goto :eof

REM Function to check if command exists
:command_exists
where %1 >nul 2>&1
exit /b %errorlevel%

REM Check prerequisites
call :print_section "Checking Prerequisites"

call :command_exists go
if %errorlevel% neq 0 (
    echo ❌ Go is not installed. Please install Go to continue.
    exit /b 1
)

echo ✅ Go is installed: 
go version

REM Check if the performance test directory exists
if not exist "%PERFORMANCE_DIR%" (
    echo ❌ Performance test directory not found: %PERFORMANCE_DIR%
    exit /b 1
)

REM Change to the performance test directory
cd /d "%PERFORMANCE_DIR%"

REM Run the performance tests
call :print_section "Running Performance Tests"

echo Running performance tests...
go run main.go

REM Check if the performance tests ran successfully
if %errorlevel% equ 0 (
    echo ✅ Performance tests completed successfully
) else (
    echo ❌ Performance tests failed
    exit /b 1
)

REM Move the generated reports to the reports directory
call :print_section "Organizing Reports"

if exist "performance-report.md" (
    move performance-report.md "%REPORTS_DIR%\" >nul
    echo ✅ Performance report moved to: %REPORTS_DIR%\performance-report.md
)

if exist "performance-results.json" (
    move performance-results.json "%REPORTS_DIR%\" >nul
    echo ✅ Performance results moved to: %REPORTS_DIR%\performance-results.json
)

if exist "cpu.prof" (
    move cpu.prof "%REPORTS_DIR%\" >nul
    echo ✅ CPU profile moved to: %REPORTS_DIR%\cpu.prof
)

REM Display final status
call :print_section "Performance Testing Complete"

echo 🎉 Performance testing completed successfully!
echo 📊 Reports are available in: %REPORTS_DIR%\
echo 📄 View the performance report: %REPORTS_DIR%\performance-report.md
echo 📈 View the detailed results: %REPORTS_DIR%\performance-results.json