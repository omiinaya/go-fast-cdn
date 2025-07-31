# Media Hooks Documentation

This directory contains hooks for managing media files in the application. The hooks have been updated to work with the unified media system while maintaining backward compatibility.

## Hook Categories

### 1. Legacy Hooks (Backward Compatible)

These hooks maintain the original API while using the unified media infrastructure internally:

- `useGetFilesQuery` - Fetch files (images or documents)
- `useUploadFileMutation` - Upload files
- `useDeleteFileMutation` - Delete files
- `useRenameFileMutation` - Rename files
- `useResizeImageMutation` - Resize images
- `useGetFileDataQuery` - Get file metadata
- `useResizeModalQuery` - Get image data for resize modal

### 2. Media Hooks (Partially Unified)

These hooks work with specific media types:

- `useGetMediaQuery` - Fetch media by type
- `useUploadMediaMutation` - Upload media by type
- `useDeleteMediaMutation` - Delete media by type
- `useRenameMediaMutation` - Rename media by type
- `useResizeMediaMutation` - Resize media (images only)
- `useGetMediaMetadataQuery` - Get media metadata

### 3. Unified Media Hooks (Fully Unified)

These hooks work with the unified media system:

- `useGetAllMediaQuery` - Fetch all media (optionally filtered by type)
- `useUploadUnifiedMediaMutation` - Upload any media type
- `useDeleteUnifiedMediaMutation` - Delete any media type
- `useRenameUnifiedMediaMutation` - Rename any media type
- `useResizeUnifiedMediaMutation` - Resize images (type-specific operation)

### 4. Comprehensive Unified Media Hook

This hook provides all media operations in a single hook:

- `useUnifiedMedia` (also exported as `useMedia`) - Comprehensive hook for all media operations

### 5. Utility Hooks

- `useGetSizeQuery` - Get CDN size
- `useMediaMigrationHelper` - Helper for migrating to unified media

## Migration Guide

### For Existing Components

Existing components can continue to use the legacy hooks without any changes:

```typescript
// Before (still works)
import { useGetFilesQuery, useUploadFileMutation } from '@/modules/content/hooks';

const MyComponent = () => {
  const { data: files } = useGetFilesQuery({ type: 'images' });
  const uploadMutation = useUploadFileMutation();
  
  // ... component logic
};
```

### For New Components

New components should use the unified media hooks:

```typescript
// Recommended for new components
import { useUnifiedMedia } from '@/modules/content/hooks';

const MyComponent = () => {
  const { useGetAllMedia, useUploadMedia } = useUnifiedMedia({ mediaType: 'image' });
  const { data: media } = useGetAllMedia();
  const uploadMutation = useUploadMedia();
  
  // ... component logic
};
```

### Gradual Migration

Components can gradually migrate using the migration helper:

```typescript
import { useMediaMigrationHelper } from '@/modules/content/hooks';

const MyComponent = () => {
  const { legacy, shouldUseUnifiedMedia } = useMediaMigrationHelper();
  
  if (shouldUseUnifiedMedia) {
    // Use unified hooks
    const { useGetAllMedia } = useUnifiedMedia();
    const { data: media } = useGetAllMedia();
  } else {
    // Use legacy hooks
    const { data: files } = legacy.useGetFilesQuery({ type: 'images' });
  }
  
  // ... component logic
};
```

## Media Types

The unified media system supports the following media types:

- `image` - Image files (JPEG, PNG, GIF, etc.)
- `document` - Document files (PDF, DOC, TXT, etc.)
- `video` - Video files (MP4, AVI, MOV, etc.)
- `audio` - Audio files (MP3, WAV, etc.)
- `other` - Other file types

## Type-Specific Operations

Some operations are specific to certain media types:

- **Resize**: Only available for image media
- **Thumbnail**: Only available for image and video media
- **Metadata**: Available for all media types, with type-specific properties

## Error Handling

All hooks include proper error handling with toast notifications. Errors are logged and displayed to users with appropriate messages.

## Query Invalidation

The hooks automatically invalidate related queries when mutations succeed:

- Upload operations invalidate size and media queries
- Delete operations invalidate size and media queries
- Rename operations invalidate media queries
- Resize operations invalidate media queries

## Examples

### Basic Usage

```typescript
import { useUnifiedMedia } from '@/modules/content/hooks';

const MediaGallery = () => {
  const { 
    useGetAllMedia, 
    useUploadMedia, 
    useDeleteMedia,
    getMediaDownloadUrl 
  } = useUnifiedMedia();
  
  const { data: media, isLoading } = useGetAllMedia();
  const uploadMutation = useUploadMedia();
  const deleteMutation = useDeleteMedia();
  
  const handleUpload = (file: File) => {
    uploadMutation.mutate({ 
      file, 
      mediaType: 'image' 
    });
  };
  
  const handleDelete = (mediaItem: Media) => {
    deleteMutation.mutate({ 
      fileName: mediaItem.fileName, 
      mediaType: mediaItem.mediaType 
    });
  };
  
  if (isLoading) return <div>Loading...</div>;
  
  return (
    <div>
      {/* Upload button */}
      <input type="file" onChange={(e) => handleUpload(e.target.files[0])} />
      
      {/* Media grid */}
      <div className="media-grid">
        {media?.map((item) => (
          <div key={item.id} className="media-item">
            <img src={getMediaDownloadUrl(item.fileName, item.mediaType)} alt={item.fileName} />
            <button onClick={() => handleDelete(item)}>Delete</button>
          </div>
        ))}
      </div>
    </div>
  );
};
```

### Type-Specific Operations

```typescript
import { useUnifiedMedia } from '@/modules/content/hooks';
import { isImageMedia } from '@/types/media';

const ImageEditor = () => {
  const { 
    useGetAllMedia, 
    useResizeImage,
    checkIsImageMedia 
  } = useUnifiedMedia({ mediaType: 'image' });
  
  const { data: images } = useGetAllMedia();
  const resizeMutation = useResizeImage();
  
  const handleResize = (image: Media, width: number, height: number) => {
    if (checkIsImageMedia(image)) {
      resizeMutation.mutate({
        filename: image.fileName,
        width,
        height
      });
    }
  };
  
  return (
    <div>
      {images?.map((image) => (
        <div key={image.id}>
          <img src={image.downloadUrl} alt={image.fileName} />
          {checkIsImageMedia(image) && (
            <button onClick={() => handleResize(image, 300, 300)}>
              Resize to 300x300
            </button>
          )}
        </div>
      ))}
    </div>
  );
};