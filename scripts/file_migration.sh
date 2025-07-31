#!/bin/bash

# File Migration Script
# This script migrates files from separate directories to a unified media directory

# Parse command line arguments
ROLLBACK=false
CLEANUP=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --rollback)
            ROLLBACK=true
            shift
            ;;
        --cleanup)
            CLEANUP=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--rollback] [--cleanup]"
            echo "  --rollback  : Rollback file migration to legacy directories"
            echo "  --cleanup   : Clean up legacy files after successful migration"
            exit 1
            ;;
    esac
done

# Build the file migration tool
echo "Building file migration tool..."
go build -o bin/file_migration cmd/file_migration/main.go
if [ $? -ne 0 ]; then
    echo "Failed to build file migration tool"
    exit 1
fi

# Run the file migration
if [ "$ROLLBACK" = true ]; then
    echo "Rolling back file migration..."
    ./bin/file_migration --rollback
elif [ "$CLEANUP" = true ]; then
    echo "Cleaning up legacy files..."
    ./bin/file_migration --cleanup
else
    echo "Running file migration to unified media directory..."
    ./bin/file_migration
fi

# Check if migration was successful
if [ $? -eq 0 ]; then
    echo "File migration operation completed successfully!"
else
    echo "File migration operation failed!"
    exit 1
fi