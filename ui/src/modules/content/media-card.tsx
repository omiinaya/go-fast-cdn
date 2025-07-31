import { Media, isImageMedia, isDocumentMedia, isVideoMedia, isAudioMedia, getMediaUrl } from "@/types/media";
import { DownloadCloud, FileText, Files, Trash2, Film, Music, File } from "lucide-react";
import { toast } from "react-hot-toast";
import MediaDataModal from "./media-data-modal";
import RenameModal from "./rename-modal";
import ResizeModal from "./resize-modal";
import useDeleteUnifiedMediaMutation from "./hooks/use-delete-unified-media-mutation";
import { Tooltip, TooltipContent } from "@/components/ui/tooltip";
import { TooltipTrigger } from "@radix-ui/react-tooltip";
import { Dialog, DialogTrigger } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";

interface TMediaCardProps {
  media: Media;
  disabled?: boolean;
  isSelected?: boolean;
  onSelect?: (fileName: string) => void;
  isSelecting?: boolean;
}

const MediaCard: React.FC<TMediaCardProps> = ({
  media,
  disabled = false,
  isSelected,
  onSelect,
  isSelecting,
}) => {
  const url = getMediaUrl(media);

  const deleteMedia = useDeleteUnifiedMediaMutation();

  const handleDeleteMedia = () => {
    deleteMedia.mutate({ fileName: media.fileName, mediaType: media.mediaType });
  };

  const getMediaIcon = () => {
    if (isImageMedia(media)) {
      return (
        <img
          src={url}
          alt={media.fileName}
          width={224}
          height={150}
          className="object-cover max-h-[150px] max-w-[224px]"
        />
      );
    } else if (isDocumentMedia(media)) {
      return <FileText size="128" />;
    } else if (isVideoMedia(media)) {
      return <Film size="128" />;
    } else if (isAudioMedia(media)) {
      return <Music size="128" />;
    } else {
      return <File size="128" />;
    }
  };

  const getMediaTypeName = () => {
    if (isImageMedia(media)) return "Image";
    if (isDocumentMedia(media)) return "Document";
    if (isVideoMedia(media)) return "Video";
    if (isAudioMedia(media)) return "Audio";
    return "File";
  };

  return (
    <div className="border rounded-lg shadow-lg flex flex-col min-h-[264px] w-64 max-w-[256px] justify-between items-center gap-4 p-4 relative">
      {isSelecting && (
        <Checkbox
          className="absolute top-2 right-2 bg-background"
          checked={isSelected}
          onCheckedChange={() => onSelect && onSelect(media.fileName)}
          disabled={disabled}
          aria-label="Select file"
        />
      )}
      <Dialog>
        <DialogTrigger disabled={isSelecting}>
          {getMediaIcon()}
        </DialogTrigger>
        <MediaDataModal media={media} />
      </Dialog>
      <div className="w-full flex flex-col gap-2">
        <p className="truncate" title={media.fileName}>{media.fileName}</p>
        <p className="text-xs text-muted-foreground">{getMediaTypeName()}</p>
        {/* Non-destructive buttons */}
        <div className={`flex w-full justify-between ${disabled && "sr-only"}`}>
          <div className="flex">
            <Tooltip>
              <TooltipTrigger>
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-sky-600"
                  onClick={() => {
                    navigator.clipboard.writeText(url);
                    toast.success("Link copied to clipboard");
                  }}
                  aria-label="Copy Link"
                  disabled={isSelecting}
                >
                  <Files />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="bottom">
                <p>Copy Link to clipboard</p>
              </TooltipContent>
            </Tooltip>
            <Tooltip>
              <TooltipTrigger>
                <Button
                  variant="ghost"
                  size="icon"
                  className="text-sky-600"
                  disabled={isSelecting}
                  asChild={!isSelecting}
                >
                  <a href={url} download aria-label="Download file">
                    <DownloadCloud />
                  </a>
                </Button>
              </TooltipTrigger>
              <TooltipContent side="bottom">
                <p>Download file</p>
              </TooltipContent>
            </Tooltip>
            <RenameModal
              filename={media.fileName}
              type={media.mediaType === "image" ? "images" : "documents"}
              isSelecting={isSelecting}
            />
            {isImageMedia(media) && (
              <ResizeModal
                filename={media.fileName}
                isSelecting={isSelecting}
              />
            )}
          </div>
          {/* Destructive buttons */}
          <div className="flex gap-2">
            <Tooltip>
              <TooltipTrigger>
                <Button
                  variant="destructive"
                  size="icon"
                  onClick={() => handleDeleteMedia()}
                  aria-label="Delete file"
                  disabled={isSelecting}
                >
                  <Trash2 className="inline" size="24" />
                </Button>
              </TooltipTrigger>
              <TooltipContent side="bottom">
                <p>Delete file</p>
              </TooltipContent>
            </Tooltip>
          </div>
        </div>
      </div>
    </div>
  );
};

export default MediaCard;