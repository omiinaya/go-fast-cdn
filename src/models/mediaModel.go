package models

import (
	"fmt"

	"gorm.io/gorm"
)

// MediaType defines the type of media stored
type MediaType string

const (
	MediaTypeImage    MediaType = "image"
	MediaTypeDocument MediaType = "document"
	MediaTypeVideo    MediaType = "video"
	MediaTypeAudio    MediaType = "audio"
)

// Media represents a unified media entity that can be either an image or a document
type Media struct {
	gorm.Model

	FileName string    `json:"fileName" gorm:"column:file_name;not null"`
	Checksum []byte    `json:"checksum"`
	Type     MediaType `json:"mediaType" gorm:"type:varchar(20);not null;default:'document';column:type"`

	// Image-specific fields (will be empty/null for non-image media)
	Width  *int `json:"width,omitempty" gorm:"default:null"`
	Height *int `json:"height,omitempty" gorm:"default:null"`
}

// TableName specifies the table name for the Media model
func (Media) TableName() string {
	return "media"
}

// AddCompositeUniqueIndex adds a composite unique index on FileName and Type
// This ensures that the combination of FileName and Type is unique across all media records
func (Media) AddCompositeUniqueIndex(db *gorm.DB) error {
	// First, try to create the index directly
	err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_media_file_name_type ON media(file_name, type)").Error
	if err == nil {
		return nil
	}

	// If index creation fails due to constraint violation, we need to handle duplicates
	// This can happen when there are existing records with duplicate (file_name, type) combinations

	// Log the issue for debugging
	fmt.Printf("Warning: Failed to create composite unique index due to existing duplicates: %v\n", err)

	// Find and handle duplicate records
	var duplicates []Media
	err = db.Raw(`
		SELECT file_name, type, COUNT(*) as count 
		FROM media 
		GROUP BY file_name, type 
		HAVING COUNT(*) > 1
	`).Scan(&duplicates).Error

	if err != nil {
		return fmt.Errorf("failed to find duplicate records: %w", err)
	}

	if len(duplicates) > 0 {
		fmt.Printf("Found %d duplicate (file_name, type) combinations. Cleaning up...\n", len(duplicates))

		// Remove duplicates by keeping the oldest record (lowest ID) for each combination
		err = db.Exec(`
			DELETE FROM media 
			WHERE id NOT IN (
				SELECT MIN(id) 
				FROM media 
				GROUP BY file_name, type
			)
		`).Error

		if err != nil {
			return fmt.Errorf("failed to remove duplicate records: %w", err)
		}

		fmt.Printf("Successfully removed duplicate records\n")
	}

	// Now try to create the index again
	err = db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS idx_media_file_name_type ON media(file_name, type)").Error
	if err != nil {
		return fmt.Errorf("failed to create composite unique index after cleanup: %w", err)
	}

	fmt.Printf("Successfully created composite unique index\n")
	return nil
}

// ToImage converts a Media entity to an Image entity for backward compatibility
func (m *Media) ToImage() Image {
	return Image{
		Model:    m.Model,
		FileName: m.FileName,
		Checksum: m.Checksum,
	}
}

// ToDoc converts a Media entity to a Doc entity for backward compatibility
func (m *Media) ToDoc() Doc {
	return Doc{
		Model:    m.Model,
		FileName: m.FileName,
		Checksum: m.Checksum,
	}
}

// ImageFromMedia creates an Image from a Media entity for backward compatibility
func ImageFromMedia(media Media) Image {
	return media.ToImage()
}

// DocFromMedia creates a Doc from a Media entity for backward compatibility
func DocFromMedia(media Media) Doc {
	return media.ToDoc()
}

// MediaFromImage creates a Media entity from an Image for migration
func MediaFromImage(image Image) Media {
	return Media{
		FileName: image.FileName,
		Checksum: image.Checksum,
		Type:     MediaTypeImage,
	}
}

// MediaFromDoc creates a Media entity from a Doc for migration
func MediaFromDoc(doc Doc) Media {
	return Media{
		FileName: doc.FileName,
		Checksum: doc.Checksum,
		Type:     MediaTypeDocument,
	}
}

// MediaRepository defines the interface for media operations
type MediaRepository interface {
	// Generic media operations
	GetAllMedia() []Media
	GetMediaByCheckSum(checksum []byte) Media
	GetMediaByFileName(fileName string) Media
	GetDocByFileName(fileName string) Media
	GetMediaByType(mediaType MediaType) []Media
	AddMedia(media Media) (string, error)
	DeleteMedia(fileName string) (string, bool)
	RenameMedia(oldFileName, newFileName string) error

	// Backward compatibility methods for images
	GetAllImages() []Media
	GetImageByCheckSum(checksum []byte) Media
	AddImage(image Media) (string, error)
	DeleteImage(fileName string) (string, bool)
	RenameImage(oldFileName, newFileName string) error

	// Backward compatibility methods for documents
	GetAllDocs() []Media
	GetDocByCheckSum(checksum []byte) Media
	AddDoc(doc Media) (string, error)
	DeleteDoc(fileName string) (string, bool)
	RenameDoc(oldFileName, newFileName string) error

	// Type-specific methods
	UpdateImageDimensions(fileName string, width, height int) error
}
