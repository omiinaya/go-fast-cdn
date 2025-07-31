import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import UploadMediaForm from '../upload-media-form';
import { MediaType } from '@/types/media';

// Mock the toast notifications
jest.mock('react-hot-toast', () => ({
  toast: {
    success: jest.fn(),
    error: jest.fn(),
    dismiss: jest.fn(),
  },
}));

// Import the mocked toast
import { toast } from 'react-hot-toast';

describe('File Type Validation', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  const createMockFile = (name: string, type: string, size: number): File => {
    const file = new File([''], name, { type });
    Object.defineProperty(file, 'size', { value: size });
    return file;
  };

  it('should accept JS files for document upload', () => {
    const mockOnChangeFiles = jest.fn();
    render(
      <UploadMediaForm
        files={[]}
        onChangeFiles={mockOnChangeFiles}
        isLoading={false}
        mediaType="document"
        onChangeMediaType={jest.fn()}
      />
    );

    const fileInput = screen.getByLabelText('Select media files');
    expect(fileInput).toHaveAttribute('accept', expect.stringContaining('application/javascript'));
  });

  it('should accept CSS files for document upload', () => {
    const mockOnChangeFiles = jest.fn();
    render(
      <UploadMediaForm
        files={[]}
        onChangeFiles={mockOnChangeFiles}
        isLoading={false}
        mediaType="document"
        onChangeMediaType={jest.fn()}
      />
    );

    const fileInput = screen.getByLabelText('Select media files');
    expect(fileInput).toHaveAttribute('accept', expect.stringContaining('text/css'));
  });

  it('should accept JSON files for document upload', () => {
    const mockOnChangeFiles = jest.fn();
    render(
      <UploadMediaForm
        files={[]}
        onChangeFiles={mockOnChangeFiles}
        isLoading={false}
        mediaType="document"
        onChangeMediaType={jest.fn()}
      />
    );

    const fileInput = screen.getByLabelText('Select media files');
    expect(fileInput).toHaveAttribute('accept', expect.stringContaining('application/json'));
  });

  it('should accept YAML files for document upload', () => {
    const mockOnChangeFiles = jest.fn();
    render(
      <UploadMediaForm
        files={[]}
        onChangeFiles={mockOnChangeFiles}
        isLoading={false}
        mediaType="document"
        onChangeMediaType={jest.fn()}
      />
    );

    const fileInput = screen.getByLabelText('Select media files');
    expect(fileInput).toHaveAttribute('accept', expect.stringContaining('application/x-yaml'));
  });

  it('should successfully validate JS files when dropped', async () => {
    const mockOnChangeFiles = jest.fn();
    render(
      <UploadMediaForm
        files={[]}
        onChangeFiles={mockOnChangeFiles}
        isLoading={false}
        mediaType="document"
        onChangeMediaType={jest.fn()}
      />
    );

    const dropZone = screen.getByText('Drop your media files here, or click to select files.').closest('div') as HTMLDivElement;
    const mockJsFile = createMockFile('test-script.js', 'application/javascript', 1024);

    // Simulate file drop
    fireEvent.drop(dropZone!, {
      dataTransfer: {
        files: [mockJsFile],
      },
    });

    await waitFor(() => {
      expect(mockOnChangeFiles).toHaveBeenCalledWith([mockJsFile]);
      expect(toast.error).not.toHaveBeenCalled();
    });
  });

  it('should successfully validate CSS files when dropped', async () => {
    const mockOnChangeFiles = jest.fn();
    render(
      <UploadMediaForm
        files={[]}
        onChangeFiles={mockOnChangeFiles}
        isLoading={false}
        mediaType="document"
        onChangeMediaType={jest.fn()}
      />
    );

    const dropZone = screen.getByText('Drop your media files here, or click to select files.').closest('div') as HTMLDivElement;
    const mockCssFile = createMockFile('test-styles.css', 'text/css', 1024);

    // Simulate file drop
    fireEvent.drop(dropZone!, {
      dataTransfer: {
        files: [mockCssFile],
      },
    });

    await waitFor(() => {
      expect(mockOnChangeFiles).toHaveBeenCalledWith([mockCssFile]);
      expect(toast.error).not.toHaveBeenCalled();
    });
  });

  it('should successfully validate JSON files when dropped', async () => {
    const mockOnChangeFiles = jest.fn();
    render(
      <UploadMediaForm
        files={[]}
        onChangeFiles={mockOnChangeFiles}
        isLoading={false}
        mediaType="document"
        onChangeMediaType={jest.fn()}
      />
    );

    const dropZone = screen.getByText('Drop your media files here, or click to select files.').closest('div') as HTMLDivElement;
    const mockJsonFile = createMockFile('test-data.json', 'application/json', 1024);

    // Simulate file drop
    fireEvent.drop(dropZone!, {
      dataTransfer: {
        files: [mockJsonFile],
      },
    });

    await waitFor(() => {
      expect(mockOnChangeFiles).toHaveBeenCalledWith([mockJsonFile]);
      expect(toast.error).not.toHaveBeenCalled();
    });
  });

  it('should successfully validate YAML files when dropped', async () => {
    const mockOnChangeFiles = jest.fn();
    render(
      <UploadMediaForm
        files={[]}
        onChangeFiles={mockOnChangeFiles}
        isLoading={false}
        mediaType="document"
        onChangeMediaType={jest.fn()}
      />
    );

    const dropZone = screen.getByText('Drop your media files here, or click to select files.').closest('div') as HTMLDivElement;
    const mockYamlFile = createMockFile('test-config.yaml', 'application/x-yaml', 1024);

    // Simulate file drop
    fireEvent.drop(dropZone!, {
      dataTransfer: {
        files: [mockYamlFile],
      },
    });

    await waitFor(() => {
      expect(mockOnChangeFiles).toHaveBeenCalledWith([mockYamlFile]);
      expect(toast.error).not.toHaveBeenCalled();
    });
  });

  it('should accept all supported file types for document upload', () => {
    const mockOnChangeFiles = jest.fn();
    render(
      <UploadMediaForm
        files={[]}
        onChangeFiles={mockOnChangeFiles}
        isLoading={false}
        mediaType="document"
        onChangeMediaType={jest.fn()}
      />
    );

    const fileInput = screen.getByLabelText('Select media files');
    const acceptAttribute = fileInput.getAttribute('accept');
    
    // Check that all the new file types are included
    expect(acceptAttribute).toContain('application/javascript');
    expect(acceptAttribute).toContain('text/css');
    expect(acceptAttribute).toContain('application/json');
    expect(acceptAttribute).toContain('application/x-yaml');
    
    // Check that some of the existing file types are still there
    expect(acceptAttribute).toContain('application/pdf');
    expect(acceptAttribute).toContain('text/plain');
    expect(acceptAttribute).toContain('application/msword');
  });
});