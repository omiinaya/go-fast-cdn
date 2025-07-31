import { constant } from "@/lib/constant";
import { mediaService, MediaDeleteParams } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { MediaType } from "@/types/media";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import toast from "react-hot-toast";

interface DeleteMediaParams {
  fileName: string;
  mediaType: MediaType;
}

const useDeleteMediaMutation = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ fileName, mediaType }: DeleteMediaParams) => {
      // Use the unified media service to delete media
      const deleteParams: MediaDeleteParams = {
        fileName,
        mediaType,
      };
      
      return mediaService.deleteMedia(deleteParams);
    },
    onSuccess: (_, { mediaType }) => {
      toast.dismiss();
      toast.success("Successfully deleted media!");
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.size(),
      });
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.media(mediaType),
      });
    },
    onError: (error: unknown) => {
      const err = error as AxiosError<IErrorResponse>;
      toast.dismiss();
      const message =
        err.response?.data?.error || err.message || "Delete failed";
      toast.error(message);
    },
  });
};

export default useDeleteMediaMutation;