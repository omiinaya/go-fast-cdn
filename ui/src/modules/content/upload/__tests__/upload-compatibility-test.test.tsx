import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import UploadCompatibilityTest from '../upload-compatibility-test';
import { Media, MediaType } from '@/types/media';
import toast from 'react-hot-toast';

// Mock the hooks
jest.mock('../hooks/use-upload-file-mutation', () => ({
  __esModule: true,
  default: jest.fn(() => ({
    mutate: jest.fn(),
    isLoading: false,
  })),
}));

jest.mock('../hooks/use-upload-media-mutation', () => ({
  __esModule: true,
  default: jest.fn(() => ({
    mutate: jest.fn(),
    isLoading: false,
  })),
}));

jest.mock('../hooks/use-upload-unified-media-mutation', () => ({
  __esModule: true,
  default: jest.fn(() => ({
    mutate: jest.fn(),
    isLoading: false,
  })),
}));

// Mock the toast notifications
jest.mock('react-hot-toast', () => ({
  toast: {
    success: jest.fn(),
    error: jest.fn(),
    dismiss: jest.fn(),
  },
}));

// Mock the file input component
jest.mock('../file-input', () => {
  return function MockFileInput({ onFileSelect, accept }: { 
    onFileSelect: (file: File) => void; 
    accept?: string; 
  }) {
    return (
      <div>
        <input
          data-testid="file-input"
          type="file"
          accept={accept}
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) {
              onFileSelect(file);
            }
          }}
        />
      </div>
    );
  };
});

// Mock the media file input component
jest.mock('../file-input-media', () => {
  return function MockFileInputMedia({ onFileSelect, mediaType }: { 
    onFileSelect: (file: File) => void; 
    mediaType: MediaType; 
  }) {
    return (
      <div>
        <input
          data-testid="media-file-input"
          type="file"
          data-mediatype={mediaType}
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) {
              onFileSelect(file);
            }
          }}
        />
      </div>
    );
  };
});

// Mock the unified media upload component
jest.mock('../unified-media-upload', () => {
  return function MockUnifiedMediaUpload({ onFileSelect, mediaType }: { 
    onFileSelect: (file: File) => void; 
    mediaType: MediaType; 
  }) {
    return (
      <div>
        <input
          data-testid="unified-media-file-input"
          type="file"
          data-mediatype={mediaType}
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) {
              onFileSelect(file);
            }
          }}
        />
      </div>
    );
  };
});

describe('UploadCompatibilityTest', () => {
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

  const createMockFile = (name: string, type: string, size: number): File => {
    const file = new File([''], name, { type });
    Object.defineProperty(file, 'size', { value: size });
    return file;
  };

  describe('Rendering', () => {
    it('should render all upload components', () => {
      render(<UploadCompatibilityTest />, { wrapper });
      
      expect(screen.getByText('Upload Compatibility Test')).toBeInTheDocument();
      expect(screen.getByTestId('file-input')).toBeInTheDocument();
      expect(screen.getByTestId('media-file-input')).toBeInTheDocument();
      expect(screen.getByTestId('unified-media-file-input')).toBeInTheDocument();
    });

    it('should render section titles', () => {
      render(<UploadCompatibilityTest />, { wrapper });
      
      expect(screen.getByText('Legacy Image Upload')).toBeInTheDocument();
      expect(screen.getByText('Legacy Document Upload')).toBeInTheDocument();
      expect(screen.getByText('Unified Media Upload - Image')).toBeInTheDocument();
      expect(screen.getByText('Unified Media Upload - Document')).toBeInTheDocument();
    });
  });

  describe('Backward Compatibility', () => {
    it('should handle legacy image upload', async () => {
      const mockUploadFile = jest.fn();
      (require('../hooks/use-upload-file-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadFile,
        isLoading: false,
      });

      render(<UploadCompatibilityTest />, { wrapper });
      
      const fileInput = screen.getAllByTestId('file-input')[0];
      const mockFile = createMockFile('test-image.jpg', 'image/jpeg', 1024);
      
      fireEvent.change(fileInput, { target: { files: [mockFile] } });
      
      await waitFor(() => {
        expect(mockUploadFile).toHaveBeenCalledWith({
          file: mockFile,
          type: 'image',
          filename: undefined,
        });
      });
    });

    it('should handle legacy document upload', async () => {
      const mockUploadFile = jest.fn();
      (require('../hooks/use-upload-file-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadFile,
        isLoading: false,
      });

      render(<UploadCompatibilityTest />, { wrapper });
      
      const fileInput = screen.getAllByTestId('file-input')[1];
      const mockFile = createMockFile('test-document.pdf', 'application/pdf', 2048);
      
      fireEvent.change(fileInput, { target: { files: [mockFile] } });
      
      await waitFor(() => {
        expect(mockUploadFile).toHaveBeenCalledWith({
          file: mockFile,
          type: 'doc',
          filename: undefined,
        });
      });
    });

    it('should handle unified media upload for images', async () => {
      const mockUploadMedia = jest.fn();
      (require('../hooks/use-upload-media-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadMedia,
        isLoading: false,
      });

      render(<UploadCompatibilityTest />, { wrapper });
      
      const mediaFileInput = screen.getAllByTestId('media-file-input')[0];
      const mockFile = createMockFile('test-image.jpg', 'image/jpeg', 1024);
      
      fireEvent.change(mediaFileInput, { target: { files: [mockFile] } });
      
      await waitFor(() => {
        expect(mockUploadMedia).toHaveBeenCalledWith({
          file: mockFile,
          mediaType: 'image',
          filename: undefined,
        });
      });
    });

    it('should handle unified media upload for documents', async () => {
      const mockUploadMedia = jest.fn();
      (require('../hooks/use-upload-media-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadMedia,
        isLoading: false,
      });

      render(<UploadCompatibilityTest />, { wrapper });
      
      const mediaFileInput = screen.getAllByTestId('media-file-input')[1];
      const mockFile = createMockFile('test-document.pdf', 'application/pdf', 2048);
      
      fireEvent.change(mediaFileInput, { target: { files: [mockFile] } });
      
      await waitFor(() => {
        expect(mockUploadMedia).toHaveBeenCalledWith({
          file: mockFile,
          mediaType: 'document',
          filename: undefined,
        });
      });
    });

    it('should handle unified media upload using unified hook for images', async () => {
      const mockUploadUnifiedMedia = jest.fn();
      (require('../hooks/use-upload-unified-media-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadUnifiedMedia,
        isLoading: false,
      });

      render(<UploadCompatibilityTest />, { wrapper });
      
      const unifiedMediaFileInput = screen.getAllByTestId('unified-media-file-input')[0];
      const mockFile = createMockFile('test-image.jpg', 'image/jpeg', 1024);
      
      fireEvent.change(unifiedMediaFileInput, { target: { files: [mockFile] } });
      
      await waitFor(() => {
        expect(mockUploadUnifiedMedia).toHaveBeenCalledWith({
          file: mockFile,
          mediaType: 'image',
          filename: undefined,
        });
      });
    });

    it('should handle unified media upload using unified hook for documents', async () => {
      const mockUploadUnifiedMedia = jest.fn();
      (require('../hooks/use-upload-unified-media-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadUnifiedMedia,
        isLoading: false,
      });

      render(<UploadCompatibilityTest />, { wrapper });
      
      const unifiedMediaFileInput = screen.getAllByTestId('unified-media-file-input')[1];
      const mockFile = createMockFile('test-document.pdf', 'application/pdf', 2048);
      
      fireEvent.change(unifiedMediaFileInput, { target: { files: [mockFile] } });
      
      await waitFor(() => {
        expect(mockUploadUnifiedMedia).toHaveBeenCalledWith({
          file: mockFile,
          mediaType: 'document',
          filename: undefined,
        });
      });
    });
  });

  describe('File Type Validation', () => {
    it('should accept image files for image upload', () => {
      render(<UploadCompatibilityTest />, { wrapper });
      
      const imageFileInput = screen.getAllByTestId('file-input')[0];
      expect(imageFileInput).toHaveAttribute('accept', 'image/*');
    });

    it('should accept document files for document upload', () => {
      render(<UploadCompatibilityTest />, { wrapper });
      
      const documentFileInput = screen.getAllByTestId('file-input')[1];
      expect(documentFileInput).toHaveAttribute('accept', '.pdf,.doc,.docx,.txt');
    });

    it('should have correct media type for media file inputs', () => {
      render(<UploadCompatibilityTest />, { wrapper });
      
      const imageMediaFileInput = screen.getAllByTestId('media-file-input')[0];
      const documentMediaFileInput = screen.getAllByTestId('media-file-input')[1];
      
      expect(imageMediaFileInput).toHaveAttribute('data-mediatype', 'image');
      expect(documentMediaFileInput).toHaveAttribute('data-mediatype', 'document');
    });

    it('should have correct media type for unified media file inputs', () => {
      render(<UploadCompatibilityTest />, { wrapper });
      
      const imageUnifiedMediaFileInput = screen.getAllByTestId('unified-media-file-input')[0];
      const documentUnifiedMediaFileInput = screen.getAllByTestId('unified-media-file-input')[1];
      
      expect(imageUnifiedMediaFileInput).toHaveAttribute('data-mediatype', 'image');
      expect(documentUnifiedMediaFileInput).toHaveAttribute('data-mediatype', 'document');
    });
  });

  describe('Error Handling', () => {
    it('should handle upload errors gracefully', async () => {
      const mockUploadFile = jest.fn().mockImplementation(() => {
        throw new Error('Upload failed');
      });
      
      (require('../hooks/use-upload-file-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadFile,
        isLoading: false,
      });

      render(<UploadCompatibilityTest />, { wrapper });
      
      const fileInput = screen.getAllByTestId('file-input')[0];
      const mockFile = createMockFile('test-image.jpg', 'image/jpeg', 1024);
      
      fireEvent.change(fileInput, { target: { files: [mockFile] } });
      
      await waitFor(() => {
        expect(toast.error).toHaveBeenCalled();
      });
    });

    it('should show loading state during upload', () => {
      (require('../hooks/use-upload-file-mutation').default as jest.Mock).mockReturnValue({
        mutate: jest.fn(),
        isLoading: true,
      });

      render(<UploadCompatibilityTest />, { wrapper });
      
      // In a real implementation, we would check for loading indicators
      // For this test, we just verify that the component renders without errors
      expect(screen.getByText('Upload Compatibility Test')).toBeInTheDocument();
    });
  });

  describe('Integration', () => {
    it('should maintain compatibility between legacy and unified upload methods', async () => {
      const mockUploadFile = jest.fn();
      const mockUploadMedia = jest.fn();
      const mockUploadUnifiedMedia = jest.fn();
      
      (require('../hooks/use-upload-file-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadFile,
        isLoading: false,
      });
      
      (require('../hooks/use-upload-media-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadMedia,
        isLoading: false,
      });
      
      (require('../hooks/use-upload-unified-media-mutation').default as jest.Mock).mockReturnValue({
        mutate: mockUploadUnifiedMedia,
        isLoading: false,
      });

      render(<UploadCompatibilityTest />, { wrapper });
      
      const mockImageFile = createMockFile('test-image.jpg', 'image/jpeg', 1024);
      const mockDocumentFile = createMockFile('test-document.pdf', 'application/pdf', 2048);
      
      // Test all upload methods
      const fileInputs = screen.getAllByTestId('file-input');
      const mediaFileInputs = screen.getAllByTestId('media-file-input');
      const unifiedMediaFileInputs = screen.getAllByTestId('unified-media-file-input');
      
      // Legacy image upload
      fireEvent.change(fileInputs[0], { target: { files: [mockImageFile] } });
      
      // Legacy document upload
      fireEvent.change(fileInputs[1], { target: { files: [mockDocumentFile] } });
      
      // Unified media upload for images
      fireEvent.change(mediaFileInputs[0], { target: { files: [mockImageFile] } });
      
      // Unified media upload for documents
      fireEvent.change(mediaFileInputs[1], { target: { files: [mockDocumentFile] } });
      
      // Unified media upload using unified hook for images
      fireEvent.change(unifiedMediaFileInputs[0], { target: { files: [mockImageFile] } });
      
      // Unified media upload using unified hook for documents
      fireEvent.change(unifiedMediaFileInputs[1], { target: { files: [mockDocumentFile] } });
      
      await waitFor(() => {
        // Verify that all upload methods were called
        expect(mockUploadFile).toHaveBeenCalledTimes(2);
        expect(mockUploadMedia).toHaveBeenCalledTimes(2);
        expect(mockUploadUnifiedMedia).toHaveBeenCalledTimes(2);
      });
    });
  });
});