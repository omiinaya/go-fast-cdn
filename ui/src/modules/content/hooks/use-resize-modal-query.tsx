import { constant } from "@/lib/constant";
import { mediaService, MediaMetadataResponse } from "@/services/mediaService";
import { FileMetadata } from "@/types/fileMetadata";
import { useQuery } from "@tanstack/react-query";

const useResizeModalQuery = (filename: string) => {
  return useQuery({
    queryKey: constant.queryKeys.image(filename),
    queryFn: async (): Promise<FileMetadata> => {
      // Use the unified media service to get image metadata
      const metadataResponse: MediaMetadataResponse = await mediaService.getMediaMetadata(filename, "image");
      
      // Convert the MediaMetadataResponse to FileMetadata for backward compatibility
      return {
        download_url: metadataResponse.download_url,
        file_size: metadataResponse.file_size,
        filename: metadataResponse.filename,
        width: metadataResponse.width,
        height: metadataResponse.height,
      } as FileMetadata;
    },
    enabled: !!filename,
  });
};

export default useResizeModalQuery;
