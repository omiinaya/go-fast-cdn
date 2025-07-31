#!/bin/bash

# Database Backup Script
# This script provides a convenient way to run the database backup tool

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Build the backup tool if it doesn't exist or if source files are newer
if [ ! -f "$PROJECT_ROOT/bin/db_backup" ] || [ "$PROJECT_ROOT/cmd/db_backup/main.go" -nt "$PROJECT_ROOT/bin/db_backup" ]; then
    echo "Building database backup tool..."
    cd "$PROJECT_ROOT"
    go build -o bin/db_backup cmd/db_backup/main.go
    if [ $? -ne 0 ]; then
        echo "Error: Failed to build database backup tool"
        exit 1
    fi
    echo "Build completed successfully."
fi

# Run the backup tool with all provided arguments
echo "Running database backup tool..."
"$PROJECT_ROOT/bin/db_backup" "$@"