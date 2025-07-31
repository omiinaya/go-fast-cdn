import { constant } from "@/lib/constant";
import { mediaService, MediaRenameParams } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { MediaType } from "@/types/media";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import toast from "react-hot-toast";

interface RenameMediaParams {
  fileName: string;
  newFileName: string;
  mediaType: MediaType;
}

const useRenameMediaMutation = (options?: {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
}) => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ fileName, newFileName, mediaType }: RenameMediaParams) => {
      // Use the unified media service to rename media
      const renameParams: MediaRenameParams = {
        fileName,
        newFileName,
        mediaType,
      };
      
      return mediaService.renameMedia(renameParams);
    },
    onSuccess: (_, { mediaType }) => {
      toast.dismiss();
      toast.success("Successfully renamed media!");
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.media(mediaType),
      });
      options?.onSuccess?.();
    },
    onError: (error: unknown) => {
      const err = error as AxiosError<IErrorResponse>;
      toast.dismiss();
      const message =
        err.response?.data?.error || err.message || "Rename failed";
      toast.error(message);
      options?.onError?.(new Error(message));
    },
  });
};

export default useRenameMediaMutation;