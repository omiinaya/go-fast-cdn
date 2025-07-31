import { sanitizeFileName } from "@/utils";
import { useCallback, useRef, useState } from "react";
import { cn } from "@/lib/utils";
import ImageCardUpload from "./image-card-upload";
import FileInput from "./file-input";
import toast from "react-hot-toast";
import DocCardUpload from "./doc-card-upload";
import MediaCardUpload from "./media-card-upload";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { MediaType } from "@/types/media";

interface UploadProps {
  files: File[];
  onChangeFiles: (files: File[]) => void;
  isLoading: boolean;
  tab: MediaType;
  onChangeTab: (tab: MediaType) => void;
  disableTabSwitching?: boolean;
  useUnifiedMedia?: boolean; // Flag to enable unified media handling
}

const UploadForm = ({
  isLoading,
  files,
  onChangeFiles,
  onChangeTab,
  tab,
  disableTabSwitching = false,
  useUnifiedMedia = false,
}: UploadProps) => {
  const fileRef = useRef<HTMLInputElement>(null);
  const [isDragOver, setIsDragOver] = useState(false);

  const handleOnChangeFiles = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const selectedFiles = e.target.files;
      if (selectedFiles) {
        const fileArray = Array.from(selectedFiles);
        onChangeFiles(fileArray);
      }
    },
    [onChangeFiles]
  );

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
          "text/plain; charset=utf-8",
          "text/csv",
          "application/json",
          "application/xml",
          "text/xml",
          "text/html",
          "text/html; charset=utf-8",
          "text/css",
          "text/css; charset=utf-8",
          "application/javascript",
          "text/javascript",
          "application/javascript; charset=utf-8",
          "text/javascript; charset=utf-8",
          "text/markdown",
          "text/yaml",
          "text/yaml; charset=utf-8",
          "application/x-yaml",
          "application/yaml",
          "application/zip",
          "application/msword",
          "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
          "application/vnd.ms-excel",
          "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
          "application/vnd.ms-powerpoint",
          "application/vnd.openxmlformats-officedocument.presentationml.presentation",
          "application/rtf",
          "application/vnd.oasis.opendocument.text",
          "application/vnd.oasis.opendocument.spreadsheet",
          "application/vnd.oasis.opendocument.presentation",
          "application/pdf",
          "application/x-rar-compressed",
          "application/x-7z-compressed",
          "application/x-tar",
          "application/gzip",
          "application/x-gzip",
          "application/x-bzip2",
          "application/x-xz",
          "application/octet-stream",
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
      const acceptedTypes = getAcceptedTypes(tab);
      const isValidFiles = droppedFiles.every((file) =>
        acceptedTypes.includes(file.type)
      );
      
      if (!isValidFiles) {
        toast.error(
          `Invalid file type. Please upload ${tab} files only.`
        );
        return;
      }
      
      if (droppedFiles.length === 0) return;
      onChangeFiles([...files, ...droppedFiles]);
    },
    [isLoading, onChangeFiles, files, tab]
  );

  const handleClick = useCallback(() => {
    if (fileRef.current && files.length === 0 && !isLoading) {
      fileRef.current.click();
    }
  }, [files.length, isLoading]);

  const getAcceptAttribute = (type: MediaType): string => {
    switch (type) {
      case "image":
        return "image/jpeg,image/png,image/jpg,image/webp,image/gif,image/bmp,image/svg+xml";
      case "document":
        return "text/plain,text/plain; charset=utf-8,text/csv,application/json,application/xml,text/xml,text/html,text/html; charset=utf-8,text/css,text/css; charset=utf-8,application/javascript,text/javascript,application/javascript; charset=utf-8,text/javascript; charset=utf-8,text/markdown,text/yaml,text/yaml; charset=utf-8,application/x-yaml,application/yaml,application/zip,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.ms-powerpoint,application/vnd.openxmlformats-officedocument.presentationml.presentation,application/rtf,application/vnd.oasis.opendocument.text,application/vnd.oasis.opendocument.spreadsheet,application/vnd.oasis.opendocument.presentation,application/pdf,application/x-rar-compressed,application/x-7z-compressed,application/x-tar,application/gzip,application/x-gzip,application/x-bzip2,application/x-xz,application/octet-stream";
      case "video":
        return "video/mp4,video/webm,video/ogg,video/quicktime,video/x-msvideo";
      case "audio":
        return "audio/mpeg,audio/ogg,audio/wav,audio/webm,audio/aac";
      default:
        return "*/*";
    }
  };

  const getMediaTypeName = (type: MediaType): string => {
    switch (type) {
      case "image": return "images";
      case "document": return "documents";
      case "video": return "videos";
      case "audio": return "audio files";
      default: return "files";
    }
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
        <div className="flex gap-2 flex-wrap p-2 overflow-y-auto max-h-80">
          {useUnifiedMedia ? (
            // Use unified media card for all media types
            files.map((file, index) => (
              <MediaCardUpload
                key={file.name + index}
                media={file}
                fileName={sanitizeFileName(file).name}
                mediaType={tab}
                onClickDelete={() => handleDeleteFile(index)}
              />
            ))
          ) : (
            // Use legacy components for backward compatibility
            <>
              {tab === "document" ? (
                files.map((file, index) => (
                  <DocCardUpload
                    key={file.name + index}
                    media={file}
                    fileName={sanitizeFileName(file).name}
                    onClickDelete={() => handleDeleteFile(index)}
                  />
                ))
              ) : (
                files.map((file, index) => (
                  <ImageCardUpload
                    key={file.name + index}
                    media={file}
                    fileName={sanitizeFileName(file).name}
                    onClickDelete={() => handleDeleteFile(index)}
                  />
                ))
              )}
            </>
          )}
        </div>
      ) : (
        <div className="flex flex-col justify-center items-center h-full">
          {disableTabSwitching ? (
            // When tab switching is disabled, show content without tabs
            <>
              {isDragOver ? (
                <div className="text-center">
                  <p className="text-blue-600 font-medium">
                    Drop files here to upload
                  </p>
                </div>
              ) : (
                <div className="text-center">
                  <p className="text-gray-500 text-center">
                    Drop your {getMediaTypeName(tab)} here, or click to select files.
                  </p>
                  <input
                    type="file"
                    accept={getAcceptAttribute(tab)}
                    multiple
                    name={tab}
                    id={tab}
                    aria-label={`Select ${getMediaTypeName(tab)}`}
                    ref={fileRef}
                    className="hidden"
                    onChange={handleOnChangeFiles}
                  />
                </div>
              )}
            </>
          ) : (
            // When tab switching is enabled, show tabs with content
            <Tabs
              onValueChange={(value) => {
                onChangeTab(value as MediaType);
              }}
              value={tab}
            >
              <TabsList
                className="self-center mb-4"
                onClick={(e) => e.stopPropagation()}
              >
                <TabsTrigger value="document">Documents</TabsTrigger>
                <TabsTrigger value="image">Images</TabsTrigger>
                {useUnifiedMedia && (
                  <>
                    <TabsTrigger value="video">Videos</TabsTrigger>
                    <TabsTrigger value="audio">Audio</TabsTrigger>
                  </>
                )}
              </TabsList>

              {isDragOver ? (
                <div className="text-center">
                  <p className="text-blue-600 font-medium">
                    Drop files here to upload
                  </p>
                </div>
              ) : (
                <>
                  <FileInput
                    type="document"
                    fileRef={fileRef}
                    onFileChange={handleOnChangeFiles}
                  />
                  <FileInput
                    type="image"
                    fileRef={fileRef}
                    onFileChange={handleOnChangeFiles}
                  />
                  {useUnifiedMedia && (
                    <>
                      <FileInput
                        type="video"
                        fileRef={fileRef}
                        onFileChange={handleOnChangeFiles}
                      />
                      <FileInput
                        type="audio"
                        fileRef={fileRef}
                        onFileChange={handleOnChangeFiles}
                      />
                    </>
                  )}
                </>
              )}
            </Tabs>
          )}
        </div>
      )}
    </div>
  );
};

export default UploadForm;
