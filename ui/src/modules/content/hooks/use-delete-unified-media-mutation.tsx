import { constant } from "@/lib/constant";
import { mediaService, MediaDeleteParams } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import toast from "react-hot-toast";

const useDeleteUnifiedMediaMutation = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: async (params: MediaDeleteParams) => {
      return mediaService.deleteMedia(params);
    },
    onSuccess: (_data: { message: string; fileName: string }, { mediaType }: MediaDeleteParams) => {
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

export default useDeleteUnifiedMediaMutation;