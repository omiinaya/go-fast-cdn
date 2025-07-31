import { constant } from "@/lib/constant";
import { mediaService } from "@/services/mediaService";
import { MediaType, convertMediaToFile } from "@/types/media";
import { TFile } from "@/types/file";
import { useQuery } from "@tanstack/react-query";

type GetFilesParams = {
  type: "documents" | "images";
};

const useGetFilesQuery = ({ type }: GetFilesParams) => {
  // Map the old type to the new MediaType
  const mediaType: MediaType = type === "images" ? "image" : "document";
  
  return useQuery({
    queryKey: constant.queryKeys.images(type),
    queryFn: async (): Promise<TFile[]> => {
      // Use the unified media service to get media
      const media = await mediaService.getAllMedia({ mediaType });
      
      // Convert Media objects back to TFile for backward compatibility
      return media.map(convertMediaToFile);
    },
  });
};

export default useGetFilesQuery;
