import { constant } from "@/lib/constant";
import { mediaService } from "@/services/mediaService";
import { IErrorResponse } from "@/types/response";
import { MediaType, convertMediaToFile } from "@/types/media";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { AxiosError } from "axios";
import { TFile } from "@/types/file";
import toast from "react-hot-toast";

interface UploadFileParams {
  file: File;
  type: "doc" | "image";
  filename?: string;
}

const useUploadFileMutation = () => {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ file, type, filename }: UploadFileParams): Promise<TFile> => {
      // Map the old type to the new MediaType
      const mediaType: MediaType = type === "image" ? "image" : "document";
      
      // Use the unified media service to upload media
      await mediaService.uploadMedia(file, mediaType, filename);
      
      // Get the uploaded media to convert to TFile for backward compatibility
      const mediaList = await mediaService.getAllMedia({ mediaType });
      const uploadedMedia = mediaList.find(m => m.fileName === (filename || file.name));
      
      if (!uploadedMedia) {
        throw new Error("Failed to retrieve uploaded media");
      }
      
      return convertMediaToFile(uploadedMedia);
    },
    onSuccess: (_data: TFile, { type }: UploadFileParams) => {
      toast.dismiss();
      toast.success("Successfully uploaded file!");
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.size(),
      });
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.images(type === "image" ? "images" : "documents"),
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

export default useUploadFileMutation;
