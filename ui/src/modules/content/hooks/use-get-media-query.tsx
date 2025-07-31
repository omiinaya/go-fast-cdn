import { constant } from "@/lib/constant";
import { mediaService } from "@/services/mediaService";
import { Media, MediaType } from "@/types/media";
import { useQuery } from "@tanstack/react-query";

type GetMediaParams = {
  mediaType: MediaType;
};

const useGetMediaQuery = ({ mediaType }: GetMediaParams) => {
  return useQuery({
    queryKey: constant.queryKeys.media(mediaType),
    queryFn: async (): Promise<Media[]> => {
      // Use the unified media service to get media
      return mediaService.getAllMedia({ mediaType });
    },
  });
};

export default useGetMediaQuery;