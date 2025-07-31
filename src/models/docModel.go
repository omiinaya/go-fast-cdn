package models

import (
	"gorm.io/gorm"
)

type Doc struct {
	gorm.Model

	FileName string `json:"file_name" gorm:"unique;not null"`
	Checksum []byte `json:"checksum"`
}

type DocRepository interface {
	GetAllDocs() []Doc
	GetDocByCheckSum(checksum []byte) Doc
	GetDocByFileName(fileName string) Doc
	AddDoc(doc Doc) (string, error)
	DeleteDoc(fileName string) (string, bool)
	RenameDoc(oldFileName, newFileName string) error
}
