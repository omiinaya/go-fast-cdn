package database

import (
	"github.com/kevinanielsen/go-fast-cdn/src/models"
	"gorm.io/gorm"
)

type mediaRepo struct {
	DB *gorm.DB
}

func NewMediaRepo(db *gorm.DB) models.MediaRepository {
	return &mediaRepo{DB: db}
}

// Generic media operations

func (repo *mediaRepo) GetAllMedia() []models.Media {
	var entries []models.Media
	repo.DB.Find(&entries, &models.Media{})
	return entries
}

func (repo *mediaRepo) GetMediaByCheckSum(checksum []byte) models.Media {
	var entry models.Media
	repo.DB.Where("checksum = ?", checksum).First(&entry)
	return entry
}

func (repo *mediaRepo) GetMediaByFileName(fileName string) models.Media {
	var entry models.Media
	repo.DB.Where("file_name = ?", fileName).First(&entry)
	return entry
}

func (repo *mediaRepo) GetMediaByType(mediaType models.MediaType) []models.Media {
	var entries []models.Media
	repo.DB.Where("type = ?", mediaType).Find(&entries)
	return entries
}

func (repo *mediaRepo) AddMedia(media models.Media) (string, error) {
	result := repo.DB.Create(&media)
	if result.Error != nil {
		return "", result.Error
	}
	return media.FileName, nil
}

func (repo *mediaRepo) DeleteMedia(fileName string) (string, bool) {
	var media models.Media
	result := repo.DB.Where("file_name = ?", fileName).First(&media)

	if result.Error == nil {
		repo.DB.Delete(&media)
		return fileName, true
	} else {
		return "", false
	}
}

func (repo *mediaRepo) RenameMedia(oldFileName, newFileName string) error {
	media := models.Media{}
	return repo.DB.Model(&media).Where("file_name = ?", oldFileName).Update("file_name", newFileName).Error
}

// Backward compatibility methods for images

func (repo *mediaRepo) GetAllImages() []models.Media {
	return repo.GetMediaByType(models.MediaTypeImage)
}

func (repo *mediaRepo) GetImageByCheckSum(checksum []byte) models.Media {
	var media models.Media
	repo.DB.Where("checksum = ? AND type = ?", checksum, models.MediaTypeImage).First(&media)
	return media
}

func (repo *mediaRepo) AddImage(image models.Media) (string, error) {
	// Ensure the media type is set to image
	image.Type = models.MediaTypeImage
	return repo.AddMedia(image)
}

func (repo *mediaRepo) DeleteImage(fileName string) (string, bool) {
	var media models.Media
	result := repo.DB.Where("file_name = ? AND type = ?", fileName, models.MediaTypeImage).First(&media)

	if result.Error == nil {
		repo.DB.Delete(&media)
		return fileName, true
	} else {
		return "", false
	}
}

func (repo *mediaRepo) RenameImage(oldFileName, newFileName string) error {
	media := models.Media{}
	return repo.DB.Model(&media).Where("file_name = ? AND type = ?", oldFileName, models.MediaTypeImage).Update("file_name", newFileName).Error
}

// Backward compatibility methods for documents

func (repo *mediaRepo) GetAllDocs() []models.Media {
	return repo.GetMediaByType(models.MediaTypeDocument)
}

func (repo *mediaRepo) GetDocByCheckSum(checksum []byte) models.Media {
	var media models.Media
	repo.DB.Where("checksum = ? AND type = ?", checksum, models.MediaTypeDocument).First(&media)
	return media
}

func (repo *mediaRepo) AddDoc(doc models.Media) (string, error) {
	// Ensure the media type is set to document
	doc.Type = models.MediaTypeDocument
	return repo.AddMedia(doc)
}

func (repo *mediaRepo) DeleteDoc(fileName string) (string, bool) {
	var media models.Media
	result := repo.DB.Where("file_name = ? AND type = ?", fileName, models.MediaTypeDocument).First(&media)

	if result.Error == nil {
		repo.DB.Delete(&media)
		return fileName, true
	} else {
		return "", false
	}
}

func (repo *mediaRepo) RenameDoc(oldFileName, newFileName string) error {
	media := models.Media{}
	return repo.DB.Model(&media).Where("file_name = ? AND type = ?", oldFileName, models.MediaTypeDocument).Update("file_name", newFileName).Error
}

// UpdateImageDimensions updates the width and height of an image in the database
func (repo *mediaRepo) UpdateImageDimensions(fileName string, width, height int) error {
	media := models.Media{}
	return repo.DB.Model(&media).
		Where("file_name = ? AND type = ?", fileName, models.MediaTypeImage).
		Updates(map[string]interface{}{
			"width":  width,
			"height": height,
		}).Error
}
