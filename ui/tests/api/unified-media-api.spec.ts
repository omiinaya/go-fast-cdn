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
 * Test suite for the unified media API endpoints
 */
test.describe("Unified Media API", () => {
  
  /**
   * Tests that the /api/cdn/media endpoint returns a 200 OK status
   * when retrieving all media without type filter
   */
  test("should return all media types", async ({ request }) => {
    const response = await request.get("/api/cdn/media");
    expect(response.status()).toBe(200);
    expect(response.ok()).toBe(true);
    
    const media = await response.json();
    expect(Array.isArray(media)).toBe(true);
  });

  /**
   * Tests that the /api/cdn/media endpoint returns filtered media
   * when type parameter is provided
   */
  test("should return filtered media by type", async ({ request }) => {
    const types = ["image", "document", "video", "audio"];
    
    for (const type of types) {
      const response = await request.get(`/api/cdn/media?type=${type}`);
      expect(response.status()).toBe(200);
      expect(response.ok()).toBe(true);
      
      const media = await response.json();
      expect(Array.isArray(media)).toBe(true);
      
      // Verify all returned media are of the requested type
      media.forEach((item: any) => {
        expect(item.type).toBe(type);
      });
    }
  });

  /**
   * Tests that the /api/cdn/media/upload endpoint handles image uploads
   */
  test("should upload image media", async ({ request }) => {
    const formData = new FormData();
    formData.append("file", testFiles.image);
    formData.append("filename", "test-upload-image");
    
    const response = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(response.status()).toBe(200);
    expect(response.ok()).toBe(true);
    
    const result = await response.json();
    expect(result).toHaveProperty("file_url");
    expect(result).toHaveProperty("type");
    expect(result.type).toBe("image");
  });

  /**
   * Tests that the /api/cdn/media/upload endpoint handles document uploads
   */
  test("should upload document media", async ({ request }) => {
    const formData = new FormData();
    formData.append("file", testFiles.document);
    formData.append("filename", "test-upload-document");
    
    const response = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(response.status()).toBe(200);
    expect(response.ok()).toBe(true);
    
    const result = await response.json();
    expect(result).toHaveProperty("file_url");
    expect(result).toHaveProperty("type");
    expect(result.type).toBe("document");
  });

  /**
   * Tests that the /api/cdn/media/upload endpoint handles video uploads
   */
  test("should upload video media", async ({ request }) => {
    const formData = new FormData();
    formData.append("file", testFiles.video);
    formData.append("filename", "test-upload-video");
    
    const response = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(response.status()).toBe(200);
    expect(response.ok()).toBe(true);
    
    const result = await response.json();
    expect(result).toHaveProperty("file_url");
    expect(result).toHaveProperty("type");
    expect(result.type).toBe("video");
  });

  /**
   * Tests that the /api/cdn/media/upload endpoint handles audio uploads
   */
  test("should upload audio media", async ({ request }) => {
    const formData = new FormData();
    formData.append("file", testFiles.audio);
    formData.append("filename", "test-upload-audio");
    
    const response = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(response.status()).toBe(200);
    expect(response.ok()).toBe(true);
    
    const result = await response.json();
    expect(result).toHaveProperty("file_url");
    expect(result).toHaveProperty("type");
    expect(result.type).toBe("audio");
  });

  /**
   * Tests that the /api/cdn/media/upload endpoint rejects duplicate files
   */
  test("should reject duplicate file upload", async ({ request }) => {
    // First upload
    const formData = new FormData();
    formData.append("file", testFiles.image);
    formData.append("filename", "test-duplicate-image");
    
    const firstResponse = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(firstResponse.status()).toBe(200);
    
    // Second upload of the same file
    const secondResponse = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(secondResponse.status()).toBe(409);
    const error = await secondResponse.json();
    expect(error).toHaveProperty("error");
    expect(error.error).toBe("File already exists");
  });

  /**
   * Tests that the /api/cdn/media/upload endpoint rejects invalid file types
   */
  test("should reject invalid file type", async ({ request }) => {
    const formData = new FormData();
    // Create a fake file with invalid extension
    const invalidFile = new Blob(["invalid content"], { type: "application/octet-stream" });
    formData.append("file", invalidFile, "test.invalid");
    
    const response = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(response.status()).toBe(400);
    const error = await response.json();
    expect(error).toHaveProperty("error");
  });

  /**
   * Tests that the /api/cdn/media/:filename endpoint returns metadata for images
   */
  test("should return image metadata", async ({ request }) => {
    // First upload an image
    const formData = new FormData();
    formData.append("file", testFiles.image);
    formData.append("filename", "test-metadata-image");
    
    const uploadResponse = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(uploadResponse.status()).toBe(200);
    const uploadResult = await uploadResponse.json();
    const fileName = uploadResult.file_url.split("/").pop() || "test-metadata-image.jpg";
    
    // Get metadata
    const metadataResponse = await request.get(`/api/cdn/media/${fileName}?type=image`);
    expect(metadataResponse.status()).toBe(200);
    
    const metadata = await metadataResponse.json();
    expect(metadata).toHaveProperty("filename");
    expect(metadata).toHaveProperty("download_url");
    expect(metadata).toHaveProperty("file_size");
    expect(metadata).toHaveProperty("type");
    expect(metadata.type).toBe("image");
    expect(metadata).toHaveProperty("width");
    expect(metadata).toHaveProperty("height");
  });

  /**
   * Tests that the /api/cdn/media/:filename endpoint returns metadata for documents
   */
  test("should return document metadata", async ({ request }) => {
    // First upload a document
    const formData = new FormData();
    formData.append("file", testFiles.document);
    formData.append("filename", "test-metadata-document");
    
    const uploadResponse = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(uploadResponse.status()).toBe(200);
    const uploadResult = await uploadResponse.json();
    const fileName = uploadResult.file_url.split("/").pop() || "test-metadata-document.pdf";
    
    // Get metadata
    const metadataResponse = await request.get(`/api/cdn/media/${fileName}?type=document`);
    expect(metadataResponse.status()).toBe(200);
    
    const metadata = await metadataResponse.json();
    expect(metadata).toHaveProperty("filename");
    expect(metadata).toHaveProperty("download_url");
    expect(metadata).toHaveProperty("file_size");
    expect(metadata).toHaveProperty("type");
    expect(metadata.type).toBe("document");
  });

  /**
   * Tests that the /api/cdn/media/:filename endpoint returns 404 for non-existent files
   */
  test("should return 404 for non-existent file metadata", async ({ request }) => {
    const response = await request.get("/api/cdn/media/non-existent-file.jpg?type=image");
    expect(response.status()).toBe(404);
    
    const error = await response.json();
    expect(error).toHaveProperty("error");
    expect(error.error).toBe("Media not found");
  });

  /**
   * Tests that the /api/cdn/media/:filename endpoint returns 400 for missing type parameter
   */
  test("should return 400 for missing type parameter", async ({ request }) => {
    const response = await request.get("/api/cdn/media/test-file.jpg");
    expect(response.status()).toBe(400);
    
    const error = await response.json();
    expect(error).toHaveProperty("error");
    expect(error.error).toBe("Media type is required");
  });

  /**
   * Tests that the /api/cdn/media/delete endpoint deletes media files
   */
  test("should delete media file", async ({ request }) => {
    // First upload an image
    const formData = new FormData();
    formData.append("file", testFiles.image);
    formData.append("filename", "test-delete-image");
    
    const uploadResponse = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(uploadResponse.status()).toBe(200);
    const uploadResult = await uploadResponse.json();
    const fileName = uploadResult.file_url.split("/").pop() || "test-delete-image.jpg";
    
    // Delete the file
    const deleteResponse = await request.delete(`/api/cdn/media/delete/${fileName}?type=image`);
    expect(deleteResponse.status()).toBe(200);
    
    const deleteResult = await deleteResponse.json();
    expect(deleteResult).toHaveProperty("message");
    expect(deleteResult).toHaveProperty("fileName");
    expect(deleteResult.message).toBe("Media deleted successfully");
    
    // Verify the file is deleted
    const metadataResponse = await request.get(`/api/cdn/media/${fileName}?type=image`);
    expect(metadataResponse.status()).toBe(404);
  });

  /**
   * Tests that the /api/cdn/media/delete endpoint returns 404 for non-existent files
   */
  test("should return 404 when deleting non-existent file", async ({ request }) => {
    const response = await request.delete("/api/cdn/media/delete/non-existent-file.jpg?type=image");
    expect(response.status()).toBe(404);
    
    const error = await response.json();
    expect(error).toHaveProperty("error");
    expect(error.error).toBe("Media not found");
  });

  /**
   * Tests that the /api/cdn/media/delete endpoint returns 400 for missing type parameter
   */
  test("should return 400 for missing type parameter in delete", async ({ request }) => {
    const response = await request.delete("/api/cdn/media/delete/test-file.jpg");
    expect(response.status()).toBe(400);
    
    const error = await response.json();
    expect(error).toHaveProperty("error");
    expect(error.error).toBe("Media type is required");
  });

  /**
   * Tests that the /api/cdn/media/rename endpoint renames media files
   */
  test("should rename media file", async ({ request }) => {
    // First upload an image
    const formData = new FormData();
    formData.append("file", testFiles.image);
    formData.append("filename", "test-rename-image-original");
    
    const uploadResponse = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(uploadResponse.status()).toBe(200);
    const uploadResult = await uploadResponse.json();
    const originalFileName = uploadResult.file_url.split("/").pop() || "test-rename-image-original.jpg";
    
    // Rename the file
    const renameFormData = new FormData();
    renameFormData.append("filename", originalFileName);
    renameFormData.append("newname", "test-rename-image-new");
    renameFormData.append("type", "image");
    
    const renameResponse = await request.post("/api/cdn/media/rename", {
      multipart: renameFormData,
    });
    
    expect(renameResponse.status()).toBe(200);
    
    const renameResult = await renameResponse.json();
    expect(renameResult).toHaveProperty("status");
    expect(renameResult.status).toBe("File renamed successfully");
    
    // Verify the file is accessible with the new name
    const newFileName = "test-rename-image-new.jpg";
    const metadataResponse = await request.get(`/api/cdn/media/${newFileName}?type=image`);
    expect(metadataResponse.status()).toBe(200);
    
    // Clean up - delete the renamed file
    await request.delete(`/api/cdn/media/delete/${newFileName}?type=image`);
  });

  /**
   * Tests that the /api/cdn/media/resize endpoint resizes images
   */
  test("should resize image", async ({ request }) => {
    // First upload an image
    const formData = new FormData();
    formData.append("file", testFiles.image);
    formData.append("filename", "test-resize-image");
    
    const uploadResponse = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(uploadResponse.status()).toBe(200);
    const uploadResult = await uploadResponse.json();
    const fileName = uploadResult.file_url.split("/").pop() || "test-resize-image.jpg";
    
    // Get original dimensions
    const originalMetadataResponse = await request.get(`/api/cdn/media/${fileName}?type=image`);
    expect(originalMetadataResponse.status()).toBe(200);
    const originalMetadata = await originalMetadataResponse.json();
    
    // Resize the image
    const resizeData = {
      filename: fileName,
      width: 100,
      height: 100,
    };
    
    const resizeResponse = await request.post("/api/cdn/media/resize", {
      data: resizeData,
    });
    
    expect(resizeResponse.status()).toBe(200);
    
    const resizeResult = await resizeResponse.json();
    expect(resizeResult).toHaveProperty("status");
    expect(resizeResult).toHaveProperty("width");
    expect(resizeResult).toHaveProperty("height");
    expect(resizeResult.width).toBe(100);
    expect(resizeResult.height).toBe(100);
    
    // Verify the new dimensions
    const newMetadataResponse = await request.get(`/api/cdn/media/${fileName}?type=image`);
    expect(newMetadataResponse.status()).toBe(200);
    const newMetadata = await newMetadataResponse.json();
    expect(newMetadata.width).toBe(100);
    expect(newMetadata.height).toBe(100);
    
    // Clean up - delete the resized image
    await request.delete(`/api/cdn/media/delete/${fileName}?type=image`);
  });

  /**
   * Tests that the /api/cdn/media/resize endpoint returns 400 for non-image media
   */
  test("should return 400 when trying to resize non-image media", async ({ request }) => {
    // First upload a document
    const formData = new FormData();
    formData.append("file", testFiles.document);
    formData.append("filename", "test-resize-document");
    
    const uploadResponse = await request.post("/api/cdn/media/upload", {
      multipart: formData,
    });
    
    expect(uploadResponse.status()).toBe(200);
    const uploadResult = await uploadResponse.json();
    const fileName = uploadResult.file_url.split("/").pop() || "test-resize-document.pdf";
    
    // Try to resize the document
    const resizeData = {
      filename: fileName,
      width: 100,
      height: 100,
    };
    
    const resizeResponse = await request.post("/api/cdn/media/resize", {
      data: resizeData,
    });
    
    expect(resizeResponse.status()).toBe(400);
    
    const error = await resizeResponse.json();
    expect(error).toHaveProperty("error");
    expect(error.error).toContain("Cannot resize media of type");
    
    // Clean up - delete the document
    await request.delete(`/api/cdn/media/delete/${fileName}?type=document`);
  });

  /**
   * Tests backward compatibility with image endpoints
   */
  test.describe("Backward Compatibility", () => {
    
    /**
     * Tests that the legacy /api/cdn/images endpoint still works
     */
    test("should return all images using legacy endpoint", async ({ request }) => {
      const response = await request.get("/api/cdn/images");
      expect(response.status()).toBe(200);
      expect(response.ok()).toBe(true);
      
      const images = await response.json();
      expect(Array.isArray(images)).toBe(true);
    });

    /**
     * Tests that the legacy /api/cdn/docs endpoint still works
     */
    test("should return all documents using legacy endpoint", async ({ request }) => {
      const response = await request.get("/api/cdn/docs");
      expect(response.status()).toBe(200);
      expect(response.ok()).toBe(true);
      
      const docs = await response.json();
      expect(Array.isArray(docs)).toBe(true);
    });

    /**
     * Tests that the legacy image upload endpoint still works
     */
    test("should upload image using legacy endpoint", async ({ request }) => {
      const formData = new FormData();
      formData.append("image", testFiles.image);
      formData.append("filename", "test-legacy-image");
      
      const response = await request.post("/api/cdn/upload/image", {
        multipart: formData,
      });
      
      expect(response.status()).toBe(200);
      expect(response.ok()).toBe(true);
      
      const result = await response.json();
      expect(result).toHaveProperty("file_url");
      
      // Clean up - delete the uploaded image
      const fileName = result.file_url.split("/").pop() || "test-legacy-image.jpg";
      await request.delete(`/api/cdn/image/delete/${fileName}`);
    });

    /**
     * Tests that the legacy document upload endpoint still works
     */
    test("should upload document using legacy endpoint", async ({ request }) => {
      const formData = new FormData();
      formData.append("doc", testFiles.document);
      formData.append("filename", "test-legacy-document");
      
      const response = await request.post("/api/cdn/upload/doc", {
        multipart: formData,
      });
      
      expect(response.status()).toBe(200);
      expect(response.ok()).toBe(true);
      
      const result = await response.json();
      expect(result).toHaveProperty("file_url");
      
      // Clean up - delete the uploaded document
      const fileName = result.file_url.split("/").pop() || "test-legacy-document.pdf";
      await request.delete(`/api/cdn/doc/delete/${fileName}`);
    });

    /**
     * Tests that the legacy image metadata endpoint still works
     */
    test("should return image metadata using legacy endpoint", async ({ request }) => {
      // First upload an image using the legacy endpoint
      const formData = new FormData();
      formData.append("image", testFiles.image);
      formData.append("filename", "test-legacy-metadata-image");
      
      const uploadResponse = await request.post("/api/cdn/upload/image", {
        multipart: formData,
      });
      
      expect(uploadResponse.status()).toBe(200);
      const uploadResult = await uploadResponse.json();
      const fileName = uploadResult.file_url.split("/").pop() || "test-legacy-metadata-image.jpg";
      
      // Get metadata using legacy endpoint
      const metadataResponse = await request.get(`/api/cdn/image/metadata/${fileName}`);
      expect(metadataResponse.status()).toBe(200);
      
      const metadata = await metadataResponse.json();
      expect(metadata).toHaveProperty("filename");
      expect(metadata).toHaveProperty("download_url");
      expect(metadata).toHaveProperty("file_size");
      expect(metadata).toHaveProperty("width");
      expect(metadata).toHaveProperty("height");
      
      // Clean up - delete the uploaded image
      await request.delete(`/api/cdn/image/delete/${fileName}`);
    });

    /**
     * Tests that the legacy document metadata endpoint still works
     */
    test("should return document metadata using legacy endpoint", async ({ request }) => {
      // First upload a document using the legacy endpoint
      const formData = new FormData();
      formData.append("doc", testFiles.document);
      formData.append("filename", "test-legacy-metadata-document");
      
      const uploadResponse = await request.post("/api/cdn/upload/doc", {
        multipart: formData,
      });
      
      expect(uploadResponse.status()).toBe(200);
      const uploadResult = await uploadResponse.json();
      const fileName = uploadResult.file_url.split("/").pop() || "test-legacy-metadata-document.pdf";
      
      // Get metadata using legacy endpoint
      const metadataResponse = await request.get(`/api/cdn/doc/metadata/${fileName}`);
      expect(metadataResponse.status()).toBe(200);
      
      const metadata = await metadataResponse.json();
      expect(metadata).toHaveProperty("filename");
      expect(metadata).toHaveProperty("download_url");
      expect(metadata).toHaveProperty("file_size");
      
      // Clean up - delete the uploaded document
      await request.delete(`/api/cdn/doc/delete/${fileName}`);
    });
  });
});