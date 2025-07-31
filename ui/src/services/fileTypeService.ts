import { MediaType } from '@/types/media';
import { cdnApiClient } from './authService';

export interface FileTypeInfo {
  extension: string;
  mimeTypes: string[];
  category: MediaType;
}

export interface FileTypeConfig {
  fileTypes: Record<string, FileTypeInfo>;
}

export interface SupportedFileTypes {
  image: string[];
  document: string[];
  video: string[];
  audio: string[];
  other?: string[];
}

export interface SupportedMimeTypes {
  image: string[];
  document: string[];
  video: string[];
  audio: string[];
  other?: string[];
}

class FileTypeService {
  private config: FileTypeConfig | null = null;
  private supportedFileTypes: SupportedFileTypes | null = null;
  private supportedMimeTypes: SupportedMimeTypes | null = null;
  private loadingPromise: Promise<void> | null = null;

  /**
   * Load the file type configuration from the backend
   */
  async loadConfig(): Promise<void> {
    if (this.loadingPromise) {
      return this.loadingPromise;
    }

    this.loadingPromise = (async () => {
      try {
        const response = await cdnApiClient.get<FileTypeConfig>('/config/file-types');
        this.config = response.data;
        
        // Load supported file types and mime types
        const [fileTypesResponse, mimeTypesResponse] = await Promise.all([
          cdnApiClient.get<SupportedFileTypes>('/config/file-types/extensions'),
          cdnApiClient.get<SupportedMimeTypes>('/config/file-types/mime-types')
        ]);
        
        this.supportedFileTypes = fileTypesResponse.data;
        this.supportedMimeTypes = mimeTypesResponse.data;
      } catch (error) {
        console.error('Failed to load file type configuration:', error);
        throw error;
      }
    })();

    return this.loadingPromise;
  }

  /**
   * Get the file type configuration
   */
  getConfig(): FileTypeConfig | null {
    return this.config;
  }

  /**
   * Get supported file extensions by category
   */
  getSupportedExtensions(category: MediaType): string[] {
    if (!this.supportedFileTypes) {
      return [];
    }
    // Handle 'other' category which might not exist in the backend response
    if (category === 'other') {
      return this.supportedFileTypes.other || [];
    }
    return this.supportedFileTypes[category as keyof SupportedFileTypes] || [];
  }

  /**
   * Get supported MIME types by category
   */
  getSupportedMimeTypes(category: MediaType): string[] {
    if (!this.supportedMimeTypes) {
      return [];
    }
    // Handle 'other' category which might not exist in the backend response
    if (category === 'other') {
      return this.supportedMimeTypes.other || [];
    }
    return this.supportedMimeTypes[category as keyof SupportedMimeTypes] || [];
  }

  /**
   * Get all supported file extensions
   */
  getAllSupportedExtensions(): string[] {
    if (!this.config) {
      return [];
    }
    return Object.keys(this.config.fileTypes);
  }

  /**
   * Get all supported MIME types
   */
  getAllSupportedMimeTypes(): string[] {
    if (!this.config) {
      return [];
    }
    const mimeTypes: string[] = [];
    Object.values(this.config.fileTypes).forEach(info => {
      mimeTypes.push(...info.mimeTypes);
    });
    return mimeTypes;
  }

  /**
   * Get file info for a specific extension
   */
  getFileInfo(extension: string): FileTypeInfo | null {
    if (!this.config) {
      return null;
    }
    
    // Ensure extension starts with a dot and is lowercase
    const ext = extension.toLowerCase();
    const normalizedExt = ext.startsWith('.') ? ext : `.${ext}`;
    
    return this.config.fileTypes[normalizedExt] || null;
  }

  /**
   * Check if an extension is supported
   */
  isSupportedExtension(extension: string): boolean {
    return this.getFileInfo(extension) !== null;
  }

  /**
   * Check if a MIME type is supported
   */
  isSupportedMimeType(mimeType: string): boolean {
    if (!this.config) {
      return false;
    }
    
    return Object.values(this.config.fileTypes).some(info => 
      info.mimeTypes.includes(mimeType)
    );
  }

  /**
   * Get the category for a file extension
   */
  getCategoryFromExtension(extension: string): MediaType | null {
    const fileInfo = this.getFileInfo(extension);
    return fileInfo ? fileInfo.category : null;
  }

  /**
   * Get the category for a MIME type
   */
  getCategoryFromMimeType(mimeType: string): MediaType | null {
    if (!this.config) {
      return null;
    }
    
    for (const [ext, info] of Object.entries(this.config.fileTypes)) {
      if (info.mimeTypes.includes(mimeType)) {
        return info.category;
      }
    }
    
    return null;
  }

  /**
   * Get accept attribute value for file input
   */
  getAcceptAttribute(category: MediaType): string {
    const mimeTypes = this.getSupportedMimeTypes(category);
    return mimeTypes.join(',');
  }

  /**
   * Reset the service (useful for testing or forced refresh)
   */
  reset(): void {
    this.config = null;
    this.supportedFileTypes = null;
    this.supportedMimeTypes = null;
    this.loadingPromise = null;
  }
}

// Export singleton instance
export const fileTypeService = new FileTypeService();