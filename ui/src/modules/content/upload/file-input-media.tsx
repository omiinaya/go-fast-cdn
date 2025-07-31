import { useRef } from "react";
import { Upload } from "lucide-react";
import { MediaType } from "@/types/media";

interface FileInputMediaProps {
  mediaType: MediaType;
  fileRef: React.RefObject<HTMLInputElement>;
  onFileChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  accept?: string;
}

const FileInputMedia = ({
  mediaType,
  fileRef,
  onFileChange,
  accept,
}: FileInputMediaProps) => {
  const getMediaTypeName = (type: MediaType): string => {
    switch (type) {
      case "image": return "Images";
      case "document": return "Documents";
      case "video": return "Videos";
      case "audio": return "Audio";
      default: return "Files";
    }
  };

  const getDefaultAcceptTypes = (type: MediaType): string => {
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
        return "*/*";
    }
  };

  const acceptTypes = accept || getDefaultAcceptTypes(mediaType);

  return (
    <div className="hidden">
      <input
        type="file"
        accept={acceptTypes}
        multiple
        name={mediaType}
        id={mediaType}
        aria-label={`Select ${getMediaTypeName(mediaType)}`}
        ref={fileRef}
        className="hidden"
        onChange={onFileChange}
      />
      <label htmlFor={mediaType} className="sr-only">
        Select {getMediaTypeName(mediaType)}
      </label>
    </div>
  );
};

export default FileInputMedia;