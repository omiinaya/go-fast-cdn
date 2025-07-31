package handlers

import "github.com/kevinanielsen/go-fast-cdn/src/models"

// ImageHandler is deprecated and should not be used for new code.
// Use MediaHandler instead.
// This struct is kept for backward compatibility only.
type ImageHandler struct {
	repo models.ImageRepository
}

// NewImageHandler is deprecated and should not be used for new code.
// Use NewMediaHandler instead.
// This function is kept for backward compatibility only.
func NewImageHandler(repo models.ImageRepository) *ImageHandler {
	return &ImageHandler{repo}
}
