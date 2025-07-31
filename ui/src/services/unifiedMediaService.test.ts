import { unifiedMediaService } from './unifiedMediaService.ts';
import { MediaType } from '@/types/media';

// Mock the cdnApiClient
jest.mock('./authService', () => ({
  cdnApiClient: {
    get: jest.fn(),
    post: jest.fn(),
    put: jest.fn(),
    delete: jest.fn(),
  },
}));

describe('UnifiedMediaService', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    
  });

  describe('getAllMedia', () => {
    it('should get all media of a specific type', async () => {
      const mockMediaData = [
        {
          id: 1,
          fileName: 'test-image.jpg',
          mediaType: 'image' as MediaType,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'abc123',
          downloadUrl: 'http://example.com/download/test-image.jpg',
          fileSize: 1024,
          width: 800,
          height: 600,
        },
      ];

      const { cdnApiClient } = require('./authService');
      (cdnApiClient.get as jest.Mock).mockResolvedValue({ data: mockMediaData });

      const result = await unifiedMediaService.getAllMedia('image');

      expect(cdnApiClient.get).toHaveBeenCalledWith('/media/all/image');
      expect(result).toEqual(mockMediaData);
    });

    it('should handle errors when getting media', async () => {
      const { cdnApiClient } = require('./authService');
      (cdnApiClient.get as jest.Mock).mockRejectedValue(new Error('Network error'));

      await expect(unifiedMediaService.getAllMedia('image')).rejects.toThrow('Network error');
    });
  });

  describe('uploadMedia', () => {
    it('should upload media of a specific type', async () => {
      const mockFile = new File(['test'], 'test.jpg', { type: 'image/jpeg' });
      const mockResponse = {
        file_url: 'http://example.com/download/test.jpg',
        type: 'image',
      };

      const { cdnApiClient } = require('./authService');
      (cdnApiClient.post as jest.Mock).mockResolvedValue({ data: mockResponse });

      const result = await unifiedMediaService.uploadMedia(mockFile, 'image');

      expect(cdnApiClient.post).toHaveBeenCalledWith(
        '/media/upload',
        expect.any(FormData),
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        }
      );
      expect(result).toEqual(mockResponse);
    });

    it('should handle errors when uploading media', async () => {
      const mockFile = new File(['test'], 'test.jpg', { type: 'image/jpeg' });

      const { cdnApiClient } = require('./authService');
      (cdnApiClient.post as jest.Mock).mockRejectedValue(new Error('Upload failed'));

      await expect(unifiedMediaService.uploadMedia(mockFile, 'image')).rejects.toThrow('Upload failed');
    });
  });

  describe('deleteMedia', () => {
    it('should delete media of a specific type', async () => {
      const mockResponse = { message: 'Media deleted successfully' };

      const { cdnApiClient } = require('./authService');
      (cdnApiClient.delete as jest.Mock).mockResolvedValue({ data: mockResponse });

      const result = await unifiedMediaService.deleteMedia('test.jpg', 'image');

      expect(cdnApiClient.delete).toHaveBeenCalledWith('/media/delete/image/test.jpg');
      expect(result).toEqual(mockResponse);
    });

    it('should handle errors when deleting media', async () => {
      const { cdnApiClient } = require('./authService');
      (cdnApiClient.delete as jest.Mock).mockRejectedValue(new Error('Delete failed'));

      await expect(unifiedMediaService.deleteMedia('test.jpg', 'image')).rejects.toThrow('Delete failed');
    });
  });

  describe('renameMedia', () => {
    it('should rename media of a specific type', async () => {
      const mockResponse = { message: 'Media renamed successfully' };

      const { cdnApiClient } = require('./authService');
      (cdnApiClient.put as jest.Mock).mockResolvedValue({ data: mockResponse });

      const result = await unifiedMediaService.renameMedia('test.jpg', 'new-test.jpg', 'image');

      expect(cdnApiClient.put).toHaveBeenCalledWith('/media/rename/image/test.jpg', {
        new_filename: 'new-test.jpg',
      });
      expect(result).toEqual(mockResponse);
    });

    it('should handle errors when renaming media', async () => {
      const { cdnApiClient } = require('./authService');
      (cdnApiClient.put as jest.Mock).mockRejectedValue(new Error('Rename failed'));

      await expect(unifiedMediaService.renameMedia('test.jpg', 'new-test.jpg', 'image')).rejects.toThrow('Rename failed');
    });
  });

  describe('resizeMedia', () => {
    it('should resize image media', async () => {
      const mockResponse = { message: 'Image resized successfully' };

      const { cdnApiClient } = require('./authService');
      (cdnApiClient.put as jest.Mock).mockResolvedValue({ data: mockResponse });

      const result = await unifiedMediaService.resizeMedia('test.jpg', 800, 600);

      expect(cdnApiClient.put).toHaveBeenCalledWith('/media/resize/image', {
        filename: 'test.jpg',
        width: 800,
        height: 600,
      });
      expect(result).toEqual(mockResponse);
    });

    it('should handle errors when resizing media', async () => {
      const { cdnApiClient } = require('./authService');
      (cdnApiClient.put as jest.Mock).mockRejectedValue(new Error('Resize failed'));

      await expect(unifiedMediaService.resizeMedia('test.jpg', 800, 600)).rejects.toThrow('Resize failed');
    });
  });

  describe('getMediaMetadata', () => {
    it('should get metadata for media of a specific type', async () => {
      const mockMetadata = {
        download_url: 'http://example.com/download/test.jpg',
        file_size: 1024,
        filename: 'test.jpg',
        width: 800,
        height: 600,
      };

      const { cdnApiClient } = require('./authService');
      (cdnApiClient.get as jest.Mock).mockResolvedValue({ data: mockMetadata });

      const result = await unifiedMediaService.getMediaMetadata('test.jpg', 'image');

      expect(cdnApiClient.get).toHaveBeenCalledWith('/media/metadata/image/test.jpg');
      expect(result).toEqual(mockMetadata);
    });

    it('should handle errors when getting media metadata', async () => {
      const { cdnApiClient } = require('./authService');
      (cdnApiClient.get as jest.Mock).mockRejectedValue(new Error('Metadata fetch failed'));

      await expect(unifiedMediaService.getMediaMetadata('test.jpg', 'image')).rejects.toThrow('Metadata fetch failed');
    });
  });
});