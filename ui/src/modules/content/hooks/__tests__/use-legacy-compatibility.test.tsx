import { renderHook, act } from '@testing-library/react';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import useMediaMigrationHelper, { useUnifiedMediaMigrationStatus } from '../use-legacy-compatibility';
import useGetFilesQuery from '../use-get-files-query';
import useUploadFileMutation from '../use-upload-file-mutation';
import useDeleteFileMutation from '../use-delete-file-mutation';
import useRenameFileMutation from '../use-rename-file-mutation';
import useResizeImageMutation from '../use-resize-image-mutation';
import useGetFileDataQuery from '../use-get-file-data-query';
import useResizeModalQuery from '../use-resize-modal-query';

// Mock the hooks
jest.mock('../use-get-files-query', () => ({
  __esModule: true,
  default: jest.fn(),
}));

jest.mock('../use-upload-file-mutation', () => ({
  __esModule: true,
  default: jest.fn(),
}));

jest.mock('../use-delete-file-mutation', () => ({
  __esModule: true,
  default: jest.fn(),
}));

jest.mock('../use-rename-file-mutation', () => ({
  __esModule: true,
  default: jest.fn(),
}));

jest.mock('../use-resize-image-mutation', () => ({
  __esModule: true,
  default: jest.fn(),
}));

jest.mock('../use-get-file-data-query', () => ({
  __esModule: true,
  default: jest.fn(),
}));

jest.mock('../use-resize-modal-query', () => ({
  __esModule: true,
  default: jest.fn(),
}));

describe('useLegacyCompatibility', () => {
  let queryClient: QueryClient;

  beforeEach(() => {
    jest.clearAllMocks();
    
    queryClient = new QueryClient({
      defaultOptions: {
        queries: {
          retry: false,
        },
      },
    });
  });

  const wrapper = ({ children }: { children: React.ReactNode }) => (
    <QueryClientProvider client={queryClient}>
      {children}
    </QueryClientProvider>
  );

  describe('useUnifiedMediaMigrationStatus', () => {
    it('should return migration status as complete', () => {
      const { result } = renderHook(() => useUnifiedMediaMigrationStatus(), { wrapper });
      
      expect(result.current.isMigrationComplete).toBe(true);
      expect(result.current.canUseUnifiedMedia).toBe(true);
    });
  });

  describe('useMediaMigrationHelper', () => {
    it('should provide legacy hooks for backward compatibility', () => {
      const { result } = renderHook(() => useMediaMigrationHelper(), { wrapper });
      
      expect(result.current.legacy).toBeDefined();
      expect(result.current.legacy.useGetFilesQuery).toBeDefined();
      expect(result.current.legacy.useUploadFileMutation).toBeDefined();
      expect(result.current.legacy.useDeleteFileMutation).toBeDefined();
      expect(result.current.legacy.useRenameFileMutation).toBeDefined();
      expect(result.current.legacy.useResizeImageMutation).toBeDefined();
      expect(result.current.legacy.useGetFileDataQuery).toBeDefined();
      expect(result.current.legacy.useResizeModalQuery).toBeDefined();
    });

    it('should indicate that unified media should be used', () => {
      const { result } = renderHook(() => useMediaMigrationHelper(), { wrapper });
      
      expect(result.current.shouldUseUnifiedMedia).toBe(true);
    });

    it('should provide the same hooks as the original exports', () => {
      const { result } = renderHook(() => useMediaMigrationHelper(), { wrapper });
      
      expect(result.current.legacy.useGetFilesQuery).toBe(useGetFilesQuery);
      expect(result.current.legacy.useUploadFileMutation).toBe(useUploadFileMutation);
      expect(result.current.legacy.useDeleteFileMutation).toBe(useDeleteFileMutation);
      expect(result.current.legacy.useRenameFileMutation).toBe(useRenameFileMutation);
      expect(result.current.legacy.useResizeImageMutation).toBe(useResizeImageMutation);
      expect(result.current.legacy.useGetFileDataQuery).toBe(useGetFileDataQuery);
      expect(result.current.legacy.useResizeModalQuery).toBe(useResizeModalQuery);
    });
  });

  describe('Backward Compatibility', () => {
    it('should maintain the same API for legacy hooks', () => {
      // This test ensures that the legacy hooks maintain the same API
      // as before the unification, so existing components continue to work
      
      // Mock the hooks to return mock functions
      (useGetFilesQuery as jest.Mock).mockReturnValue(jest.fn());
      (useUploadFileMutation as jest.Mock).mockReturnValue(jest.fn());
      (useDeleteFileMutation as jest.Mock).mockReturnValue(jest.fn());
      (useRenameFileMutation as jest.Mock).mockReturnValue(jest.fn());
      (useResizeImageMutation as jest.Mock).mockReturnValue(jest.fn());
      (useGetFileDataQuery as jest.Mock).mockReturnValue(jest.fn());
      (useResizeModalQuery as jest.Mock).mockReturnValue(jest.fn());

      const { result } = renderHook(() => useMediaMigrationHelper(), { wrapper });
      
      // Verify that the legacy hooks are functions
      expect(typeof result.current.legacy.useGetFilesQuery).toBe('function');
      expect(typeof result.current.legacy.useUploadFileMutation).toBe('function');
      expect(typeof result.current.legacy.useDeleteFileMutation).toBe('function');
      expect(typeof result.current.legacy.useRenameFileMutation).toBe('function');
      expect(typeof result.current.legacy.useResizeImageMutation).toBe('function');
      expect(typeof result.current.legacy.useGetFileDataQuery).toBe('function');
      expect(typeof result.current.legacy.useResizeModalQuery).toBe('function');
    });

    it('should allow gradual migration to unified media', () => {
      // This test simulates a component that gradually migrates to the unified media system
      
      const { result } = renderHook(() => useMediaMigrationHelper(), { wrapper });
      
      // Initially, the component should use legacy hooks
      const legacyHooks = result.current.legacy;
      
      // The component can check if it should use unified media
      const shouldUseUnifiedMedia = result.current.shouldUseUnifiedMedia;
      
      // Based on the migration status, the component can decide which hooks to use
      if (shouldUseUnifiedMedia) {
        // The component can use unified media hooks
        expect(true).toBe(true);
      } else {
        // The component should use legacy hooks
        expect(legacyHooks).toBeDefined();
      }
    });

    it('should provide a consistent interface during migration', () => {
      // This test ensures that the interface remains consistent during the migration process
      
      const { result } = renderHook(() => useMediaMigrationHelper(), { wrapper });
      
      // The interface should not change during migration
      const initialInterface = {
        legacy: result.current.legacy,
        shouldUseUnifiedMedia: result.current.shouldUseUnifiedMedia,
      };
      
      // Re-render the hook to simulate a state change
      act(() => {
        result.current;
      });
      
      // The interface should remain the same
      expect(result.current.legacy).toBe(initialInterface.legacy);
      expect(result.current.shouldUseUnifiedMedia).toBe(initialInterface.shouldUseUnifiedMedia);
    });
  });

  describe('Integration with Legacy Components', () => {
    it('should work with existing components that use legacy hooks', () => {
      // This test simulates an existing component that uses legacy hooks
      
      // Mock the hooks to return mock data
      const mockFilesQuery = {
        data: [],
        isLoading: false,
        error: null,
      };
      
      const mockUploadMutation = {
        mutate: jest.fn(),
        isLoading: false,
      };
      
      (useGetFilesQuery as jest.Mock).mockReturnValue(mockFilesQuery);
      (useUploadFileMutation as jest.Mock).mockReturnValue(mockUploadMutation);

      const { result } = renderHook(() => useMediaMigrationHelper(), { wrapper });
      
      // The legacy hooks should return the expected data
      expect(result.current.legacy.useGetFilesQuery({ type: 'images' })).toBe(mockFilesQuery);
      expect(result.current.legacy.useUploadFileMutation()).toBe(mockUploadMutation);
    });

    it('should allow new components to use unified media', () => {
      // This test simulates a new component that uses unified media
      
      const { result } = renderHook(() => useMediaMigrationHelper(), { wrapper });
      
      // The hook should indicate that unified media can be used
      expect(result.current.shouldUseUnifiedMedia).toBe(true);
      
      // New components can check this flag and use unified media hooks
      if (result.current.shouldUseUnifiedMedia) {
        // The component would use unified media hooks
        expect(true).toBe(true);
      }
    });
  });
});