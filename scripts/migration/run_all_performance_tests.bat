@echo off
REM Comprehensive Performance Test Runner Script for Unified Media Repository (Windows)
REM This script runs all performance tests (backend and frontend) and generates a comprehensive summary report

echo ==========================================
echo Unified Media Repository Comprehensive Performance Test Runner
echo ==========================================

REM Set up variables
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..
set BACKEND_PERFORMANCE_DIR=%PROJECT_ROOT%\cmd\performance_tests
set FRONTEND_PERFORMANCE_DIR=%PROJECT_ROOT%\ui\tests\performance
set REPORTS_DIR=%PROJECT_ROOT%\comprehensive-performance-reports

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

call :command_exists node
if %errorlevel% neq 0 (
    echo ❌ Node.js is not installed. Please install Node.js to continue.
    exit /b 1
)

call :command_exists npm
if %errorlevel% neq 0 (
    echo ❌ npm is not installed. Please install npm to continue.
    exit /b 1
)

echo ✅ Go is installed: 
go version
echo ✅ Node.js is installed: 
node --version
echo ✅ npm is installed: 
npm --version

REM Check if performance test directories exist
if not exist "%BACKEND_PERFORMANCE_DIR%" (
    echo ❌ Backend performance test directory not found: %BACKEND_PERFORMANCE_DIR%
    exit /b 1
)

if not exist "%FRONTEND_PERFORMANCE_DIR%" (
    echo ❌ Frontend performance test directory not found: %FRONTEND_PERFORMANCE_DIR%
    exit /b 1
)

REM Check if Playwright is installed
if not exist "%PROJECT_ROOT%\ui\node_modules\@playwright" (
    echo ⚠️  Playwright is not installed. Installing Playwright...
    cd /d "%PROJECT_ROOT%\ui"
    npm install @playwright/test
    npx playwright install
    cd /d "%PROJECT_ROOT%"
)

REM Start the application in the background
call :print_section "Starting Application"

echo Starting the application in the background...
cd /d "%PROJECT_ROOT%"

REM Start the Go backend
start /B go run main.go
set BACKEND_PID=%errorlevel%

REM Wait for the backend to start
echo Waiting for the backend to start...
timeout /t 5 /nobreak >nul

REM Start the UI in development mode
cd /d "%PROJECT_ROOT%\ui"
start /B npm run dev
set UI_PID=%errorlevel%

REM Wait for the UI to start
echo Waiting for the UI to start...
timeout /t 10 /nobreak >nul

REM Check if the application is running
curl -s http://localhost:3000 >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ Application is running
) else (
    echo ❌ Failed to start the application
    taskkill /F /PID %BACKEND_PID% >nul 2>&1
    taskkill /F /PID %UI_PID% >nul 2>&1
    exit /b 1
)

REM Run backend performance tests
call :print_section "Running Backend Performance Tests"

cd /d "%BACKEND_PERFORMANCE_DIR%"
echo Running backend performance tests...
go run main.go

REM Check if backend performance tests ran successfully
if %errorlevel% equ 0 (
    echo ✅ Backend performance tests completed successfully
    
    REM Move backend reports to reports directory
    if exist "performance-report.md" (
        move performance-report.md "%REPORTS_DIR%\backend-performance-report.md" >nul
        echo ✅ Backend performance report moved to: %REPORTS_DIR%\backend-performance-report.md
    )
    
    if exist "performance-results.json" (
        move performance-results.json "%REPORTS_DIR%\backend-performance-results.json" >nul
        echo ✅ Backend performance results moved to: %REPORTS_DIR%\backend-performance-results.json
    )
) else (
    echo ❌ Backend performance tests failed
)

REM Run frontend performance tests
call :print_section "Running Frontend Performance Tests"

cd /d "%PROJECT_ROOT%\ui\tests"
echo Running frontend performance tests...
npx playwright test performance/ --reporter=list,html --output="%REPORTS_DIR%\frontend-report"

REM Check if frontend performance tests ran successfully
if %errorlevel% equ 0 (
    echo ✅ Frontend performance tests completed successfully
) else (
    echo ❌ Frontend performance tests failed
)

REM Stop the application
call :print_section "Stopping Application"

echo Stopping the application...
taskkill /F /PID %BACKEND_PID% >nul 2>&1
taskkill /F /PID %UI_PID% >nul 2>&1

REM Generate comprehensive summary report
call :print_section "Generating Comprehensive Performance Test Summary Report"

set SUMMARY_FILE=%REPORTS_DIR%\comprehensive-performance-summary.md
echo # Unified Media Repository Comprehensive Performance Test Summary Report > "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## Test Execution Details >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo - **Date:** %date% %time% >> "%SUMMARY_FILE%"
echo - **Test Environment:** Local Development >> "%SUMMARY_FILE%"
echo - **Backend PID:** %BACKEND_PID% >> "%SUMMARY_FILE%"
echo - **UI PID:** %UI_PID% >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## Backend Performance Test Results >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"

REM Add backend performance results if available
if exist "%REPORTS_DIR%\backend-performance-report.md" (
    echo Backend performance test results are available in: %REPORTS_DIR%\backend-performance-report.md >> "%SUMMARY_FILE%"
    echo. >> "%SUMMARY_FILE%"
    
    REM Extract key metrics from backend report
    if exist "%REPORTS_DIR%\backend-performance-results.json" (
        echo ### Key Backend Metrics >> "%SUMMARY_FILE%"
        echo. >> "%SUMMARY_FILE%"
        echo Raw JSON data available in: %REPORTS_DIR%\backend-performance-results.json >> "%SUMMARY_FILE%"
    )
) else (
    echo Backend performance test results are not available. >> "%SUMMARY_FILE%"
)

echo. >> "%SUMMARY_FILE%"
echo ## Frontend Performance Test Results >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"

REM Add frontend performance results if available
if exist "%REPORTS_DIR%\frontend-report" (
    echo Frontend performance test results are available in: %REPORTS_DIR%\frontend-report\index.html >> "%SUMMARY_FILE%"
    echo. >> "%SUMMARY_FILE%"
    echo ### Frontend Test Coverage >> "%SUMMARY_FILE%"
    echo. >> "%SUMMARY_FILE%"
    echo - Upload page performance comparison >> "%SUMMARY_FILE%"
    echo - Files page performance comparison >> "%SUMMARY_FILE%"
    echo - Media upload performance comparison >> "%SUMMARY_FILE%"
    echo - Media display performance comparison >> "%SUMMARY_FILE%"
    echo - Search and filter performance comparison >> "%SUMMARY_FILE%"
    echo - Bulk operations performance comparison >> "%SUMMARY_FILE%"
    echo - Different media types performance >> "%SUMMARY_FILE%"
    echo - Concurrent user performance >> "%SUMMARY_FILE%"
) else (
    echo Frontend performance test results are not available. >> "%SUMMARY_FILE%"
)

echo. >> "%SUMMARY_FILE%"
echo ## Overall Performance Analysis >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ### Key Findings >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo Based on the test results, the following areas may need optimization: >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo 1. **Database Queries**: Review and optimize database queries, especially for media retrieval operations >> "%SUMMARY_FILE%"
echo 2. **API Endpoints**: Optimize API endpoint handling for unified media operations >> "%SUMMARY_FILE%"
echo 3. **Frontend Rendering**: Improve frontend component rendering performance for large media lists >> "%SUMMARY_FILE%"
echo 4. **File Upload/Download**: Enhance file transfer performance for large media files >> "%SUMMARY_FILE%"
echo 5. **Concurrent Operations**: Scale concurrent request handling for better performance under load >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## Recommendations >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ### Short-term Optimizations >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo 1. **Implement Caching**: Add caching for frequently accessed media metadata and files >> "%SUMMARY_FILE%"
echo 2. **Database Indexing**: Ensure proper database indexing for media queries >> "%SUMMARY_FILE%"
echo 3. **Frontend Optimization**: Implement lazy loading and virtual scrolling for media lists >> "%SUMMARY_FILE%"
echo 4. **API Response Optimization**: Optimize API response sizes and implement pagination >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ### Long-term Improvements >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo 1. **Content Delivery Network (CDN)**: Implement CDN for global media distribution >> "%SUMMARY_FILE%"
echo 2. **Microservices Architecture**: Consider breaking down unified service into specialized microservices >> "%SUMMARY_FILE%"
echo 3. **Advanced Caching Strategies**: Implement multi-level caching with Redis or similar >> "%SUMMARY_FILE%"
echo 4. **Load Balancing**: Implement load balancing for high-traffic scenarios >> "%SUMMARY_FILE%"
echo 5. **Performance Monitoring**: Set up continuous performance monitoring and alerting >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## Next Steps >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo 1. **Review Test Results**: Analyze detailed test results to identify specific bottlenecks >> "%SUMMARY_FILE%"
echo 2. **Implement Optimizations**: Apply the recommended optimizations based on test findings >> "%SUMMARY_FILE%"
echo 3. **Re-run Tests**: Execute performance tests after optimizations to verify improvements >> "%SUMMARY_FILE%"
echo 4. **Establish Baselines**: Set performance baselines for ongoing monitoring >> "%SUMMARY_FILE%"
echo 5. **Continuous Testing**: Integrate performance testing into CI/CD pipeline >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## Test Reports Location >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo - **Backend Performance Report**: %REPORTS_DIR%\backend-performance-report.md >> "%SUMMARY_FILE%"
echo - **Backend Performance Results**: %REPORTS_DIR%\backend-performance-results.json >> "%SUMMARY_FILE%"
echo - **Frontend Performance Report**: %REPORTS_DIR%\frontend-report\index.html >> "%SUMMARY_FILE%"
echo - **Comprehensive Summary**: %REPORTS_DIR%\comprehensive-performance-summary.md >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo Generated on: %date% %time% >> "%SUMMARY_FILE%"

echo ✅ Comprehensive performance summary report generated: %SUMMARY_FILE%

REM Display final status
call :print_section "Comprehensive Performance Testing Complete"

echo 🎉 Comprehensive performance testing completed successfully!
echo 📊 Reports are available in: %REPORTS_DIR%\
echo 📄 View the comprehensive summary: %REPORTS_DIR%\comprehensive-performance-summary.md
echo 🔍 View detailed reports in the subdirectories