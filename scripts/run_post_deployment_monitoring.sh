#!/bin/bash

# Post-Deployment Monitoring Runner Script for Unified Media Repository
# This script runs the comprehensive post-deployment monitoring for the unified media repository

echo "=========================================="
echo "Unified Media Repository Post-Deployment Monitoring Runner"
echo "=========================================="

# Set up variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
MONITORING_DIR="$PROJECT_ROOT/cmd/post_deployment_monitoring"
REPORTS_DIR="$PROJECT_ROOT/post-deployment-monitoring-reports"

# Create reports directory if it doesn't exist
mkdir -p "$REPORTS_DIR"

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

# Check if the monitoring directory exists
if [ ! -d "$MONITORING_DIR" ]; then
    echo "❌ Monitoring directory not found: $MONITORING_DIR"
    exit 1
fi

# Change to the monitoring directory
cd "$MONITORING_DIR"

# Run the monitoring tests
print_section "Running Post-Deployment Monitoring"

echo "Running post-deployment monitoring..."
go run main.go

# Check if the monitoring ran successfully
if [ $? -eq 0 ]; then
    echo "✅ Post-deployment monitoring completed successfully"
else
    echo "❌ Post-deployment monitoring failed"
    exit 1
fi

# Move the generated reports to the reports directory
print_section "Organizing Reports"

if [ -f "post-deployment-monitoring-report.md" ]; then
    mv post-deployment-monitoring-report.md "$REPORTS_DIR/"
    echo "✅ Monitoring report moved to: $REPORTS_DIR/post-deployment-monitoring-report.md"
fi

if [ -f "post-deployment-monitoring-results.json" ]; then
    mv post-deployment-monitoring-results.json "$REPORTS_DIR/"
    echo "✅ Monitoring results moved to: $REPORTS_DIR/post-deployment-monitoring-results.json"
fi

if [ -f "monitoring-cpu.prof" ]; then
    mv monitoring-cpu.prof "$REPORTS_DIR/"
    echo "✅ CPU profile moved to: $REPORTS_DIR/monitoring-cpu.prof"
fi

# Display final status
print_section "Post-Deployment Monitoring Complete"

echo "🎉 Post-deployment monitoring completed successfully!"
echo "📊 Reports are available in: $REPORTS_DIR/"
echo "📄 View the monitoring report: $REPORTS_DIR/post-deployment-monitoring-report.md"
echo "📈 View the detailed results: $REPORTS_DIR/post-deployment-monitoring-results.json"

# Optional: Open the report if the user is on a system with an open command
if command_exists open; then
    echo "📖 Opening monitoring report..."
    open "$REPORTS_DIR/post-deployment-monitoring-report.md"
elif command_exists xdg-open; then
    echo "📖 Opening monitoring report..."
    xdg-open "$REPORTS_DIR/post-deployment-monitoring-report.md"
fi