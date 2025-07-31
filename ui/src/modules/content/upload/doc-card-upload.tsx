import { X, FileText } from "lucide-react";
import { Media, isDocumentMedia } from "@/types/media";

interface DocCardUploadProps {
  media: Media | File;
  onClickDelete: () => void;
  fileName?: string;
}

const DocCardUpload = ({ media, onClickDelete, fileName }: DocCardUploadProps) => {
  const displayName = fileName || (media instanceof File ? media.name : media.fileName);
  
  return (
    <div className="bg-zinc-200 py-1 px-2 rounded-sm text-xs truncate inline-flex items-center">
      <FileText size={14} className="mr-1 flex-shrink-0" />
      <span className="truncate">{displayName}</span>
      <button
        onClick={onClickDelete}
        className="ml-2 text-muted-foreground hover:text-foreground transition-colors"
        aria-label="Remove document"
      >
        <X size={16} />
      </button>
    </div>
  );
};

export default DocCardUpload;
