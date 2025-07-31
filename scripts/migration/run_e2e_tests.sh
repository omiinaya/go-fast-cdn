#!/bin/bash

# End-to-End Test Runner Script for Unified Media Repository
# This script runs the comprehensive E2E tests for the unified media repository

echo "=========================================="
echo "Unified Media Repository E2E Test Runner"
echo "=========================================="

# Set up variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
UI_DIR="$PROJECT_ROOT/ui"
TESTS_DIR="$UI_DIR/tests"
REPORTS_DIR="$PROJECT_ROOT/test-reports"

# Create reports directory if it doesn't exist
mkdir -p "$REPORTS_DIR"

# Function to print section headers
print_section() {
    echo ""
    echo "=========================================="
    echo "$1"
    echo "=========================================="
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
print_section "Checking Prerequisites"

if ! command_exists node; then
    echo "❌ Node.js is not installed. Please install Node.js to continue."
    exit 1
fi

if ! command_exists npm; then
    echo "❌ npm is not installed. Please install npm to continue."
    exit 1
fi

echo "✅ Node.js is installed: $(node --version)"
echo "✅ npm is installed: $(npm --version)"

# Check if Playwright is installed
if [ ! -d "$UI_DIR/node_modules/@playwright" ]; then
    echo "⚠️  Playwright is not installed. Installing Playwright..."
    cd "$UI_DIR"
    npm install @playwright/test
    npx playwright install
    cd "$PROJECT_ROOT"
fi

# Create test fixture files if they don't exist
print_section "Setting Up Test Fixtures"

FIXTURES_DIR="$TESTS_DIR/fixtures"
mkdir -p "$FIXTURES_DIR"

# Create minimal test files if they don't exist
if [ ! -f "$FIXTURES_DIR/test-image.jpg" ]; then
    echo "Creating test image fixture..."
    # Create a minimal 1x1 pixel JPEG file
    echo -e "\xFF\xD8\xFF\xE0\x00\x10JFIF\x00\x01\x01\x01\x00H\x00H\x00\x00\xFF\xDB\x00C\x00\xFF\xD9" > "$FIXTURES_DIR/test-image.jpg"
fi

if [ ! -f "$FIXTURES_DIR/test-document.pdf" ]; then
    echo "Creating test document fixture..."
    # Create a minimal PDF file
    echo -e "%PDF-1.4\n1 0 obj\n<<\n/Type /Catalog\n/Pages 2 0 R\n>>\nendobj\n2 0 obj\n<<\n/Type /Pages\n/Kids [3 0 R]\n/Count 1\n>>\nendobj\n3 0 obj\n<<\n/Type /Page\n/Parent 2 0 R\n/MediaBox [0 0 612 792]\n>>\nendobj\nxref\n0 4\n0000000000 65535 f \n0000000009 00000 n \n0000000058 00000 n \n0000000115 00000 n \ntrailer\n<<\n/Size 4\n/Root 1 0 R\n>>\nstartxref\n174\n%%EOF" > "$FIXTURES_DIR/test-document.pdf"
fi

if [ ! -f "$FIXTURES_DIR/test-video.mp4" ]; then
    echo "Creating test video fixture..."
    # Create a minimal MP4 file (just a header)
    echo -e "\x00\x00\x00\x20ftypmp41\x00\x00\x00\x00mp41isom\x00\x00\x00\x08free" > "$FIXTURES_DIR/test-video.mp4"
fi

if [ ! -f "$FIXTURES_DIR/test-audio.mp3" ]; then
    echo "Creating test audio fixture..."
    # Create a minimal MP3 file (just a header)
    echo -e "ID3\x03\x00\x00\x00\x00\x00\x1FTALB\x00\x00\x00\x0CTest Album\x00\x00\x00" > "$FIXTURES_DIR/test-audio.mp3"
fi

if [ ! -f "$FIXTURES_DIR/test-file.txt" ]; then
    echo "Creating test text file fixture..."
    echo "This is a test text file for E2E testing." > "$FIXTURES_DIR/test-file.txt"
fi

echo "✅ Test fixtures are ready"

# Start the application in the background
print_section "Starting Application"

echo "Starting the application in the background..."
cd "$PROJECT_ROOT"

# Start the Go backend
go run main.go &
BACKEND_PID=$!

# Wait for the backend to start
echo "Waiting for the backend to start..."
sleep 5

# Start the UI in development mode
cd "$UI_DIR"
npm run dev &
UI_PID=$!

# Wait for the UI to start
echo "Waiting for the UI to start..."
sleep 10

# Check if the application is running
if curl -s http://localhost:3000 > /dev/null; then
    echo "✅ Application is running"
else
    echo "❌ Failed to start the application"
    kill $BACKEND_PID $UI_PID 2>/dev/null
    exit 1
fi

# Run the tests
print_section "Running E2E Tests"

cd "$TESTS_DIR"

# Run API tests
echo "Running API tests..."
npx playwright test api/ --reporter=list,html --output="$REPORTS_DIR/api-report"
API_TEST_EXIT_CODE=$?

# Run UI tests
echo "Running UI tests..."
npx playwright test ui/ --reporter=list,html --output="$REPORTS_DIR/ui-report"
UI_TEST_EXIT_CODE=$?

# Stop the application
print_section "Stopping Application"

echo "Stopping the application..."
kill $BACKEND_PID $UI_PID 2>/dev/null

# Generate summary report
print_section "Generating Test Summary Report"

SUMMARY_FILE="$REPORTS_DIR/test-summary.md"
cat > "$SUMMARY_FILE" << EOF
# Unified Media Repository E2E Test Summary Report

## Test Execution Details

- **Date:** $(date)
- **Test Environment:** Local Development
- **Backend PID:** $BACKEND_PID
- **UI PID:** $UI_PID

## API Test Results

- **Exit Code:** $API_TEST_EXIT_CODE
- **Status:** $([ $API_TEST_EXIT_CODE -eq 0 ] && echo "✅ PASSED" || echo "❌ FAILED")
- **Report Location:** $REPORTS_DIR/api-report/index.html

## UI Test Results

- **Exit Code:** $UI_TEST_EXIT_CODE
- **Status:** $([ $UI_TEST_EXIT_CODE -eq 0 ] && echo "✅ PASSED" || echo "❌ FAILED")
- **Report Location:** $REPORTS_DIR/ui-report/index.html

## Overall Test Results

- **Overall Status:** $([ $API_TEST_EXIT_CODE -eq 0 ] && [ $UI_TEST_EXIT_CODE -eq 0 ] && echo "✅ ALL TESTS PASSED" || echo "❌ SOME TESTS FAILED")

## Test Coverage

The E2E tests cover the following areas:

### API Tests
- Unified media endpoints (/api/cdn/media/*)
- Media upload for all types (image, document, video, audio)
- Media metadata retrieval
- Media deletion
- Media renaming
- Image resizing
- Error handling and validation
- Backward compatibility with legacy endpoints

### UI Tests
- Unified media upload page
- Unified media files pages for all media types
- Media upload functionality
- Media display and metadata viewing
- File operations (copy link, download, rename, resize, delete)
- Search and filtering
- Bulk selection and deletion
- Backward compatibility with legacy pages

## Issues Found

EOF

# Add issues to the summary report if tests failed
if [ $API_TEST_EXIT_CODE -ne 0 ] || [ $UI_TEST_EXIT_CODE -ne 0 ]; then
    cat >> "$SUMMARY_FILE" << EOF
The following issues were identified during testing:

EOF

    # Check API test results
    if [ $API_TEST_EXIT_CODE -ne 0 ]; then
        cat >> "$SUMMARY_FILE" << EOF
### API Test Failures
- Some API tests failed. Please check the API test report for details.
- Common issues include:
  - Backend server not running
  - API endpoints not implemented correctly
  - Database connection issues
  - File permission issues

EOF
    fi

    # Check UI test results
    if [ $UI_TEST_EXIT_CODE -ne 0 ]; then
        cat >> "$SUMMARY_FILE" << EOF
### UI Test Failures
- Some UI tests failed. Please check the UI test report for details.
- Common issues include:
  - UI server not running
  - UI components not rendering correctly
  - Missing test IDs or selectors
  - Timing issues with asynchronous operations

EOF
    fi
else
    cat >> "$SUMMARY_FILE" << EOF
No issues were identified during testing. All tests passed successfully.

EOF
fi

# Add recommendations to the summary report
cat >> "$SUMMARY_FILE" << EOF
## Recommendations

EOF

if [ $API_TEST_EXIT_CODE -eq 0 ] && [ $UI_TEST_EXIT_CODE -eq 0 ]; then
    cat >> "$SUMMARY_FILE" << EOF
- The unified media repository is ready for deployment.
- Continue to run these tests as part of the CI/CD pipeline.
- Consider adding additional edge case tests for production environments.

EOF
else
    cat >> "$SUMMARY_FILE" << EOF
- Review and fix the failed tests before deployment.
- Ensure all prerequisites are properly installed and configured.
- Check the application logs for any error messages.
- Run the tests again after fixing the identified issues.

EOF
fi

echo "✅ Test summary report generated: $SUMMARY_FILE"

# Display final status
print_section "Test Execution Complete"

if [ $API_TEST_EXIT_CODE -eq 0 ] && [ $UI_TEST_EXIT_CODE -eq 0 ]; then
    echo "🎉 All tests passed! The unified media repository is working correctly."
    exit 0
else
    echo "⚠️  Some tests failed. Please check the reports for details."
    exit 1
fi