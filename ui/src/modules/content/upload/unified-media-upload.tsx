import { useCallback, useRef, useState, useEffect } from "react";
import { cn } from "@/lib/utils";
import MediaCardUpload from "./media-card-upload";
import toast from "react-hot-toast";
import { MediaType, Media } from "@/types/media";
import { useDropzone } from "react-dropzone";
import { fileTypeService } from "@/services/fileTypeService";

interface UnifiedMediaUploadProps {
  files: File[];
  onChangeFiles: (files: File[]) => void;
  isLoading?: boolean;
  mediaType?: MediaType;
  onChangeMediaType?: (mediaType: MediaType) => void;
  disableMediaTypeSwitching?: boolean;
  maxFileSize?: number; // in bytes
  maxFiles?: number;
  className?: string;
  onUpload?: (files: File[], mediaType: MediaType) => Promise<void>;
}

const UnifiedMediaUpload = ({
  files,
  onChangeFiles,
  isLoading = false,
  mediaType = "image",
  onChangeMediaType,
  disableMediaTypeSwitching = false,
  maxFileSize = 50 * 1024 * 1024, // 50MB default
  maxFiles = 10, // 10 files default
  className,
  onUpload,
}: UnifiedMediaUploadProps) => {
  const [isDragActive, setIsDragActive] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [fileTypesLoaded, setFileTypesLoaded] = useState(false);

  // Load file type configuration on component mount
  useEffect(() => {
    const loadFileTypes = async () => {
      try {
        await fileTypeService.loadConfig();
        setFileTypesLoaded(true);
      } catch (err) {
        console.error('Failed to load file type configuration:', err);
        toast.error('Failed to load file type configuration');
      }
    };

    loadFileTypes();
  }, []);

  const getAcceptedTypes = (type: MediaType): string[] => {
    if (!fileTypesLoaded) {
      // Fallback to basic types while loading
      return ["*/*"];
    }
    
    return fileTypeService.getSupportedMimeTypes(type);
  };

  const getAcceptAttribute = (type: MediaType): string => {
    if (!fileTypesLoaded) {
      // Fallback to basic types while loading
      return "*/*";
    }
    
    return fileTypeService.getAcceptAttribute(type);
  };

  const validateFiles = (newFiles: File[], currentFileCount: number, maxSize: number, maxCount: number): string | null => {
    // Check max files limit
    if (currentFileCount + newFiles.length > maxCount) {
      return `You can only upload up to ${maxCount} files at a time.`;
    }
    
    // Check file sizes
    for (const file of newFiles) {
      if (file.size > maxSize) {
        const maxSizeMB = Math.round(maxSize / (1024 * 1024));
        return `File "${file.name}" is too large. Maximum size is ${maxSizeMB}MB.`;
      }
    }
    
    return null;
  };

  const handleDeleteFile = useCallback(
    (index: number) => {
      onChangeFiles(files.filter((_, i) => i !== index));
    },
    [files, onChangeFiles]
  );

  const handleFileInputChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const selectedFiles = e.target.files;
      if (selectedFiles) {
        const fileArray = Array.from(selectedFiles);
        
        // Validate files
        const validationError = validateFiles(fileArray, files.length, maxFileSize, maxFiles);
        if (validationError) {
          toast.error(validationError);
          return;
        }
        
        onChangeFiles([...files, ...fileArray]);
      }
    },
    [onChangeFiles, files, maxFileSize, maxFiles]
  );

  const handleUpload = useCallback(async () => {
    if (files.length === 0) {
      toast.error("Please select at least one file to upload");
      return;
    }

    if (onUpload) {
      try {
        await onUpload(files, mediaType);
      } catch (error) {
        toast.error("Upload failed");
      }
    }
  }, [files, mediaType, onUpload]);

  const onDrop = useCallback(
    (acceptedFiles: File[]) => {
      setIsDragActive(false);
      
      // Validate files
      const validationError = validateFiles(acceptedFiles, files.length, maxFileSize, maxFiles);
      if (validationError) {
        toast.error(validationError);
        return;
      }
      
      onChangeFiles([...files, ...acceptedFiles]);
    },
    [onChangeFiles, files, maxFileSize, maxFiles]
  );

  const onDropRejected = useCallback((rejectedFiles: any[]) => {
    setIsDragActive(false);
    
    if (rejectedFiles.length > 0) {
      const error = rejectedFiles[0].errors[0];
      if (error.code === "file-too-large") {
        const maxSizeMB = Math.round(maxFileSize / (1024 * 1024));
        toast.error(`File is too large. Maximum size is ${maxSizeMB}MB.`);
      } else if (error.code === "file-invalid-type") {
        toast.error(`Invalid file type. Please upload ${mediaType} files only.`);
      } else {
        toast.error("File validation failed");
      }
    }
  }, [maxFileSize, mediaType]);

  const { getRootProps, getInputProps } = useDropzone({
    onDrop,
    onDropRejected,
    accept: mediaType && fileTypesLoaded ? { [mediaType]: getAcceptedTypes(mediaType) } : undefined,
    maxSize: maxFileSize,
    maxFiles: maxFiles - files.length,
    disabled: isLoading || !fileTypesLoaded,
    noClick: files.length > 0,
  });

  const getMediaTypeName = (type: MediaType): string => {
    switch (type) {
      case "image": return "Images";
      case "document": return "Documents";
      case "video": return "Videos";
      case "audio": return "Audio";
      default: return "Files";
    }
  };

  const getFileCountText = () => {
    if (files.length === 0) return "";
    return `${files.length} file${files.length !== 1 ? 's' : ''} selected`;
  };

  return (
    <div className={cn("w-full", className)}>
      <div
        {...getRootProps()}
        className={cn(
          "w-full h-96 border-2 border-dashed rounded-md transition-colors",
          {
            "cursor-pointer": files.length === 0 && !isLoading,
            "border-zinc-300": !isDragActive,
            "border-blue-400 bg-blue-50": isDragActive,
            "opacity-50": isLoading,
          }
        )}
      >
        <input {...getInputProps()} ref={fileInputRef} onChange={handleFileInputChange} />
        
        {files.length > 0 ? (
          <div className="flex flex-col h-full">
            <div className="text-xs text-gray-500 px-2 py-1 text-right">
              {getFileCountText()}
            </div>
            <div className="flex gap-2 flex-wrap p-2 overflow-y-auto flex-grow">
              {files.map((file, index) => (
                <MediaCardUpload
                  key={file.name + index}
                  media={file}
                  fileName={file.name}
                  mediaType={mediaType}
                  onClickDelete={() => handleDeleteFile(index)}
                />
              ))}
            </div>
          </div>
        ) : (
          <div className="flex flex-col justify-center items-center h-full p-4 text-center">
            <div className="mb-4">
              <svg
                className="w-12 h-12 mx-auto text-gray-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 48 48"
                aria-hidden="true"
              >
                <path
                  d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"
                  strokeWidth={2}
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
              </svg>
            </div>
            <div>
              <p className="text-lg font-medium text-gray-700 mb-1">
                {isDragActive ? "Drop files here" : "Drag and drop files here"}
              </p>
              <p className="text-sm text-gray-500 mb-3">
                or click to browse files
              </p>
              <p className="text-xs text-gray-400">
                Max file size: {Math.round(maxFileSize / (1024 * 1024))}MB | Max files: {maxFiles}
                {!fileTypesLoaded && <span> | Loading file types...</span>}
              </p>
            </div>
          </div>
        )}
      </div>
      
      {onUpload && (
        <div className="mt-4 flex justify-end">
          <button
            type="button"
            onClick={handleUpload}
            disabled={isLoading || files.length === 0}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {isLoading ? "Uploading..." : `Upload ${files.length > 0 ? `(${files.length})` : ''}`}
          </button>
        </div>
      )}
    </div>
  );
};

export default UnifiedMediaUpload;