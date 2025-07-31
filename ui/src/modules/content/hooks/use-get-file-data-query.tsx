import { constant } from "@/lib/constant";
import { mediaService, MediaMetadataResponse } from "@/services/mediaService";
import { MediaType } from "@/types/media";
import { FileMetadata } from "@/types/fileMetadata";
import { useQuery } from "@tanstack/react-query";

type FileDataParams = {
  filename: string;
  type: "documents" | "images";
};

const useGetFileDataQuery = ({ filename, type }: FileDataParams) => {
  // Map the old type to the new MediaType
  const mediaType: MediaType = type === "images" ? "image" : "document";
  
  return useQuery({
    queryKey: constant.queryKeys.image(filename),
    queryFn: async (): Promise<FileMetadata> => {
      // Use the unified media service to get media metadata
      const metadataResponse: MediaMetadataResponse = await mediaService.getMediaMetadata(filename, mediaType);
      
      // Convert the MediaMetadataResponse to FileMetadata for backward compatibility
      return {
        download_url: metadataResponse.download_url,
        file_size: metadataResponse.file_size,
        filename: metadataResponse.filename,
        width: metadataResponse.width,
        height: metadataResponse.height,
      } as FileMetadata;
    },
  });
};

export default useGetFileDataQuery;
