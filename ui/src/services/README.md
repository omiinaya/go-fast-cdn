# Unified Media Services

This directory contains the unified media services for the frontend that replace the separate image and document API services.

## Services

### Unified Media Service (`unifiedMediaService.ts`)

The unified media service provides a consistent interface for communicating with the backend's unified media repository. It handles all media types (images, documents, videos, audio, and other files) with a single API.

#### Methods

- `getAllMedia(mediaType: MediaType)`: Get all media of a specific type
- `uploadMedia(file: File, mediaType: MediaType)`: Upload media of a specific type
- `deleteMedia(fileName: string, mediaType: MediaType)`: Delete media of a specific type
- `renameMedia(fileName: string, newFileName: string, mediaType: MediaType)`: Rename media of a specific type
- `resizeMedia(fileName: string, width: number, height: number)`: Resize image media
- `getMediaMetadata(fileName: string, mediaType: MediaType)`: Get metadata for media of a specific type

### Legacy Services

The following legacy services are still available for backward compatibility:

- `authService.ts`: Authentication service with CDN API client
- `adminUserService.ts`: Admin user service
- `configService.ts`: Configuration service

## Hooks

### Unified Media Hooks

The following hooks are available for use with React components:

- `useGetUnifiedMediaQuery({ mediaType })`: Get all media of a specific type
- `useUploadUnifiedMediaMutation()`: Upload media of a specific type
- `useDeleteUnifiedMediaMutation()`: Delete media of a specific type
- `useRenameUnifiedMediaMutation()`: Rename media of a specific type
- `useResizeUnifiedMediaMutation()`: Resize image media

### Legacy Hooks

The following legacy hooks are still available for backward compatibility:

- `useGetFilesQuery({ type })`: Get files of a specific type (images or documents)
- `useGetMediaQuery({ mediaType })`: Get media of a specific type
- `useUploadFileMutation()`: Upload files of a specific type
- `useUploadMediaMutation()`: Upload media of a specific type
- `useDeleteFileMutation(type)`: Delete files of a specific type
- `useDeleteMediaMutation()`: Delete media of a specific type
- `useRenameFileMutation(type)`: Rename files of a specific type
- `useRenameMediaMutation()`: Rename media of a specific type
- `useResizeImageMutation()`: Resize images
- `useResizeMediaMutation()`: Resize media

## Components

### Unified Media Components

The following components have been updated to use the unified media services:

- `MediaCard`: Displays media items with actions (delete, rename, resize, download)
- `UploadMediaModal`: Modal for uploading media files
- `RenameMediaModal`: Modal for renaming media files
- `ResizeMediaModal`: Modal for resizing image media
- `MediaFiles`: Page for displaying media files

### Legacy Components

The following components still use the legacy services for backward compatibility:

- `ContentCard`: Displays content items with actions
- `UploadModal`: Modal for uploading files
- `RenameModal`: Modal for renaming files
- `ResizeModal`: Modal for resizing images
- `Files`: Page for displaying files

## Migration Guide

### Migrating from Legacy Services to Unified Services

1. Replace `useGetFilesQuery` with `useGetUnifiedMediaQuery`
   ```typescript
   // Before
   const files = useGetFilesQuery({ type: "images" });
   
   // After
   const media = useGetUnifiedMediaQuery({ mediaType: "image" });
   ```

2. Replace `useUploadFileMutation` with `useUploadUnifiedMediaMutation`
   ```typescript
   // Before
   const uploadFile = useUploadFileMutation();
   uploadFile.mutate({ file, type: "image" });
   
   // After
   const uploadMedia = useUploadUnifiedMediaMutation();
   uploadMedia.mutate({ file, mediaType: "image" });
   ```

3. Replace `useDeleteFileMutation` with `useDeleteUnifiedMediaMutation`
   ```typescript
   // Before
   const deleteFile = useDeleteFileMutation("image");
   deleteFile.mutate(fileName);
   
   // After
   const deleteMedia = useDeleteUnifiedMediaMutation();
   deleteMedia.mutate({ fileName, mediaType: "image" });
   ```

4. Replace `useRenameFileMutation` with `useRenameUnifiedMediaMutation`
   ```typescript
   // Before
   const renameFile = useRenameFileMutation("image");
   renameFile.mutate(formData);
   
   // After
   const renameMedia = useRenameUnifiedMediaMutation();
   renameMedia.mutate({ fileName, newFileName, mediaType: "image" });
   ```

5. Replace `useResizeImageMutation` with `useResizeUnifiedMediaMutation`
   ```typescript
   // Before
   const resizeImage = useResizeImageMutation();
   resizeImage.mutate({ filename, width, height });
   
   // After
   const resizeMedia = useResizeUnifiedMediaMutation();
   resizeMedia.mutate({ media: imageMedia, width, height });
   ```

## Testing

The unified media services are tested with Jest. To run the tests:

```bash
npm test -- unifiedMediaService.test.ts
```

## Backward Compatibility

The legacy services and components are still available and functional. This ensures that existing functionality continues to work while new features can be developed using the unified services.

## Future Enhancements

1. Add support for more media types (3D models, VR content, etc.)
2. Implement batch operations for media management
3. Add media transcoding capabilities
4. Implement media versioning
5. Add media analytics and usage statistics