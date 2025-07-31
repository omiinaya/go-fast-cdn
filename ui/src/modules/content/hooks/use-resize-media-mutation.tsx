import { constant } from "@/lib/constant";
import { mediaService, MediaResizeParams } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { useMutation, useQueryClient, UseMutationOptions } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { Media, MediaType, isImageMedia } from "@/types/media";
import toast from "react-hot-toast";

interface ResizeMediaParams {
  media: Media;
  width: number;
  height: number;
}

const useResizeMediaMutation = (
  options?: UseMutationOptions<
    { status: string; width: number; height: number; type: MediaType; message: string },
    Error,
    ResizeMediaParams
  >
) => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (data: ResizeMediaParams) => {
      // Only allow resizing for image media
      if (!isImageMedia(data.media)) {
        throw new Error("Only images can be resized");
      }
      
      // Use the unified media service to resize image
      const resizeParams: MediaResizeParams = {
        filename: data.media.fileName,
        width: data.width,
        height: data.height,
      };
      
      return mediaService.resizeImage(resizeParams);
    },
    onSuccess: (data, variables) => {
      toast.dismiss();
      toast.success(`Successfully resized image to ${data.width}x${data.height}!`);
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.media("image"),
      });
      options?.onSuccess?.(data, variables, undefined);
    },
    onError: (error: unknown, variables) => {
      const err = error as AxiosError<IErrorResponse>;
      toast.dismiss();
      const message =
        err.response?.data?.error || err.message || "Resize failed";
      toast.error(message);
      options?.onError?.(error as Error, variables, undefined);
    },
    ...options,
  });
};

export default useResizeMediaMutation;