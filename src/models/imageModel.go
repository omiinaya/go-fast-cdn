package models

import "gorm.io/gorm"

type Image struct {
	gorm.Model

	FileName string `json:"file_name" gorm:"unique;not null"`
	Checksum []byte `json:"checksum"`
}

type ImageRepository interface {
	GetAllImages() []Image
	GetImageByCheckSum(checksum []byte) Image
	GetImageByFileName(fileName string) Image
	AddImage(image Image) (string, error)
	DeleteImage(fileName string) (string, bool)
	RenameImage(oldFileName, newFileName string) error
}
