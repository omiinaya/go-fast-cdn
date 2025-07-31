// Package util contains utility functions used throughout the application.
package util

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// countVal counts the number of occurrences of val in str.
// It splits str into a slice of strings, iterates over the slice
// checking each string for equality with val, and increments a counter
// each time a match is found. The final count is returned.
func countVal(str string, val string) int {
	var count int
	arr := strings.Split(str, "")
	for _, v := range arr {
		if v == val {
			count++
		}
	}

	return count
}

// FilterFilename removes illegal characters from a filename string.
// It ensures there is at most one period in the filename,
// replaces any '/' and '\' characters,
// and returns the filtered string.
func FilterFilename(filename string) (string, error) {
	if countVal(filename, ".") > 1 {
		return filename, errors.New("filename cannot contain more than one period character")
	}

	var filteredStr string

	filteredStr = strings.Replace(filename, "/", "", -1)
	filteredStr = strings.Replace(filteredStr, `\`, "", -1)

	return filteredStr, nil
}

// SanitizeFilename removes or replaces characters that are not safe for filenames
// while preserving the file extension. This is a more comprehensive version
// of FilterFilename for use with the unified media system.
func SanitizeFilename(filename string) (string, error) {
	if filename == "" {
		return "", errors.New("filename cannot be empty")
	}

	// Extract the file extension
	ext := filepath.Ext(filename)
	nameWithoutExt := filename[:len(filename)-len(ext)]

	// Replace spaces with underscores
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, " ", "_")

	// Remove or replace special characters
	// Allow alphanumeric characters, underscores, hyphens, and periods
	reg := regexp.MustCompile(`[^\w\-\.]`)
	nameWithoutExt = reg.ReplaceAllString(nameWithoutExt, "")

	// Remove consecutive non-alphanumeric characters
	reg = regexp.MustCompile(`[^\w\-]+`)
	nameWithoutExt = reg.ReplaceAllString(nameWithoutExt, "_")

	// Remove leading and trailing special characters
	nameWithoutExt = strings.Trim(nameWithoutExt, "_.-")

	// Ensure the filename is not empty after sanitization
	if nameWithoutExt == "" {
		return "", errors.New("filename is empty after sanitization")
	}

	// Reassemble the filename
	sanitizedFilename := nameWithoutExt + ext

	// Check for reserved filenames (Windows)
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	for _, reserved := range reservedNames {
		if strings.ToUpper(sanitizedFilename) == reserved || strings.ToUpper(nameWithoutExt) == reserved {
			return "", errors.New("filename is a reserved name")
		}
	}

	// Check for filename length (limit to 255 characters)
	if len(sanitizedFilename) > 255 {
		return "", errors.New("filename is too long (maximum 255 characters)")
	}

	return sanitizedFilename, nil
}

// ValidateFilename checks if a filename is valid according to common filesystem rules
func ValidateFilename(filename string) error {
	if filename == "" {
		return errors.New("filename cannot be empty")
	}

	// Check for invalid characters
	invalidChars := `<>:"/\|?*`
	for _, char := range filename {
		if strings.ContainsRune(invalidChars, char) {
			return errors.New("filename contains invalid characters")
		}
	}

	// Check for control characters
	for _, char := range filename {
		if unicode.IsControl(char) {
			return errors.New("filename contains control characters")
		}
	}

	// Check for reserved filenames (Windows)
	reservedNames := []string{
		"CON", "PRN", "AUX", "NUL",
		"COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9",
		"LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9",
	}

	nameWithoutExt := filename[:len(filename)-len(filepath.Ext(filename))]
	for _, reserved := range reservedNames {
		if strings.ToUpper(nameWithoutExt) == reserved {
			return errors.New("filename is a reserved name")
		}
	}

	// Check for filename length (limit to 255 characters)
	if len(filename) > 255 {
		return errors.New("filename is too long (maximum 255 characters)")
	}

	return nil
}

// GenerateUniqueFilename generates a unique filename by appending a number
// if the filename already exists
func GenerateUniqueFilename(baseFilename string, checkExists func(string) bool) string {
	if !checkExists(baseFilename) {
		return baseFilename
	}

	ext := filepath.Ext(baseFilename)
	nameWithoutExt := baseFilename[:len(baseFilename)-len(ext)]

	counter := 1
	for {
		newFilename := fmt.Sprintf("%s_%d%s", nameWithoutExt, counter, ext)
		if !checkExists(newFilename) {
			return newFilename
		}
		counter++
	}
}
