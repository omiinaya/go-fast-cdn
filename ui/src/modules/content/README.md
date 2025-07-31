# Unified Media Components

This directory contains the updated content components that handle all media types in a unified way, replacing the separate image and document components.

## Components

### MediaCard
A unified card component that can display any type of media (images, documents, videos, audio, etc.).

**Props:**
- `media`: Media object (required)
- `disabled`: boolean (optional)
- `isSelected`: boolean (optional)
- `onSelect`: function (optional)
- `isSelecting`: boolean (optional)

**Usage:**
```tsx
import MediaCard from "@/modules/content/media-card";

<MediaCard 
  media={mediaObject} 
  isSelecting={isSelecting}
  isSelected={selectedMedia.includes(mediaObject.fileName)}
  onSelect={handleOnSelectMedia}
/>
```

### MediaFiles
A component that displays a grid of media files with search and selection functionality.

**Props:**
- `mediaType`: MediaType (required) - 'image' | 'document' | 'video' | 'audio' | 'other'

**Usage:**
```tsx
import MediaFiles from "@/modules/content/media-files";

<MediaFiles mediaType="image" />
```

### MediaDataModal
A modal that displays detailed information about a media file, with type-specific properties.

**Props:**
- `media`: Media object (required)

**Usage:**
```tsx
import MediaDataModal from "@/modules/content/media-data-modal";

<MediaDataModal media={mediaObject} />
```

### RenameMediaModal
A modal that allows renaming a media file.

**Props:**
- `media`: Media object (required)
- `isSelecting`: boolean (optional)

**Usage:**
```tsx
import RenameMediaModal from "@/modules/content/rename-media-modal";

<RenameMediaModal media={mediaObject} />
```

### ResizeMediaModal
A modal that allows resizing images (only available for image media).

**Props:**
- `media`: Media object (required)
- `isSelecting`: boolean (optional)

**Usage:**
```tsx
import ResizeMediaModal from "@/modules/content/resize-media-modal";

<ResizeMediaModal media={mediaObject} />
```

### UploadMediaModal
A modal for uploading media files with support for all media types.

**Props:**
- `placement`: 'header' | 'sidebar' (optional, default: 'sidebar')
- `mediaType`: MediaType (optional, required when placement is 'header')

**Usage:**
```tsx
import UploadMediaModal from "@/modules/content/upload/upload-media-modal";

// In header with specific media type
<UploadMediaModal placement="header" mediaType="image" />

// In sidebar with media type selection
<UploadMediaModal />
```

### UploadMediaForm
A form component for uploading media files with drag-and-drop support.

**Props:**
- `files`: File[] (required)
- `onChangeFiles`: function (required)
- `isLoading`: boolean (required)
- `mediaType`: MediaType (required)
- `onChangeMediaType`: function (required)
- `disableMediaTypeSwitching`: boolean (optional, default: false)

**Usage:**
```tsx
import UploadMediaForm from "@/modules/content/upload/upload-media-form";

<UploadMediaForm
  isLoading={isUploadPending}
  mediaType={selectedMediaType}
  onChangeMediaType={setSelectedMediaType}
  files={files}
  onChangeFiles={setFiles}
  disableMediaTypeSwitching={placement === "header"}
/>
```

### MediaCardUpload
A card component for displaying files in the upload queue.

**Props:**
- `file`: File object (required)
- `onClickDelete`: function (required)
- `fileName`: string (optional)
- `mediaType`: MediaType (required)

**Usage:**
```tsx
import MediaCardUpload from "@/modules/content/upload/media-card-upload";

<MediaCardUpload
  file={file}
  fileName={sanitizeFileName(file).name}
  mediaType={mediaType}
  onClickDelete={() => handleDeleteFile(index)}
/>
```

## Hooks

### useGetMediaQuery
A hook for fetching media files of a specific type.

**Parameters:**
- `mediaType`: MediaType (required)

**Usage:**
```tsx
import useGetMediaQuery from "@/modules/content/hooks/use-get-media-query";

const media = useGetMediaQuery({ mediaType: "image" });
```

### useDeleteMediaMutation
A hook for deleting media files.

**Usage:**
```tsx
import useDeleteMediaMutation from "@/modules/content/hooks/use-delete-media-mutation";

const deleteMedia = useDeleteMediaMutation();

const handleDelete = () => {
  deleteMedia.mutate({ fileName: media.fileName, mediaType: media.mediaType });
};
```

### useUploadMediaMutation
A hook for uploading media files.

**Usage:**
```tsx
import useUploadMediaMutation from "@/modules/content/hooks/use-upload-media-mutation";

const uploadMedia = useUploadMediaMutation();

const handleUpload = () => {
  uploadMedia.mutate({ file, mediaType: "image" });
};
```

### useResizeMediaMutation
A hook for resizing image media.

**Usage:**
```tsx
import useResizeMediaMutation from "@/modules/content/hooks/use-resize-media-mutation";

const resizeMedia = useResizeMediaMutation({
  onSuccess: () => {
    toast.success("Image resized!");
  },
});

const handleResize = () => {
  resizeMedia.mutate({
    media: imageMedia,
    width: 800,
    height: 600,
  });
};
```

### useRenameMediaMutation
A hook for renaming media files.

**Usage:**
```tsx
import useRenameMediaMutation from "@/modules/content/hooks/use-rename-media-mutation";

const renameMedia = useRenameMediaMutation({
  onSuccess: () => {
    toast.success("File renamed!");
  },
});

const handleRename = () => {
  renameMedia.mutate({
    fileName: media.fileName,
    newFileName: "new-name.jpg",
    mediaType: media.mediaType,
    apiEndpoint: getMediaApiEndpoint(media),
  });
};
```

## Migration from Legacy Components

### For Images
- Replace `Files` component with `MediaFiles` component with `mediaType="image"`
- Replace `ContentCard` component with `MediaCard` component
- Replace `UploadModal` component with `UploadMediaModal` component

### For Documents
- Replace `Files` component with `MediaFiles` component with `mediaType="document"`
- Replace `ContentCard` component with `MediaCard` component
- Replace `UploadModal` component with `UploadMediaModal` component

## Backward Compatibility

The existing image and document components remain functional to ensure backward compatibility. The new unified components can be gradually adopted as needed.

## Type Definitions

The components use the unified `Media` type defined in `@/types/media.ts`, which includes:

- `BaseMedia`: Common properties for all media types
- `ImageMedia`: Image-specific properties (width, height, altText, thumbnailUrl)
- `DocumentMedia`: Document-specific properties (pageCount, author, subject, keywords)
- `VideoMedia`: Video-specific properties (duration, width, height, thumbnailUrl)
- `AudioMedia`: Audio-specific properties (duration, artist, album, genre)
- `OtherMedia`: Properties for other media types (mimeType)

Type guard functions are provided to check the media type:
- `isImageMedia(media: Media)`
- `isDocumentMedia(media: Media)`
- `isVideoMedia(media: Media)`
- `isAudioMedia(media: Media)`
- `isOtherMedia(media: Media)`