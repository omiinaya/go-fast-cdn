# Media Repository Unification Analysis

## 1. Introduction

### Purpose of the Analysis
This document provides a comprehensive analysis of the current implementation of separate image and document repositories in the go-fast-cdn project and outlines the strategy for unifying them into a single media repository. The goal is to evaluate the benefits, challenges, and implementation approach for this architectural change.

### Current State (Separate Image and Document Repositories)
The go-fast-cdn project currently implements two separate repositories for handling different types of media files:
- **Image Repository**: Handles image files (JPEG, PNG, GIF, WebP, BMP)
- **Document Repository**: Handles document files (Plain text, Word documents, PDFs, etc.)

These repositories have parallel but separate implementations across all layers of the application, resulting in code duplication and maintenance overhead.

### Goal (Unified Media Repository)
The objective is to create a unified media repository that can handle all types of media files through a single, extensible framework. This unified approach will:
- Reduce code duplication
- Improve maintainability
- Simplify the addition of new media types
- Provide a consistent API for all media operations

## 2. Current Implementation Analysis

### Backend Implementation for Images
The image repository implementation consists of the following components:

**Model** ([`src/models/imageModel.go`](src/models/imageModel.go:5-18)):
```go
type Image struct {
    gorm.Model
    FileName string `json:"file_name"`
    Checksum []byte `json:"checksum"`
}

type ImageRepository interface {
    GetAllImages() []Image
    GetImageByCheckSum(checksum []byte) Image
    AddImage(image Image) (string, error)
    DeleteImage(fileName string) (string, bool)
    RenameImage(oldFileName, newFileName string) error
}
```

**Database Operations** ([`src/database/image.go`](src/database/image.go:8-57)):
- `imageRepo` struct implements the `ImageRepository` interface
- Standard CRUD operations for image entities
- MD5 checksum-based duplicate detection

**Handlers** ([`src/handlers/image/`](src/handlers/image/)):
- `ImageHandler` struct for coordinating operations
- Specialized handlers for each operation:
  - [`handleImageUpload.go`](src/handlers/image/handleImageUpload.go:13-100): Handles image upload with MIME type validation
  - [`handleImageDelete.go`](src/handlers/image/handleImageDelete.go): Handles image deletion
  - [`handleImageRename.go`](src/handlers/image/handleImageRename.go): Handles image renaming
  - [`handleImageResize.go`](src/handlers/image/handleImageResize.go): Handles image resizing (image-specific functionality)
  - [`handleAllImages.go`](src/handlers/image/handleAllImages.go): Retrieves all images
  - [`handleImageMetadata.go`](src/handlers/image/handleImageMetadata.go): Retrieves image metadata

**Image-Specific Features**:
- MIME type validation for images only
- Image resizing functionality
- Storage in `uploads/images/` directory

### Backend Implementation for Documents
The document repository implementation mirrors the image repository with document-specific variations:

**Model** ([`src/models/docModel.go`](src/models/docModel.go:7-20)):
```go
type Doc struct {
    gorm.Model
    FileName string `json:"file_name"`
    Checksum []byte `json:"checksum"`
}

type DocRepository interface {
    GetAllDocs() []Doc
    GetDocByCheckSum(checksum []byte) Doc
    AddDoc(doc Doc) (string, error)
    DeleteDoc(fileName string) (string, bool)
    RenameDoc(oldFileName, newFileName string) error
}
```

**Database Operations** ([`src/database/doc.go`](src/database/doc.go:8-57)):
- `DocRepo` struct implements the `DocRepository` interface
- Nearly identical implementation to `imageRepo`
- MD5 checksum-based duplicate detection

**Handlers** ([`src/handlers/docs/`](src/handlers/docs/)):
- `DocHandler` struct for coordinating operations
- Specialized handlers for each operation:
  - [`handleDocUpload.go`](src/handlers/docs/handleDocUpload.go:13-97): Handles document upload with MIME type validation
  - [`handleDocDelete.go`](src/handlers/docs/handleDocDelete.go): Handles document deletion
  - [`handleDocsRename.go`](src/handlers/docs/handleDocsRename.go): Handles document renaming
  - [`handleAllDocs.go`](src/handlers/docs/handleAllDocs.go): Retrieves all documents
  - [`handleDocMetadata.go`](src/handlers/docs/handleDocMetadata.go): Retrieves document metadata

**Document-Specific Features**:
- MIME type validation for documents only
- Storage in `uploads/docs/` directory
- No document-specific processing equivalent to image resizing

### Frontend Implementation for Images
The frontend implementation for images includes:

**Type Definitions** ([`ui/src/types/file.ts`](ui/src/types/file.ts:1-8), [`ui/src/types/fileMetadata.ts`](ui/src/types/fileMetadata.ts:1-7)):
```typescript
export type TFile = {
  ID: number;
  CreatedAt: string;
  UpdatedAt: string;
  DeletedAt: string | null;
  file_name: string;
  checksum: string;
};

export type FileMetadata = {
  download_url: string;
  file_size: number;
  filename: string;
  height?: number;  // Image-specific
  width?: number;   // Image-specific
};
```

**Components**:
- [`ImageCardUpload`](ui/src/modules/content/upload/image-card-upload.tsx:8-26): Component for displaying uploaded images
- [`Files`](ui/src/modules/content/files.tsx:33-227): Main component for displaying files with type parameter ("images" | "documents")

**Hooks**:
- [`useGetFilesQuery`](ui/src/modules/content/hooks/use-get-files-query.tsx:9-21): Fetches files with type parameter
- [`useUploadFileMutation`](ui/src/modules/content/hooks/use-upload-file-mutation.tsx:4-19): Handles file uploads with type parameter
- [`useDeleteFileMutation`](ui/src/modules/content/hooks/use-delete-file-mutation.tsx): Handles file deletion with type parameter
- [`useRenameFileMutation`](ui/src/modules/content/hooks/use-rename-file-mutation.tsx): Handles file renaming with type parameter
- [`useResizeImageMutation`](ui/src/modules/content/hooks/use-resize-image-mutation.tsx): Image-specific hook for resizing

### Frontend Implementation for Documents
The frontend implementation for documents is similar to images but with document-specific components:

**Components**:
- [`DocCardUpload`](ui/src/modules/content/upload/doc-card-upload.tsx:7-18): Component for displaying uploaded documents
- [`Files`](ui/src/modules/content/files.tsx:33-227): Reused main component with type parameter set to "documents"

**Hooks**:
- Same hooks as images but with type parameter set to "doc" or "documents"
- No document-specific processing equivalent to image resizing

### Database Schema for Images and Documents
The current database schema maintains separate tables for images and documents:

**Images Table** (defined in [`src/models/imageModel.go`](src/models/imageModel.go:5-10)):
```sql
CREATE TABLE images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME NULL,
    file_name VARCHAR(255),
    checksum BLOB
);
```

**Documents Table** (defined in [`src/models/docModel.go`](src/models/docModel.go:7-12)):
```sql
CREATE TABLE docs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME NULL,
    file_name VARCHAR(255),
    checksum BLOB
);
```

Both tables have identical structures, differing only in the table name.

### API Endpoints for Images and Documents
The API endpoints are organized in separate groups for images and documents:

**Image Endpoints** (defined in [`src/router/api.go`](src/router/api.go:56-93)):
```go
// Public endpoints
cdn.GET("/image/all", imageHandler.HandleAllImages)
cdn.GET("/image/:filename", iHandlers.HandleImageMetadata)
cdn.Static("/download/images", util.ExPath+"/uploads/images")

// Protected endpoints
upload.POST("/image", imageHandler.HandleImageUpload)
delete.DELETE("/image/:filename", imageHandler.HandleImageDelete)
rename.PUT("/image", imageHandler.HandleImageRename)
resize.PUT("/image", iHandlers.HandleImageResize)  // Image-specific
```

**Document Endpoints** (defined in [`src/router/api.go`](src/router/api.go:54-88)):
```go
// Public endpoints
cdn.GET("/doc/all", docHandler.HandleAllDocs)
cdn.GET("/doc/:filename", dHandlers.HandleDocMetadata)
cdn.Static("/download/docs", util.ExPath+"/uploads/docs")

// Protected endpoints
upload.POST("/doc", docHandler.HandleDocUpload)
delete.DELETE("/doc/:filename", docHandler.HandleDocDelete)
rename.PUT("/doc", docHandler.HandleDocsRename)
// No document-specific equivalent to image resize
```

## 3. Components Requiring Changes

### Backend Components

#### Models
**Current State**:
- Separate models: [`Image`](src/models/imageModel.go:5-10) and [`Doc`](src/models/docModel.go:7-12)
- Separate interfaces: [`ImageRepository`](src/models/imageModel.go:12-18) and [`DocRepository`](src/models/docModel.go:14-20)

**Required Changes**:
- Create a unified `Media` model with a `MediaType` field to distinguish between different media types
- Create a unified `MediaRepository` interface that handles all media types
- Implement type-specific methods within the unified interface for operations like image resizing

#### Database Operations
**Current State**:
- Separate implementations: [`imageRepo`](src/database/image.go:8-57) and [`DocRepo`](src/database/doc.go:8-57)
- Nearly identical code with only model type differences

**Required Changes**:
- Create a unified `mediaRepo` implementation
- Implement type-specific handling where necessary (e.g., image processing)
- Update the migration function in [`src/database/migrate.go`](src/database/migrate.go:7-9) to create the unified media table

#### Handlers
**Current State**:
- Separate handlers: [`ImageHandler`](src/handlers/image/ImageHandler.go:5-11) and [`DocHandler`](src/handlers/docs/DocHandler.go:7-13)
- Separate handler files for each operation (upload, delete, rename, etc.)

**Required Changes**:
- Create a unified `MediaHandler` that can handle all media types
- Implement type-specific logic within handlers where necessary
- Consolidate handler files to reduce duplication
- Create a strategy pattern for handling media-specific operations

#### API Routes
**Current State**:
- Separate route groups for images and documents
- Type-specific endpoints (e.g., `/cdn/upload/image`, `/cdn/upload/doc`)

**Required Changes**:
- Create unified endpoints that accept a media type parameter
- Maintain backward compatibility with existing endpoints
- Implement type-specific routing for operations like image resizing

#### Utilities
**Current State**:
- File storage in separate directories: `uploads/images/` and `uploads/docs/`
- Separate file path handling logic

**Required Changes**:
- Implement a unified file storage structure with media type subdirectories
- Create utility functions for handling media type-specific paths
- Update file deletion and renaming utilities to work with the unified structure

### Frontend Components

#### Type Definitions
**Current State**:
- Generic [`TFile`](ui/src/types/file.ts:1-8) type used for both images and documents
- [`FileMetadata`](ui/src/types/fileMetadata.ts:1-7) with optional image-specific fields (height, width)

**Required Changes**:
- Create a more comprehensive type definition that explicitly includes media type
- Define type-specific metadata interfaces
- Update all components to use the new type definitions

#### Components
**Current State**:
- Separate components for image and document upload cards
- Generic [`Files`](ui/src/modules/content/files.tsx:33-227) component with type parameter

**Required Changes**:
- Create a unified upload card component that can handle different media types
- Update the Files component to work with the unified media type
- Implement conditional rendering based on media type where necessary

#### Services
**Current State**:
- API calls with type parameters (e.g., `/upload/${type}`)
- Separate handling for image-specific operations like resizing

**Required Changes**:
- Update API service functions to work with unified endpoints
- Implement media type-specific API calls where necessary
- Maintain backward compatibility with existing API calls

#### Hooks
**Current State**:
- Hooks with type parameters (e.g., [`useGetFilesQuery`](ui/src/modules/content/hooks/use-get-files-query.tsx:9-21))
- Image-specific hooks like [`useResizeImageMutation`](ui/src/modules/content/hooks/use-resize-image-mutation.tsx)

**Required Changes**:
- Update hooks to work with unified media types
- Implement conditional logic for media-specific operations
- Create a more flexible hook system that can handle different media types

### Database Schema Changes

**Current State**:
- Separate tables: `images` and `docs`
- Identical schema structure

**Required Changes**:
- Create a unified `media` table with a `media_type` field
- Implement a migration strategy to preserve existing data
- Update all database queries to work with the unified table
- Consider indexing strategies for the media_type field for optimal performance

**Proposed Unified Schema**:
```sql
CREATE TABLE media (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at DATETIME,
    updated_at DATETIME,
    deleted_at DATETIME NULL,
    file_name VARCHAR(255),
    checksum BLOB,
    media_type VARCHAR(50),  -- 'image', 'document', etc.
    -- Additional metadata fields as needed
);
```

### File Storage Structure Changes

**Current State**:
- Separate directories: `uploads/images/` and `uploads/docs/`

**Required Changes**:
- Implement a unified structure: `uploads/media/{media_type}/`
- Create a migration strategy for existing files
- Update file path handling in all components
- Implement media type-specific subdirectories for organization

## 4. Benefits of Unification

### Reduced Code Duplication
The current implementation has significant code duplication between image and document repositories:

**Backend Duplication**:
- Model definitions are nearly identical
- Database repository implementations are almost the same
- Handler patterns are duplicated with only minor variations
- API route definitions follow the same pattern

**Frontend Duplication**:
- Type definitions are shared but with conditional fields
- Components have similar structures with media-specific variations
- Hooks and services implement the same patterns with type parameters

**Estimated Reduction**:
- Approximately 40-50% reduction in backend code related to media handling
- Approximately 30-40% reduction in frontend code related to media handling
- Simplified maintenance and testing efforts

### Improved Maintainability
A unified media repository would significantly improve maintainability:

**Single Point of Change**:
- Bug fixes and improvements only need to be implemented once
- Feature additions automatically apply to all media types
- Reduced risk of inconsistencies between implementations

**Simplified Testing**:
- Fewer test cases needed for common functionality
- Easier to ensure consistent behavior across media types
- Reduced testing overhead for new features

**Code Organization**:
- Clearer separation of concerns
- More logical code structure
- Easier for new developers to understand and contribute

### Easier Extensibility for Future Media Types
The unified approach would make it significantly easier to add support for new media types:

**Current Process for Adding a New Media Type**:
1. Create a new model struct
2. Implement a new repository interface
3. Create handler files for all operations
4. Add API routes
5. Update frontend types and components
6. Create frontend hooks and services
7. Update database schema

**Process with Unified Repository**:
1. Add the new media type to the media type enumeration
2. Implement media-specific validation and processing (if needed)
3. Update frontend display logic (if needed)

**Examples of Future Media Types**:
- Video files (MP4, WebM, etc.)
- Audio files (MP3, WAV, etc.)
- Archive files (ZIP, RAR, etc.)
- 3D model files (OBJ, STL, etc.)

### Simplified User Experience
A unified media repository would provide a more consistent user experience:

**Consistent Interface**:
- Same upload, management, and interaction patterns for all media types
- Reduced user confusion when working with different file types
- Streamlined workflow for managing mixed media content

**Enhanced Features**:
- Unified search and filtering across all media types
- Consistent metadata handling
- Easier implementation of cross-media features like galleries or collections

**API Consistency**:
- Simplified API client implementation
- Reduced need for type-specific API calls
- More predictable API behavior

## 5. Challenges and Considerations

### Data Migration
Migrating from separate image and document repositories to a unified media repository presents several challenges:

**Database Migration**:
- Need to migrate existing data from `images` and `docs` tables to a unified `media` table
- Must preserve all existing metadata and relationships
- Requires careful planning to avoid data loss

**File System Migration**:
- Need to reorganize file storage from `uploads/images/` and `uploads/docs/` to `uploads/media/{media_type}/`
- Must update all file references in the database
- Requires significant disk I/O operations for large file repositories

**Migration Strategy**:
1. Create the unified `media` table alongside existing tables
2. Implement a migration script to copy data from existing tables
3. Update file paths in the database
4. Physically move files to the new directory structure
5. Update application code to use the unified repository
6. Remove old tables and directories after verification

**Downtime Considerations**:
- Plan for maintenance window during migration
- Implement read-only mode during critical migration phases
- Consider a phased approach to minimize disruption

### Backward Compatibility
Maintaining backward compatibility is crucial for existing clients and integrations:

**API Compatibility**:
- Existing endpoints must continue to work
- Response formats should remain consistent
- Need to implement endpoint redirection or proxying

**Frontend Compatibility**:
- Existing UI components must continue to function
- User workflows should not be disrupted
- Need to maintain consistent behavior during transition

**Compatibility Strategies**:
- Implement proxy endpoints that redirect old URLs to new unified endpoints
- Maintain old database tables as views during transition period
- Create compatibility layers in the codebase
- Implement feature flags to enable gradual migration

### Type-Specific Functionality
Different media types have unique requirements and processing needs:

**Image-Specific Functionality**:
- Image resizing and thumbnail generation
- EXIF data extraction
- Color space manipulation
- Format conversion

**Document-Specific Functionality**:
- Text extraction and indexing
- Page count and document metadata
- Format conversion (PDF to text, etc.)
- Document preview generation

**Handling Strategies**:
- Implement a strategy pattern for media-specific operations
- Create plugin architecture for media processors
- Use factory pattern to create appropriate handlers based on media type
- Define clear interfaces for media-specific operations

**Implementation Example**:
```go
type MediaProcessor interface {
    Process(file *os.File) error
    GetMetadata(file *os.File) (map[string]interface{}, error)
}

type ImageProcessor struct{}
type DocumentProcessor struct{}

func GetProcessor(mediaType string) MediaProcessor {
    switch mediaType {
    case "image":
        return &ImageProcessor{}
    case "document":
        return &DocumentProcessor{}
    default:
        return &DefaultProcessor{}
    }
}
```

### File Type Validation and Handling
Different media types require different validation and handling approaches:

**Current Validation**:
- Images: MIME type checking for specific image formats
- Documents: MIME type checking for specific document formats
- Separate validation logic in each handler

**Unified Validation Challenges**:
- Need to support a wider range of file types
- Must handle type-specific validation rules
- Security considerations for different file types

**Validation Strategies**:
- Create a comprehensive file type registry
- Implement MIME type detection with magic number verification
- Define validation rules for each media type
- Implement file content scanning for security

**Security Considerations**:
- Malicious file detection
- File size limits per media type
- Virus scanning for executable content
- Sanitization of user-provided metadata

**Implementation Example**:
```go
type MediaValidator interface {
    Validate(file *os.File) error
    GetAllowedTypes() []string
}

type ImageValidator struct{}
type DocumentValidator struct{}

func GetValidator(mediaType string) MediaValidator {
    switch mediaType {
    case "image":
        return &ImageValidator{}
    case "document":
        return &DocumentValidator{}
    default:
        return &DefaultValidator{}
    }
}
```

## Conclusion

The unification of image and document repositories into a single media repository represents a significant architectural improvement for the go-fast-cdn project. While the implementation requires careful planning and execution, the benefits in terms of reduced code duplication, improved maintainability, and enhanced extensibility make this a worthwhile endeavor.

The key to success lies in:
1. Thorough planning of the migration strategy
2. Maintaining backward compatibility throughout the transition
3. Implementing flexible patterns for handling media-specific functionality
4. Ensuring robust validation and security measures

This analysis provides a foundation for implementing the unified media repository and should serve as a reference for the development team throughout the process.