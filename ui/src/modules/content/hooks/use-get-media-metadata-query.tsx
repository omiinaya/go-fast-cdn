import { mediaService, MediaMetadataResponse } from "@/services/mediaService";
import { MediaType } from "@/types/media";
import { useQuery } from "@tanstack/react-query";

interface GetMediaMetadataParams {
  fileName: string;
  mediaType: MediaType;
}

const useGetMediaMetadataQuery = ({ fileName, mediaType }: GetMediaMetadataParams) => {
  return useQuery({
    queryKey: ['media-metadata', fileName, mediaType],
    queryFn: async (): Promise<MediaMetadataResponse> => {
      return mediaService.getMediaMetadata(fileName, mediaType);
    },
    enabled: !!fileName && !!mediaType,
  });
};

export default useGetMediaMetadataQuery;