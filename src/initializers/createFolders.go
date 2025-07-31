package initializers

import (
	"fmt"
	"os"

	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

func CreateFolders() {
	uploadsFolder := fmt.Sprintf("%v/uploads", util.ExPath)
	os.Mkdir(uploadsFolder, 0o755)

	// Create legacy directories for backward compatibility
	os.Mkdir(fmt.Sprintf("%v/docs", uploadsFolder), 0o755)
	os.Mkdir(fmt.Sprintf("%v/images", uploadsFolder), 0o755)

	// Create unified media directory for all media files
	os.Mkdir(fmt.Sprintf("%v/media", uploadsFolder), 0o755)
}
