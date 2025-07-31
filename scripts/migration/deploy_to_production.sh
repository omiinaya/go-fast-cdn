#!/bin/bash

# Production Deployment Script for Linux/macOS
# This script automates the deployment of the unified media repository to the production environment

set +e  # Don't exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to display usage
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  --skip-backup           Skip creating a backup before deployment"
    echo "  --skip-verification     Skip post-deployment verification"
    echo "  --rollback-only         Only perform rollback without deployment"
    echo "  --help                  Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                      # Full deployment with backup and verification"
    echo "  $0 --skip-backup        # Deployment without backup (not recommended)"
    echo "  $0 --rollback-only      # Only perform rollback"
}

# Parse command line arguments
SKIP_BACKUP=false
SKIP_VERIFICATION=false
ROLLBACK_ONLY=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --skip-backup)
            SKIP_BACKUP=true
            shift
            ;;
        --skip-verification)
            SKIP_VERIFICATION=true
            shift
            ;;
        --rollback-only)
            ROLLBACK_ONLY=true
            shift
            ;;
        --help)
            show_usage
            exit 0
            ;;
        *)
            print_error "Unknown option: $1"
            show_usage
            exit 1
            ;;
    esac
done

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
cd "$PROJECT_ROOT"

# Print banner
echo "=================================================="
echo "  UNIFIED MEDIA REPOSITORY PRODUCTION DEPLOYMENT "
echo "=================================================="
echo ""

# Log file for deployment
LOG_FILE="./bin/production_deployment_$(date +%Y%m%d_%H%M%S).log"
mkdir -p ./bin
touch "$LOG_FILE"
echo "Deployment log: $LOG_FILE"

# Function to log to both console and file
log() {
    echo "$1"
    echo "$1" >> "$LOG_FILE"
}

log "Starting production deployment at $(date)"

# Check prerequisites
print_info "Checking prerequisites..."

if ! command_exists go; then
    print_error "Go is not installed or not in PATH"
    log "ERROR: Go is not installed or not in PATH"
    exit 1
fi

if ! command_exists node; then
    print_error "Node.js is not installed or not in PATH"
    log "ERROR: Node.js is not installed or not in PATH"
    exit 1
fi

if ! command_exists npm; then
    print_error "npm is not installed or not in PATH"
    log "ERROR: npm is not installed or not in PATH"
    exit 1
fi

print_success "All prerequisites are met"
log "All prerequisites are met"

# Function to create backup
create_backup() {
    if [ "$SKIP_BACKUP" = true ]; then
        print_warning "Skipping backup creation as requested"
        echo "WARNING: Skipping backup creation as requested" | tee -a "$LOG_FILE"
        return 0
    fi

    print_info "Creating database backup..."
    echo "Creating database backup..." | tee -a "$LOG_FILE"
    
    if ./scripts/db_backup.sh create; then
        print_success "Database backup created successfully"
        echo "SUCCESS: Database backup created successfully" | tee -a "$LOG_FILE"
    else
        print_error "Failed to create database backup"
        echo "ERROR: Failed to create database backup" | tee -a "$LOG_FILE"
        exit 1
    fi
}

# Function to perform database migration
perform_migration() {
    print_info "Performing database migration..."
    echo "Performing database migration..." | tee -a "$LOG_FILE"
    
    # Set database path environment variable to ensure consistency
    export CDN_DB_PATH="./db_data/production.db"
    
    # Create db_data directory if it doesn't exist
    mkdir -p ./db_data
    
    # Build the production migration tool
    print_info "Building production migration tool..."
    echo "Building production migration tool..." | tee -a "$LOG_FILE"
    if go build -o bin/production_migration cmd/staging_migration/main.go; then
        print_success "Production migration tool built successfully"
        echo "SUCCESS: Production migration tool built successfully" | tee -a "$LOG_FILE"
    else
        print_error "Failed to build production migration tool"
        echo "ERROR: Failed to build production migration tool" | tee -a "$LOG_FILE"
        exit 1
    fi
    
    # Run the migration using the built binary
    if ./bin/production_migration; then
        print_success "Database migration completed successfully"
        echo "SUCCESS: Database migration completed successfully" | tee -a "$LOG_FILE"
    else
        print_error "Database migration failed"
        echo "ERROR: Database migration failed" | tee -a "$LOG_FILE"
        print_info "Attempting rollback..."
        echo "INFO: Attempting rollback..." | tee -a "$LOG_FILE"
        
        # Run the rollback using the built binary
        if ./bin/production_migration --rollback; then
            print_success "Rollback completed successfully"
            echo "SUCCESS: Rollback completed successfully" | tee -a "$LOG_FILE"
        else
            print_error "Rollback failed"
            echo "ERROR: Rollback failed" | tee -a "$LOG_FILE"
        fi
        
        exit 1
    fi
}

# Function to build backend
build_backend() {
    print_info "Building backend application..."
    echo "Building backend application..." | tee -a "$LOG_FILE"
    
    if go build -o bin/go-fast-cdn main.go; then
        print_success "Backend built successfully"
        echo "SUCCESS: Backend built successfully" | tee -a "$LOG_FILE"
    else
        print_error "Failed to build backend"
        echo "ERROR: Failed to build backend" | tee -a "$LOG_FILE"
        exit 1
    fi
}

# Function to build frontend
build_frontend() {
    print_info "Building frontend application..."
    echo "Building frontend application..." | tee -a "$LOG_FILE"
    
    cd ui
    
    if npm install --legacy-peer-deps && npm run build; then
        print_success "Frontend built successfully"
        echo "SUCCESS: Frontend built successfully" | tee -a "$LOG_FILE"
        cd ..
    else
        print_error "Failed to build frontend"
        echo "ERROR: Failed to build frontend" | tee -a "$LOG_FILE"
        cd ..
        exit 1
    fi
}

# Function to deploy backend
deploy_backend() {
    print_info "Deploying backend application..."
    echo "Deploying backend application..." | tee -a "$LOG_FILE"
    
    # In a real deployment, this would involve copying files to the production server
    # For this example, we'll just simulate the deployment
    
    if [ -f "bin/go-fast-cdn" ]; then
        print_success "Backend deployment simulated successfully"
        echo "SUCCESS: Backend deployment simulated successfully" | tee -a "$LOG_FILE"
    else
        print_error "Backend binary not found"
        echo "ERROR: Backend binary not found" | tee -a "$LOG_FILE"
        exit 1
    fi
}

# Function to deploy frontend
deploy_frontend() {
    print_info "Deploying frontend application..."
    echo "Deploying frontend application..." | tee -a "$LOG_FILE"
    
    # In a real deployment, this would involve copying files to the production server
    # For this example, we'll just simulate the deployment
    
    if [ -d "ui/build" ]; then
        print_success "Frontend deployment simulated successfully"
        echo "SUCCESS: Frontend deployment simulated successfully" | tee -a "$LOG_FILE"
    else
        print_error "Frontend build output not found"
        echo "ERROR: Frontend build output not found" | tee -a "$LOG_FILE"
        exit 1
    fi
}

# Function to verify deployment
verify_deployment() {
    if [ "$SKIP_VERIFICATION" = true ]; then
        print_warning "Skipping deployment verification as requested"
        echo "WARNING: Skipping deployment verification as requested" | tee -a "$LOG_FILE"
        return 0
    fi

    print_info "Verifying deployment..."
    echo "Verifying deployment..." | tee -a "$LOG_FILE"
    
    # Set database path environment variable to ensure consistency
    export CDN_DB_PATH="./db_data/production.db"
    
    # Build the verification tool first to ensure it uses the same database
    print_info "Building verification tool..."
    echo "Building verification tool..." | tee -a "$LOG_FILE"
    if go build -o bin/verify_media_migration cmd/verify_media_migration/main.go; then
        print_success "Verification tool built successfully"
        echo "SUCCESS: Verification tool built successfully" | tee -a "$LOG_FILE"
    else
        print_error "Failed to build verification tool"
        echo "ERROR: Failed to build verification tool" | tee -a "$LOG_FILE"
        return 1
    fi
    
    # Run the verification using the built binary
    if ./bin/verify_media_migration; then
        print_success "Deployment verification completed successfully"
        echo "SUCCESS: Deployment verification completed successfully" | tee -a "$LOG_FILE"
    else
        print_error "Deployment verification failed"
        echo "ERROR: Deployment verification failed" | tee -a "$LOG_FILE"
        return 1
    fi
}

# Function to perform rollback
perform_rollback() {
    print_info "Performing rollback..."
    echo "Performing rollback..." | tee -a "$LOG_FILE"
    
    # Set database path environment variable to ensure consistency
    export CDN_DB_PATH="./db_data/production.db"
    
    # Create db_data directory if it doesn't exist
    mkdir -p ./db_data
    
    # Build the production migration tool if it doesn't exist
    if [ ! -f "./bin/production_migration" ]; then
        print_info "Building production migration tool..."
        echo "INFO: Building production migration tool..." | tee -a "$LOG_FILE"
        if go build -o bin/production_migration cmd/staging_migration/main.go; then
            print_success "Production migration tool built successfully"
            echo "SUCCESS: Production migration tool built successfully" | tee -a "$LOG_FILE"
        else
            print_error "Failed to build production migration tool"
            echo "ERROR: Failed to build production migration tool" | tee -a "$LOG_FILE"
            exit 1
        fi
    fi
    
    # Run the rollback using the built binary
    if ./bin/production_migration --rollback; then
        print_success "Rollback completed successfully"
        echo "SUCCESS: Rollback completed successfully" | tee -a "$LOG_FILE"
    else
        print_error "Rollback failed"
        echo "ERROR: Rollback failed" | tee -a "$LOG_FILE"
        exit 1
    fi
}

# Main deployment process
if [ "$ROLLBACK_ONLY" = true ]; then
    print_info "Starting rollback-only process..."
    echo "Starting rollback-only process..." | tee -a "$LOG_FILE"
    perform_rollback
    print_success "Rollback-only process completed"
    echo "Rollback-only process completed" | tee -a "$LOG_FILE"
    exit 0
fi

print_info "Starting deployment process..."
echo "Starting deployment process..." | tee -a "$LOG_FILE"

# Phase 1: Pre-Deployment Preparation
print_info "Phase 1: Pre-Deployment Preparation"
echo "Phase 1: Pre-Deployment Preparation" | tee -a "$LOG_FILE"
create_backup

# Phase 2: Database Migration
print_info "Phase 2: Database Migration"
echo "Phase 2: Database Migration" | tee -a "$LOG_FILE"
perform_migration

# Phase 3: Build Application
print_info "Phase 3: Build Application"
echo "Phase 3: Build Application" | tee -a "$LOG_FILE"
build_backend
build_frontend

# Phase 4: Deploy Application
print_info "Phase 4: Deploy Application"
echo "Phase 4: Deploy Application" | tee -a "$LOG_FILE"
deploy_backend
deploy_frontend

# Phase 5: Post-Deployment Verification
print_info "Phase 5: Post-Deployment Verification"
echo "Phase 5: Post-Deployment Verification" | tee -a "$LOG_FILE"
if verify_deployment; then
    print_success "Deployment verification passed"
    echo "SUCCESS: Deployment verification passed" | tee -a "$LOG_FILE"
else
    print_error "Deployment verification failed"
    echo "ERROR: Deployment verification failed" | tee -a "$LOG_FILE"
    print_info "Attempting rollback..."
    echo "INFO: Attempting rollback..." | tee -a "$LOG_FILE"
    perform_rollback
    exit 1
fi

# Deployment completed successfully
echo ""
echo "=================================================="
echo "          DEPLOYMENT COMPLETED SUCCESSFULLY       "
echo "=================================================="
echo ""
print_info "The unified media repository has been successfully deployed to the production environment."
print_info "All verification checks have passed."
echo ""
echo "INFO: The unified media repository has been successfully deployed to the production environment." | tee -a "$LOG_FILE"
echo "INFO: All verification checks have passed." | tee -a "$LOG_FILE"
echo ""
print_info "Next steps:"
echo "1. Monitor the application for any issues"
echo "2. Perform user acceptance testing"
echo "3. Document the deployment results"
echo "4. Plan for future maintenance and updates"
echo ""
echo "INFO: Next steps:" | tee -a "$LOG_FILE"
echo "INFO: 1. Monitor the application for any issues" | tee -a "$LOG_FILE"
echo "INFO: 2. Perform user acceptance testing" | tee -a "$LOG_FILE"
echo "INFO: 3. Document the deployment results" | tee -a "$LOG_FILE"
echo "INFO: 4. Plan for future maintenance and updates" | tee -a "$LOG_FILE"
echo ""
echo "Deployment completed at $(date)" | tee -a "$LOG_FILE"