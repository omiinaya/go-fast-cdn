import { TabsContent } from "@/components/ui/tabs";
import React, { useState, useEffect } from "react";
import { MediaType } from "@/types/media";
import { fileTypeService } from "@/services/fileTypeService";

interface FileInputProps {
  type: MediaType;
  fileRef: React.RefObject<HTMLInputElement>;
  onFileChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
}

const FileInput = ({ type, fileRef, onFileChange }: FileInputProps) => {
  const [acceptType, setAcceptType] = useState<string>("*/*");
  const [fileTypesLoaded, setFileTypesLoaded] = useState(false);

  // Load file type configuration on component mount
  useEffect(() => {
    const loadFileTypes = async () => {
      try {
        await fileTypeService.loadConfig();
        setFileTypesLoaded(true);
        setAcceptType(fileTypeService.getAcceptAttribute(type));
      } catch (err) {
        console.error('Failed to load file type configuration:', err);
        // Keep default accept type on error
      }
    };

    loadFileTypes();
  }, [type]);

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
      <p className="text-gray-500 text-center">
        {UPLOAD_MESSAGES[type]}
        {!fileTypesLoaded && <span className="text-xs text-gray-400 ml-2">Loading file types...</span>}
      </p>
      <input
        type="file"
        accept={acceptType}
        multiple
        name={type}
        id={type}
        aria-label={`Select ${getTypeName(type)}`}
        ref={fileRef}
        className="hidden"
        onChange={onFileChange}
        disabled={!fileTypesLoaded}
      />
    </TabsContent>
  );
};

export default FileInput;
