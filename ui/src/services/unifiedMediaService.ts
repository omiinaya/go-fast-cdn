import { MediaType } from '@/types/media';
import { cdnApiClient } from './authService';

export const unifiedMediaService = {
  /**
   * Get all media of a specific type
   * @param mediaType The type of media to retrieve
   * @returns Promise resolving to an array of media objects
   */
  getAllMedia: async (mediaType: MediaType) => {
    const response = await cdnApiClient.get(`/media/all/${mediaType}`);
    return response.data;
  },

  /**
   * Upload media of a specific type
   * @param file The file to upload
   * @param mediaType The type of media being uploaded
   * @returns Promise resolving to the upload response
   */
  uploadMedia: async (file: File, mediaType: MediaType) => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('type', mediaType);

    const response = await cdnApiClient.post('/media/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  /**
   * Delete media of a specific type
   * @param fileName The name of the file to delete
   * @param mediaType The type of media being deleted
   * @returns Promise resolving to the delete response
   */
  deleteMedia: async (fileName: string, mediaType: MediaType) => {
    const response = await cdnApiClient.delete(`/media/delete/${mediaType}/${fileName}`);
    return response.data;
  },

  /**
   * Rename media of a specific type
   * @param fileName The current name of the file
   * @param newFileName The new name for the file
   * @param mediaType The type of media being renamed
   * @returns Promise resolving to the rename response
   */
  renameMedia: async (fileName: string, newFileName: string, mediaType: MediaType) => {
    const response = await cdnApiClient.put(`/media/rename/${mediaType}/${fileName}`, {
      new_filename: newFileName,
    });
    return response.data;
  },

  /**
   * Resize image media
   * @param fileName The name of the image file to resize
   * @param width The new width for the image
   * @param height The new height for the image
   * @returns Promise resolving to the resize response
   */
  resizeMedia: async (fileName: string, width: number, height: number) => {
    const response = await cdnApiClient.put('/media/resize/image', {
      filename: fileName,
      width,
      height,
    });
    return response.data;
  },

  /**
   * Get metadata for media of a specific type
   * @param fileName The name of the file
   * @param mediaType The type of media
   * @returns Promise resolving to the metadata response
   */
  getMediaMetadata: async (fileName: string, mediaType: MediaType) => {
    const response = await cdnApiClient.get(`/media/metadata/${mediaType}/${fileName}`);
    return response.data;
  },
};