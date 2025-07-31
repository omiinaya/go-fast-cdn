import { TabsContent } from "@/components/ui/tabs";
import React from "react";
import { MediaType } from "@/types/media";

interface FileInputProps {
  type: MediaType;
  fileRef: React.RefObject<HTMLInputElement>;
  onFileChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

const FileInput = ({ type, fileRef, onFileChange }: FileInputProps) => {
  const ACCEPT_TYPES = {
    image: "image/jpeg,image/png,image/jpg,image/webp,image/gif,image/bmp,image/svg+xml",
    document: "text/plain,application/zip,application/msword,application/vnd.openxmlformats-officedocument.wordprocessingml.document,application/vnd.openxmlformats-officedocument.spreadsheetml.sheet,application/vnd.openxmlformats-officedocument.presentationml.presentation,application/pdf,application/rtf,application/x-freearc",
    video: "video/mp4,video/webm,video/ogg,video/quicktime,video/x-msvideo",
    audio: "audio/mpeg,audio/ogg,audio/wav,audio/webm,audio/aac",
    other: "*/*",
  } as const;

  const UPLOAD_MESSAGES = {
    image: "Drop your images here, or click to select files.",
    document: "Drop your documents here, or click to select files.",
    video: "Drop your videos here, or click to select files.",
    audio: "Drop your audio files here, or click to select files.",
    other: "Drop your files here, or click to select files.",
  } as const;

  const getTypeName = (type: MediaType): string => {
    switch (type) {
      case "image": return "images";
      case "document": return "documents";
      case "video": return "videos";
      case "audio": return "audio files";
      default: return "files";
    }
  };

  return (
    <TabsContent value={type}>
      <p className="text-gray-500 text-center">{UPLOAD_MESSAGES[type]}</p>
      <input
        type="file"
        accept={ACCEPT_TYPES[type]}
        multiple
        name={type}
        id={type}
        aria-label={`Select ${getTypeName(type)}`}
        ref={fileRef}
        className="hidden"
        onChange={onFileChange}
      />
    </TabsContent>
  );
};

export default FileInput;
