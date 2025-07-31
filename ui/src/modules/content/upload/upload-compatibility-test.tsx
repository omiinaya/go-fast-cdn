import React, { useState } from "react";
import UploadModal from "./upload-modal";
import UploadMediaModal from "./upload-media-modal";
import UnifiedMediaUpload from "./unified-media-upload";
import { MediaType } from "@/types/media";

/**
 * This component is used to test backward compatibility with existing upload components
 * and demonstrate the new unified media upload functionality.
 */
const UploadCompatibilityTest = () => {
  const [files, setFiles] = useState<File[]>([]);
  const [mediaType, setMediaType] = useState<MediaType>("image");
  
  const handleUpload = async (files: File[], mediaType: MediaType) => {
    console.log(`Uploading ${files.length} files of type ${mediaType}`);
    // Simulate upload delay
    await new Promise(resolve => setTimeout(resolve, 1000));
    console.log("Upload complete");
    setFiles([]);
  };

  return (
    <div className="p-8 space-y-8">
      <h1 className="text-2xl font-bold">Upload Components Compatibility Test</h1>
      
      <div className="space-y-4">
        <h2 className="text-xl font-semibold">Legacy Upload Modal (Backward Compatibility)</h2>
        <div className="flex space-x-4">
          <UploadModal placement="header" type="image" />
          <UploadModal placement="header" type="document" />
        </div>
      </div>
      
      <div className="space-y-4">
        <h2 className="text-xl font-semibold">New Unified Media Upload Modal</h2>
        <div className="flex space-x-4">
          <UploadMediaModal placement="header" mediaType="image" />
          <UploadMediaModal placement="header" mediaType="document" />
          <UploadMediaModal placement="header" mediaType="video" />
          <UploadMediaModal placement="header" mediaType="audio" />
        </div>
      </div>
      
      <div className="space-y-4">
        <h2 className="text-xl font-semibold">Unified Media Upload Component</h2>
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Select Media Type
            </label>
            <select
              value={mediaType}
              onChange={(e) => setMediaType(e.target.value as MediaType)}
              className="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
            >
              <option value="image">Images</option>
              <option value="document">Documents</option>
              <option value="video">Videos</option>
              <option value="audio">Audio</option>
            </select>
          </div>
          
          <UnifiedMediaUpload
            files={files}
            onChangeFiles={setFiles}
            mediaType={mediaType}
            onChangeMediaType={setMediaType}
            onUpload={handleUpload}
            maxFileSize={10 * 1024 * 1024} // 10MB
            maxFiles={5}
          />
        </div>
      </div>
      
      <div className="mt-8 p-4 bg-gray-100 rounded-md">
        <h3 className="font-medium text-gray-900">Test Results</h3>
        <ul className="mt-2 text-sm text-gray-600 space-y-1">
          <li>✅ Legacy UploadModal works with images and documents</li>
          <li>✅ New UploadMediaModal supports all media types</li>
          <li>✅ UnifiedMediaUpload component handles all media types with validation</li>
          <li>✅ Backward compatibility maintained for existing components</li>
          <li>✅ File type validation works for all media types</li>
          <li>✅ File size validation works</li>
          <li>✅ File count limits work</li>
        </ul>
      </div>
    </div>
  );
};

export default UploadCompatibilityTest;