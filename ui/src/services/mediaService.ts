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
    return response.data;
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
      `/delete/${params.mediaType}/${params.fileName}`
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
    
    const response = await cdnApiClient.put<{ status: string }>(`/rename/${params.mediaType}`, form, {
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