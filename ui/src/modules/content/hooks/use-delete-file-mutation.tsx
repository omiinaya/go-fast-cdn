import { constant } from "@/lib/constant";
import { mediaService, MediaDeleteParams } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { MediaType } from "@/types/media";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import toast from "react-hot-toast";

const useDeleteFileMutation = (type: "doc" | "image") => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (filename: string) => {
      // Map the old type to the new MediaType
      const mediaType: MediaType = type === "image" ? "image" : "document";
      
      // Use the unified media service to delete media
      const deleteParams: MediaDeleteParams = {
        fileName: filename,
        mediaType,
      };
      
      return mediaService.deleteMedia(deleteParams);
    },
    onSuccess: () => {
      toast.dismiss();
      toast.success("Successfully deleted file!");
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.size(),
      });
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.images(
          type === "image" ? "images" : "documents"
        ),
      });
    },
    onError: (error) => {
      const err = error as AxiosError<IErrorResponse>;
      toast.dismiss();
      const message =
        err.response?.data?.error || err.message || "Delete failed";
      toast.error(message);
    },
  });
};

export default useDeleteFileMutation;
