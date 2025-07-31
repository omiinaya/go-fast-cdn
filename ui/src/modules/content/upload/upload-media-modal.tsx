import { Loader2Icon, Plus } from "lucide-react";
import { useCallback, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { sanitizeFileName } from "@/utils";
import UploadMediaForm from "./upload-media-form.tsx";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { SidebarGroupAction } from "@/components/ui/sidebar";
import useUploadUnifiedMediaMutation from "../hooks/use-upload-unified-media-mutation";
import { AxiosError } from "axios";
import { IErrorResponse } from "@/types/response";
import toast from "react-hot-toast";
import { constant } from "@/lib/constant";
import { MediaType } from "@/types/media";

type ConditionalUploadMediaModalProps =
  | { placement: "header"; mediaType: MediaType; maxFileSize?: number; maxFiles?: number }
  | { placement?: "sidebar"; mediaType?: MediaType; maxFileSize?: number; maxFiles?: number };

const UploadMediaModal = ({
  placement = "sidebar",
  mediaType,
  maxFileSize = 50 * 1024 * 1024, // 50MB default
  maxFiles = 10, // 10 files default
}: ConditionalUploadMediaModalProps) => {
  const [open, setOpen] = useState(false);
  const [files, setFiles] = useState<File[]>([]);

  // Set initial tab based on mediaType when placement is header, otherwise default to image
  const [selectedMediaType, setSelectedMediaType] = useState<MediaType>(
    placement === "header" && mediaType ? mediaType : "image"
  );

  const uploadMediaMutation = useUploadUnifiedMediaMutation();

  const handleReset = useCallback(() => {
    setFiles([]);

    // Reset media type to initial value based on placement and mediaType
    const initialMediaType = placement === "header" && mediaType ? mediaType : "image";
    setSelectedMediaType(initialMediaType);
    setOpen(false);
    uploadMediaMutation.reset();
  }, [uploadMediaMutation, placement, mediaType]);

  const queryClient = useQueryClient();

  const { mutate: uploadMediaMutate, isPending: isUploadPending } = useMutation({
    mutationFn: async () => {
      if (files.length === 0) {
        throw new Error("No files to upload");
      }
      
      return Promise.all(
        files.map((file) => {
          const sanitizedFile = sanitizeFileName(file);
          return uploadMediaMutation.mutateAsync({
            file: sanitizedFile,
            mediaType: selectedMediaType,
          });
        })
      );
    },
    onSuccess: async () => {
      toast.success(`Successfully uploaded ${files.length} media file${files.length !== 1 ? 's' : ''}!`);
      Promise.all([
        queryClient.invalidateQueries({ queryKey: constant.queryKeys.all }),
        queryClient.invalidateQueries({
          queryKey: [constant.queryKeys.dashboard],
        }),
        queryClient.invalidateQueries({
          queryKey: constant.queryKeys.media(selectedMediaType),
        }),
      ]);
      handleReset();
    },
    onError: (error) => {
      const err = error as AxiosError<IErrorResponse>;
      const message = err.response?.data?.error || err.message || "Upload failed";
      toast.error(message);
    },
  });

  const handleUpload = useCallback(() => {
    if (files.length === 0) {
      toast.error("Please select at least one file to upload");
      return;
    }
    uploadMediaMutate();
  }, [files, uploadMediaMutate]);

  const getMediaTypeName = (type: MediaType): string => {
    switch (type) {
      case "image": return "Image";
      case "document": return "Document";
      case "video": return "Video";
      case "audio": return "Audio";
      default: return "File";
    }
  };

  const getMediaDescription = (type: MediaType): string => {
    switch (type) {
      case "image": return "Upload images (JPEG, PNG, GIF, etc.)";
      case "document": return "Upload documents (PDF, DOC, TXT, etc.)";
      case "video": return "Upload videos (MP4, WebM, etc.)";
      case "audio": return "Upload audio files (MP3, WAV, etc.)";
      default: return "Upload files";
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen} modal>
      <DialogTrigger asChild>
        {placement === "header" ? (
          <Button onClick={() => {}} variant="default">
            <Plus />
            Add {mediaType ? getMediaTypeName(mediaType) : "Media"}
          </Button>
        ) : (
          <SidebarGroupAction title="Add Content">
            <Plus /> <span className="sr-only">Add Content</span>
          </SidebarGroupAction>
        )}
      </DialogTrigger>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>Upload Media</DialogTitle>
          <DialogDescription>
            {getMediaDescription(selectedMediaType)}
          </DialogDescription>
        </DialogHeader>

        <UploadMediaForm
          isLoading={isUploadPending}
          mediaType={selectedMediaType}
          onChangeMediaType={setSelectedMediaType}
          files={files}
          onChangeFiles={setFiles}
          disableMediaTypeSwitching={placement === "header"}
          maxFileSize={maxFileSize}
          maxFiles={maxFiles}
        />

        <DialogFooter className="sm:justify-end">
          <DialogClose asChild>
            <Button
              onClick={handleReset}
              type="button"
              variant="secondary"
              disabled={isUploadPending}
            >
              Cancel
            </Button>
          </DialogClose>
          <Button
            disabled={isUploadPending || files.length === 0}
            type="button"
            variant="default"
            onClick={handleUpload}
          >
            {isUploadPending ? (
              <>
                <Loader2Icon className="animate-spin mr-2" />
                Uploading...
              </>
            ) : (
              `Upload ${files.length > 0 ? `(${files.length})` : ''}`
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default UploadMediaModal;