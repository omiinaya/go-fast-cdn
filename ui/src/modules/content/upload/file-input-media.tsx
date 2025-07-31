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
        return "text/plain,text/plain; charset=utf-8,text/csv,application/json,application/xml,text/html,text/css,application/javascript,text/javascript,text/markdown,text/yaml,application/x-yaml,application/yaml,application/zip,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/vnd.ms-excel,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.ms-powerpoint,application/vnd.openxmlformats-officedocument.presentationml.presentation,application/rtf,application/vnd.oasis.opendocument.text,application/vnd.oasis.opendocument.spreadsheet,application/vnd.oasis.opendocument.presentation,application/pdf,application/x-rar-compressed,application/x-7z-compressed,application/x-tar,application/gzip,application/x-gzip,application/x-bzip2,application/x-xz,application/octet-stream";
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