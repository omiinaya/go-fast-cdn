import { constant } from "@/lib/constant";
import { mediaService, MediaRenameParams } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { MediaType } from "@/types/media";
import { useMutation, useQueryClient, UseMutationOptions } from "@tanstack/react-query";
import { AxiosError } from "axios";
import toast from "react-hot-toast";

interface RenameFileParams {
  fileName: string;
  newFileName: string;
}

const useRenameFileMutation = (
  type: "doc" | "image",
  options?: UseMutationOptions<string, Error, RenameFileParams>
) => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ fileName, newFileName }: RenameFileParams) => {
      // Map the old type to the new MediaType
      const mediaType: MediaType = type === "image" ? "image" : "document";
      
      // Use the unified media service to rename media
      const renameParams: MediaRenameParams = {
        fileName,
        newFileName,
        mediaType,
      };
      
      const result = await mediaService.renameMedia(renameParams);
      return result.status;
    },
    onSuccess: (data, { fileName, newFileName }) => {
      toast.dismiss();
      toast.success(`Successfully renamed ${fileName} to ${newFileName}!`);
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.images(
          type === "image" ? "images" : "documents"
        ),
      });
      options?.onSuccess?.(data, { fileName, newFileName }, undefined);
    },
    onError: (error: unknown, variables) => {
      const err = error as AxiosError<IErrorResponse>;
      toast.dismiss();
      const message =
        err.response?.data?.error || err.message || "Rename failed";
      toast.error(message);
      options?.onError?.(error as Error, variables, undefined);
    },
  });
};

export default useRenameFileMutation;
