import { renderHook } from '@testing-library/react';
import { waitFor } from '@testing-library/dom';
import { useUnifiedMedia } from '../use-unified-media';
import { mediaService } from '@/services/mediaService';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

// Mock the mediaService
jest.mock('@/services/mediaService', () => ({
  mediaService: {
    getAllMedia: jest.fn(),
    getMediaMetadata: jest.fn(),
    uploadMedia: jest.fn(),
    deleteMedia: jest.fn(),
    renameMedia: jest.fn(),
    resizeImage: jest.fn(),
    getMediaDownloadUrl: jest.fn((fileName: string, mediaType: string, baseUrl?: string) => {
      if (baseUrl) {
        return `${baseUrl}/${mediaType}/${fileName}`;
      }
      return `http://localhost:8080/api/cdn/download/${mediaType}/${fileName}`;
    }),
  },
}));

// Mock the toast notifications
jest.mock('react-hot-toast', () => ({
  toast: {
    success: jest.fn(),
    error: jest.fn(),
    dismiss: jest.fn(),
  },
}));

// Mock the constant
jest.mock('@/lib/constant', () => ({
  constant: {
    queryKeys: {
      media: jest.fn().mockReturnValue(['media']),
      size: jest.fn().mockReturnValue(['size']),
    },
  },
}));

describe('useUnifiedMedia', () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    jest.clearAllMocks();
    
    
    queryClient = new QueryClient({
      defaultOptions: {
        queries: {
          retry: false,
        },
      },
    });
  });

  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );

  describe('useGetAllMedia', () => {
    it('should fetch all media successfully', async () => {
      const mockMedia = [
        {
          id: 1,
          fileName: 'test.jpg',
          mediaType: 'image' as const,
          downloadUrl: 'http://example.com/test.jpg',
          fileSize: 1024,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'abc123',
          width: 800,
          height: 600,
        },
      ];

      (mediaService.getAllMedia as jest.Mock).mockResolvedValue(mockMedia);

      const { result } = renderHook(() => useUnifiedMedia(), { wrapper });
      const { useGetAllMedia } = result.current;
      const { result: mediaResult } = renderHook(() => useGetAllMedia(), { wrapper });

      await waitFor(() => {
        expect(mediaResult.current.data).toEqual(mockMedia);
        expect(mediaService.getAllMedia).toHaveBeenCalledWith({});
      });
    });

    it('should fetch media by type', async () => {
      const mockMedia = [
        {
          id: 1,
          fileName: 'test.pdf',
          mediaType: 'document' as const,
          downloadUrl: 'http://example.com/test.pdf',
          fileSize: 2048,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'def456',
        },
      ];

      (mediaService.getAllMedia as jest.Mock).mockResolvedValue(mockMedia);

      const { result } = renderHook(() => useUnifiedMedia({ mediaType: 'document' }), { wrapper });
      const { useGetAllMedia } = result.current;
      const { result: mediaResult } = renderHook(() => useGetAllMedia(), { wrapper });

      await waitFor(() => {
        expect(mediaResult.current.data).toEqual(mockMedia);
        expect(mediaService.getAllMedia).toHaveBeenCalledWith({ mediaType: 'document' });
      });
    });
  });

  describe('useUploadMedia', () => {
    it('should upload media successfully', async () => {
      const mockResponse = {
        file_url: 'http://example.com/uploaded.jpg',
        type: 'image' as const,
      };

      const mockFile = new File([''], 'test.jpg', { type: 'image/jpeg' });

      (mediaService.uploadMedia as jest.Mock).mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUnifiedMedia(), { wrapper });
      const { useUploadMedia } = result.current;
      const { result: uploadResult } = renderHook(() => useUploadMedia(), { wrapper });

      await waitFor(() => {
        uploadResult.current.mutate({ file: mockFile, mediaType: 'image' });
      });

      await waitFor(() => {
        expect(mediaService.uploadMedia).toHaveBeenCalledWith(mockFile, 'image', undefined);
      });
    });
  });

  describe('useDeleteMedia', () => {
    it('should delete media successfully', async () => {
      const mockResponse = {
        message: 'File deleted successfully',
        fileName: 'test.jpg',
      };

      (mediaService.deleteMedia as jest.Mock).mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUnifiedMedia(), { wrapper });
      const { useDeleteMedia } = result.current;
      const { result: deleteResult } = renderHook(() => useDeleteMedia(), { wrapper });

      await waitFor(() => {
        deleteResult.current.mutate({ fileName: 'test.jpg', mediaType: 'image' });
      });

      await waitFor(() => {
        expect(mediaService.deleteMedia).toHaveBeenCalledWith({
          fileName: 'test.jpg',
          mediaType: 'image',
        });
      });
    });
  });

  describe('useRenameMedia', () => {
    it('should rename media successfully', async () => {
      const mockResponse = {
        status: 'success',
      };

      (mediaService.renameMedia as jest.Mock).mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUnifiedMedia(), { wrapper });
      const { useRenameMedia } = result.current;
      const { result: renameResult } = renderHook(() => useRenameMedia(), { wrapper });

      await waitFor(() => {
        renameResult.current.mutate({
          fileName: 'old-name.jpg',
          newFileName: 'new-name.jpg',
          mediaType: 'image',
        });
      });

      await waitFor(() => {
        expect(mediaService.renameMedia).toHaveBeenCalledWith({
          fileName: 'old-name.jpg',
          newFileName: 'new-name.jpg',
          mediaType: 'image',
        });
      });
    });
  });

  describe('useResizeImage', () => {
    it('should resize image successfully', async () => {
      const mockResponse = {
        status: 'success',
        width: 300,
        height: 300,
        type: 'image' as const,
        message: 'Image resized successfully',
      };

      (mediaService.resizeImage as jest.Mock).mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUnifiedMedia(), { wrapper });
      const { useResizeImage } = result.current;
      const { result: resizeResult } = renderHook(() => useResizeImage(), { wrapper });

      await waitFor(() => {
        resizeResult.current.mutate({
          filename: 'test.jpg',
          width: 300,
          height: 300,
        });
      });

      await waitFor(() => {
        expect(mediaService.resizeImage).toHaveBeenCalledWith({
          filename: 'test.jpg',
          width: 300,
          height: 300,
        });
      });
    });
  });

  describe('getMediaDownloadUrl', () => {
    it('should return correct download URL', () => {
      const { result } = renderHook(() => useUnifiedMedia(), { wrapper });
      const { getMediaDownloadUrl } = result.current;

      const url = getMediaDownloadUrl('test.jpg', 'image', 'http://localhost:8080/api/cdn/download');

      expect(url).toBe('http://localhost:8080/api/cdn/download/image/test.jpg');
    });
  });

  describe('checkIsImageMedia', () => {
    it('should return true for image media', () => {
      const { result } = renderHook(() => useUnifiedMedia(), { wrapper });
      const { checkIsImageMedia } = result.current;

      const imageMedia = {
        id: 1,
        fileName: 'test.jpg',
        mediaType: 'image' as const,
        downloadUrl: 'http://example.com/test.jpg',
        fileSize: 1024,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
        deletedAt: null,
        checksum: 'abc123',
        width: 800,
        height: 600,
      };

      expect(checkIsImageMedia(imageMedia)).toBe(true);
    });

    it('should return false for non-image media', () => {
      const { result } = renderHook(() => useUnifiedMedia(), { wrapper });
      const { checkIsImageMedia } = result.current;

      const documentMedia = {
        id: 1,
        fileName: 'test.pdf',
        mediaType: 'document' as const,
        downloadUrl: 'http://example.com/test.pdf',
        fileSize: 1024,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
        deletedAt: null,
        checksum: 'abc123',
      };

      expect(checkIsImageMedia(documentMedia)).toBe(false);
    });
  });
});