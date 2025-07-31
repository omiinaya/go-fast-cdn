import { 
  Media, 
  MediaType, 
  ImageMedia, 
  DocumentMedia, 
  VideoMedia, 
  AudioMedia, 
  OtherMedia,
  isImageMedia,
  isDocumentMedia,
  isVideoMedia,
  isAudioMedia,
  isOtherMedia,
  convertFileToMedia,
  convertMediaToFile,
  convertMediaToFileMetadata,
  convertMediaToImageDimensions,
  getMediaUrl,
  getMediaApiEndpoint
} from '../media';
import { TFile } from '../file';
import { FileMetadata } from '../fileMetadata';
import { ImageDimensions } from '../imageDimensions';


describe('Media Types', () => {
  describe('Type Guards', () => {
    it('should correctly identify image media', () => {
      const imageMedia: ImageMedia = {
        id: 1,
        fileName: 'test.jpg',
        mediaType: 'image',
        downloadUrl: 'http://example.com/test.jpg',
        fileSize: 1024,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
        deletedAt: null,
        checksum: 'abc123',
        width: 800,
        height: 600,
      };

      expect(isImageMedia(imageMedia)).toBe(true);
      expect(isDocumentMedia(imageMedia)).toBe(false);
      expect(isVideoMedia(imageMedia)).toBe(false);
      expect(isAudioMedia(imageMedia)).toBe(false);
      expect(isOtherMedia(imageMedia)).toBe(false);
    });

    it('should correctly identify document media', () => {
      const documentMedia: DocumentMedia = {
        id: 2,
        fileName: 'test.pdf',
        mediaType: 'document',
        downloadUrl: 'http://example.com/test.pdf',
        fileSize: 2048,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
        deletedAt: null,
        checksum: 'def456',
        pageCount: 10,
        author: 'Test Author',
      };

      expect(isDocumentMedia(documentMedia)).toBe(true);
      expect(isImageMedia(documentMedia)).toBe(false);
      expect(isVideoMedia(documentMedia)).toBe(false);
      expect(isAudioMedia(documentMedia)).toBe(false);
      expect(isOtherMedia(documentMedia)).toBe(false);
    });

    it('should correctly identify video media', () => {
      const videoMedia: VideoMedia = {
        id: 3,
        fileName: 'test.mp4',
        mediaType: 'video',
        downloadUrl: 'http://example.com/test.mp4',
        fileSize: 1048576,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
        deletedAt: null,
        checksum: 'ghi789',
        duration: 120,
        width: 1920,
        height: 1080,
      };

      expect(isVideoMedia(videoMedia)).toBe(true);
      expect(isImageMedia(videoMedia)).toBe(false);
      expect(isDocumentMedia(videoMedia)).toBe(false);
      expect(isAudioMedia(videoMedia)).toBe(false);
      expect(isOtherMedia(videoMedia)).toBe(false);
    });

    it('should correctly identify audio media', () => {
      const audioMedia: AudioMedia = {
        id: 4,
        fileName: 'test.mp3',
        mediaType: 'audio',
        downloadUrl: 'http://example.com/test.mp3',
        fileSize: 524288,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
        deletedAt: null,
        checksum: 'jkl012',
        duration: 180,
        artist: 'Test Artist',
        album: 'Test Album',
        genre: 'Test Genre',
      };

      expect(isAudioMedia(audioMedia)).toBe(true);
      expect(isImageMedia(audioMedia)).toBe(false);
      expect(isDocumentMedia(audioMedia)).toBe(false);
      expect(isVideoMedia(audioMedia)).toBe(false);
      expect(isOtherMedia(audioMedia)).toBe(false);
    });

    it('should correctly identify other media', () => {
      const otherMedia: OtherMedia = {
        id: 5,
        fileName: 'test.bin',
        mediaType: 'other',
        downloadUrl: 'http://example.com/test.bin',
        fileSize: 4096,
        createdAt: '2023-01-01T00:00:00Z',
        updatedAt: '2023-01-01T00:00:00Z',
        deletedAt: null,
        checksum: 'mno345',
        mimeType: 'application/octet-stream',
      };

      expect(isOtherMedia(otherMedia)).toBe(true);
      expect(isImageMedia(otherMedia)).toBe(false);
      expect(isDocumentMedia(otherMedia)).toBe(false);
      expect(isVideoMedia(otherMedia)).toBe(false);
      expect(isAudioMedia(otherMedia)).toBe(false);
    });
  });

  describe('Conversion Functions', () => {
    const mockFile: TFile = {
      ID: 1,
      CreatedAt: '2023-01-01T00:00:00Z',
      UpdatedAt: '2023-01-01T00:00:00Z',
      DeletedAt: null,
      file_name: 'test.jpg',
      checksum: 'abc123',
    };

    const mockMetadata: FileMetadata = {
      download_url: 'http://example.com/test.jpg',
      file_size: 1024,
      filename: 'test.jpg',
      width: 800,
      height: 600,
    };

    describe('convertFileToMedia', () => {
      it('should convert to image media', () => {
        const result = convertFileToMedia(mockFile, mockMetadata, 'image');
        
        expect(result.mediaType).toBe('image');
        expect(result.id).toBe(mockFile.ID);
        expect(result.fileName).toBe(mockFile.file_name);
        expect(result.checksum).toBe(mockFile.checksum);
        expect(result.downloadUrl).toBe(mockMetadata.download_url);
        expect(result.fileSize).toBe(mockMetadata.file_size);
        
        if (isImageMedia(result)) {
          expect(result.width).toBe(mockMetadata.width);
          expect(result.height).toBe(mockMetadata.height);
        }
      });

      it('should convert to document media', () => {
        const result = convertFileToMedia(mockFile, mockMetadata, 'document');
        
        expect(result.mediaType).toBe('document');
        expect(isDocumentMedia(result)).toBe(true);
      });

      it('should convert to video media', () => {
        const result = convertFileToMedia(mockFile, mockMetadata, 'video');
        
        expect(result.mediaType).toBe('video');
        expect(isVideoMedia(result)).toBe(true);
      });

      it('should convert to audio media', () => {
        const result = convertFileToMedia(mockFile, mockMetadata, 'audio');
        
        expect(result.mediaType).toBe('audio');
        expect(isAudioMedia(result)).toBe(true);
      });

      it('should convert to other media', () => {
        const result = convertFileToMedia(mockFile, mockMetadata, 'other');
        
        expect(result.mediaType).toBe('other');
        expect(isOtherMedia(result)).toBe(true);
        
        if (isOtherMedia(result)) {
          expect(result.mimeType).toBe('application/octet-stream');
        }
      });
    });

    describe('convertMediaToFile', () => {
      it('should convert media back to TFile', () => {
        const media: ImageMedia = {
          id: 1,
          fileName: 'test.jpg',
          mediaType: 'image',
          downloadUrl: 'http://example.com/test.jpg',
          fileSize: 1024,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'abc123',
          width: 800,
          height: 600,
        };

        const result = convertMediaToFile(media);
        
        expect(result.ID).toBe(media.id);
        expect(result.file_name).toBe(media.fileName);
        expect(result.checksum).toBe(media.checksum);
        expect(result.CreatedAt).toBe(media.createdAt);
        expect(result.UpdatedAt).toBe(media.updatedAt);
        expect(result.DeletedAt).toBe(media.deletedAt);
      });
    });

    describe('convertMediaToFileMetadata', () => {
      it('should convert image media to FileMetadata with dimensions', () => {
        const media: ImageMedia = {
          id: 1,
          fileName: 'test.jpg',
          mediaType: 'image',
          downloadUrl: 'http://example.com/test.jpg',
          fileSize: 1024,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'abc123',
          width: 800,
          height: 600,
        };

        const result = convertMediaToFileMetadata(media);
        
        expect(result.download_url).toBe(media.downloadUrl);
        expect(result.file_size).toBe(media.fileSize);
        expect(result.filename).toBe(media.fileName);
        expect(result.width).toBe(media.width);
        expect(result.height).toBe(media.height);
      });

      it('should convert document media to FileMetadata without dimensions', () => {
        const media: DocumentMedia = {
          id: 1,
          fileName: 'test.pdf',
          mediaType: 'document',
          downloadUrl: 'http://example.com/test.pdf',
          fileSize: 2048,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'def456',
        };

        const result = convertMediaToFileMetadata(media);
        
        expect(result.download_url).toBe(media.downloadUrl);
        expect(result.file_size).toBe(media.fileSize);
        expect(result.filename).toBe(media.fileName);
        expect(result.width).toBeUndefined();
        expect(result.height).toBeUndefined();
      });
    });

    describe('convertMediaToImageDimensions', () => {
      it('should convert image media to ImageDimensions', () => {
        const media: ImageMedia = {
          id: 1,
          fileName: 'test.jpg',
          mediaType: 'image',
          downloadUrl: 'http://example.com/test.jpg',
          fileSize: 1024,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'abc123',
          width: 800,
          height: 600,
        };

        const result = convertMediaToImageDimensions(media);
        
        expect(result).toEqual({
          width: media.width,
          height: media.height,
        });
      });

      it('should return null for non-image media', () => {
        const media: DocumentMedia = {
          id: 1,
          fileName: 'test.pdf',
          mediaType: 'document',
          downloadUrl: 'http://example.com/test.pdf',
          fileSize: 2048,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'def456',
        };

        const result = convertMediaToImageDimensions(media);
        
        expect(result).toBeNull();
      });
    });
  });

  describe('Helper Functions', () => {
    describe('getMediaUrl', () => {
      const testBaseUrl = 'http://localhost:8080/api/cdn/download';
      
      it('should return correct URL for image media', () => {
        const media: ImageMedia = {
          id: 1,
          fileName: 'test.jpg',
          mediaType: 'image',
          downloadUrl: 'http://example.com/test.jpg',
          fileSize: 1024,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'abc123',
          width: 800,
          height: 600,
        };

        const result = getMediaUrl(media, testBaseUrl);
        
        expect(result).toBe('http://localhost:8080/api/cdn/download/images/test.jpg');
      });

      it('should return correct URL for document media', () => {
        const media: DocumentMedia = {
          id: 1,
          fileName: 'test.pdf',
          mediaType: 'document',
          downloadUrl: 'http://example.com/test.pdf',
          fileSize: 2048,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'def456',
        };

        const result = getMediaUrl(media, testBaseUrl);
        
        expect(result).toBe('http://localhost:8080/api/cdn/download/docs/test.pdf');
      });

      it('should return correct URL for video media', () => {
        const media: VideoMedia = {
          id: 1,
          fileName: 'test.mp4',
          mediaType: 'video',
          downloadUrl: 'http://example.com/test.mp4',
          fileSize: 1048576,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'ghi789',
        };

        const result = getMediaUrl(media, testBaseUrl);
        
        expect(result).toBe('http://localhost:8080/api/cdn/download/videos/test.mp4');
      });

      it('should return correct URL for audio media', () => {
        const media: AudioMedia = {
          id: 1,
          fileName: 'test.mp3',
          mediaType: 'audio',
          downloadUrl: 'http://example.com/test.mp3',
          fileSize: 524288,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'jkl012',
        };

        const result = getMediaUrl(media, testBaseUrl);
        
        expect(result).toBe('http://localhost:8080/api/cdn/download/audio/test.mp3');
      });

      it('should return correct URL for other media', () => {
        const media: OtherMedia = {
          id: 1,
          fileName: 'test.bin',
          mediaType: 'other',
          downloadUrl: 'http://example.com/test.bin',
          fileSize: 4096,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'mno345',
          mimeType: 'application/octet-stream',
        };

        const result = getMediaUrl(media, testBaseUrl);
        
        expect(result).toBe('http://localhost:8080/api/cdn/download/other/test.bin');
      });
    });

    describe('getMediaApiEndpoint', () => {
      it('should return correct endpoint for image media', () => {
        const media: ImageMedia = {
          id: 1,
          fileName: 'test.jpg',
          mediaType: 'image',
          downloadUrl: 'http://example.com/test.jpg',
          fileSize: 1024,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'abc123',
          width: 800,
          height: 600,
        };

        const result = getMediaApiEndpoint(media);
        
        expect(result).toBe('image');
      });

      it('should return correct endpoint for document media', () => {
        const media: DocumentMedia = {
          id: 1,
          fileName: 'test.pdf',
          mediaType: 'document',
          downloadUrl: 'http://example.com/test.pdf',
          fileSize: 2048,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'def456',
        };

        const result = getMediaApiEndpoint(media);
        
        expect(result).toBe('doc');
      });

      it('should return correct endpoint for video media', () => {
        const media: VideoMedia = {
          id: 1,
          fileName: 'test.mp4',
          mediaType: 'video',
          downloadUrl: 'http://example.com/test.mp4',
          fileSize: 1048576,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'ghi789',
        };

        const result = getMediaApiEndpoint(media);
        
        expect(result).toBe('video');
      });

      it('should return correct endpoint for audio media', () => {
        const media: AudioMedia = {
          id: 1,
          fileName: 'test.mp3',
          mediaType: 'audio',
          downloadUrl: 'http://example.com/test.mp3',
          fileSize: 524288,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'jkl012',
        };

        const result = getMediaApiEndpoint(media);
        
        expect(result).toBe('audio');
      });

      it('should return correct endpoint for other media', () => {
        const media: OtherMedia = {
          id: 1,
          fileName: 'test.bin',
          mediaType: 'other',
          downloadUrl: 'http://example.com/test.bin',
          fileSize: 4096,
          createdAt: '2023-01-01T00:00:00Z',
          updatedAt: '2023-01-01T00:00:00Z',
          deletedAt: null,
          checksum: 'mno345',
          mimeType: 'application/octet-stream',
        };

        const result = getMediaApiEndpoint(media);
        
        expect(result).toBe('other');
      });
    });
  });
});