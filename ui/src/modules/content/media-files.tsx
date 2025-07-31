import MediaCard from "./media-card";
import useGetAllMediaQuery from "./hooks/use-get-all-media-query";
import useDeleteUnifiedMediaMutation from "./hooks/use-delete-unified-media-mutation";
import { Input } from "@/components/ui/input";
import { useCallback, useEffect, useMemo, useState } from "react";
import { Button, buttonVariants } from "@/components/ui/button";
import { List, Trash, X } from "lucide-react";
import { cn } from "@/lib/utils";
import MainContentWrapper from "@/components/layouts/main-content-wrapper";
import { Skeleton } from "@/components/ui/skeleton";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { constant } from "@/lib/constant";
import { AxiosError } from "axios";
import { IErrorResponse } from "@/types/response";
import { MediaType } from "@/types/media";
import toast from "react-hot-toast";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import UploadMediaModal from "./upload/upload-media-modal";

type TMediaFilesProps = {
  mediaType: MediaType;
};

const MediaFiles: React.FC<TMediaFilesProps> = ({ mediaType }) => {
  const media = useGetAllMediaQuery({ mediaType });
  const [search, setSearch] = useState("");
  const [debounceSearch, setDebounceSearch] = useState("");

  const [selectedMedia, setSelectedMedia] = useState<string[]>([]);
  const [isSelecting, setIsSelecting] = useState(false);

  const deleteMutation = useDeleteUnifiedMediaMutation();

  const queryClient = useQueryClient();

  const {
    mutateAsync: deleteMediaMutation,
    isPending: isDeletingMediaLoading,
  } = useMutation({
    mutationFn: () =>
      Promise.all(
        selectedMedia.map((fileName: string) =>
          deleteMutation.mutateAsync({ fileName, mediaType })
        )
      ),
    onSuccess: () => {
      setSelectedMedia([]);
      setIsSelecting(false);
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.media(mediaType),
      });
    },
    onError: (error: unknown) => {
      const errorResponse = error as AxiosError<IErrorResponse>;
      const message = errorResponse.response?.data?.error || "Delete failed";
      toast.error(message);
    },
  });

  const filteredMedia = useMemo(() => {
    return media.data?.filter((item) =>
      item.fileName.toLowerCase().includes(search.toLowerCase())
    );
  }, [media.data, search]);

  const handleSearchChange = useCallback(
    (event: React.ChangeEvent<HTMLInputElement>) => {
      setDebounceSearch(event.target.value);
    },
    []
  );

  const handleOnSelectMedia = useCallback(
    (fileName: string) => {
      if (selectedMedia.includes(fileName)) {
        setSelectedMedia((prev) => prev.filter((name) => name !== fileName));
      } else {
        setSelectedMedia((prev) => [...prev, fileName]);
      }
    },
    [selectedMedia]
  );

  useEffect(() => {
    const handler = setTimeout(() => {
      setSearch(debounceSearch);
    }, 300);

    return () => {
      clearTimeout(handler);
    };
  }, [debounceSearch]);

  const getMediaTypeName = () => {
    switch (mediaType) {
      case "image": return "Images";
      case "document": return "Documents";
      case "video": return "Videos";
      case "audio": return "Audio";
      default: return "Files";
    }
  };

  return (
    <MainContentWrapper title={getMediaTypeName()}>
      <div className="flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <section className="flex items-center gap-2">
            <Input
              className="w-full max-w-md min-w-xs"
              placeholder="Search files by name"
              value={debounceSearch}
              onChange={handleSearchChange}
              aria-label="Search files"
            />
            <Button
              variant="outline"
              className={cn({
                hidden: !search,
              })}
              onClick={() => {
                setSearch("");
                setDebounceSearch("");
              }}
              aria-label="Clear search"
            >
              <X />
              Clear Search
            </Button>
          </section>
          <section className="flex items-center gap-2">
            {isSelecting ? (
              <>
                <Button
                  onClick={() => {
                    setIsSelecting(false);
                    setSelectedMedia([]);
                  }}
                  variant="outline"
                  disabled={isDeletingMediaLoading}
                >
                  <X />
                  {selectedMedia.length}
                  {selectedMedia.length === 1
                    ? " File Selected"
                    : " Files Selected"}
                </Button>
                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button
                      variant="destructive"
                      disabled={
                        selectedMedia.length === 0 || isDeletingMediaLoading
                      }
                    >
                      <Trash />
                      Delete Selected Files
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>
                        Are you sure you want to delete these files?
                      </AlertDialogTitle>
                      <AlertDialogDescription>
                        This action cannot be undone. All selected files will be
                        permanently deleted.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction
                        className={buttonVariants({
                          variant: "destructive",
                        })}
                        onClick={() => deleteMediaMutation()}
                      >
                        Continue
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </>
            ) : (
              <>
                <Button onClick={() => setIsSelecting(true)} variant="outline">
                  <List />
                  Select
                </Button>
                <UploadMediaModal placement="header" mediaType={mediaType} />
              </>
            )}
          </section>
        </div>
        <div className="flex flex-wrap gap-4">
          {media.isLoading ? (
            <>
              <Skeleton className="min-h-[264px] w-64 max-w-[256px]" />
              <Skeleton className="min-h-[264px] w-64 max-w-[256px]" />
              <Skeleton className="min-h-[264px] w-64 max-w-[256px]" />
              <Skeleton className="min-h-[264px] w-64 max-w-[256px]" />
              <Skeleton className="min-h-[264px] w-64 max-w-[256px]" />
              <Skeleton className="min-h-[264px] w-64 max-w-[256px]" />
            </>
          ) : (
            <>
              {filteredMedia?.map((item) => (
                <MediaCard
                  media={item}
                  key={item.id}
                  isSelecting={isSelecting}
                  isSelected={selectedMedia.includes(item.fileName)}
                  onSelect={handleOnSelectMedia}
                />
              ))}
            </>
          )}
        </div>
      </div>
    </MainContentWrapper>
  );
};

export default MediaFiles;