#!/bin/bash

# Media Migration Test Script
# This script creates a test database, populates it with sample data,
# runs the media migration, verifies it, tests rollback, and cleans up.

echo "Media Migration Test Suite"
echo "=========================="

# Check if cleanup-only flag is provided
if [ "$1" = "--cleanup-only" ]; then
    echo "Running cleanup only..."
    go run cmd/test_media_migration/main.go --cleanup-only
    exit $?
fi

# Run the test suite
go run cmd/test_media_migration/main.go

# Check the exit status
if [ $? -eq 0 ]; then
    echo ""
    echo "All tests passed successfully!"
    echo "The migration script is ready for staging environment."
else
    echo ""
    echo "Tests failed! Please check the output above for details."
    exit 1
fi