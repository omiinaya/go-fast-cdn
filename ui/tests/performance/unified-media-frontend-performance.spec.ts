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
 * Performance test suite for unified media frontend components vs legacy components
 */
test.describe("Unified Media Frontend Performance", () => {
  
  /**
   * Tests the performance of unified media upload page vs legacy upload pages
   */
  test("upload page performance - unified vs legacy", async ({ page }) => {
    // Test legacy image upload page performance
    await page.goto("/upload/image");
    const legacyImageLoadStart = Date.now();
    await page.waitForSelector('[data-testid="upload-form"]');
    const legacyImageLoadTime = Date.now() - legacyImageLoadStart;
    
    // Test legacy document upload page performance
    await page.goto("/upload/doc");
    const legacyDocLoadStart = Date.now();
    await page.waitForSelector('[data-testid="upload-form"]');
    const legacyDocLoadTime = Date.now() - legacyDocLoadStart;
    
    // Test unified media upload page performance
    await page.goto("/upload/media");
    const unifiedLoadStart = Date.now();
    await page.waitForSelector('[data-testid="upload-form"]');
    const unifiedLoadTime = Date.now() - unifiedLoadStart;
    
    // Log performance metrics
    console.log(`Legacy Image Upload Page Load Time: ${legacyImageLoadTime}ms`);
    console.log(`Legacy Document Upload Page Load Time: ${legacyDocLoadTime}ms`);
    console.log(`Unified Media Upload Page Load Time: ${unifiedLoadTime}ms`);
    
    // Performance assertions
    expect(unifiedLoadTime).toBeLessThan(legacyImageLoadTime + legacyDocLoadTime);
  });

  /**
   * Tests the performance of unified media files page vs legacy files pages
   */
  test("files page performance - unified vs legacy", async ({ page }) => {
    // Test legacy images page performance
    await page.goto("/files/images");
    const legacyImagesLoadStart = Date.now();
    await page.waitForSelector('[data-testid="files-grid"]');
    const legacyImagesLoadTime = Date.now() - legacyImagesLoadStart;
    
    // Test legacy documents page performance
    await page.goto("/files/docs");
    const legacyDocsLoadStart = Date.now();
    await page.waitForSelector('[data-testid="files-grid"]');
    const legacyDocsLoadTime = Date.now() - legacyDocsLoadStart;
    
    // Test unified media files page performance
    await page.goto("/files/media");
    const unifiedLoadStart = Date.now();
    await page.waitForSelector('[data-testid="files-grid"]');
    const unifiedLoadTime = Date.now() - unifiedLoadStart;
    
    // Log performance metrics
    console.log(`Legacy Images Page Load Time: ${legacyImagesLoadTime}ms`);
    console.log(`Legacy Documents Page Load Time: ${legacyDocsLoadTime}ms`);
    console.log(`Unified Media Files Page Load Time: ${unifiedLoadTime}ms`);
    
    // Performance assertions
    expect(unifiedLoadTime).toBeLessThan(legacyImagesLoadTime + legacyDocsLoadTime);
  });

  /**
   * Tests the performance of media upload operations
   */
  test("media upload performance - unified vs legacy", async ({ page }) => {
    // Test legacy image upload performance
    await page.goto("/upload/image");
    const legacyImageUploadStart = Date.now();
    
    // Set up file input change listener
    const fileInput = await page.locator('input[type="file"]').first();
    await fileInput.setInputFiles(testFiles.image);
    
    // Wait for upload to complete
    await page.waitForSelector('[data-testid="upload-success"]');
    const legacyImageUploadTime = Date.now() - legacyImageUploadStart;
    
    // Test legacy document upload performance
    await page.goto("/upload/doc");
    const legacyDocUploadStart = Date.now();
    
    // Set up file input change listener
    const docFileInput = await page.locator('input[type="file"]').first();
    await docFileInput.setInputFiles(testFiles.document);
    
    // Wait for upload to complete
    await page.waitForSelector('[data-testid="upload-success"]');
    const legacyDocUploadTime = Date.now() - legacyDocUploadStart;
    
    // Test unified media upload performance
    await page.goto("/upload/media");
    const unifiedUploadStart = Date.now();
    
    // Set up file input change listener
    const unifiedFileInput = await page.locator('input[type="file"]').first();
    await unifiedFileInput.setInputFiles(testFiles.image);
    
    // Wait for upload to complete
    await page.waitForSelector('[data-testid="upload-success"]');
    const unifiedUploadTime = Date.now() - unifiedUploadStart;
    
    // Log performance metrics
    console.log(`Legacy Image Upload Time: ${legacyImageUploadTime}ms`);
    console.log(`Legacy Document Upload Time: ${legacyDocUploadTime}ms`);
    console.log(`Unified Media Upload Time: ${unifiedUploadTime}ms`);
    
    // Performance assertions
    expect(unifiedUploadTime).toBeLessThan(legacyImageUploadTime * 1.5); // Allow some overhead for unified handling
  });

  /**
   * Tests the performance of media display and rendering
   */
  test("media display performance - unified vs legacy", async ({ page }) => {
    // Test legacy image display performance
    await page.goto("/files/images");
    const legacyImageDisplayStart = Date.now();
    
    // Wait for images to load
    await page.waitForSelector('[data-testid="media-card"] img');
    const legacyImageDisplayTime = Date.now() - legacyImageDisplayStart;
    
    // Test unified media display performance
    await page.goto("/files/media?type=image");
    const unifiedDisplayStart = Date.now();
    
    // Wait for media to load
    await page.waitForSelector('[data-testid="media-card"] img');
    const unifiedDisplayTime = Date.now() - unifiedDisplayStart;
    
    // Log performance metrics
    console.log(`Legacy Image Display Time: ${legacyImageDisplayTime}ms`);
    console.log(`Unified Media Display Time: ${unifiedDisplayTime}ms`);
    
    // Performance assertions
    expect(unifiedDisplayTime).toBeLessThan(legacyImageDisplayTime * 1.2); // Allow some overhead for unified handling
  });

  /**
   * Tests the performance of search and filtering operations
   */
  test("search and filter performance - unified vs legacy", async ({ page }) => {
    // Test legacy image search performance
    await page.goto("/files/images");
    const legacySearchStart = Date.now();
    
    // Perform search
    await page.fill('[data-testid="search-input"]', "test");
    await page.press('[data-testid="search-input"]', "Enter");
    
    // Wait for search results
    await page.waitForSelector('[data-testid="search-results"]');
    const legacySearchTime = Date.now() - legacySearchStart;
    
    // Test unified media search performance
    await page.goto("/files/media");
    const unifiedSearchStart = Date.now();
    
    // Perform search
    await page.fill('[data-testid="search-input"]', "test");
    await page.press('[data-testid="search-input"]', "Enter");
    
    // Wait for search results
    await page.waitForSelector('[data-testid="search-results"]');
    const unifiedSearchTime = Date.now() - unifiedSearchStart;
    
    // Log performance metrics
    console.log(`Legacy Image Search Time: ${legacySearchTime}ms`);
    console.log(`Unified Media Search Time: ${unifiedSearchTime}ms`);
    
    // Performance assertions
    expect(unifiedSearchTime).toBeLessThan(legacySearchTime * 1.3); // Allow some overhead for unified handling
  });

  /**
   * Tests the performance of bulk operations
   */
  test("bulk operations performance - unified vs legacy", async ({ page }) => {
    // Test legacy bulk delete performance
    await page.goto("/files/images");
    const legacyBulkStart = Date.now();
    
    // Select multiple items
    const checkboxes = await page.locator('[data-testid="select-checkbox"]').first();
    await checkboxes.check();
    
    // Perform bulk delete
    await page.click('[data-testid="bulk-delete-button"]');
    await page.click('[data-testid="confirm-delete-button"]');
    
    // Wait for operation to complete
    await page.waitForSelector('[data-testid="delete-success"]');
    const legacyBulkTime = Date.now() - legacyBulkStart;
    
    // Test unified bulk delete performance
    await page.goto("/files/media");
    const unifiedBulkStart = Date.now();
    
    // Select multiple items
    const unifiedCheckboxes = await page.locator('[data-testid="select-checkbox"]').first();
    await unifiedCheckboxes.check();
    
    // Perform bulk delete
    await page.click('[data-testid="bulk-delete-button"]');
    await page.click('[data-testid="confirm-delete-button"]');
    
    // Wait for operation to complete
    await page.waitForSelector('[data-testid="delete-success"]');
    const unifiedBulkTime = Date.now() - unifiedBulkStart;
    
    // Log performance metrics
    console.log(`Legacy Bulk Delete Time: ${legacyBulkTime}ms`);
    console.log(`Unified Bulk Delete Time: ${unifiedBulkTime}ms`);
    
    // Performance assertions
    expect(unifiedBulkTime).toBeLessThan(legacyBulkTime * 1.3); // Allow some overhead for unified handling
  });

  /**
   * Tests the performance with different media types and sizes
   */
  test("different media types performance", async ({ page }) => {
    const mediaTypes = [
      { type: "image", file: testFiles.image, expectedMaxTime: 5000 },
      { type: "document", file: testFiles.document, expectedMaxTime: 3000 },
      { type: "video", file: testFiles.video, expectedMaxTime: 10000 },
      { type: "audio", file: testFiles.audio, expectedMaxTime: 5000 },
    ];

    for (const media of mediaTypes) {
      await page.goto("/upload/media");
      const uploadStart = Date.now();
      
      // Upload media file
      const fileInput = await page.locator('input[type="file"]').first();
      await fileInput.setInputFiles(media.file);
      
      // Wait for upload to complete
      await page.waitForSelector('[data-testid="upload-success"]');
      const uploadTime = Date.now() - uploadStart;
      
      // Log performance metrics
      console.log(`${media.type} Upload Time: ${uploadTime}ms`);
      
      // Performance assertions
      expect(uploadTime).toBeLessThan(media.expectedMaxTime);
    }
  });

  /**
   * Tests the performance under load (simulating multiple concurrent users)
   */
  test("concurrent user performance", async ({ browser }) => {
    const numUsers = 5;
    const promises: Promise<number>[] = [];
    
    // Simulate multiple concurrent users
    for (let i = 0; i < numUsers; i++) {
      const context = await browser.newContext();
      const page = await context.newPage();
      
      promises.push((async () => {
        const start = Date.now();
        
        // Navigate to unified media page
        await page.goto("/files/media");
        
        // Perform search
        await page.fill('[data-testid="search-input"]', "test");
        await page.press('[data-testid="search-input"]', "Enter");
        
        // Wait for results
        await page.waitForSelector('[data-testid="search-results"]');
        
        return Date.now() - start;
      })());
    }
    
    // Wait for all users to complete
    const results = await Promise.all(promises);
    
    // Calculate average performance
    const avgTime = results.reduce((sum, time) => sum + time, 0) / results.length;
    const maxTime = Math.max(...results);
    const minTime = Math.min(...results);
    
    // Log performance metrics
    console.log(`Concurrent User Performance:`);
    console.log(`Average Time: ${avgTime}ms`);
    console.log(`Max Time: ${maxTime}ms`);
    console.log(`Min Time: ${minTime}ms`);
    
    // Performance assertions
    expect(avgTime).toBeLessThan(5000); // Average should be under 5 seconds
    expect(maxTime).toBeLessThan(10000); // No user should take more than 10 seconds
  });
});