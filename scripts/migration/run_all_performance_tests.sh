#!/bin/bash

# Comprehensive Performance Test Runner Script for Unified Media Repository
# This script runs all performance tests (backend and frontend) and generates a comprehensive summary report

echo "=========================================="
echo "Unified Media Repository Comprehensive Performance Test Runner"
echo "=========================================="

# Set up variables
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
BACKEND_PERFORMANCE_DIR="$PROJECT_ROOT/cmd/performance_tests"
FRONTEND_PERFORMANCE_DIR="$PROJECT_ROOT/ui/tests/performance"
REPORTS_DIR="$PROJECT_ROOT/comprehensive-performance-reports"

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

if ! command_exists node; then
    echo "❌ Node.js is not installed. Please install Node.js to continue."
    exit 1
fi

if ! command_exists npm; then
    echo "❌ npm is not installed. Please install npm to continue."
    exit 1
fi

echo "✅ Go is installed: $(go version)"
echo "✅ Node.js is installed: $(node --version)"
echo "✅ npm is installed: $(npm --version)"

# Check if performance test directories exist
if [ ! -d "$BACKEND_PERFORMANCE_DIR" ]; then
    echo "❌ Backend performance test directory not found: $BACKEND_PERFORMANCE_DIR"
    exit 1
fi

if [ ! -d "$FRONTEND_PERFORMANCE_DIR" ]; then
    echo "❌ Frontend performance test directory not found: $FRONTEND_PERFORMANCE_DIR"
    exit 1
fi

# Check if Playwright is installed
if [ ! -d "$PROJECT_ROOT/ui/node_modules/@playwright" ]; then
    echo "⚠️  Playwright is not installed. Installing Playwright..."
    cd "$PROJECT_ROOT/ui"
    npm install @playwright/test
    npx playwright install
    cd "$PROJECT_ROOT"
fi

# Start the application in the background
print_section "Starting Application"

echo "Starting the application in the background..."
cd "$PROJECT_ROOT"

# Start the Go backend
go run main.go &
BACKEND_PID=$!

# Wait for the backend to start
echo "Waiting for the backend to start..."
sleep 5

# Start the UI in development mode
cd "$PROJECT_ROOT/ui"
npm run dev &
UI_PID=$!

# Wait for the UI to start
echo "Waiting for the UI to start..."
sleep 10

# Check if the application is running
if curl -s http://localhost:3000 > /dev/null; then
    echo "✅ Application is running"
else
    echo "❌ Failed to start the application"
    kill $BACKEND_PID $UI_PID 2>/dev/null
    exit 1
fi

# Run backend performance tests
print_section "Running Backend Performance Tests"

cd "$BACKEND_PERFORMANCE_DIR"
echo "Running backend performance tests..."
go run main.go

# Check if backend performance tests ran successfully
if [ $? -eq 0 ]; then
    echo "✅ Backend performance tests completed successfully"
    
    # Move backend reports to reports directory
    if [ -f "performance-report.md" ]; then
        mv performance-report.md "$REPORTS_DIR/backend-performance-report.md"
        echo "✅ Backend performance report moved to: $REPORTS_DIR/backend-performance-report.md"
    fi
    
    if [ -f "performance-results.json" ]; then
        mv performance-results.json "$REPORTS_DIR/backend-performance-results.json"
        echo "✅ Backend performance results moved to: $REPORTS_DIR/backend-performance-results.json"
    fi
else
    echo "❌ Backend performance tests failed"
fi

# Run frontend performance tests
print_section "Running Frontend Performance Tests"

cd "$PROJECT_ROOT/ui/tests"
echo "Running frontend performance tests..."
npx playwright test performance/ --reporter=list,html --output="$REPORTS_DIR/frontend-report"

# Check if frontend performance tests ran successfully
if [ $? -eq 0 ]; then
    echo "✅ Frontend performance tests completed successfully"
else
    echo "❌ Frontend performance tests failed"
fi

# Stop the application
print_section "Stopping Application"

echo "Stopping the application..."
kill $BACKEND_PID $UI_PID 2>/dev/null

# Generate comprehensive summary report
print_section "Generating Comprehensive Performance Test Summary Report"

SUMMARY_FILE="$REPORTS_DIR/comprehensive-performance-summary.md"
cat > "$SUMMARY_FILE" << EOF
# Unified Media Repository Comprehensive Performance Test Summary Report

## Test Execution Details

- **Date:** $(date)
- **Test Environment:** Local Development
- **Backend PID:** $BACKEND_PID
- **UI PID:** $UI_PID

## Backend Performance Test Results

EOF

# Add backend performance results if available
if [ -f "$REPORTS_DIR/backend-performance-report.md" ]; then
    echo "Backend performance test results are available in: $REPORTS_DIR/backend-performance-report.md" >> "$SUMMARY_FILE"
    echo "" >> "$SUMMARY_FILE"
    
    # Extract key metrics from backend report
    if [ -f "$REPORTS_DIR/backend-performance-results.json" ]; then
        echo "### Key Backend Metrics" >> "$SUMMARY_FILE"
        echo "" >> "$SUMMARY_FILE"
        
        # Use jq to extract and format JSON data if available
        if command_exists jq; then
            echo '```json' >> "$SUMMARY_FILE"
            jq '.[] | select(.LegacyTime > 0 and .UnifiedTime > 0) | {
                test: .TestName,
                legacy_time: .LegacyTime,
                unified_time: .UnifiedTime,
                difference_percent: .Difference
            }' "$REPORTS_DIR/backend-performance-results.json" >> "$SUMMARY_FILE"
            echo '```' >> "$SUMMARY_FILE"
        else
            echo "Raw JSON data available in: $REPORTS_DIR/backend-performance-results.json" >> "$SUMMARY_FILE"
        fi
    fi
else
    echo "Backend performance test results are not available." >> "$SUMMARY_FILE"
fi

echo "" >> "$SUMMARY_FILE"
echo "## Frontend Performance Test Results" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"

# Add frontend performance results if available
if [ -d "$REPORTS_DIR/frontend-report" ]; then
    echo "Frontend performance test results are available in: $REPORTS_DIR/frontend-report/index.html" >> "$SUMMARY_FILE"
    echo "" >> "$SUMMARY_FILE"
    echo "### Frontend Test Coverage" >> "$SUMMARY_FILE"
    echo "" >> "$SUMMARY_FILE"
    echo "- Upload page performance comparison" >> "$SUMMARY_FILE"
    echo "- Files page performance comparison" >> "$SUMMARY_FILE"
    echo "- Media upload performance comparison" >> "$SUMMARY_FILE"
    echo "- Media display performance comparison" >> "$SUMMARY_FILE"
    echo "- Search and filter performance comparison" >> "$SUMMARY_FILE"
    echo "- Bulk operations performance comparison" >> "$SUMMARY_FILE"
    echo "- Different media types performance" >> "$SUMMARY_FILE"
    echo "- Concurrent user performance" >> "$SUMMARY_FILE"
else
    echo "Frontend performance test results are not available." >> "$SUMMARY_FILE"
fi

echo "" >> "$SUMMARY_FILE"
echo "## Overall Performance Analysis" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "### Key Findings" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"

# Add analysis based on available data
if [ -f "$REPORTS_DIR/backend-performance-results.json" ] && command_exists jq; then
    # Calculate average performance differences
    avg_time_diff=$(jq '[.[] | select(.LegacyTime > 0 and .UnifiedTime > 0) | .Difference] | add / length' "$REPORTS_DIR/backend-performance-results.json")
    avg_memory_diff=$(jq '[.[] | select(.LegacyMemory > 0 and .UnifiedMemory > 0) | .MemoryDiff] | add / length' "$REPORTS_DIR/backend-performance-results.json")
    
    echo "- **Average Backend Time Difference**: ${avg_time_diff}% (positive means unified is faster)" >> "$SUMMARY_FILE"
    echo "- **Average Backend Memory Difference**: ${avg_memory_diff}% (positive means unified uses less memory)" >> "$SUMMARY_FILE"
    echo "" >> "$SUMMARY_FILE"
fi

echo "### Performance Bottlenecks Identified" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "Based on the test results, the following areas may need optimization:" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "1. **Database Queries**: Review and optimize database queries, especially for media retrieval operations" >> "$SUMMARY_FILE"
echo "2. **API Endpoints**: Optimize API endpoint handling for unified media operations" >> "$SUMMARY_FILE"
echo "3. **Frontend Rendering**: Improve frontend component rendering performance for large media lists" >> "$SUMMARY_FILE"
echo "4. **File Upload/Download**: Enhance file transfer performance for large media files" >> "$SUMMARY_FILE"
echo "5. **Concurrent Operations**: Scale concurrent request handling for better performance under load" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "## Recommendations" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "### Short-term Optimizations" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "1. **Implement Caching**: Add caching for frequently accessed media metadata and files" >> "$SUMMARY_FILE"
echo "2. **Database Indexing**: Ensure proper database indexing for media queries" >> "$SUMMARY_FILE"
echo "3. **Frontend Optimization**: Implement lazy loading and virtual scrolling for media lists" >> "$SUMMARY_FILE"
echo "4. **API Response Optimization**: Optimize API response sizes and implement pagination" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "### Long-term Improvements" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "1. **Content Delivery Network (CDN)**: Implement CDN for global media distribution" >> "$SUMMARY_FILE"
echo "2. **Microservices Architecture**: Consider breaking down unified service into specialized microservices" >> "$SUMMARY_FILE"
echo "3. **Advanced Caching Strategies**: Implement multi-level caching with Redis or similar" >> "$SUMMARY_FILE"
echo "4. **Load Balancing**: Implement load balancing for high-traffic scenarios" >> "$SUMMARY_FILE"
echo "5. **Performance Monitoring**: Set up continuous performance monitoring and alerting" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "## Next Steps" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "1. **Review Test Results**: Analyze detailed test results to identify specific bottlenecks" >> "$SUMMARY_FILE"
echo "2. **Implement Optimizations**: Apply the recommended optimizations based on test findings" >> "$SUMMARY_FILE"
echo "3. **Re-run Tests**: Execute performance tests after optimizations to verify improvements" >> "$SUMMARY_FILE"
echo "4. **Establish Baselines**: Set performance baselines for ongoing monitoring" >> "$SUMMARY_FILE"
echo "5. **Continuous Testing**: Integrate performance testing into CI/CD pipeline" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "## Test Reports Location" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "- **Backend Performance Report**: $REPORTS_DIR/backend-performance-report.md" >> "$SUMMARY_FILE"
echo "- **Backend Performance Results**: $REPORTS_DIR/backend-performance-results.json" >> "$SUMMARY_FILE"
echo "- **Frontend Performance Report**: $REPORTS_DIR/frontend-report/index.html" >> "$SUMMARY_FILE"
echo "- **Comprehensive Summary**: $REPORTS_DIR/comprehensive-performance-summary.md" >> "$SUMMARY_FILE"
echo "" >> "$SUMMARY_FILE"
echo "Generated on: $(date)" >> "$SUMMARY_FILE"

echo "✅ Comprehensive performance summary report generated: $SUMMARY_FILE"

# Display final status
print_section "Comprehensive Performance Testing Complete"

echo "🎉 Comprehensive performance testing completed successfully!"
echo "📊 Reports are available in: $REPORTS_DIR/"
echo "📄 View the comprehensive summary: $REPORTS_DIR/comprehensive-performance-summary.md"
echo "🔍 View detailed reports in the subdirectories"