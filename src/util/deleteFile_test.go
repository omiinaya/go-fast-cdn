package util

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeleteFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	ExPath = tempDir // Set the execution path to temp directory for testing

	// Create test directories
	testDir := filepath.Join(tempDir, "uploads", "images")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Create a test file
	testFileName := "test.jpg"
	testFilePath := filepath.Join(testDir, testFileName)
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(testFilePath)
	require.NoError(t, err)

	// Act
	err = DeleteFile(testFileName, "images")

	// Assert
	require.NoError(t, err)

	// Verify file was deleted
	_, err = os.Stat(testFilePath)
	require.True(t, os.IsNotExist(err))
}

func TestDeleteMediaFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	ExPath = tempDir // Set the execution path to temp directory for testing

	// Create test directories
	testDir := filepath.Join(tempDir, "uploads", "images")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Create a test file
	testFileName := "test.jpg"
	testFilePath := filepath.Join(testDir, testFileName)
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(testFilePath)
	require.NoError(t, err)

	// Act
	err = DeleteMediaFile(testFileName, "images")

	// Assert
	require.NoError(t, err)

	// Verify file was deleted
	_, err = os.Stat(testFilePath)
	require.True(t, os.IsNotExist(err))
}

func TestDeleteUnifiedMediaFile(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	ExPath = tempDir // Set the execution path to temp directory for testing

	// Create test directories
	testDir := filepath.Join(tempDir, "uploads", "media")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Create a test file
	testFileName := "test.jpg"
	testFilePath := filepath.Join(testDir, testFileName)
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	require.NoError(t, err)

	// Verify file exists
	_, err = os.Stat(testFilePath)
	require.NoError(t, err)

	// Act
	err = DeleteUnifiedMediaFile(testFileName)

	// Assert
	require.NoError(t, err)

	// Verify file was deleted
	_, err = os.Stat(testFilePath)
	require.True(t, os.IsNotExist(err))
}

func TestGetMediaFilePath(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	ExPath = tempDir // Set the execution path to temp directory for testing
	testFileName := "test.jpg"

	// Act
	filePath := GetMediaFilePath(testFileName)

	// Assert
	expectedPath := filepath.Join(tempDir, "uploads", "media", testFileName)
	require.Equal(t, expectedPath, filePath)
}

func TestGetLegacyMediaFilePath(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	ExPath = tempDir // Set the execution path to temp directory for testing
	testFileName := "test.jpg"
	mediaType := "images"

	// Act
	filePath := GetLegacyMediaFilePath(testFileName, mediaType)

	// Assert
	expectedPath := filepath.Join(tempDir, "uploads", mediaType, testFileName)
	require.Equal(t, expectedPath, filePath)
}

func TestFileExists(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()

	// Create a test file
	testFilePath := filepath.Join(tempDir, "test.jpg")
	err := os.WriteFile(testFilePath, []byte("test content"), 0644)
	require.NoError(t, err)

	nonexistentFilePath := filepath.Join(tempDir, "nonexistent.jpg")

	// Act & Assert
	require.True(t, FileExists(testFilePath))
	require.False(t, FileExists(nonexistentFilePath))
}

func TestMediaFileExists(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	ExPath = tempDir // Set the execution path to temp directory for testing

	// Create test directories
	testDir := filepath.Join(tempDir, "uploads", "media")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Create a test file
	testFileName := "test.jpg"
	testFilePath := filepath.Join(testDir, testFileName)
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	require.NoError(t, err)

	nonexistentFileName := "nonexistent.jpg"

	// Act & Assert
	require.True(t, MediaFileExists(testFileName))
	require.False(t, MediaFileExists(nonexistentFileName))
}

func TestLegacyMediaFileExists(t *testing.T) {
	// Arrange
	tempDir := t.TempDir()
	ExPath = tempDir // Set the execution path to temp directory for testing

	// Create test directories
	testDir := filepath.Join(tempDir, "uploads", "images")
	err := os.MkdirAll(testDir, 0755)
	require.NoError(t, err)

	// Create a test file
	testFileName := "test.jpg"
	testFilePath := filepath.Join(testDir, testFileName)
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	require.NoError(t, err)

	nonexistentFileName := "nonexistent.jpg"

	// Act & Assert
	require.True(t, LegacyMediaFileExists(testFileName, "images"))
	require.False(t, LegacyMediaFileExists(nonexistentFileName, "images"))
}

func TestGetFileExtension(t *testing.T) {
	// Arrange
	testCases := []struct {
		input    string
		expected string
	}{
		{"file.jpg", ".jpg"},
		{"file.JPG", ".jpg"},
		{"file.jpeg", ".jpeg"},
		{"file.png", ".png"},
		{"file.pdf", ".pdf"},
		{"file.docx", ".docx"},
		{"file", ""},
		{".hiddenfile", ".hiddenfile"},
		{"file.with.multiple.dots.jpg", ".jpg"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			// Act
			result := GetFileExtension(tc.input)

			// Assert
			require.Equal(t, tc.expected, result)
		})
	}
}
