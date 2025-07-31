import { useState, useEffect } from "react";
import { Scaling } from "lucide-react";
import { Media, isImageMedia } from "@/types/media";
import toast from "react-hot-toast";
import useResizeUnifiedMediaMutation from "./hooks/use-resize-unified-media-mutation";
import { useQueryClient } from "@tanstack/react-query";
import { constant } from "@/lib/constant";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";

interface ResizeMediaModalProps {
  media: Media;
  isSelecting?: boolean;
}

const ResizeMediaModal: React.FC<ResizeMediaModalProps> = ({ media, isSelecting }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [resizeFormData, setResizeFormData] = useState({ width: 0, height: 0 });
  const queryClient = useQueryClient();

  const resizeMediaMutation = useResizeUnifiedMediaMutation({
    onSuccess: () => {
      setIsOpen(false);
      toast.dismiss();
      const toastId = toast.loading("Processing...");
      toast.success("Media resized!", { id: toastId, duration: 1500 });
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.media(media.mediaType),
      });
    },
    onError: (err: Error) => {
      toast.dismiss();
      const toastId = toast.loading("Processing...");
      toast.error(err.message, { id: toastId, duration: 4000 });
    },
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;

    const formattedValue = value.replace(/\D/g, "");
    const positiveIntegerValue =
      formattedValue === "" ? 0 : parseInt(formattedValue, 10);

    setResizeFormData({
      ...resizeFormData,
      [name]: positiveIntegerValue,
    });
  };

  const handleResizeMedia = () => {
    const { width, height } = resizeFormData;

    if (!width || !height) {
      toast.error("Width and height are required!");
      return;
    }

    resizeMediaMutation.mutate({
      media,
      width: Math.abs(Math.floor(width)),
      height: Math.abs(Math.floor(height)),
    });
  };

  useEffect(() => {
    if (isImageMedia(media)) {
      toast.dismiss();
      setResizeFormData({
        width: media.width || 0,
        height: media.height || 0,
      });
    }
  }, [media]);

  // Only show resize option for images
  if (!isImageMedia(media)) {
    return null;
  }

  return (
    <Dialog open={isOpen} onOpenChange={setIsOpen}>
      <Tooltip>
        <TooltipTrigger>
          <DialogTrigger asChild>
            <Button
              size="icon"
              variant="ghost"
              className="text-sky-600"
              disabled={isSelecting}
            >
              <Scaling className="inline" size="24" />
            </Button>
          </DialogTrigger>
        </TooltipTrigger>
        <TooltipContent side="bottom">
          <p>Resize Image</p>
        </TooltipContent>
      </Tooltip>
      <DialogContent className="sm:max-w-[425px]">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleResizeMedia();
          }}
        >
          <DialogHeader>
            <DialogTitle>Resize image</DialogTitle>
            <DialogDescription>
              Change the height and width of the image. Click save when you're
              done.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="height" className="text-right">
                Height
              </Label>
              <Input
                name="height"
                value={resizeFormData.height}
                onChange={handleInputChange}
                className="col-span-3"
              />
            </div>

            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="width" className="text-right">
                Width
              </Label>
              <Input
                name="width"
                value={resizeFormData.width}
                onChange={handleInputChange}
                className="col-span-3"
              />
            </div>
          </div>
          <DialogFooter>
            <Button type="submit">Save changes</Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
};

export default ResizeMediaModal;