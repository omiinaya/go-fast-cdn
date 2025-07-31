#!/bin/bash

# Emergency Media Migration Rollback Script
# This script provides a quick reference for rolling back the media migration in an emergency situation.

echo "================================================"
echo "EMERGENCY MEDIA MIGRATION ROLLBACK SCRIPT"
echo "================================================"
echo ""

# Check if being run with root/sudo privileges
if [[ $EUID -eq 0 ]]; then
   echo "WARNING: This script is running as root. This is not recommended."
   echo "Please run as a regular user with appropriate permissions."
   echo ""
   read -p "Continue anyway? (y/N): " -n 1 -r
   echo ""
   if [[ ! $REPLY =~ ^[Yy]$ ]]; then
       exit 1
   fi
fi

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

echo "This script will guide you through the emergency rollback process."
echo "Please follow the instructions carefully."
echo ""

# Step 1: Stop the application
echo "STEP 1: STOP THE APPLICATION"
echo "============================"
echo "Before proceeding, you must stop any running instances of the application."
echo ""
echo "Common ways to stop the application:"
echo "- If running as a service: sudo systemctl stop go-fast-cdn"
echo "- If running in a terminal: Press Ctrl+C"
echo "- If running as a background process: pkill -f 'go run main.go'"
echo ""
read -p "Have you stopped the application? (y/N): " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Please stop the application before continuing."
    exit 1
fi

# Step 2: Choose rollback method
echo ""
echo "STEP 2: CHOOSE ROLLBACK METHOD"
echo "==============================="
echo "Select the rollback method to use:"
echo ""
echo "1. Built-in rollback functionality (RECOMMENDED)"
echo "   - Uses the migration script's built-in rollback"
echo "   - Fast and safe if the migration script is working"
echo ""
echo "2. Restore from backup"
echo "   - Restores the database from a pre-migration backup"
echo "   - Use if built-in rollback fails"
echo ""
echo "3. Exit this script"
echo ""
read -p "Enter your choice (1-3): " -n 1 -r
echo ""

case $REPLY in
    1)
        echo "You selected: Built-in rollback functionality"
        echo ""
        echo "Running built-in rollback..."
        echo ""
        
        # Run the rollback script
        if ./scripts/media_migration.sh --rollback; then
            echo ""
            echo "✓ Rollback completed successfully!"
        else
            echo ""
            echo "✗ Rollback failed!"
            echo ""
            echo "Please try restoring from backup instead."
            exit 1
        fi
        ;;
    2)
        echo "You selected: Restore from backup"
        echo ""
        
        # List available backups
        echo "Available backups:"
        echo "=================="
        if ! ./scripts/db_backup.sh list; then
            echo "Failed to list backups. Please check the backup system."
            exit 1
        fi
        
        echo ""
        read -p "Enter the full path to the pre-migration backup file: " backup_path
        
        if [ -z "$backup_path" ]; then
            echo "No backup path provided. Exiting."
            exit 1
        fi
        
        if [ ! -f "$backup_path" ]; then
            echo "Backup file not found: $backup_path"
            exit 1
        fi
        
        echo ""
        echo "WARNING: This will overwrite the current database with the backup."
        echo "All data added after the backup was created will be lost."
        echo ""
        read -p "Continue with restore? (y/N): " -n 1 -r
        echo ""
        
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            echo "Restoring from backup..."
            echo ""
            
            if ./scripts/db_backup.sh restore -backup "$backup_path" -force; then
                echo ""
                echo "✓ Restore completed successfully!"
            else
                echo ""
                echo "✗ Restore failed!"
                exit 1
            fi
        else
            echo "Restore cancelled."
            exit 0
        fi
        ;;
    3)
        echo "Exiting script."
        exit 0
        ;;
    *)
        echo "Invalid choice. Please run the script again."
        exit 1
        ;;
esac

# Step 3: Verify rollback
echo ""
echo "STEP 3: VERIFY ROLLBACK"
echo "========================"
echo "Verifying that the rollback was successful..."
echo ""

if ./scripts/verify_media_migration.sh; then
    echo ""
    echo "✓ Verification completed successfully!"
    echo "The rollback was successful and the system is in a consistent state."
else
    echo ""
    echo "✗ Verification failed!"
    echo ""
    echo "The rollback may not have been completed successfully."
    echo "Please check the verification output and take appropriate action."
    exit 1
fi

# Step 4: Restart the application
echo ""
echo "STEP 4: RESTART THE APPLICATION"
echo "==============================="
echo "The rollback has been completed and verified."
echo "You can now restart the application."
echo ""
echo "Common ways to start the application:"
echo "- As a service: sudo systemctl start go-fast-cdn"
echo "- In a terminal: go run main.go"
echo "- As a background process: nohup go run main.go > app.log 2>&1 &"
echo ""
read -p "Do you want to start the application now? (y/N): " -n 1 -r
echo ""

if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Starting the application..."
    echo ""
    
    # Try to start the application
    if go run main.go & then
        echo ""
        echo "✓ Application started successfully!"
        echo "Check the application logs to ensure it's running correctly."
    else
        echo ""
        echo "✗ Failed to start the application!"
        echo "Please start it manually using the appropriate command for your setup."
    fi
else
    echo "Please start the application manually when ready."
fi

# Completion message
echo ""
echo "================================================"
echo "EMERGENCY ROLLBACK COMPLETED"
echo "================================================"
echo ""
echo "The media migration rollback has been completed."
echo "The system has been restored to its pre-migration state."
echo ""
echo "Next steps:"
echo "1. Monitor the application for any issues"
echo "2. Investigate the cause of the rollback"
echo "3. Notify stakeholders of the rollback"
echo "4. Plan for re-migration after fixing the issues"
echo ""
echo "For more detailed information, see the rollback plan documentation:"
echo "docs/MEDIA_MIGRATION_ROLLBACK_PLAN.md"
echo ""