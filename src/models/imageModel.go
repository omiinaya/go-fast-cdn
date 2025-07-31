package models

import "gorm.io/gorm"

// Image is deprecated and should not be used for new code.
// Use Media with Type = MediaTypeImage instead.
// This struct is kept for backward compatibility only.
type Image struct {
	gorm.Model

	FileName string `json:"file_name" gorm:"unique;not null"`
	Checksum []byte `json:"checksum"`
}

// ImageRepository is deprecated and should not be used for new code.
// Use MediaRepository instead.
// This interface is kept for backward compatibility only.
type ImageRepository interface {
	GetAllImages() []Image
	GetImageByCheckSum(checksum []byte) Image
	GetImageByFileName(fileName string) Image
	AddImage(image Image) (string, error)
	DeleteImage(fileName string) (string, bool)
	RenameImage(oldFileName, newFileName string) error
}
