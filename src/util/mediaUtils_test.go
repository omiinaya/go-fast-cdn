package util

import (
	"testing"
)

func TestGetMediaTypeFromExtension(t *testing.T) {
	tests := []struct {
		name      string
		extension string
		want      MediaType
		wantErr   bool
	}{
		{"Image JPG", ".jpg", MediaTypeImage, false},
		{"Image PNG", ".png", MediaTypeImage, false},
		{"Document PDF", ".pdf", MediaTypeDocument, false},
		{"Document DOCX", ".docx", MediaTypeDocument, false},
		{"Video MP4", ".mp4", MediaTypeVideo, false},
		{"Audio MP3", ".mp3", MediaTypeAudio, false},
		{"Unsupported extension", ".xyz", "", true},
		{"Empty extension", "", "", true},
		{"Extension without dot", "jpg", MediaTypeImage, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMediaTypeFromExtension(tt.extension)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMediaTypeFromExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMediaTypeFromExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMediaTypeFromMIME(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		want     MediaType
		wantErr  bool
	}{
		{"Image JPEG", "image/jpeg", MediaTypeImage, false},
		{"Image PNG", "image/png", MediaTypeImage, false},
		{"Document PDF", "application/pdf", MediaTypeDocument, false},
		{"Document DOCX", "application/vnd.openxmlformats-officedocument.wordprocessingml.document", MediaTypeDocument, false},
		{"Video MP4", "video/mp4", MediaTypeVideo, false},
		{"Audio MP3", "audio/mpeg", MediaTypeAudio, false},
		{"Unsupported MIME type", "application/xyz", "", true},
		{"Empty MIME type", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMediaTypeFromMIME(tt.mimeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMediaTypeFromMIME() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMediaTypeFromMIME() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMediaTypeFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     MediaType
		wantErr  bool
	}{
		{"Image filename", "image.jpg", MediaTypeImage, false},
		{"Document filename", "document.pdf", MediaTypeDocument, false},
		{"Video filename", "video.mp4", MediaTypeVideo, false},
		{"Audio filename", "audio.mp3", MediaTypeAudio, false},
		{"Unsupported filename", "file.xyz", "", true},
		{"Empty filename", "", "", true},
		{"Filename without extension", "file", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMediaTypeFromFilename(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMediaTypeFromFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMediaTypeFromFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectMediaType(t *testing.T) {
	tests := []struct {
		name       string
		fileBuffer []byte
		want       MediaType
		wantErr    bool
	}{
		{"Empty buffer", []byte{}, "", true},
		{"PNG header", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, MediaTypeImage, false},
		{"JPEG header", []byte{0xFF, 0xD8, 0xFF}, MediaTypeImage, false},
		{"PDF header", []byte{0x25, 0x50, 0x44, 0x46}, MediaTypeDocument, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DetectMediaType(tt.fileBuffer)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectMediaType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DetectMediaType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSupportedExtension(t *testing.T) {
	tests := []struct {
		name      string
		extension string
		want      bool
	}{
		{"Supported image", ".jpg", true},
		{"Supported document", ".pdf", true},
		{"Supported video", ".mp4", true},
		{"Supported audio", ".mp3", true},
		{"Unsupported extension", ".xyz", false},
		{"Empty extension", "", false},
		{"Extension without dot", "jpg", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSupportedExtension(tt.extension); got != tt.want {
				t.Errorf("IsSupportedExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsSupportedMIMEType(t *testing.T) {
	tests := []struct {
		name     string
		mimeType string
		want     bool
	}{
		{"Supported image", "image/jpeg", true},
		{"Supported document", "application/pdf", true},
		{"Supported video", "video/mp4", true},
		{"Supported audio", "audio/mpeg", true},
		{"Unsupported MIME type", "application/xyz", false},
		{"Empty MIME type", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsSupportedMIMEType(tt.mimeType); got != tt.want {
				t.Errorf("IsSupportedMIMEType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMIMETypeFromExtension(t *testing.T) {
	tests := []struct {
		name      string
		extension string
		want      string
		wantErr   bool
	}{
		{"JPG extension", ".jpg", "image/jpeg", false},
		{"PNG extension", ".png", "image/png", false},
		{"PDF extension", ".pdf", "application/pdf", false},
		{"Unsupported extension", ".xyz", "", true},
		{"Empty extension", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMIMETypeFromExtension(tt.extension)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMIMETypeFromExtension() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetMIMETypeFromExtension() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateMediaFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		mimeType string
		want     MediaType
		wantErr  bool
	}{
		{"Valid image", "image.jpg", "image/jpeg", MediaTypeImage, false},
		{"Valid document", "document.pdf", "application/pdf", MediaTypeDocument, false},
		{"Valid video", "video.mp4", "video/mp4", MediaTypeVideo, false},
		{"Valid audio", "audio.mp3", "audio/mpeg", MediaTypeAudio, false},
		{"Mismatched type", "image.jpg", "application/pdf", "", true},
		{"Unsupported filename", "file.xyz", "application/xyz", "", true},
		{"Empty filename", "", "image/jpeg", "", true},
		{"Empty MIME type", "image.jpg", "", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateMediaFile(tt.filename, tt.mimeType)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMediaFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ValidateMediaFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMediaTypeInfo(t *testing.T) {
	tests := []struct {
		name      string
		mediaType MediaType
		wantErr   bool
	}{
		{"Image type", MediaTypeImage, false},
		{"Document type", MediaTypeDocument, false},
		{"Video type", MediaTypeVideo, false},
		{"Audio type", MediaTypeAudio, false},
		{"Unsupported type", MediaType("unsupported"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetMediaTypeInfo(tt.mediaType)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMediaTypeInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
