import { TFile } from './file';
import { FileMetadata } from './fileMetadata';
import { ImageDimensions } from './imageDimensions';

/**
 * Media type discriminator to differentiate between different media types
 */
export type MediaType = 'image' | 'document' | 'video' | 'audio' | 'other';

/**
 * Base interface for all media types
 */
export interface BaseMedia {
  id: number;
  createdAt: string;
  updatedAt: string;
  deletedAt: string | null;
  fileName: string;
  checksum: string;
  mediaType: MediaType;
  downloadUrl: string;
  fileSize: number;
}

/**
 * Extended interface for image-specific properties
 */
export interface ImageMedia extends BaseMedia {
  mediaType: 'image';
  width: number;
  height: number;
  altText?: string;
  thumbnailUrl?: string;
}

/**
 * Extended interface for document-specific properties
 */
export interface DocumentMedia extends BaseMedia {
  mediaType: 'document';
  pageCount?: number;
  author?: string;
  subject?: string;
  keywords?: string[];
}

/**
 * Extended interface for video-specific properties
 */
export interface VideoMedia extends BaseMedia {
  mediaType: 'video';
  duration?: number; // in seconds
  width?: number;
  height?: number;
  thumbnailUrl?: string;
}

/**
 * Extended interface for audio-specific properties
 */
export interface AudioMedia extends BaseMedia {
  mediaType: 'audio';
  duration?: number; // in seconds
  artist?: string;
  album?: string;
  genre?: string;
}

/**
 * Extended interface for other media types
 */
export interface OtherMedia extends BaseMedia {
  mediaType: 'other';
  mimeType: string;
}

/**
 * Unified media type that can represent any media type
 */
export type Media = ImageMedia | DocumentMedia | VideoMedia | AudioMedia | OtherMedia;

/**
 * Generic media type with additional metadata
 */
export interface MediaWithMetadata<T extends Media = Media> {
  media: T;
  metadata: {
    uploadDate: string;
    lastAccessedDate?: string;
    accessCount: number;
    tags?: string[];
    description?: string;
  };
}

/**
 * Type guard to check if media is an image
 */
export function isImageMedia(media: Media): media is ImageMedia {
  return media.mediaType === 'image';
}

/**
 * Type guard to check if media is a document
 */
export function isDocumentMedia(media: Media): media is DocumentMedia {
  return media.mediaType === 'document';
}

/**
 * Type guard to check if media is a video
 */
export function isVideoMedia(media: Media): media is VideoMedia {
  return media.mediaType === 'video';
}

/**
 * Type guard to check if media is audio
 */
export function isAudioMedia(media: Media): media is AudioMedia {
  return media.mediaType === 'audio';
}

/**
 * Type guard to check if media is other type
 */
export function isOtherMedia(media: Media): media is OtherMedia {
  return media.mediaType === 'other';
}

/**
 * Conversion function to convert legacy TFile to Media
 */
export function convertFileToMedia(file: TFile, metadata: FileMetadata, type: MediaType): Media {
  const baseMedia: BaseMedia = {
    id: file.ID,
    createdAt: file.CreatedAt,
    updatedAt: file.UpdatedAt,
    deletedAt: file.DeletedAt,
    fileName: file.file_name,
    checksum: file.checksum,
    mediaType: type,
    downloadUrl: metadata.download_url,
    fileSize: metadata.file_size,
  };

  switch (type) {
    case 'image':
      return {
        ...baseMedia,
        mediaType: 'image',
        width: metadata.width || 0,
        height: metadata.height || 0,
      } as ImageMedia;
    case 'document':
      return {
        ...baseMedia,
        mediaType: 'document',
      } as DocumentMedia;
    case 'video':
      return {
        ...baseMedia,
        mediaType: 'video',
      } as VideoMedia;
    case 'audio':
      return {
        ...baseMedia,
        mediaType: 'audio',
      } as AudioMedia;
    default:
      return {
        ...baseMedia,
        mediaType: 'other',
        mimeType: 'application/octet-stream',
      } as OtherMedia;
  }
}

/**
 * Conversion function to convert Media to legacy TFile
 */
export function convertMediaToFile(media: Media): TFile {
  return {
    ID: media.id,
    CreatedAt: media.createdAt,
    UpdatedAt: media.updatedAt,
    DeletedAt: media.deletedAt,
    file_name: media.fileName,
    checksum: media.checksum,
  };
}

/**
 * Conversion function to convert Media to legacy FileMetadata
 */
export function convertMediaToFileMetadata(media: Media): FileMetadata {
  const baseMetadata: FileMetadata = {
    download_url: media.downloadUrl,
    file_size: media.fileSize,
    filename: media.fileName,
  };

  if (isImageMedia(media)) {
    return {
      ...baseMetadata,
      width: media.width,
      height: media.height,
    };
  }

  return baseMetadata;
}

/**
 * Conversion function to convert Media to legacy ImageDimensions (if applicable)
 */
export function convertMediaToImageDimensions(media: Media): ImageDimensions | null {
  if (isImageMedia(media)) {
    return {
      width: media.width,
      height: media.height,
    };
  }
  return null;
}

/**
 * Helper function to get media URL based on media type
 */
export function getMediaUrl(media: Media, baseUrl?: string): string {
  // Use the downloadUrl from the media object if available
  if (media.downloadUrl) {
    // If downloadUrl doesn't include protocol, add it
    if (media.downloadUrl.startsWith('localhost:')) {
      return `${window.location.protocol}//${media.downloadUrl}`;
    }
    return media.downloadUrl;
  }
  
  // Fallback to legacy URL structure if downloadUrl is not available
  const resolvedBaseUrl = baseUrl || `${window.location.protocol}//${window.location.host}/api/cdn/download`;
  
  switch (media.mediaType) {
    case 'image':
      return `${resolvedBaseUrl}/images/${media.fileName}`;
    case 'document':
      return `${resolvedBaseUrl}/docs/${media.fileName}`;
    case 'video':
      return `${resolvedBaseUrl}/videos/${media.fileName}`;
    case 'audio':
      return `${resolvedBaseUrl}/audio/${media.fileName}`;
    default:
      return `${resolvedBaseUrl}/other/${media.fileName}`;
  }
}

/**
 * Helper function to get media API endpoint based on media type
 */
export function getMediaApiEndpoint(media: Media): string {
  switch (media.mediaType) {
    case 'image':
      return 'image';
    case 'document':
      return 'doc';
    case 'video':
      return 'video';
    case 'audio':
      return 'audio';
    default:
      return 'other';
  }
}