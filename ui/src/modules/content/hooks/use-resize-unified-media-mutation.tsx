import { mediaService, MediaResizeParams } from "@/services/mediaService";
import { useMutation, UseMutationOptions } from "@tanstack/react-query";
import { Media, isImageMedia } from "@/types/media";

interface ResizeUnifiedMediaParams {
  media: Media;
  width: number;
  height: number;
}

const useResizeUnifiedMediaMutation = (
  options?: UseMutationOptions<
    { status: string; width: number; height: number; type: string; message: string },
    Error,
    ResizeUnifiedMediaParams
  >
) => {
  return useMutation({
    mutationFn: async (data: ResizeUnifiedMediaParams) => {
      // Only allow resizing for image media
      if (!isImageMedia(data.media)) {
        throw new Error("Only images can be resized");
      }
      
      const resizeParams: MediaResizeParams = {
        filename: data.media.fileName,
        width: data.width,
        height: data.height,
      };
      
      return mediaService.resizeImage(resizeParams);
    },
    ...options,
  });
};

export default useResizeUnifiedMediaMutation;