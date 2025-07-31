import { test, expect } from "@playwright/test";
import path from "path";

// Test file paths for different media types
const testFiles = {
  image: path.join(__dirname, "..", "..", "fixtures", "test-image.jpg"),
  document: path.join(__dirname, "..", "..", "fixtures", "test-document.pdf"),
  video: path.join(__dirname, "..", "..", "fixtures", "test-video.mp4"),
  audio: path.join(__dirname, "..", "..", "fixtures", "test-audio.mp3"),
};

/**
 * Test suite for the unified media UI components
 */
test.describe("Unified Media UI", () => {
  
  /**
   * Tests that the unified media upload page loads correctly
   */
  test("should load unified media upload page", async ({ page }) => {
    await page.goto("/upload/media");
    
    await expect(page).toHaveTitle(/Go-Fast CDN/);
    await expect(page.getByRole("heading", { name: "Upload Files" })).toBeVisible();
    await expect(page.getByTestId("unified-media-upload")).toBeVisible();
  });

  /**
   * Tests that the unified media files page loads correctly for images
   */
  test("should load unified media files page for images", async ({ page }) => {
    await page.goto("/media/images");
    
    await expect(page).toHaveTitle(/Go-Fast CDN/);
    await expect(page.getByRole("heading", { name: "Images" })).toBeVisible();
    await expect(page.getByPlaceholder("Search files by name")).toBeVisible();
  });

  /**
   * Tests that the unified media files page loads correctly for documents
   */
  test("should load unified media files page for documents", async ({ page }) => {
    await page.goto("/media/documents");
    
    await expect(page).toHaveTitle(/Go-Fast CDN/);
    await expect(page.getByRole("heading", { name: "Documents" })).toBeVisible();
    await expect(page.getByPlaceholder("Search files by name")).toBeVisible();
  });

  /**
   * Tests that the unified media files page loads correctly for videos
   */
  test("should load unified media files page for videos", async ({ page }) => {
    await page.goto("/media/videos");
    
    await expect(page).toHaveTitle(/Go-Fast CDN/);
    await expect(page.getByRole("heading", { name: "Videos" })).toBeVisible();
    await expect(page.getByPlaceholder("Search files by name")).toBeVisible();
  });

  /**
   * Tests that the unified media files page loads correctly for audio
   */
  test("should load unified media files page for audio", async ({ page }) => {
    await page.goto("/media/audio");
    
    await expect(page).toHaveTitle(/Go-Fast CDN/);
    await expect(page.getByRole("heading", { name: "Audio" })).toBeVisible();
    await expect(page.getByPlaceholder("Search files by name")).toBeVisible();
  });

  /**
   * Tests that the unified media upload component works for images
   */
  test("should upload image using unified media upload", async ({ page }) => {
    await page.goto("/upload/media");
    
    // Set media type to image
    await page.getByRole("button", { name: "Images" }).click();
    
    // Upload file
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.image);
    
    // Wait for file to appear in upload area
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Click upload button
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Verify success message
    await expect(page.getByText("Successfully uploaded media!")).toBeVisible();
  });

  /**
   * Tests that the unified media upload component works for documents
   */
  test("should upload document using unified media upload", async ({ page }) => {
    await page.goto("/upload/media");
    
    // Set media type to document
    await page.getByRole("button", { name: "Documents" }).click();
    
    // Upload file
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.document);
    
    // Wait for file to appear in upload area
    await expect(page.getByText("test-document.pdf")).toBeVisible();
    
    // Click upload button
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Verify success message
    await expect(page.getByText("Successfully uploaded media!")).toBeVisible();
  });

  /**
   * Tests that the unified media upload component works for videos
   */
  test("should upload video using unified media upload", async ({ page }) => {
    await page.goto("/upload/media");
    
    // Set media type to video
    await page.getByRole("button", { name: "Videos" }).click();
    
    // Upload file
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.video);
    
    // Wait for file to appear in upload area
    await expect(page.getByText("test-video.mp4")).toBeVisible();
    
    // Click upload button
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Verify success message
    await expect(page.getByText("Successfully uploaded media!")).toBeVisible();
  });

  /**
   * Tests that the unified media upload component works for audio
   */
  test("should upload audio using unified media upload", async ({ page }) => {
    await page.goto("/upload/media");
    
    // Set media type to audio
    await page.getByRole("button", { name: "Audio" }).click();
    
    // Upload file
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.audio);
    
    // Wait for file to appear in upload area
    await expect(page.getByText("test-audio.mp3")).toBeVisible();
    
    // Click upload button
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Verify success message
    await expect(page.getByText("Successfully uploaded media!")).toBeVisible();
  });

  /**
   * Tests that the unified media upload component handles file rejection
   */
  test("should reject invalid file type", async ({ page }) => {
    await page.goto("/upload/media");
    
    // Set media type to image
    await page.getByRole("button", { name: "Images" }).click();
    
    // Create a fake file with invalid extension
    const invalidFile = path.join(__dirname, "..", "..", "fixtures", "test-file.txt");
    
    // Upload file
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(invalidFile);
    
    // Wait for error message
    await expect(page.getByText("Invalid file type")).toBeVisible();
  });

  /**
   * Tests that the media card displays correctly for images
   */
  test("should display image media card correctly", async ({ page }) => {
    // First upload an image
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Images" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.image);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to images page
    await page.goto("/media/images");
    
    // Wait for images to load
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Click on the image to open modal
    await page.getByText("test-image.jpg").click();
    
    // Verify modal content
    await expect(page.getByRole("heading", { name: "test-image.jpg" })).toBeVisible();
    await expect(page.getByText("Filename")).toBeVisible();
    await expect(page.getByText("File Size")).toBeVisible();
    await expect(page.getByText("Media Type")).toBeVisible();
    await expect(page.getByText("Created")).toBeVisible();
    await expect(page.getByText("Checksum")).toBeVisible();
    await expect(page.getByText("Width")).toBeVisible();
    await expect(page.getByText("Height")).toBeVisible();
    
    // Close modal
    await page.getByRole("button", { name: "Close" }).click();
  });

  /**
   * Tests that the media card displays correctly for documents
   */
  test("should display document media card correctly", async ({ page }) => {
    // First upload a document
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Documents" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.document);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to documents page
    await page.goto("/media/documents");
    
    // Wait for documents to load
    await expect(page.getByText("test-document.pdf")).toBeVisible();
    
    // Click on the document to open modal
    await page.getByText("test-document.pdf").click();
    
    // Verify modal content
    await expect(page.getByRole("heading", { name: "test-document.pdf" })).toBeVisible();
    await expect(page.getByText("Filename")).toBeVisible();
    await expect(page.getByText("File Size")).toBeVisible();
    await expect(page.getByText("Media Type")).toBeVisible();
    await expect(page.getByText("Created")).toBeVisible();
    await expect(page.getByText("Checksum")).toBeVisible();
    
    // Close modal
    await page.getByRole("button", { name: "Close" }).click();
  });

  /**
   * Tests that the media card displays correctly for videos
   */
  test("should display video media card correctly", async ({ page }) => {
    // First upload a video
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Videos" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.video);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to videos page
    await page.goto("/media/videos");
    
    // Wait for videos to load
    await expect(page.getByText("test-video.mp4")).toBeVisible();
    
    // Click on the video to open modal
    await page.getByText("test-video.mp4").click();
    
    // Verify modal content
    await expect(page.getByRole("heading", { name: "test-video.mp4" })).toBeVisible();
    await expect(page.getByText("Filename")).toBeVisible();
    await expect(page.getByText("File Size")).toBeVisible();
    await expect(page.getByText("Media Type")).toBeVisible();
    await expect(page.getByText("Created")).toBeVisible();
    await expect(page.getByText("Checksum")).toBeVisible();
    
    // Close modal
    await page.getByRole("button", { name: "Close" }).click();
  });

  /**
   * Tests that the media card displays correctly for audio
   */
  test("should display audio media card correctly", async ({ page }) => {
    // First upload an audio file
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Audio" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.audio);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to audio page
    await page.goto("/media/audio");
    
    // Wait for audio files to load
    await expect(page.getByText("test-audio.mp3")).toBeVisible();
    
    // Click on the audio file to open modal
    await page.getByText("test-audio.mp3").click();
    
    // Verify modal content
    await expect(page.getByRole("heading", { name: "test-audio.mp3" })).toBeVisible();
    await expect(page.getByText("Filename")).toBeVisible();
    await expect(page.getByText("File Size")).toBeVisible();
    await expect(page.getByText("Media Type")).toBeVisible();
    await expect(page.getByText("Created")).toBeVisible();
    await expect(page.getByText("Checksum")).toBeVisible();
    
    // Close modal
    await page.getByRole("button", { name: "Close" }).click();
  });

  /**
   * Tests that the search functionality works correctly
   */
  test("should filter media files by search", async ({ page }) => {
    // First upload an image
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Images" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.image);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to images page
    await page.goto("/media/images");
    
    // Wait for images to load
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Search for the file
    await page.getByPlaceholder("Search files by name").fill("test-image");
    
    // Verify the file is still visible
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Search for something else
    await page.getByPlaceholder("Search files by name").fill("non-existent");
    
    // Verify no files are shown
    await expect(page.getByText("test-image.jpg")).not.toBeVisible();
    
    // Clear search
    await page.getByRole("button", { name: "Clear Search" }).click();
    
    // Verify the file is visible again
    await expect(page.getByText("test-image.jpg")).toBeVisible();
  });

  /**
   * Tests that the copy link functionality works
   */
  test("should copy link to clipboard", async ({ page }) => {
    // First upload an image
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Images" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.image);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to images page
    await page.goto("/media/images");
    
    // Wait for images to load
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Click on copy link button
    await page.getByTitle("Copy Link to clipboard").click();
    
    // Verify success message
    await expect(page.getByText("Link copied to clipboard")).toBeVisible();
  });

  /**
   * Tests that the download functionality works
   */
  test("should download file", async ({ page }) => {
    // First upload an image
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Images" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.image);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to images page
    await page.goto("/media/images");
    
    // Wait for images to load
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Click on download button
    const downloadPromise = page.waitForEvent('download');
    await page.getByTitle("Download file").click();
    const download = await downloadPromise;
    
    // Verify download
    expect(download.suggestedFilename()).toBe("test-image.jpg");
  });

  /**
   * Tests that the rename functionality works
   */
  test("should rename file", async ({ page }) => {
    // First upload an image
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Images" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.image);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to images page
    await page.goto("/media/images");
    
    // Wait for images to load
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Click on rename button
    await page.getByTitle("Rename file").click();
    
    // Enter new name
    await page.getByLabel("New Filename").fill("renamed-image.jpg");
    
    // Click rename button
    await page.getByRole("button", { name: "Rename" }).click();
    
    // Verify success message
    await expect(page.getByText("Successfully renamed media!")).toBeVisible();
    
    // Verify the file is renamed
    await expect(page.getByText("renamed-image.jpg")).toBeVisible();
    await expect(page.getByText("test-image.jpg")).not.toBeVisible();
  });

  /**
   * Tests that the resize functionality works for images
   */
  test("should resize image", async ({ page }) => {
    // First upload an image
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Images" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.image);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to images page
    await page.goto("/media/images");
    
    // Wait for images to load
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Click on resize button
    await page.getByTitle("Resize image").click();
    
    // Enter new dimensions
    await page.getByLabel("Width").fill("100");
    await page.getByLabel("Height").fill("100");
    
    // Click resize button
    await page.getByRole("button", { name: "Resize" }).click();
    
    // Verify success message
    await expect(page.getByText("Successfully resized image to 100x100!")).toBeVisible();
  });

  /**
   * Tests that the delete functionality works
   */
  test("should delete file", async ({ page }) => {
    // First upload an image
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Images" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles(testFiles.image);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to images page
    await page.goto("/media/images");
    
    // Wait for images to load
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Click on delete button
    await page.getByTitle("Delete file").click();
    
    // Confirm deletion
    await page.getByRole("button", { name: "Continue" }).click();
    
    // Verify success message
    await expect(page.getByText("Successfully deleted media!")).toBeVisible();
    
    // Verify the file is deleted
    await expect(page.getByText("test-image.jpg")).not.toBeVisible();
  });

  /**
   * Tests that the bulk selection and deletion works
   */
  test("should select and delete multiple files", async ({ page }) => {
    // Upload multiple images
    await page.goto("/upload/media");
    await page.getByRole("button", { name: "Images" }).click();
    
    const fileInput = page.locator('input[type="file"]');
    await fileInput.setInputFiles([testFiles.image]);
    await page.getByRole("button", { name: "Upload (1)" }).click();
    
    // Wait for upload to complete
    await expect(page.getByText("Upload (1)")).not.toBeVisible();
    
    // Navigate to images page
    await page.goto("/media/images");
    
    // Wait for images to load
    await expect(page.getByText("test-image.jpg")).toBeVisible();
    
    // Enable selection mode
    await page.getByRole("button", { name: "Select" }).click();
    
    // Select the file
    await page.locator("input[type='checkbox']").check();
    
    // Verify selection counter
    await expect(page.getByText("1 File Selected")).toBeVisible();
    
    // Delete selected files
    await page.getByRole("button", { name: "Delete Selected Files" }).click();
    
    // Confirm deletion
    await page.getByRole("button", { name: "Continue" }).click();
    
    // Verify success message
    await expect(page.getByText("Successfully deleted media!")).toBeVisible();
    
    // Verify the file is deleted
    await expect(page.getByText("test-image.jpg")).not.toBeVisible();
    
    // Verify selection mode is disabled
    await expect(page.getByRole("button", { name: "Select" })).toBeVisible();
  });

  /**
   * Tests backward compatibility with legacy upload pages
   */
  test.describe("Backward Compatibility", () => {
    
    /**
     * Tests that the legacy image upload page still works
     */
    test("should load legacy image upload page", async ({ page }) => {
      await page.goto("/upload/images");
      
      await expect(page).toHaveTitle(/Go-Fast CDN/);
      await expect(page.getByRole("heading", { name: "Upload Files" })).toBeVisible();
      await expect(page.getByRole("button", { name: "Images" })).toBeVisible();
      await expect(page.getByRole("button", { name: "Documents" })).toBeVisible();
    });

    /**
     * Tests that the legacy document upload page still works
     */
    test("should load legacy document upload page", async ({ page }) => {
      await page.goto("/upload/docs");
      
      await expect(page).toHaveTitle(/Go-Fast CDN/);
      await expect(page.getByRole("heading", { name: "Upload Files" })).toBeVisible();
      await expect(page.getByRole("button", { name: "Images" })).toBeVisible();
      await expect(page.getByRole("button", { name: "Documents" })).toBeVisible();
    });

    /**
     * Tests that the legacy images page still works
     */
    test("should load legacy images page", async ({ page }) => {
      await page.goto("/images");
      
      await expect(page).toHaveTitle(/Go-Fast CDN/);
      await expect(page.getByRole("heading", { name: "Images" })).toBeVisible();
      await expect(page.getByPlaceholder("Search files by name")).toBeVisible();
    });

    /**
     * Tests that the legacy documents page still works
     */
    test("should load legacy documents page", async ({ page }) => {
      await page.goto("/docs");
      
      await expect(page).toHaveTitle(/Go-Fast CDN/);
      await expect(page.getByRole("heading", { name: "Documents" })).toBeVisible();
      await expect(page.getByPlaceholder("Search files by name")).toBeVisible();
    });
  });
});