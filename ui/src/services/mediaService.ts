import { cdnApiClient } from "@/services/authService";
import { Media, MediaType } from "@/types/media";

export interface MediaUploadResponse {
  file_url: string;
  type: MediaType;
}

export interface MediaMetadataResponse {
  filename: string;
  download_url: string;
  file_size: number;
  type: MediaType;
  width?: number;
  height?: number;
}

export interface MediaResizeParams {
  filename: string;
  width: number;
  height: number;
}

export interface MediaRenameParams {
  fileName: string;
  newFileName: string;
  mediaType: MediaType;
}

export interface MediaDeleteParams {
  fileName: string;
  mediaType: MediaType;
}

export interface GetAllMediaParams {
  mediaType?: MediaType;
}

/**
 * Unified Media Service that handles all media types (images, documents, videos, audio, other)
 */
export const mediaService = {
  /**
   * Get all media files, optionally filtered by media type
   */
  async getAllMedia(params?: GetAllMediaParams): Promise<Media[]> {
    const queryParams = params?.mediaType ? `?type=${params.mediaType}` : '';
    const response = await cdnApiClient.get<Media[]>(`/media/all${queryParams}`);
    
    // For each media item, fetch the metadata to get the download URL
    const mediaWithMetadata = await Promise.all(
      response.data.map(async (media) => {
        try {
          const metadata = await this.getMediaMetadata(media.fileName, media.mediaType);
          
          // Create a new media object with the metadata
          const enhancedMedia = {
            ...media,
            downloadUrl: metadata.download_url,
            fileSize: metadata.file_size,
          };
          
          // Add dimensions for images if available
          if (media.mediaType === 'image' && metadata.width && metadata.height) {
            return {
              ...enhancedMedia,
              width: metadata.width,
              height: metadata.height,
            } as Media;
          }
          
          return enhancedMedia as Media;
        } catch (error) {
          // If metadata fetch fails, use the media object as-is
          console.warn(`Failed to fetch metadata for ${media.fileName}:`, error);
          return media;
        }
      })
    );
    
    return mediaWithMetadata;
  },

  /**
   * Get metadata for a specific media file
   */
  async getMediaMetadata(fileName: string, mediaType: MediaType): Promise<MediaMetadataResponse> {
    const response = await cdnApiClient.get<MediaMetadataResponse>(`/media/${fileName}?type=${mediaType}`);
    return response.data;
  },

  /**
   * Upload a media file
   */
  async uploadMedia(file: File, _mediaType: MediaType, filename?: string): Promise<MediaUploadResponse> {
    const form = new FormData();
    form.append('file', file);
    if (filename) {
      form.append('filename', filename);
    }
    
    const response = await cdnApiClient.post<MediaUploadResponse>('/upload/media', form, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  /**
   * Delete a media file
   */
  async deleteMedia(params: MediaDeleteParams): Promise<{ message: string; fileName: string }> {
    const response = await cdnApiClient.delete<{ message: string; fileName: string }>(
      `/delete/media/${params.fileName}`,
      {
        params: {
          type: params.mediaType
        }
      }
    );
    return response.data;
  },

  /**
   * Rename a media file
   */
  async renameMedia(params: MediaRenameParams): Promise<{ status: string }> {
    const form = new FormData();
    form.append('filename', params.fileName);
    form.append('newname', params.newFileName);
    form.append('type', params.mediaType);
    
    const response = await cdnApiClient.put<{ status: string }>(`/rename/media`, form, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  /**
   * Resize an image (type-specific operation)
   */
  async resizeImage(params: MediaResizeParams): Promise<{ status: string; width: number; height: number; type: MediaType; message: string }> {
    const response = await cdnApiClient.put<{ status: string; width: number; height: number; type: MediaType; message: string }>(
      '/resize/media',
      params,
      {
        headers: {
          'Content-Type': 'application/json',
        },
      }
    );
    return response.data;
  },

  /**
   * Get download URL for a media file
   */
  getMediaDownloadUrl(fileName: string, mediaType: MediaType, baseUrl?: string): string {
    const resolvedBaseUrl = baseUrl || `${window.location.protocol}//${window.location.host}/api/cdn/download`;
    return `${resolvedBaseUrl}/${mediaType}/${fileName}`;
  },
};