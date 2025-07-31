# Test Fixtures

This directory contains test files used for end-to-end testing of the unified media repository.

## Files

- `test-image.jpg` - Sample image file for testing image upload, display, and manipulation
- `test-document.pdf` - Sample PDF document for testing document upload and display
- `test-video.mp4` - Sample video file for testing video upload and display
- `test-audio.mp3` - Sample audio file for testing audio upload and display
- `test-file.txt` - Sample text file for testing invalid file type rejection

## Usage

These files are referenced in the end-to-end test files:
- `ui/tests/api/unified-media-api.spec.ts`
- `ui/tests/ui/unified-media-ui.spec.ts`

The tests use these files to simulate real user uploads and test various scenarios including:
- File upload for different media types
- File validation and rejection
- File display and metadata extraction
- File manipulation (resize, rename, delete)
- Download functionality
- Search and filtering
- Bulk operations

## Note

These files should be small in size to ensure fast test execution. They don't need to be actual valid media files for most tests, as the focus is on testing the application logic rather than media processing.