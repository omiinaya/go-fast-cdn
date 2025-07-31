#!/bin/bash

echo "====================================="
echo "Media Migration Verification Script"
echo "====================================="

# Change to the project root directory
cd "$(dirname "$0")/.."

# Run the media migration verification
echo "Running media migration verification..."
go run cmd/verify_media_migration/main.go

# Check the exit code
if [ $? -eq 0 ]; then
    echo ""
    echo "Verification completed successfully!"
    exit 0
else
    echo ""
    echo "Verification failed!"
    exit 1
fi