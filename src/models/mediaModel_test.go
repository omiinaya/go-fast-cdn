package models

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMedia_Constants(t *testing.T) {
	// Assert
	require.Equal(t, MediaType("image"), MediaTypeImage)
	require.Equal(t, MediaType("document"), MediaTypeDocument)
}

func TestMedia_ToImage(t *testing.T) {
	// Arrange
	media := Media{
		Model:    gorm.Model{ID: 1},
		FileName: "test.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     MediaTypeImage,
		Width:    intPtr(800),
		Height:   intPtr(600),
	}

	// Act
	image := media.ToImage()

	// Assert
	require.Equal(t, media.Model.ID, image.Model.ID)
	require.Equal(t, media.FileName, image.FileName)
	require.Equal(t, media.Checksum, image.Checksum)
}

func TestMedia_ToDoc(t *testing.T) {
	// Arrange
	media := Media{
		Model:    gorm.Model{ID: 1},
		FileName: "test.pdf",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     MediaTypeDocument,
	}

	// Act
	doc := media.ToDoc()

	// Assert
	require.Equal(t, media.Model.ID, doc.Model.ID)
	require.Equal(t, media.FileName, doc.FileName)
	require.Equal(t, media.Checksum, doc.Checksum)
}

func TestImageFromMedia(t *testing.T) {
	// Arrange
	media := Media{
		Model:    gorm.Model{ID: 1},
		FileName: "test.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     MediaTypeImage,
	}

	// Act
	image := ImageFromMedia(media)

	// Assert
	require.Equal(t, media.Model.ID, image.Model.ID)
	require.Equal(t, media.FileName, image.FileName)
	require.Equal(t, media.Checksum, image.Checksum)
}

func TestDocFromMedia(t *testing.T) {
	// Arrange
	media := Media{
		Model:    gorm.Model{ID: 1},
		FileName: "test.pdf",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     MediaTypeDocument,
	}

	// Act
	doc := DocFromMedia(media)

	// Assert
	require.Equal(t, media.Model.ID, doc.Model.ID)
	require.Equal(t, media.FileName, doc.FileName)
	require.Equal(t, media.Checksum, doc.Checksum)
}

func TestMediaFromImage(t *testing.T) {
	// Arrange
	image := Image{
		Model:    gorm.Model{ID: 1},
		FileName: "test.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
	}

	// Act
	media := MediaFromImage(image)

	// Assert
	require.Equal(t, image.FileName, media.FileName)
	require.Equal(t, image.Checksum, media.Checksum)
	require.Equal(t, MediaTypeImage, media.Type)
	// Note: ID is not copied in MediaFromImage
}

func TestMediaFromDoc(t *testing.T) {
	// Arrange
	doc := Doc{
		Model:    gorm.Model{ID: 1},
		FileName: "test.pdf",
		Checksum: []byte{0x01, 0x02, 0x03},
	}

	// Act
	media := MediaFromDoc(doc)

	// Assert
	require.Equal(t, doc.FileName, media.FileName)
	require.Equal(t, doc.Checksum, media.Checksum)
	require.Equal(t, MediaTypeDocument, media.Type)
	// Note: ID is not copied in MediaFromDoc
}

// Mock implementation of MediaRepository for testing
type mockMediaRepo struct{}

func (m *mockMediaRepo) GetAllMedia() []Media {
	return []Media{}
}

func (m *mockMediaRepo) GetMediaByCheckSum(checksum []byte) Media {
	return Media{}
}

func (m *mockMediaRepo) GetMediaByFileName(fileName string) Media {
	return Media{}
}

func (m *mockMediaRepo) GetMediaByType(mediaType MediaType) []Media {
	return []Media{}
}

func (m *mockMediaRepo) AddMedia(media Media) (string, error) {
	return "", nil
}

func (m *mockMediaRepo) DeleteMedia(fileName string) (string, bool) {
	return "", false
}

func (m *mockMediaRepo) RenameMedia(oldFileName, newFileName string) error {
	return nil
}

func (m *mockMediaRepo) GetAllImages() []Media {
	return []Media{}
}

func (m *mockMediaRepo) GetImageByCheckSum(checksum []byte) Media {
	return Media{}
}

func (m *mockMediaRepo) AddImage(image Media) (string, error) {
	return "", nil
}

func (m *mockMediaRepo) DeleteImage(fileName string) (string, bool) {
	return "", false
}

func (m *mockMediaRepo) RenameImage(oldFileName, newFileName string) error {
	return nil
}

func (m *mockMediaRepo) GetAllDocs() []Media {
	return []Media{}
}

func (m *mockMediaRepo) GetDocByCheckSum(checksum []byte) Media {
	return Media{}
}

func (m *mockMediaRepo) AddDoc(doc Media) (string, error) {
	return "", nil
}

func (m *mockMediaRepo) DeleteDoc(fileName string) (string, bool) {
	return "", false
}

func (m *mockMediaRepo) RenameDoc(oldFileName, newFileName string) error {
	return nil
}

func (m *mockMediaRepo) UpdateImageDimensions(fileName string, width, height int) error {
	return nil
}

func TestMediaRepository_Interface(t *testing.T) {
	// This test ensures that the MediaRepository interface is properly defined
	// and can be implemented by any struct that provides all the required methods.

	// Act & Assert - Verify that the mock struct implements the MediaRepository interface
	var repo MediaRepository = &mockMediaRepo{}
	require.NotNil(t, repo)
}

// Helper function to create a pointer to an int
func intPtr(i int) *int {
	return &i
}
