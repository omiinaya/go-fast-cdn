import { useState } from "react";
import { Edit } from "lucide-react";
import { Media } from "@/types/media";
import toast from "react-hot-toast";
import useRenameUnifiedMediaMutation from "./hooks/use-rename-unified-media-mutation";
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

interface RenameMediaModalProps {
  media: Media;
  isSelecting?: boolean;
}

const RenameMediaModal: React.FC<RenameMediaModalProps> = ({ media, isSelecting }) => {
  const [isOpen, setIsOpen] = useState(false);
  const [newFileName, setNewFileName] = useState(media.fileName);
  const queryClient = useQueryClient();

  const renameMediaMutation = useRenameUnifiedMediaMutation({
    onSuccess: () => {
      setIsOpen(false);
      toast.dismiss();
      toast.success("File renamed successfully!");
      queryClient.invalidateQueries({
        queryKey: constant.queryKeys.media(media.mediaType),
      });
    },
    onError: (err: Error) => {
      toast.dismiss();
      toast.error(err.message);
    },
  });

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setNewFileName(e.target.value);
  };

  const handleRenameFile = () => {
    if (!newFileName.trim()) {
      toast.error("Filename cannot be empty!");
      return;
    }

    renameMediaMutation.mutate({
      fileName: media.fileName,
      newFileName,
      mediaType: media.mediaType,
    });
  };

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
              <Edit className="inline" size="24" />
            </Button>
          </DialogTrigger>
        </TooltipTrigger>
        <TooltipContent side="bottom">
          <p>Rename File</p>
        </TooltipContent>
      </Tooltip>
      <DialogContent className="sm:max-w-[425px]">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleRenameFile();
          }}
        >
          <DialogHeader>
            <DialogTitle>Rename file</DialogTitle>
            <DialogDescription>
              Change the filename of the media. Click save when you're done.
            </DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-4">
            <div className="grid grid-cols-4 items-center gap-4">
              <Label htmlFor="filename" className="text-right">
                Filename
              </Label>
              <Input
                id="filename"
                value={newFileName}
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

export default RenameMediaModal;