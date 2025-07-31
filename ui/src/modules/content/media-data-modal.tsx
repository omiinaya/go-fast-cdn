import {
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Media, isImageMedia, isDocumentMedia, isVideoMedia, isAudioMedia } from "@/types/media";

interface TMediaDataModalProps {
  media: Media;
}

const MediaDataModal: React.FC<TMediaDataModalProps> = ({ media }) => {
  const formatFileSize = (bytes: number): string => {
    if (bytes < 1000) return `${bytes} b`;
    if (bytes < 1000000) return `${Math.round(bytes / 100) / 10} KB`;
    if (bytes < 1000000000) return `${Math.round(bytes / 100000) / 10} MB`;
    if (bytes < 1000000000000) return `${Math.round(bytes / 100000000) / 10} GB`;
    return `${Math.round(bytes / 100000000000) / 10} TB`;
  };

  const formatDate = (dateString: string): string => {
    return new Date(dateString).toLocaleString();
  };

  const getMediaSpecificInfo = () => {
    if (isImageMedia(media)) {
      return (
        <div className="">
          <strong>Dimensions</strong>
          <p>Height: {media.height}px</p>
          <p>Width: {media.width}px</p>
          {media.altText && (
            <>
              <strong>Alt Text</strong>
              <p>{media.altText}</p>
            </>
          )}
        </div>
      );
    } else if (isDocumentMedia(media)) {
      return (
        <div className="">
          {media.pageCount && (
            <>
              <strong>Page Count</strong>
              <p>{media.pageCount}</p>
            </>
          )}
          {media.author && (
            <>
              <strong>Author</strong>
              <p>{media.author}</p>
            </>
          )}
          {media.subject && (
            <>
              <strong>Subject</strong>
              <p>{media.subject}</p>
            </>
          )}
          {media.keywords && media.keywords.length > 0 && (
            <>
              <strong>Keywords</strong>
              <p>{media.keywords.join(", ")}</p>
            </>
          )}
        </div>
      );
    } else if (isVideoMedia(media)) {
      return (
        <div className="">
          {media.duration && (
            <>
              <strong>Duration</strong>
              <p>{Math.floor(media.duration / 60)}:{Math.floor(media.duration % 60).toString().padStart(2, '0')}</p>
            </>
          )}
          {media.width && media.height && (
            <>
              <strong>Dimensions</strong>
              <p>Height: {media.height}px</p>
              <p>Width: {media.width}px</p>
            </>
          )}
        </div>
      );
    } else if (isAudioMedia(media)) {
      return (
        <div className="">
          {media.duration && (
            <>
              <strong>Duration</strong>
              <p>{Math.floor(media.duration / 60)}:{Math.floor(media.duration % 60).toString().padStart(2, '0')}</p>
            </>
          )}
          {media.artist && (
            <>
              <strong>Artist</strong>
              <p>{media.artist}</p>
            </>
          )}
          {media.album && (
            <>
              <strong>Album</strong>
              <p>{media.album}</p>
            </>
          )}
          {media.genre && (
            <>
              <strong>Genre</strong>
              <p>{media.genre}</p>
            </>
          )}
        </div>
      );
    }
    return null;
  };

  return (
    <DialogContent>
      <DialogHeader>
        <DialogTitle>{media.fileName}</DialogTitle>
      </DialogHeader>

      <div className="">
        <strong>Filename</strong>
        <p id="filename">{media.fileName}</p>
      </div>
      <div className="">
        <strong>File Size</strong>
        <p id="filesize">{formatFileSize(media.fileSize)}</p>
      </div>
      <div className="">
        <strong>Media Type</strong>
        <p id="mediatype">{media.mediaType}</p>
      </div>
      <div className="">
        <strong>Created</strong>
        <p id="created">{formatDate(media.createdAt)}</p>
      </div>
      {media.updatedAt !== media.createdAt && (
        <div className="">
          <strong>Last Modified</strong>
          <p id="modified">{formatDate(media.updatedAt)}</p>
        </div>
      )}
      <div className="">
        <strong>Checksum</strong>
        <p id="checksum" className="font-mono text-xs">{media.checksum}</p>
      </div>
      {getMediaSpecificInfo()}
    </DialogContent>
  );
};

export default MediaDataModal;