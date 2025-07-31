#!/bin/bash

# Monitoring Dashboard Runner Script for Unified Media Repository
# This script starts the monitoring dashboard for visualizing post-deployment monitoring data

echo "=========================================="
echo "Unified Media Repository Monitoring Dashboard Runner"
echo "=========================================="

# Set up variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DASHBOARD_DIR="$PROJECT_ROOT/cmd/monitoring_dashboard"
REPORTS_DIR="$PROJECT_ROOT/post-deployment-monitoring-reports"

# Function to print section headers
print_section() {
    echo ""
    echo "=========================================="
    echo "$1"
    echo "=========================================="
}

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
print_section "Checking Prerequisites"

if ! command_exists go; then
    echo "❌ Go is not installed. Please install Go to continue."
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Check if the dashboard directory exists
if [ ! -d "$DASHBOARD_DIR" ]; then
    echo "❌ Dashboard directory not found: $DASHBOARD_DIR"
    exit 1
fi

# Create reports directory if it doesn't exist
mkdir -p "$REPORTS_DIR"

# Check if monitoring data exists
if [ ! -f "$REPORTS_DIR/post-deployment-monitoring-results.json" ]; then
    echo "⚠️  Monitoring data not found. The dashboard will show sample data."
    echo "   Run the post-deployment monitoring first to generate real data:"
    echo "   ./scripts/run_post_deployment_monitoring.sh"
    echo ""
fi

# Change to the dashboard directory
cd "$DASHBOARD_DIR"

# Start the dashboard server
print_section "Starting Monitoring Dashboard"

echo "Starting monitoring dashboard server..."
echo "Monitoring reports directory: $REPORTS_DIR"
echo ""

# Set environment variables for the dashboard
export MONITORING_REPORTS_DIR="$REPORTS_DIR"
export MONITORING_DASHBOARD_PORT=":8080"

# Start the dashboard in the background
go run main.go &
DASHBOARD_PID=$!

# Wait for the dashboard to start
echo "Waiting for the dashboard to start..."
sleep 3

# Check if the dashboard is running
if curl -s http://localhost:8080 > /dev/null; then
    echo "✅ Dashboard is running"
    echo ""
    echo "🌐 Access the dashboard at: http://localhost:8080"
    echo ""
    echo "Press Ctrl+C to stop the dashboard server"
    echo ""
    
    # Open the dashboard in the default browser if possible
    if command_exists open; then
        echo "📖 Opening dashboard in default browser..."
        open http://localhost:8080
    elif command_exists xdg-open; then
        echo "📖 Opening dashboard in default browser..."
        xdg-open http://localhost:8080
    fi
    
    # Wait for user interrupt
    trap 'echo ""; echo "Stopping dashboard server..."; kill $DASHBOARD_PID 2>/dev/null; exit 0' INT
    while true; do
        sleep 1
    done
else
    echo "❌ Failed to start the dashboard"
    kill $DASHBOARD_PID 2>/dev/null
    exit 1
fi