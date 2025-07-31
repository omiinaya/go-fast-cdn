import { X, FileText, Film, Music, File, Image as ImageIcon } from "lucide-react";
import { MediaType, Media, isImageMedia, isVideoMedia, isAudioMedia, isDocumentMedia } from "@/types/media";

interface MediaCardUploadProps {
  media: Media | File;
  onClickDelete: () => void;
  fileName?: string;
  mediaType?: MediaType;
}

const MediaCardUpload = ({
  media,
  onClickDelete,
  fileName,
  mediaType,
}: MediaCardUploadProps) => {
  const getMediaIcon = () => {
    // If it's a File object, determine type from MIME type
    if (media instanceof File) {
      const file = media as File;
      if (file.type.startsWith('image/')) {
        return (
          <div className="relative w-32 h-32">
            <img
              src={URL.createObjectURL(file)}
              alt={fileName || file.name}
              className="w-full h-full object-cover"
              onError={(e) => {
                // Fallback for broken images
                e.currentTarget.src = '';
                e.currentTarget.style.display = 'none';
              }}
            />
          </div>
        );
      } else if (file.type.startsWith('video/')) {
        return (
          <div className="relative w-32 h-32">
            <video
              src={URL.createObjectURL(file)}
              className="w-full h-full object-cover"
              muted
            />
            <div className="absolute inset-0 flex items-center justify-center bg-black/20">
              <Film size={24} className="text-white" />
            </div>
          </div>
        );
      } else if (file.type.startsWith('audio/')) {
        return (
          <div className="flex flex-col items-center justify-center w-32 h-32 bg-gray-100 rounded-md">
            <Music size="48" className="text-gray-500" />
            <span className="text-xs text-gray-500 mt-2 text-center px-1">
              Audio
            </span>
          </div>
        );
      } else {
        return (
          <div className="flex flex-col items-center justify-center w-32 h-32 bg-gray-100 rounded-md">
            <FileText size="48" className="text-gray-500" />
            <span className="text-xs text-gray-500 mt-2 text-center px-1">
              Document
            </span>
          </div>
        );
      }
    }
    // If it's a Media object, use the mediaType
    else {
      const mediaObj = media as Media;
      if (isImageMedia(mediaObj)) {
        return (
          <div className="relative w-32 h-32">
            <img
              src={mediaObj.downloadUrl || mediaObj.thumbnailUrl || ''}
              alt={fileName || mediaObj.fileName}
              className="w-full h-full object-cover"
              onError={(e) => {
                // Fallback for broken images
                e.currentTarget.src = '';
                e.currentTarget.style.display = 'none';
              }}
            />
          </div>
        );
      } else if (isVideoMedia(mediaObj)) {
        return (
          <div className="flex flex-col items-center justify-center w-32 h-32 bg-gray-100 rounded-md">
            <Film size="48" className="text-gray-500" />
            <span className="text-xs text-gray-500 mt-2 text-center px-1">
              Video
            </span>
          </div>
        );
      } else if (isAudioMedia(mediaObj)) {
        return (
          <div className="flex flex-col items-center justify-center w-32 h-32 bg-gray-100 rounded-md">
            <Music size="48" className="text-gray-500" />
            <span className="text-xs text-gray-500 mt-2 text-center px-1">
              Audio
            </span>
          </div>
        );
      } else if (isDocumentMedia(mediaObj)) {
        return (
          <div className="flex flex-col items-center justify-center w-32 h-32 bg-gray-100 rounded-md">
            <FileText size="48" className="text-gray-500" />
            <span className="text-xs text-gray-500 mt-2 text-center px-1">
              Document
            </span>
          </div>
        );
      } else {
        return (
          <div className="flex flex-col items-center justify-center w-32 h-32 bg-gray-100 rounded-md">
            <File size="48" className="text-gray-500" />
            <span className="text-xs text-gray-500 mt-2 text-center px-1">
              File
            </span>
          </div>
        );
      }
    }
  };

  const displayName = fileName || (media instanceof File ? (media as File).name : (media as Media).fileName);

  return (
    <div className="relative inline-block">
      {getMediaIcon()}
      <button
        onClick={onClickDelete}
        className="absolute top-1 right-1 bg-background/30 rounded-xs text-foreground cursor-pointer hover:bg-background/50 transition-colors"
        aria-label="Remove file"
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

export default MediaCardUpload;