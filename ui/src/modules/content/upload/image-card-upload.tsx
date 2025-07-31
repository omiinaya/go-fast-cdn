import { X } from "lucide-react";
import { Media, isImageMedia } from "@/types/media";

interface ImageCardUploadProps {
  media: Media | File;
  onClickDelete: () => void;
  fileName?: string;
}

// Helper function to check if an object is a File without using instanceof
const isFileObject = (obj: any): obj is File => {
  return obj &&
         typeof obj === 'object' &&
         typeof obj.name === 'string' &&
         typeof obj.size === 'number' &&
         typeof obj.type === 'string' &&
         typeof obj.slice === 'function';
};

const ImageCardUpload = ({
  media,
  onClickDelete,
  fileName,
}: ImageCardUploadProps) => {
  const getImageUrl = (): string => {
    if (isFileObject(media)) {
      return URL.createObjectURL(media);
    } else if (isImageMedia(media)) {
      return media.downloadUrl || media.thumbnailUrl || '';
    }
    return '';
  };

  const displayName = fileName || (isFileObject(media) ? media.name : media.fileName);

  return (
    <div className="relative w-32">
      <img
        src={getImageUrl()}
        alt={displayName}
        className="w-full h-full object-cover"
        onError={(e) => {
          // Fallback for broken images
          e.currentTarget.src = '';
          e.currentTarget.style.display = 'none';
        }}
      />
      <button
        onClick={onClickDelete}
        className="absolute top-1 right-1 bg-background/30 rounded-xs text-foreground cursor-pointer hover:bg-background/50 transition-colors"
        aria-label="Remove image"
      >
        <X size={16} />
      </button>
      {displayName && (
        <div className="mt-1 text-xs text-center max-w-[128px] truncate" title={displayName}>
          {displayName}
        </div>
      )}
    </div>
  );
};

export default ImageCardUpload;
