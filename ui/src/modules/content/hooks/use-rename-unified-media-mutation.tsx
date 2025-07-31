import { constant } from "@/lib/constant";
import { mediaService, MediaRenameParams } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import toast from "react-hot-toast";

const useRenameUnifiedMediaMutation = (options?: {
  onSuccess?: () => void;
  onError?: (error: Error) => void;
}) => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (params: MediaRenameParams) => {
      return mediaService.renameMedia(params);
    },
    onSuccess: (_data: { status: string }, { mediaType }: MediaRenameParams) => {
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

export default useRenameUnifiedMediaMutation;