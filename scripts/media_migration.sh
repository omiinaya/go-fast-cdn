#!/bin/bash

# Media Migration Script
# This script runs the media unification migration to merge images and docs tables into a single media table

echo "Starting media unification migration..."

# Check if rollback flag is provided
if [ "$1" = "--rollback" ]; then
    echo "Rolling back media unification migration..."
    go run cmd/media_migration/main.go --rollback
else
    echo "Running media unification migration..."
    go run cmd/media_migration/main.go
fi

# Check the exit status
if [ $? -eq 0 ]; then
    echo "Media migration completed successfully!"
else
    echo "Media migration failed!"
    exit 1
fi