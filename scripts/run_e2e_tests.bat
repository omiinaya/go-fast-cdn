@echo off
REM End-to-End Test Runner Script for Unified Media Repository (Windows)
REM This script runs the comprehensive E2E tests for the unified media repository

echo ==========================================
echo Unified Media Repository E2E Test Runner
echo ==========================================

REM Set up variables
set SCRIPT_DIR=%~dp0
set PROJECT_ROOT=%SCRIPT_DIR%..
set UI_DIR=%PROJECT_ROOT%\ui
set TESTS_DIR=%UI_DIR%\tests
set REPORTS_DIR=%PROJECT_ROOT%\test-reports

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

echo ✅ Node.js is installed: 
node --version
echo ✅ npm is installed: 
npm --version

REM Check if Playwright is installed
if not exist "%UI_DIR%\node_modules\@playwright" (
    echo ⚠️  Playwright is not installed. Installing Playwright...
    cd /d "%UI_DIR%"
    npm install @playwright/test
    npx playwright install
    cd /d "%PROJECT_ROOT%"
)

REM Create test fixture files if they don't exist
call :print_section "Setting Up Test Fixtures"

set FIXTURES_DIR=%TESTS_DIR%\fixtures
if not exist "%FIXTURES_DIR%" mkdir "%FIXTURES_DIR%"

REM Create minimal test files if they don't exist
if not exist "%FIXTURES_DIR%\test-image.jpg" (
    echo Creating test image fixture...
    REM Create a minimal 1x1 pixel JPEG file
    echo ÿØÿà^ZJFIF^Z^Z^ZH^ZH^ZÿÛ^C^Z^ZÿÙ > "%FIXTURES_DIR%\test-image.jpg"
)

if not exist "%FIXTURES_DIR%\test-document.pdf" (
    echo Creating test document fixture...
    REM Create a minimal PDF file
    echo %%PDF-1.4 > "%FIXTURES_DIR%\test-document.pdf"
    echo 1 0 obj >> "%FIXTURES_DIR%\test-document.pdf"
    echo ^<^< >> "%FIXTURES_DIR%\test-document.pdf"
    echo /Type /Catalog >> "%FIXTURES_DIR%\test-document.pdf"
    echo /Pages 2 0 R >> "%FIXTURES_DIR%\test-document.pdf"
    echo ^>^> >> "%FIXTURES_DIR%\test-document.pdf"
    echo endobj >> "%FIXTURES_DIR%\test-document.pdf"
    echo 2 0 obj >> "%FIXTURES_DIR%\test-document.pdf"
    echo ^<^< >> "%FIXTURES_DIR%\test-document.pdf"
    echo /Type /Pages >> "%FIXTURES_DIR%\test-document.pdf"
    echo /Kids [3 0 R] >> "%FIXTURES_DIR%\test-document.pdf"
    echo /Count 1 >> "%FIXTURES_DIR%\test-document.pdf"
    echo ^>^> >> "%FIXTURES_DIR%\test-document.pdf"
    echo endobj >> "%FIXTURES_DIR%\test-document.pdf"
    echo 3 0 obj >> "%FIXTURES_DIR%\test-document.pdf"
    echo ^<^< >> "%FIXTURES_DIR%\test-document.pdf"
    echo /Type /Page >> "%FIXTURES_DIR%\test-document.pdf"
    echo /Parent 2 0 R >> "%FIXTURES_DIR%\test-document.pdf"
    echo /MediaBox [0 0 612 792] >> "%FIXTURES_DIR%\test-document.pdf"
    echo ^>^> >> "%FIXTURES_DIR%\test-document.pdf"
    echo endobj >> "%FIXTURES_DIR%\test-document.pdf"
    echo xref >> "%FIXTURES_DIR%\test-document.pdf"
    echo 0 4 >> "%FIXTURES_DIR%\test-document.pdf"
    echo 0000000000 65535 f ^Z >> "%FIXTURES_DIR%\test-document.pdf"
    echo 0000000009 00000 n ^Z >> "%FIXTURES_DIR%\test-document.pdf"
    echo 0000000058 00000 n ^Z >> "%FIXTURES_DIR%\test-document.pdf"
    echo 0000000115 00000 n ^Z >> "%FIXTURES_DIR%\test-document.pdf"
    echo trailer >> "%FIXTURES_DIR%\test-document.pdf"
    echo ^<^< >> "%FIXTURES_DIR%\test-document.pdf"
    echo /Size 4 >> "%FIXTURES_DIR%\test-document.pdf"
    echo /Root 1 0 R >> "%FIXTURES_DIR%\test-document.pdf"
    echo ^>^> >> "%FIXTURES_DIR%\test-document.pdf"
    echo startxref >> "%FIXTURES_DIR%\test-document.pdf"
    echo 174 >> "%FIXTURES_DIR%\test-document.pdf"
    echo %%%%EOF >> "%FIXTURES_DIR%\test-document.pdf"
)

if not exist "%FIXTURES_DIR%\test-video.mp4" (
    echo Creating test video fixture...
    REM Create a minimal MP4 file (just a header)
    echo ^Z^Z^Z ftypmp41^Z^Z^Z^Zmp41isom^Z^Z^Z^Zfree > "%FIXTURES_DIR%\test-video.mp4"
)

if not exist "%FIXTURES_DIR%\test-audio.mp3" (
    echo Creating test audio fixture...
    REM Create a minimal MP3 file (just a header)
    echo ID3^Z^Z^Z^Z^Z^Z^Z^Z^ZTALB^Z^Z^Z^ZTest Album^Z^Z^Z > "%FIXTURES_DIR%\test-audio.mp3"
)

if not exist "%FIXTURES_DIR%\test-file.txt" (
    echo Creating test text file fixture...
    echo This is a test text file for E2E testing. > "%FIXTURES_DIR%\test-file.txt"
)

echo ✅ Test fixtures are ready

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
cd /d "%UI_DIR%"
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

REM Run the tests
call :print_section "Running E2E Tests"

cd /d "%TESTS_DIR%"

REM Run API tests
echo Running API tests...
npx playwright test api/ --reporter=list,html --output="%REPORTS_DIR%\api-report"
set API_TEST_EXIT_CODE=%errorlevel%

REM Run UI tests
echo Running UI tests...
npx playwright test ui/ --reporter=list,html --output="%REPORTS_DIR%\ui-report"
set UI_TEST_EXIT_CODE=%errorlevel%

REM Stop the application
call :print_section "Stopping Application"

echo Stopping the application...
taskkill /F /PID %BACKEND_PID% >nul 2>&1
taskkill /F /PID %UI_PID% >nul 2>&1

REM Generate summary report
call :print_section "Generating Test Summary Report"

set SUMMARY_FILE=%REPORTS_DIR%\test-summary.md
echo # Unified Media Repository E2E Test Summary Report > "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## Test Execution Details >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo - **Date:** %date% %time% >> "%SUMMARY_FILE%"
echo - **Test Environment:** Local Development >> "%SUMMARY_FILE%"
echo - **Backend PID:** %BACKEND_PID% >> "%SUMMARY_FILE%"
echo - **UI PID:** %UI_PID% >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## API Test Results >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo - **Exit Code:** %API_TEST_EXIT_CODE% >> "%SUMMARY_FILE%"
if %API_TEST_EXIT_CODE% equ 0 (
    echo - **Status:** ✅ PASSED >> "%SUMMARY_FILE%"
) else (
    echo - **Status:** ❌ FAILED >> "%SUMMARY_FILE%"
)
echo - **Report Location:** %REPORTS_DIR%\api-report\index.html >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## UI Test Results >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo - **Exit Code:** %UI_TEST_EXIT_CODE% >> "%SUMMARY_FILE%"
if %UI_TEST_EXIT_CODE% equ 0 (
    echo - **Status:** ✅ PASSED >> "%SUMMARY_FILE%"
) else (
    echo - **Status:** ❌ FAILED >> "%SUMMARY_FILE%"
)
echo - **Report Location:** %REPORTS_DIR%\ui-report\index.html >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## Overall Test Results >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
if %API_TEST_EXIT_CODE% equ 0 if %UI_TEST_EXIT_CODE% equ 0 (
    echo - **Overall Status:** ✅ ALL TESTS PASSED >> "%SUMMARY_FILE%"
) else (
    echo - **Overall Status:** ❌ SOME TESTS FAILED >> "%SUMMARY_FILE%"
)
echo. >> "%SUMMARY_FILE%"
echo ## Test Coverage >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo The E2E tests cover the following areas: >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ### API Tests >> "%SUMMARY_FILE%"
echo - Unified media endpoints (/api/cdn/media/*) >> "%SUMMARY_FILE%"
echo - Media upload for all types (image, document, video, audio) >> "%SUMMARY_FILE%"
echo - Media metadata retrieval >> "%SUMMARY_FILE%"
echo - Media deletion >> "%SUMMARY_FILE%"
echo - Media renaming >> "%SUMMARY_FILE%"
echo - Image resizing >> "%SUMMARY_FILE%"
echo - Error handling and validation >> "%SUMMARY_FILE%"
echo - Backward compatibility with legacy endpoints >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ### UI Tests >> "%SUMMARY_FILE%"
echo - Unified media upload page >> "%SUMMARY_FILE%"
echo - Unified media files pages for all media types >> "%SUMMARY_FILE%"
echo - Media upload functionality >> "%SUMMARY_FILE%"
echo - Media display and metadata viewing >> "%SUMMARY_FILE%"
echo - File operations (copy link, download, rename, resize, delete) >> "%SUMMARY_FILE%"
echo - Search and filtering >> "%SUMMARY_FILE%"
echo - Bulk selection and deletion >> "%SUMMARY_FILE%"
echo - Backward compatibility with legacy pages >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"
echo ## Issues Found >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"

REM Add issues to the summary report if tests failed
if %API_TEST_EXIT_CODE% neq 0 if %UI_TEST_EXIT_CODE% neq 0 (
    echo The following issues were identified during testing: >> "%SUMMARY_FILE%"
    echo. >> "%SUMMARY_FILE%"
    
    REM Check API test results
    if %API_TEST_EXIT_CODE% neq 0 (
        echo ### API Test Failures >> "%SUMMARY_FILE%"
        echo - Some API tests failed. Please check the API test report for details. >> "%SUMMARY_FILE%"
        echo - Common issues include: >> "%SUMMARY_FILE%"
        echo   - Backend server not running >> "%SUMMARY_FILE%"
        echo   - API endpoints not implemented correctly >> "%SUMMARY_FILE%"
        echo   - Database connection issues >> "%SUMMARY_FILE%"
        echo   - File permission issues >> "%SUMMARY_FILE%"
        echo. >> "%SUMMARY_FILE%"
    )
    
    REM Check UI test results
    if %UI_TEST_EXIT_CODE% neq 0 (
        echo ### UI Test Failures >> "%SUMMARY_FILE%"
        echo - Some UI tests failed. Please check the UI test report for details. >> "%SUMMARY_FILE%"
        echo - Common issues include: >> "%SUMMARY_FILE%"
        echo   - UI server not running >> "%SUMMARY_FILE%"
        echo   - UI components not rendering correctly >> "%SUMMARY_FILE%"
        echo   - Missing test IDs or selectors >> "%SUMMARY_FILE%"
        echo   - Timing issues with asynchronous operations >> "%SUMMARY_FILE%"
        echo. >> "%SUMMARY_FILE%"
    )
) else (
    echo No issues were identified during testing. All tests passed successfully. >> "%SUMMARY_FILE%"
    echo. >> "%SUMMARY_FILE%"
)

REM Add recommendations to the summary report
echo ## Recommendations >> "%SUMMARY_FILE%"
echo. >> "%SUMMARY_FILE%"

if %API_TEST_EXIT_CODE% equ 0 if %UI_TEST_EXIT_CODE% equ 0 (
    echo - The unified media repository is ready for deployment. >> "%SUMMARY_FILE%"
    echo - Continue to run these tests as part of the CI/CD pipeline. >> "%SUMMARY_FILE%"
    echo - Consider adding additional edge case tests for production environments. >> "%SUMMARY_FILE%"
    echo. >> "%SUMMARY_FILE%"
) else (
    echo - Review and fix the failed tests before deployment. >> "%SUMMARY_FILE%"
    echo - Ensure all prerequisites are properly installed and configured. >> "%SUMMARY_FILE%"
    echo - Check the application logs for any error messages. >> "%SUMMARY_FILE%"
    echo - Run the tests again after fixing the identified issues. >> "%SUMMARY_FILE%"
    echo. >> "%SUMMARY_FILE%"
)

echo ✅ Test summary report generated: %SUMMARY_FILE%

REM Display final status
call :print_section "Test Execution Complete"

if %API_TEST_EXIT_CODE% equ 0 if %UI_TEST_EXIT_CODE% equ 0 (
    echo 🎉 All tests passed! The unified media repository is working correctly.
    exit /b 0
) else (
    echo ⚠️  Some tests failed. Please check the reports for details.
    exit /b 1
)