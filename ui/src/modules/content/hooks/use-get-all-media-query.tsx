import { constant } from "@/lib/constant";
import { mediaService } from "@/services/mediaService";
import { Media, MediaType } from "@/types/media";
import { useQuery } from "@tanstack/react-query";

interface GetAllMediaParams {
  mediaType?: MediaType;
}

const useGetAllMediaQuery = (params?: GetAllMediaParams) => {
  return useQuery({
    queryKey: constant.queryKeys.media(params?.mediaType || 'all'),
    queryFn: async (): Promise<Media[]> => {
      return mediaService.getAllMedia(params);
    },
  });
};

export default useGetAllMediaQuery;