/**
 * Index file for all media-related hooks
 * This file provides a centralized export point for all hooks
 */

// Legacy hooks (for backward compatibility)
export { default as useGetFilesQuery } from './use-get-files-query';
export { default as useUploadFileMutation } from './use-upload-file-mutation';
export { default as useDeleteFileMutation } from './use-delete-file-mutation';
export { default as useRenameFileMutation } from './use-rename-file-mutation';
export { default as useResizeImageMutation } from './use-resize-image-mutation';
export { default as useGetFileDataQuery } from './use-get-file-data-query';
export { default as useResizeModalQuery } from './use-resize-modal-query';

// Media hooks (partially unified)
export { default as useGetMediaQuery } from './use-get-media-query';
export { default as useUploadMediaMutation } from './use-upload-media-mutation';
export { default as useDeleteMediaMutation } from './use-delete-media-mutation';
export { default as useRenameMediaMutation } from './use-rename-media-mutation';
export { default as useResizeMediaMutation } from './use-resize-media-mutation';
export { default as useGetMediaMetadataQuery } from './use-get-media-metadata-query';

// Unified media hooks (fully unified)
export { default as useGetAllMediaQuery } from './use-get-all-media-query';
export { default as useUploadUnifiedMediaMutation } from './use-upload-unified-media-mutation';
export { default as useDeleteUnifiedMediaMutation } from './use-delete-unified-media-mutation';
export { default as useRenameUnifiedMediaMutation } from './use-rename-unified-media-mutation';
export { default as useResizeUnifiedMediaMutation } from './use-resize-unified-media-mutation';

// Comprehensive unified media hook
export { default as useUnifiedMedia, useUnifiedMedia as useMedia } from './use-unified-media';

// Legacy compatibility helpers
export { default as useMediaMigrationHelper } from './use-legacy-compatibility';

// Size query hook (used by multiple media types)
export { default as useGetSizeQuery } from './use-get-size-query';