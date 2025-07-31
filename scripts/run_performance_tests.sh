#!/bin/bash

# Performance Test Runner Script for Unified Media Repository
# This script runs the comprehensive performance tests for the unified media repository

echo "=========================================="
echo "Unified Media Repository Performance Test Runner"
echo "=========================================="

# Set up variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
PERFORMANCE_DIR="$PROJECT_ROOT/cmd/performance_tests"
REPORTS_DIR="$PROJECT_ROOT/performance-reports"

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

# Check if the performance test directory exists
if [ ! -d "$PERFORMANCE_DIR" ]; then
    echo "❌ Performance test directory not found: $PERFORMANCE_DIR"
    exit 1
fi

# Change to the performance test directory
cd "$PERFORMANCE_DIR"

# Run the performance tests
print_section "Running Performance Tests"

echo "Running performance tests..."
go run main.go

# Check if the performance tests ran successfully
if [ $? -eq 0 ]; then
    echo "✅ Performance tests completed successfully"
else
    echo "❌ Performance tests failed"
    exit 1
fi

# Move the generated reports to the reports directory
print_section "Organizing Reports"

if [ -f "performance-report.md" ]; then
    mv performance-report.md "$REPORTS_DIR/"
    echo "✅ Performance report moved to: $REPORTS_DIR/performance-report.md"
fi

if [ -f "performance-results.json" ]; then
    mv performance-results.json "$REPORTS_DIR/"
    echo "✅ Performance results moved to: $REPORTS_DIR/performance-results.json"
fi

if [ -f "cpu.prof" ]; then
    mv cpu.prof "$REPORTS_DIR/"
    echo "✅ CPU profile moved to: $REPORTS_DIR/cpu.prof"
fi

# Display final status
print_section "Performance Testing Complete"

echo "🎉 Performance testing completed successfully!"
echo "📊 Reports are available in: $REPORTS_DIR/"
echo "📄 View the performance report: $REPORTS_DIR/performance-report.md"
echo "📈 View the detailed results: $REPORTS_DIR/performance-results.json"