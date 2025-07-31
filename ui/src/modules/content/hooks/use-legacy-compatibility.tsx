import useGetFilesQueryBase from "./use-get-files-query";
import useUploadFileMutationBase from "./use-upload-file-mutation";
import useDeleteFileMutationBase from "./use-delete-file-mutation";
import useRenameFileMutationBase from "./use-rename-file-mutation";
import useResizeImageMutationBase from "./use-resize-image-mutation";
import useGetFileDataQueryBase from "./use-get-file-data-query";
import useResizeModalQueryBase from "./use-resize-modal-query";

/**
 * Legacy compatibility hooks that maintain the existing API while using the unified media infrastructure
 * These hooks ensure backward compatibility for existing components
 */

// Re-export the updated hooks with their original signatures
export const useGetFilesQuery = useGetFilesQueryBase;
export const useUploadFileMutation = useUploadFileMutationBase;
export const useDeleteFileMutation = useDeleteFileMutationBase;
export const useRenameFileMutation = useRenameFileMutationBase;
export const useResizeImageMutation = useResizeImageMutationBase;
export const useGetFileDataQuery = useGetFileDataQueryBase;
export const useResizeModalQuery = useResizeModalQueryBase;

/**
 * Helper function to check if the unified media migration is complete
 * This can be used to gradually transition components to the new unified hooks
 */
export const useUnifiedMediaMigrationStatus = () => {
  // This hook can be expanded to check migration status from the backend
  // For now, we'll assume the migration is complete
  return {
    isMigrationComplete: true,
    canUseUnifiedMedia: true,
  };
};

/**
 * Migration helper hook that provides both legacy and unified hooks
 * This allows components to gradually transition to the unified media system
 */
export const useMediaMigrationHelper = () => {
  const { isMigrationComplete } = useUnifiedMediaMigrationStatus();
  
  return {
    // Legacy hooks (for backward compatibility)
    legacy: {
      useGetFilesQuery,
      useUploadFileMutation,
      useDeleteFileMutation,
      useRenameFileMutation,
      useResizeImageMutation,
      useGetFileDataQuery,
      useResizeModalQuery,
    },
    
    // Flag to indicate if unified media should be used
    shouldUseUnifiedMedia: isMigrationComplete,
  };
};

export default useMediaMigrationHelper;