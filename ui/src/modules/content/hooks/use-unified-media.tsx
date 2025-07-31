import { constant } from "@/lib/constant";
import { mediaService, MediaUploadResponse, MediaDeleteParams, MediaRenameParams, MediaResizeParams } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { Media, MediaType, isImageMedia } from "@/types/media";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import toast from "react-hot-toast";

interface UseUnifiedMediaParams {
  mediaType?: MediaType;
  enabled?: boolean;
}

/**
 * A comprehensive hook for managing all media types in a unified way
 */
export const useUnifiedMedia = (params?: UseUnifiedMediaParams) => {
  const queryClient = useQueryClient();
  const { mediaType, enabled = true } = params || {};

  // Get all media (optionally filtered by type)
  const useGetAllMedia = () => {
    return useQuery({
      queryKey: constant.queryKeys.media(mediaType || 'all'),
      queryFn: async (): Promise<Media[]> => {
        return mediaService.getAllMedia({ mediaType });
      },
      enabled,
    });
  };

  // Get media metadata
  const useGetMediaMetadata = (fileName: string, mediaType: MediaType) => {
    return useQuery({
      queryKey: ['media-metadata', fileName, mediaType],
      queryFn: async () => {
        return mediaService.getMediaMetadata(fileName, mediaType);
      },
      enabled: !!fileName && !!mediaType && enabled,
    });
  };

  // Upload media
  const useUploadMedia = () => {
    return useMutation({
      mutationFn: async ({ file, mediaType, filename }: { 
        file: File; 
        mediaType: MediaType; 
        filename?: string; 
      }): Promise<MediaUploadResponse> => {
        return mediaService.uploadMedia(file, mediaType, filename);
      },
      onSuccess: (_data: MediaUploadResponse, { mediaType }) => {
        toast.dismiss();
        toast.success("Successfully uploaded media!");
        queryClient.invalidateQueries({
          queryKey: constant.queryKeys.size(),
        });
        queryClient.invalidateQueries({
          queryKey: constant.queryKeys.media(mediaType),
        });
      },
      onError: (error: unknown) => {
        const err = error as AxiosError<IErrorResponse>;
        toast.dismiss();
        const message =
          err.response?.data?.error || err.message || "Upload failed";
        toast.error(message);
      },
    });
  };

  // Delete media
  const useDeleteMedia = () => {
    return useMutation({
      mutationFn: async (params: MediaDeleteParams) => {
        return mediaService.deleteMedia(params);
      },
      onSuccess: (_data: { message: string; fileName: string }, { mediaType }: MediaDeleteParams) => {
        toast.dismiss();
        toast.success("Successfully deleted media!");
        queryClient.invalidateQueries({
          queryKey: constant.queryKeys.size(),
        });
        queryClient.invalidateQueries({
          queryKey: constant.queryKeys.media(mediaType),
        });
      },
      onError: (error: unknown) => {
        const err = error as AxiosError<IErrorResponse>;
        toast.dismiss();
        const message =
          err.response?.data?.error || err.message || "Delete failed";
        toast.error(message);
      },
    });
  };

  // Rename media
  const useRenameMedia = (options?: {
    onSuccess?: () => void;
    onError?: (error: Error) => void;
  }) => {
    return useMutation({
      mutationFn: async (params: MediaRenameParams) => {
        return mediaService.renameMedia(params);
      },
      onSuccess: (_data: { status: string }, { mediaType }: MediaRenameParams) => {
        toast.dismiss();
        toast.success("Successfully renamed media!");
        queryClient.invalidateQueries({
          queryKey: constant.queryKeys.media(mediaType),
        });
        options?.onSuccess?.();
      },
      onError: (error: unknown) => {
        const err = error as AxiosError<IErrorResponse>;
        toast.dismiss();
        const message =
          err.response?.data?.error || err.message || "Rename failed";
        toast.error(message);
        options?.onError?.(new Error(message));
      },
    });
  };

  // Resize image (type-specific operation)
  const useResizeImage = (options?: {
    onSuccess?: () => void;
    onError?: (error: Error) => void;
  }) => {
    return useMutation({
      mutationFn: async (params: MediaResizeParams) => {
        return mediaService.resizeImage(params);
      },
      onSuccess: (data) => {
        toast.dismiss();
        toast.success(`Successfully resized image to ${data.width}x${data.height}!`);
        queryClient.invalidateQueries({
          queryKey: constant.queryKeys.media("image"),
        });
        options?.onSuccess?.();
      },
      onError: (error: unknown) => {
        const err = error as AxiosError<IErrorResponse>;
        toast.dismiss();
        const message =
          err.response?.data?.error || err.message || "Resize failed";
        toast.error(message);
        options?.onError?.(new Error(message));
      },
    });
  };

  // Get media download URL
  const getMediaDownloadUrl = (fileName: string, mediaType: MediaType, baseUrl?: string): string => {
    return mediaService.getMediaDownloadUrl(fileName, mediaType, baseUrl);
  };

  // Check if media is an image
  const checkIsImageMedia = (media: Media): boolean => {
    return isImageMedia(media);
  };

  return {
    useGetAllMedia,
    useGetMediaMetadata,
    useUploadMedia,
    useDeleteMedia,
    useRenameMedia,
    useResizeImage,
    getMediaDownloadUrl,
    checkIsImageMedia,
  };
};

export default useUnifiedMedia;