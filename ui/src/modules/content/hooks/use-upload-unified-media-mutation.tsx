import { constant } from "@/lib/constant";
import { mediaService, MediaUploadResponse } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { MediaType } from "@/types/media";
import toast from "react-hot-toast";

interface UploadMediaParams {
  file: File;
  mediaType: MediaType;
  filename?: string;
}

const useUploadUnifiedMediaMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async ({ file, mediaType, filename }: UploadMediaParams): Promise<MediaUploadResponse> => {
      return mediaService.uploadMedia(file, mediaType, filename);
    },
    onSuccess: (_data: MediaUploadResponse, { mediaType }: UploadMediaParams) => {
      toast.dismiss();
      toast.success("Successfully uploaded media!");
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
        err.response?.data?.error || err.message || "Upload failed";
      toast.error(message);
    },
  });
};

export default useUploadUnifiedMediaMutation;