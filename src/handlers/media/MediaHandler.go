package handlers

import "github.com/kevinanielsen/go-fast-cdn/src/models"

// MediaHandler handles all media operations (images and documents)
type MediaHandler struct {
	repo models.MediaRepository
}

// NewMediaHandler creates a new MediaHandler instance
func NewMediaHandler(repo models.MediaRepository) *MediaHandler {
	return &MediaHandler{repo}
}
