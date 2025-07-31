import { Loader2Icon, Plus } from "lucide-react";
import { useCallback, useState } from "react";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { sanitizeFileName } from "@/utils";
import UploadForm from "./upload-form";
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
import useUploadFileMutation from "../hooks/use-upload-file-mutation";
import useUploadUnifiedMediaMutation from "../hooks/use-upload-unified-media-mutation";
import { AxiosError } from "axios";
import { IErrorResponse } from "@/types/response";
import toast from "react-hot-toast";
import { constant } from "@/lib/constant";
import { MediaType } from "@/types/media";

type ConditionalUploadModalProps =
  | { placement: "header"; type: MediaType; useUnifiedMedia?: boolean }
  | { placement?: "sidebar"; type?: MediaType; useUnifiedMedia?: boolean };

const UploadModal = ({
  placement = "sidebar",
  type,
  useUnifiedMedia = false,
}: ConditionalUploadModalProps) => {
  const [open, setOpen] = useState(false);
  const [files, setFiles] = useState<File[]>([]);

  // Set initial tab based on type when placement is header, otherwise default to documents
  const [tab, setTab] = useState<MediaType>(
    placement === "header" && type ? type : "document"
  );

  const uploadFileMutation = useUploadFileMutation();
  const uploadUnifiedMediaMutation = useUploadUnifiedMediaMutation();

  const handleReset = useCallback(() => {
    setFiles([]);

    // Reset tab to initial value based on placement and type
    const initialTab = placement === "header" && type ? type : "document";
    setTab(initialTab);
    setOpen(false);
    uploadFileMutation.reset();
    uploadUnifiedMediaMutation.reset();
  }, [uploadFileMutation, uploadUnifiedMediaMutation, placement, type]);

  const queryClient = useQueryClient();

  const { mutate: uploadFileMutate, isPending: isUploadPending } = useMutation({
    mutationFn: async () => {
      if (useUnifiedMedia) {
        // Use unified media upload
        return Promise.all(
          files.map((file) => {
            const sanitizedFile = sanitizeFileName(file);
            return uploadUnifiedMediaMutation.mutateAsync({
              file: sanitizedFile,
              mediaType: tab,
            });
          })
        );
      } else {
        // Use legacy file upload
        return Promise.all(
          files.map((file) => {
            const sanitizedFile = sanitizeFileName(file);
            return uploadFileMutation.mutateAsync({
              file: sanitizedFile,
              type: tab === "document" ? "doc" : "image",
            });
          })
        );
      }
    },
    onSuccess: async () => {
      toast.success(`Successfully uploaded ${useUnifiedMedia ? "media" : "file"}!`);
      Promise.all([
        queryClient.invalidateQueries({ queryKey: constant.queryKeys.all }),
        queryClient.invalidateQueries({
          queryKey: [constant.queryKeys.dashboard],
        }),
      ]);
      
      // Also invalidate media queries if using unified media
      if (useUnifiedMedia) {
        queryClient.invalidateQueries({
          queryKey: constant.queryKeys.media(tab),
        });
      }
      
      handleReset();
    },
    onError: (error) => {
      const err = error as AxiosError<IErrorResponse>;
      const message = err.response?.data?.error || "Upload failed";
      toast.error(message);
    },
  });

  const handleUpload = useCallback(() => {
    if (files.length === 0) {
      return;
    }
    uploadFileMutate();
  }, [files, uploadFileMutate]);

  const getMediaTypeName = (type: MediaType): string => {
    switch (type) {
      case "image": return "Image";
      case "document": return "Document";
      case "video": return "Video";
      case "audio": return "Audio";
      default: return "File";
    }
  };

  return (
    <Dialog open={open} onOpenChange={setOpen} modal>
      <DialogTrigger asChild>
        {placement === "header" ? (
          <Button onClick={() => {}} variant="default">
            <Plus />
            Add {type ? getMediaTypeName(type) : "File"}
          </Button>
        ) : (
          <SidebarGroupAction title="Add Content">
            <Plus /> <span className="sr-only">Add Content</span>
          </SidebarGroupAction>
        )}
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Upload {useUnifiedMedia ? "Media" : "Files"}</DialogTitle>
          <DialogDescription>
            Upload your {useUnifiedMedia ? "media files" : "files"} here and manage your content easily.
          </DialogDescription>
        </DialogHeader>

        <UploadForm
          isLoading={isUploadPending}
          tab={tab}
          onChangeTab={setTab}
          files={files}
          onChangeFiles={setFiles}
          disableTabSwitching={placement === "header"}
          useUnifiedMedia={useUnifiedMedia}
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
            disabled={isUploadPending}
            type="button"
            variant="default"
            onClick={handleUpload}
          >
            {isUploadPending ? (
              <>
                <Loader2Icon className="animate-spin" />
                Please wait
              </>
            ) : (
              "Upload"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default UploadModal;
