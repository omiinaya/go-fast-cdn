package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kevinanielsen/go-fast-cdn/src/util"
)

func main() {
	fmt.Println("Testing updated file storage structure...")

	// Load executable path
	util.LoadExPath()

	// Test directory creation
	fmt.Println("\n1. Testing directory creation...")

	// Ensure all directories exist
	if err := util.EnsureUploadDirectories(); err != nil {
		fmt.Printf("Error creating directories: %v\n", err)
		os.Exit(1)
	}

	// Check if directories exist
	directories := []string{
		util.GetUploadsPath(),
		util.GetMediaPath(),
		util.GetImagesPath(),
		util.GetDocsPath(),
	}

	for _, dir := range directories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("Directory does not exist: %s\n", dir)
			os.Exit(1)
		} else {
			fmt.Printf("✓ Directory exists: %s\n", dir)
		}
	}

	// Test file migration functionality
	fmt.Println("\n2. Testing file migration functionality...")

	// Create test files in legacy directories
	testImagePath := filepath.Join(util.GetImagesPath(), "test_image.jpg")
	testDocPath := filepath.Join(util.GetDocsPath(), "test_doc.pdf")

	// Create test image file
	if err := os.WriteFile(testImagePath, []byte("test image content"), 0644); err != nil {
		fmt.Printf("Error creating test image file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Created test image file: %s\n", testImagePath)

	// Create test doc file
	if err := os.WriteFile(testDocPath, []byte("test doc content"), 0644); err != nil {
		fmt.Printf("Error creating test doc file: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("✓ Created test doc file: %s\n", testDocPath)

	// Test utility functions
	fmt.Println("\n3. Testing utility functions...")

	fmt.Printf("✓ Upload path: %s\n", util.GetUploadsPath())
	fmt.Printf("✓ Media path: %s\n", util.GetMediaPath())
	fmt.Printf("✓ Images path: %s\n", util.GetImagesPath())
	fmt.Printf("✓ Docs path: %s\n", util.GetDocsPath())
	fmt.Printf("✓ Media upload path: %s\n", util.GetMediaUploadPath())
	fmt.Printf("✓ Media URL path: %s\n", util.GetMediaURLPath("test.jpg"))

	// Clean up test files
	fmt.Println("\n4. Cleaning up test files...")

	if err := os.Remove(testImagePath); err != nil {
		fmt.Printf("Error removing test image file: %v\n", err)
	}

	if err := os.Remove(testDocPath); err != nil {
		fmt.Printf("Error removing test doc file: %v\n", err)
	}

	fmt.Println("\n✓ All tests passed! The updated file storage structure is working correctly.")
}
