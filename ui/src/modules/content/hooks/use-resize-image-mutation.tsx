import { constant } from "@/lib/constant";
import { mediaService, MediaResizeParams } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { useMutation, useQueryClient, UseMutationOptions } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { MediaType } from "@/types/media";
import toast from "react-hot-toast";

interface ResizeImageParams {
  filename: string;
  width: number;
  height: number;
}

const useResizeImageMutation = (
  options?: UseMutationOptions<
    { status: string; width: number; height: number; type: MediaType; message: string },
    Error,
    ResizeImageParams
  >
) => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async (data: ResizeImageParams) => {
      // Use the unified media service to resize image
      const resizeParams: MediaResizeParams = {
        filename: data.filename,
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

export default useResizeImageMutation;
