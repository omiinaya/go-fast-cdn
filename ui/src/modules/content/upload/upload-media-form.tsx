import { sanitizeFileName } from "@/utils";
import { useCallback, useRef, useState } from "react";
import { cn } from "@/lib/utils";
import MediaCardUpload from "./media-card-upload.tsx";
import FileInputMedia from "./file-input-media";
import toast from "react-hot-toast";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { MediaType } from "@/types/media";

interface UploadMediaProps {
  files: File[];
  onChangeFiles: (files: File[]) => void;
  isLoading: boolean;
  mediaType: MediaType;
  onChangeMediaType: (mediaType: MediaType) => void;
  disableMediaTypeSwitching?: boolean;
  maxFileSize?: number; // in bytes
  maxFiles?: number;
}

const UploadMediaForm = ({
  isLoading,
  files,
  onChangeFiles,
  onChangeMediaType,
  mediaType,
  disableMediaTypeSwitching = false,
  maxFileSize = 50 * 1024 * 1024, // 50MB default
  maxFiles = 10, // 10 files default
}: UploadMediaProps) => {
  const fileRef = useRef<HTMLInputElement>(null);
  const [isDragOver, setIsDragOver] = useState(false);

  const handleOnChangeFiles = useCallback(
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
        
        onChangeFiles(fileArray);
      }
    },
    [onChangeFiles, files.length, maxFileSize, maxFiles]
  );

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
    
    // Check file types
    const acceptedTypes = getAcceptedTypes(mediaType);
    for (const file of newFiles) {
      if (!acceptedTypes.includes(file.type)) {
        return `File "${file.name}" is not a valid ${mediaType} file.`;
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

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setIsDragOver(false);
  }, []);

  const getAcceptedTypes = (type: MediaType): string[] => {
    switch (type) {
      case "image":
        return [
          "image/jpeg",
          "image/png",
          "image/jpg",
          "image/webp",
          "image/gif",
          "image/bmp",
          "image/svg+xml",
        ];
      case "document":
        return [
          "text/plain",
          "application/zip",
          "application/msword",
          "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
          "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
          "application/vnd.openxmlformats-officedocument.presentationml.presentation",
          "application/pdf",
          "application/rtf",
          "application/x-freearc",
        ];
      case "video":
        return [
          "video/mp4",
          "video/webm",
          "video/ogg",
          "video/quicktime",
          "video/x-msvideo",
        ];
      case "audio":
        return [
          "audio/mpeg",
          "audio/ogg",
          "audio/wav",
          "audio/webm",
          "audio/aac",
        ];
      default:
        return [];
    }
  };

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setIsDragOver(false);

      if (isLoading) return;

      const droppedFiles = Array.from(e.dataTransfer.files);
      
      // Validate files
      const validationError = validateFiles(droppedFiles, files.length, maxFileSize, maxFiles);
      if (validationError) {
        toast.error(validationError);
        return;
      }

      const acceptedTypes = getAcceptedTypes(mediaType);
      const isValidFiles = droppedFiles.every((file) =>
        acceptedTypes.includes(file.type)
      );

      if (!isValidFiles) {
        toast.error(
          `Invalid file type. Please upload ${mediaType} files only.`
        );
        return;
      }

      if (droppedFiles.length === 0) return;
      onChangeFiles([...files, ...droppedFiles]);
    },
    [isLoading, onChangeFiles, files, mediaType, maxFileSize, maxFiles]
  );

  const handleClick = useCallback(() => {
    if (fileRef.current && files.length === 0 && !isLoading) {
      fileRef.current.click();
    }
  }, [files.length, isLoading]);

  const getMediaTypeName = (type: MediaType): string => {
    switch (type) {
      case "image": return "Images";
      case "document": return "Documents";
      case "video": return "Videos";
      case "audio": return "Audio";
      default: return "Files";
    }
  };

  const getAcceptAttribute = (type: MediaType): string => {
    switch (type) {
      case "image":
        return "image/jpeg,image/png,image/jpg,image/webp,image/gif,image/bmp,image/svg+xml";
      case "document":
        return "text/plain,application/zip,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.openxmlformats-officedocument.presentationml.presentation,application/pdf,application/rtf,application/x-freearc";
      case "video":
        return "video/mp4,video/webm,video/ogg,video/quicktime,video/x-msvideo";
      case "audio":
        return "audio/mpeg,audio/ogg,audio/wav,audio/webm,audio/aac";
      default:
        return "";
    }
  };

  const getFileCountText = () => {
    if (files.length === 0) return "";
    return `${files.length} file${files.length !== 1 ? 's' : ''} selected`;
  };

  return (
    <div
      id="drop-zone"
      className={cn(
        "w-full h-96 border-2 border-dashed rounded-md transition-colors",
        {
          "cursor-pointer": files.length === 0 && !isLoading,
          "border-zinc-300": !isDragOver,
          "border-blue-400 bg-blue-50": isDragOver,
          "opacity-50": isLoading,
        }
      )}
      onClick={handleClick}
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
    >
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
                fileName={sanitizeFileName(file).name}
                mediaType={mediaType}
                onClickDelete={() => handleDeleteFile(index)}
              />
            ))}
          </div>
        </div>
      ) : (
        <div className="flex flex-col justify-center items-center h-full">
          {disableMediaTypeSwitching ? (
            // When media type switching is disabled, show content without tabs
            <>
              {isDragOver ? (
                <div className="text-center">
                  <p className="text-blue-600 font-medium">
                    Drop files here to upload
                  </p>
                </div>
              ) : (
                <div className="text-center">
                  <p className="text-gray-500 text-center mb-2">
                    Drop your {getMediaTypeName(mediaType).toLowerCase()} here, or click to select files.
                  </p>
                  <p className="text-xs text-gray-400">
                    Max file size: {Math.round(maxFileSize / (1024 * 1024))}MB | Max files: {maxFiles}
                  </p>
                  <input
                    type="file"
                    accept={getAcceptAttribute(mediaType)}
                    multiple
                    name={mediaType}
                    id={mediaType}
                    aria-label={`Select ${getMediaTypeName(mediaType)}`}
                    ref={fileRef}
                    className="hidden"
                    onChange={handleOnChangeFiles}
                  />
                </div>
              )}
            </>
          ) : (
            // When media type switching is enabled, show tabs with content
            <Tabs
              onValueChange={(value) => {
                onChangeMediaType(value as MediaType);
              }}
              value={mediaType}
            >
              <TabsList
                className="self-center mb-4"
                onClick={(e) => e.stopPropagation()}
              >
                <TabsTrigger value="image">Images</TabsTrigger>
                <TabsTrigger value="document">Documents</TabsTrigger>
                <TabsTrigger value="video">Videos</TabsTrigger>
                <TabsTrigger value="audio">Audio</TabsTrigger>
              </TabsList>

              {isDragOver ? (
                <div className="text-center">
                  <p className="text-blue-600 font-medium">
                    Drop files here to upload
                  </p>
                </div>
              ) : (
                <>
                  <FileInputMedia
                    mediaType="image"
                    fileRef={fileRef}
                    onFileChange={handleOnChangeFiles}
                    accept={getAcceptAttribute("image")}
                  />
                  <FileInputMedia
                    mediaType="document"
                    fileRef={fileRef}
                    onFileChange={handleOnChangeFiles}
                    accept={getAcceptAttribute("document")}
                  />
                  <FileInputMedia
                    mediaType="video"
                    fileRef={fileRef}
                    onFileChange={handleOnChangeFiles}
                    accept={getAcceptAttribute("video")}
                  />
                  <FileInputMedia
                    mediaType="audio"
                    fileRef={fileRef}
                    onFileChange={handleOnChangeFiles}
                    accept={getAcceptAttribute("audio")}
                  />
                </>
              )}
            </Tabs>
          )}
        </div>
      )}
    </div>
  );
};

export default UploadMediaForm;