# End-to-End Tests for Unified Media Repository

This directory contains comprehensive end-to-end tests for the unified media repository, which is the final step in Phase 5 of the unification project.

## Test Structure

The tests are organized into two main categories:

### 1. API Tests (`api/`)

These tests verify the functionality of the backend API endpoints for the unified media repository:

- `unified-media-api.spec.ts` - Comprehensive tests for all unified media API endpoints

#### API Test Coverage

- **Media Retrieval**: Testing the `/api/cdn/media` endpoint with and without type filters
- **Media Upload**: Testing upload functionality for all media types (image, document, video, audio)
- **Media Metadata**: Testing metadata retrieval for different media types
- **Media Deletion**: Testing file deletion functionality
- **Media Renaming**: Testing file renaming functionality
- **Image Resizing**: Testing image-specific resize operations
- **Error Handling**: Testing validation and error scenarios
- **Backward Compatibility**: Testing that legacy endpoints still work

### 2. UI Tests (`ui/`)

These tests verify the functionality of the frontend UI components for the unified media repository:

- `unified-media-ui.spec.ts` - Comprehensive tests for all unified media UI components

#### UI Test Coverage

- **Page Navigation**: Testing that all unified media pages load correctly
- **Media Upload**: Testing the unified media upload component for all media types
- **Media Display**: Testing that media cards display correctly for all media types
- **Media Operations**: Testing file operations (copy link, download, rename, resize, delete)
- **Search and Filtering**: Testing the search functionality
- **Bulk Operations**: Testing bulk selection and deletion
- **Backward Compatibility**: Testing that legacy UI pages still work

### 3. Test Fixtures (`fixtures/`)

This directory contains sample files used for testing:

- `test-image.jpg` - Sample image file
- `test-document.pdf` - Sample PDF document
- `test-video.mp4` - Sample video file
- `test-audio.mp3` - Sample audio file
- `test-file.txt` - Sample text file for testing invalid file types

### 4. Configuration Files

- `playwright.config.ts` - Playwright test configuration

## Running the Tests

### Prerequisites

Before running the tests, ensure you have the following installed:

- Node.js (v14 or higher)
- npm
- Go (for running the backend)

### Running Tests Manually

1. Start the backend server:
   ```bash
   cd /path/to/project/root
   go run main.go
   ```

2. Start the UI development server:
   ```bash
   cd /path/to/project/root/ui
   npm run dev
   ```

3. Install Playwright (if not already installed):
   ```bash
   cd /path/to/project/root/ui/tests
   npm install @playwright/test
   npx playwright install
   ```

4. Run the tests:
   ```bash
   # Run API tests only
   npx playwright test api/

   # Run UI tests only
   npx playwright test ui/

   # Run all tests
   npx playwright test
   ```

### Running Tests with the Test Runner Script

For convenience, you can use the provided test runner script:

#### On Linux/macOS:
```bash
cd /path/to/project/root
chmod +x scripts/run_e2e_tests.sh
./scripts/run_e2e_tests.sh
```

#### On Windows:
```cmd
cd /path/to/project/root
scripts\run_e2e_tests.bat
```

The test runner script will:
1. Check and install prerequisites
2. Create test fixtures if they don't exist
3. Start the application in the background
4. Run all tests
5. Generate a summary report
6. Stop the application

## Test Reports

After running the tests, HTML reports will be generated in the `test-reports` directory:

- `test-reports/api-report/index.html` - API test report
- `test-reports/ui-report/index.html` - UI test report
- `test-reports/test-summary.md` - Summary report with overall results and recommendations

## Test Scenarios

### API Test Scenarios

1. **Media Retrieval**
   - Retrieve all media without type filter
   - Retrieve media filtered by type (image, document, video, audio)

2. **Media Upload**
   - Upload images, documents, videos, and audio files
   - Handle duplicate file uploads
   - Reject invalid file types
   - Handle file size limits

3. **Media Metadata**
   - Retrieve metadata for different media types
   - Handle non-existent files
   - Handle missing type parameters

4. **Media Operations**
   - Delete media files
   - Rename media files
   - Resize images
   - Handle errors for invalid operations

5. **Backward Compatibility**
   - Test legacy image endpoints
   - Test legacy document endpoints
   - Ensure legacy functionality still works

### UI Test Scenarios

1. **Page Navigation**
   - Load unified media upload page
   - Load unified media files pages for all media types

2. **Media Upload**
   - Upload different media types using the unified upload component
   - Handle file validation and rejection

3. **Media Display**
   - Display media cards correctly for all media types
   - Show metadata in modal dialogs

4. **Media Operations**
   - Copy links to clipboard
   - Download files
   - Rename files
   - Resize images
   - Delete files

5. **Search and Filtering**
   - Filter media files by search term
   - Clear search and show all files

6. **Bulk Operations**
   - Select multiple files
   - Delete selected files in bulk

7. **Backward Compatibility**
   - Load legacy upload pages
   - Load legacy media pages

## Troubleshooting

### Common Issues

1. **Tests failing due to application not starting**
   - Ensure the backend and UI servers are running
   - Check that port 3000 is available
   - Verify that all dependencies are installed

2. **Tests failing due to missing fixtures**
   - Ensure the test fixtures are created in the `fixtures` directory
   - Check file permissions

3. **Tests failing due to timeouts**
   - Increase the timeout values in the test configuration
   - Check if the application is running slowly

4. **Tests failing due to element not found**
   - Verify that the test IDs and selectors are correct
   - Check if the UI has changed

### Debugging Tips

1. **Run tests in headed mode** to see the browser:
   ```bash
   npx playwright test --headed
   ```

2. **Run tests with debug mode** to pause execution:
   ```bash
   npx playwright test --debug
   ```

3. **Generate trace files** to analyze test execution:
   ```bash
   npx playwright test --trace on
   ```

4. **Take screenshots** on failure:
   Screenshots are automatically taken on failure and saved in the test results directory.

## Contributing

When adding new tests:

1. Follow the existing test structure and naming conventions
2. Add descriptive test names and comments
3. Test both success and error scenarios
4. Ensure tests are independent and can run in any order
5. Clean up any test data created during the test

## Future Enhancements

Potential improvements to the test suite:

1. **Performance Testing**: Add tests to measure upload and download speeds
2. **Accessibility Testing**: Add tests to verify accessibility compliance
3. **Cross-Browser Testing**: Expand browser coverage beyond Chromium, Firefox, and WebKit
4. **Mobile Testing**: Add tests for mobile viewports and touch interactions
5. **Integration Testing**: Add tests that verify integration with external services