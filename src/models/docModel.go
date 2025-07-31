package models

import (
	"gorm.io/gorm"
)

// Doc is deprecated and should not be used for new code.
// Use Media with Type = MediaTypeDocument instead.
// This struct is kept for backward compatibility only.
type Doc struct {
	gorm.Model

	FileName string `json:"file_name" gorm:"unique;not null"`
	Checksum []byte `json:"checksum"`
}

// DocRepository is deprecated and should not be used for new code.
// Use MediaRepository instead.
// This interface is kept for backward compatibility only.
type DocRepository interface {
	GetAllDocs() []Doc
	GetDocByCheckSum(checksum []byte) Doc
	GetDocByFileName(fileName string) Doc
	AddDoc(doc Doc) (string, error)
	DeleteDoc(fileName string) (string, bool)
	RenameDoc(oldFileName, newFileName string) error
}
