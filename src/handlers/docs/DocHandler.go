package handlers

import (
	"github.com/kevinanielsen/go-fast-cdn/src/models"
)

// DocHandler is deprecated and should not be used for new code.
// Use MediaHandler instead.
// This struct is kept for backward compatibility only.
type DocHandler struct {
	repo models.DocRepository
}

// NewDocHandler is deprecated and should not be used for new code.
// Use NewMediaHandler instead.
// This function is kept for backward compatibility only.
func NewDocHandler(repo models.DocRepository) *DocHandler {
	return &DocHandler{repo}
}
