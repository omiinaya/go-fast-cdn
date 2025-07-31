package database

import (
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	// Create a temporary directory for the test database
	tempDir := t.TempDir()
	dbPath := tempDir + "/test.db"

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.Media{})
	require.NoError(t, err)

	// Get the underlying sql.DB to close it properly
	sqlDB, err := db.DB()
	require.NoError(t, err)

	cleanup := func() {
		sqlDB.Close()
	}

	return db, cleanup
}

func TestNewMediaRepo(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Act
	repo := NewMediaRepo(db)

	// Assert
	require.NotNil(t, repo)
	require.NotNil(t, repo.(*mediaRepo).DB)
}

func TestMediaRepo_GetAllMedia(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media1 := models.Media{
		FileName: "test1.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	media2 := models.Media{
		FileName: "test2.pdf",
		Checksum: []byte{0x04, 0x05, 0x06},
		Type:     models.MediaTypeDocument,
	}
	db.Create(&media1)
	db.Create(&media2)

	// Act
	entries := repo.GetAllMedia()

	// Assert
	require.Len(t, entries, 2)
	require.Equal(t, "test1.jpg", entries[0].FileName)
	require.Equal(t, "test2.pdf", entries[1].FileName)
}

func TestMediaRepo_GetMediaByCheckSum(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	checksum := []byte{0x01, 0x02, 0x03}
	media := models.Media{
		FileName: "test.jpg",
		Checksum: checksum,
		Type:     models.MediaTypeImage,
	}
	db.Create(&media)

	// Act
	entry := repo.GetMediaByCheckSum(checksum)

	// Assert
	require.NotZero(t, entry.ID)
	require.Equal(t, "test.jpg", entry.FileName)
	require.Equal(t, checksum, entry.Checksum)
}

func TestMediaRepo_GetMediaByCheckSum_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Act
	entry := repo.GetMediaByCheckSum([]byte{0x99, 0x99, 0x99})

	// Assert
	require.Zero(t, entry.ID)
	require.Empty(t, entry.Checksum)
}

func TestMediaRepo_GetMediaByFileName(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media := models.Media{
		FileName: "test.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	db.Create(&media)

	// Act
	entry := repo.GetMediaByFileName("test.jpg")

	// Assert
	require.NotZero(t, entry.ID)
	require.Equal(t, "test.jpg", entry.FileName)
}

func TestMediaRepo_GetMediaByFileName_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Act
	entry := repo.GetMediaByFileName("nonexistent.jpg")

	// Assert
	require.Zero(t, entry.ID)
	require.Empty(t, entry.Checksum)
}

func TestMediaRepo_GetMediaByType(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media1 := models.Media{
		FileName: "test1.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	media2 := models.Media{
		FileName: "test2.pdf",
		Checksum: []byte{0x04, 0x05, 0x06},
		Type:     models.MediaTypeDocument,
	}
	media3 := models.Media{
		FileName: "test3.png",
		Checksum: []byte{0x07, 0x08, 0x09},
		Type:     models.MediaTypeImage,
	}
	db.Create(&media1)
	db.Create(&media2)
	db.Create(&media3)

	// Act
	entries := repo.GetMediaByType(models.MediaTypeImage)

	// Assert
	require.Len(t, entries, 2)
	require.Equal(t, "test1.jpg", entries[0].FileName)
	require.Equal(t, "test3.png", entries[1].FileName)
}

func TestMediaRepo_AddMedia(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	media := models.Media{
		FileName: "test.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}

	// Act
	filename, err := repo.AddMedia(media)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "test.jpg", filename)

	// Verify media was added to database
	var entry models.Media
	db.Where("file_name = ?", "test.jpg").First(&entry)
	require.NotZero(t, entry.ID)
	require.Equal(t, "test.jpg", entry.FileName)
}

func TestMediaRepo_DeleteMedia(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media := models.Media{
		FileName: "test.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	db.Create(&media)

	// Act
	filename, success := repo.DeleteMedia("test.jpg")

	// Assert
	require.True(t, success)
	require.Equal(t, "test.jpg", filename)

	// Verify media was deleted from database
	var entry models.Media
	result := db.Where("file_name = ?", "test.jpg").First(&entry)
	require.Error(t, result.Error)
	require.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestMediaRepo_DeleteMedia_NotFound(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Act
	filename, success := repo.DeleteMedia("nonexistent.jpg")

	// Assert
	require.False(t, success)
	require.Empty(t, filename)
}

func TestMediaRepo_RenameMedia(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media := models.Media{
		FileName: "old.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	db.Create(&media)

	// Act
	err := repo.RenameMedia("old.jpg", "new.jpg")

	// Assert
	require.NoError(t, err)

	// Verify media was renamed in database
	var entry models.Media
	result := db.Where("file_name = ?", "new.jpg").First(&entry)
	require.NoError(t, result.Error)
	require.Equal(t, "new.jpg", entry.FileName)
}

func TestMediaRepo_GetAllImages(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media1 := models.Media{
		FileName: "test1.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	media2 := models.Media{
		FileName: "test2.pdf",
		Checksum: []byte{0x04, 0x05, 0x06},
		Type:     models.MediaTypeDocument,
	}
	media3 := models.Media{
		FileName: "test3.png",
		Checksum: []byte{0x07, 0x08, 0x09},
		Type:     models.MediaTypeImage,
	}
	db.Create(&media1)
	db.Create(&media2)
	db.Create(&media3)

	// Act
	entries := repo.GetAllImages()

	// Assert
	require.Len(t, entries, 2)
	require.Equal(t, "test1.jpg", entries[0].FileName)
	require.Equal(t, "test3.png", entries[1].FileName)
}

func TestMediaRepo_GetImageByCheckSum(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	checksum := []byte{0x01, 0x02, 0x03}
	media := models.Media{
		FileName: "test.jpg",
		Checksum: checksum,
		Type:     models.MediaTypeImage,
	}
	db.Create(&media)

	// Act
	entry := repo.GetImageByCheckSum(checksum)

	// Assert
	require.NotZero(t, entry.ID)
	require.Equal(t, "test.jpg", entry.FileName)
	require.Equal(t, checksum, entry.Checksum)
	require.Equal(t, models.MediaTypeImage, entry.Type)
}

func TestMediaRepo_AddImage(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	media := models.Media{
		FileName: "test.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeDocument, // This should be overridden
	}

	// Act
	filename, err := repo.AddImage(media)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "test.jpg", filename)

	// Verify media was added to database with correct type
	var entry models.Media
	db.Where("file_name = ?", "test.jpg").First(&entry)
	require.NotZero(t, entry.ID)
	require.Equal(t, "test.jpg", entry.FileName)
	require.Equal(t, models.MediaTypeImage, entry.Type)
}

func TestMediaRepo_DeleteImage(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media := models.Media{
		FileName: "test.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	db.Create(&media)

	// Act
	filename, success := repo.DeleteImage("test.jpg")

	// Assert
	require.True(t, success)
	require.Equal(t, "test.jpg", filename)

	// Verify media was deleted from database
	var entry models.Media
	result := db.Where("file_name = ? AND type = ?", "test.jpg", models.MediaTypeImage).First(&entry)
	require.Error(t, result.Error)
	require.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestMediaRepo_RenameImage(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media := models.Media{
		FileName: "old.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	db.Create(&media)

	// Act
	err := repo.RenameImage("old.jpg", "new.jpg")

	// Assert
	require.NoError(t, err)

	// Verify media was renamed in database
	var entry models.Media
	result := db.Where("file_name = ? AND type = ?", "new.jpg", models.MediaTypeImage).First(&entry)
	require.NoError(t, result.Error)
	require.Equal(t, "new.jpg", entry.FileName)
}

func TestMediaRepo_GetAllDocs(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media1 := models.Media{
		FileName: "test1.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	media2 := models.Media{
		FileName: "test2.pdf",
		Checksum: []byte{0x04, 0x05, 0x06},
		Type:     models.MediaTypeDocument,
	}
	media3 := models.Media{
		FileName: "test3.docx",
		Checksum: []byte{0x07, 0x08, 0x09},
		Type:     models.MediaTypeDocument,
	}
	db.Create(&media1)
	db.Create(&media2)
	db.Create(&media3)

	// Act
	entries := repo.GetAllDocs()

	// Assert
	require.Len(t, entries, 2)
	require.Equal(t, "test2.pdf", entries[0].FileName)
	require.Equal(t, "test3.docx", entries[1].FileName)
}

func TestMediaRepo_GetDocByCheckSum(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	checksum := []byte{0x01, 0x02, 0x03}
	media := models.Media{
		FileName: "test.pdf",
		Checksum: checksum,
		Type:     models.MediaTypeDocument,
	}
	db.Create(&media)

	// Act
	entry := repo.GetDocByCheckSum(checksum)

	// Assert
	require.NotZero(t, entry.ID)
	require.Equal(t, "test.pdf", entry.FileName)
	require.Equal(t, checksum, entry.Checksum)
	require.Equal(t, models.MediaTypeDocument, entry.Type)
}

func TestMediaRepo_AddDoc(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	media := models.Media{
		FileName: "test.pdf",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage, // This should be overridden
	}

	// Act
	filename, err := repo.AddDoc(media)

	// Assert
	require.NoError(t, err)
	require.Equal(t, "test.pdf", filename)

	// Verify media was added to database with correct type
	var entry models.Media
	db.Where("file_name = ?", "test.pdf").First(&entry)
	require.NotZero(t, entry.ID)
	require.Equal(t, "test.pdf", entry.FileName)
	require.Equal(t, models.MediaTypeDocument, entry.Type)
}

func TestMediaRepo_DeleteDoc(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media := models.Media{
		FileName: "test.pdf",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeDocument,
	}
	db.Create(&media)

	// Act
	filename, success := repo.DeleteDoc("test.pdf")

	// Assert
	require.True(t, success)
	require.Equal(t, "test.pdf", filename)

	// Verify media was deleted from database
	var entry models.Media
	result := db.Where("file_name = ? AND type = ?", "test.pdf", models.MediaTypeDocument).First(&entry)
	require.Error(t, result.Error)
	require.Equal(t, gorm.ErrRecordNotFound, result.Error)
}

func TestMediaRepo_RenameDoc(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media := models.Media{
		FileName: "old.pdf",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeDocument,
	}
	db.Create(&media)

	// Act
	err := repo.RenameDoc("old.pdf", "new.pdf")

	// Assert
	require.NoError(t, err)

	// Verify media was renamed in database
	var entry models.Media
	result := db.Where("file_name = ? AND type = ?", "new.pdf", models.MediaTypeDocument).First(&entry)
	require.NoError(t, result.Error)
	require.Equal(t, "new.pdf", entry.FileName)
}

func TestMediaRepo_UpdateImageDimensions(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media
	media := models.Media{
		FileName: "test.jpg",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeImage,
	}
	db.Create(&media)

	// Act
	err := repo.UpdateImageDimensions("test.jpg", 800, 600)

	// Assert
	require.NoError(t, err)

	// Verify dimensions were updated in database
	var entry models.Media
	result := db.Where("file_name = ? AND type = ?", "test.jpg", models.MediaTypeImage).First(&entry)
	require.NoError(t, result.Error)
	require.NotNil(t, entry.Width)
	require.NotNil(t, entry.Height)
	require.Equal(t, 800, *entry.Width)
	require.Equal(t, 600, *entry.Height)
}

func TestMediaRepo_UpdateImageDimensions_NonImage(t *testing.T) {
	// Arrange
	db, cleanup := setupTestDB(t)
	defer cleanup()
	repo := NewMediaRepo(db)

	// Create test media (document)
	media := models.Media{
		FileName: "test.pdf",
		Checksum: []byte{0x01, 0x02, 0x03},
		Type:     models.MediaTypeDocument,
	}
	db.Create(&media)

	// Act
	err := repo.UpdateImageDimensions("test.pdf", 800, 600)

	// Assert
	require.NoError(t, err) // Should not error, but also should not update anything

	// Verify dimensions were not updated
	var entry models.Media
	result := db.Where("file_name = ? AND type = ?", "test.pdf", models.MediaTypeDocument).First(&entry)
	require.NoError(t, result.Error)
	require.Nil(t, entry.Width)
	require.Nil(t, entry.Height)
}
