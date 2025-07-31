package database

import (
	"fmt"

	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"gorm.io/gorm"
)

// imageRepo is deprecated and should not be used for new code.
// Use mediaRepo instead.
// This struct is kept for backward compatibility only.
type imageRepo struct {
	DB *gorm.DB
}

// NewImageRepo is deprecated and should not be used for new code.
// Use NewMediaRepo instead.
// This function is kept for backward compatibility only.
func NewImageRepo(db *gorm.DB) models.ImageRepository {
	return &imageRepo{DB: db}
}

func (repo *imageRepo) GetAllImages() []models.Image {
	var entries []models.Image

	repo.DB.Find(&entries, &models.Image{})

	return entries
}

func (repo *imageRepo) GetImageByCheckSum(checksum []byte) models.Image {
	var entries models.Image

	repo.DB.Where("checksum = ?", checksum).First(&entries)

	return entries
}

func (repo *imageRepo) GetImageByFileName(fileName string) models.Image {
	var entry models.Image
	repo.DB.Where("file_name = ?", fileName).First(&entry)
	return entry
}

func (repo *imageRepo) AddImage(image models.Image) (string, error) {
	// Check if a file with the same name already exists
	var existingImage models.Image
	result := repo.DB.Where("file_name = ?", image.FileName).First(&existingImage)
	if result.Error == nil {
		return "", fmt.Errorf("file with name '%s' already exists", image.FileName)
	}

	result = repo.DB.Create(&image)
	if result.Error != nil {
		return "", result.Error
	}

	return image.FileName, nil
}

func (repo *imageRepo) DeleteImage(fileName string) (string, bool) {
	var image models.Image

	result := repo.DB.Where("file_name = ?", fileName).First(&image)

	if result.Error == nil {
		repo.DB.Delete(&image)
		return fileName, true
	} else {
		return "", false
	}
}

func (repo *imageRepo) RenameImage(oldFileName, newFileName string) error {
	// Check if a file with the new name already exists
	var existingImage models.Image
	result := repo.DB.Where("file_name = ?", newFileName).First(&existingImage)
	if result.Error == nil {
		return fmt.Errorf("file with name '%s' already exists", newFileName)
	}

	image := models.Image{}
	return repo.DB.Model(&image).Where("file_name = ?", oldFileName).Update("file_name", newFileName).Error
}
