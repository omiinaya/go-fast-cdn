# Performance Testing for Unified Media Repository

This document describes the performance testing approach for the unified media repository, which is the third step in Phase 5 of our unification project.

## Overview

The performance testing suite compares the performance of the unified media repository with the legacy separate image and document repositories. It measures various metrics including response times, memory usage, and throughput to ensure that the unified implementation performs at least as well as the legacy implementations.

## Test Categories

### 1. Database Query Performance Tests

These tests compare the performance of database queries between the unified and legacy implementations:

- **Get All Images vs Get All Media (type=image)**: Compares retrieving all images using the legacy image handler vs retrieving all media with type=image using the unified media handler.
- **Get All Docs vs Get All Media (type=document)**: Compares retrieving all documents using the legacy document handler vs retrieving all media with type=document using the unified media handler.
- **Get Image by Checksum vs Get Media by Checksum**: Compares retrieving images by checksum using the legacy handler vs retrieving media by checksum using the unified handler.
- **Get Doc by Checksum vs Get Media by Checksum**: Compares retrieving documents by checksum using the legacy handler vs retrieving media by checksum using the unified handler.

### 2. API Endpoint Performance Tests

These tests compare the performance of API endpoints between the unified and legacy implementations:

- **GET /api/cdn/images vs GET /api/cdn/media?type=image**: Compares the performance of retrieving all images via the legacy endpoint vs the unified endpoint.
- **GET /api/cdn/docs vs GET /api/cdn/media?type=document**: Compares the performance of retrieving all documents via the legacy endpoint vs the unified endpoint.
- **POST /api/cdn/upload/image vs POST /api/cdn/upload/media (image)**: Compares the performance of uploading images via the legacy endpoint vs the unified endpoint.
- **POST /api/cdn/upload/doc vs POST /api/cdn/upload/media (document)**: Compares the performance of uploading documents via the legacy endpoint vs the unified endpoint.

### 3. Concurrent Request Performance Tests

These tests measure the throughput of the system under concurrent load:

- **Concurrent GET Requests**: Measures the throughput of concurrent GET requests to both legacy and unified endpoints.
- **Concurrent POST Requests**: Measures the throughput of concurrent POST requests (uploads) to both legacy and unified endpoints.

### 4. File Size Performance Tests

These tests measure performance with different file sizes:

- **Small Image (100x100)**: Tests performance with small image files.
- **Medium Image (500x500)**: Tests performance with medium-sized image files.
- **Large Image (1000x1000)**: Tests performance with large image files.

## Performance Metrics

The performance tests measure the following metrics:

### Time Metrics
- **Execution Time**: The time taken to complete operations.
- **Difference Percentage**: The percentage difference between legacy and unified implementation times.

### Memory Metrics
- **Memory Usage**: The amount of memory used during operations.
- **Memory Difference Percentage**: The percentage difference in memory usage between legacy and unified implementations.

### Throughput Metrics
- **Requests per Second**: The number of requests processed per second.
- **Throughput Difference Percentage**: The percentage difference in throughput between legacy and unified implementations.

## Running the Performance Tests

### Prerequisites

Before running the performance tests, ensure you have the following installed:

- Go (version 1.16 or higher)
- Git (to clone the repository)

### Running Tests Manually

1. Navigate to the performance test directory:
   ```bash
   cd cmd/performance_tests
   ```

2. Run the performance tests:
   ```bash
   go run main.go
   ```

3. The tests will generate two files:
   - `performance-report.md`: A detailed markdown report with analysis and recommendations.
   - `performance-results.json`: Raw test results in JSON format.
   - `cpu.prof`: CPU profiling data.

### Running Tests with the Test Runner Script

For convenience, you can use the provided test runner script:

#### On Linux/macOS:
```bash
cd /path/to/project/root
chmod +x scripts/run_performance_tests.sh
./scripts/run_performance_tests.sh
```

#### On Windows:
```cmd
cd /path/to/project/root
scripts\run_performance_tests.bat
```

The test runner script will:
1. Check and install prerequisites
2. Run the performance tests
3. Generate reports
4. Organize the reports in the `performance-reports` directory

## Understanding the Results

### Performance Report

The performance report (`performance-report.md`) includes:

1. **Test Results Table**: A table showing all test results with metrics for both legacy and unified implementations.
2. **Analysis**: Average performance differences across all test categories.
3. **Recommendations**: Specific recommendations based on the test results, including optimization suggestions.

### Interpreting the Metrics

#### Time Difference
- **Positive Difference**: Unified implementation is faster (better performance)
- **Negative Difference**: Legacy implementation is faster (unified needs optimization)

#### Memory Difference
- **Positive Difference**: Unified implementation uses less memory (better efficiency)
- **Negative Difference**: Legacy implementation uses less memory (unified needs optimization)

#### Throughput Difference
- **Positive Difference**: Legacy implementation has higher throughput (unified needs optimization)
- **Negative Difference**: Unified implementation has higher throughput (better performance)

## Performance Test Implementation

### Test Suite Structure

The performance test suite is implemented in `cmd/performance_tests/main.go` and consists of:

1. **PerformanceTestSuite**: A struct that manages the test environment and results.
2. **Benchmark Methods**: Methods for different test categories (database, API, concurrent, file size).
3. **Helper Methods**: Methods for measuring performance and generating reports.

### Key Components

#### PerformanceTestResult
```go
type PerformanceTestResult struct {
    TestName          string        `json:"test_name"`
    LegacyTime        time.Duration `json:"legacy_time"`
    UnifiedTime       time.Duration `json:"unified_time"`
    Difference        float64       `json:"difference_percent"`
    LegacyMemory      uint64        `json:"legacy_memory_bytes"`
    UnifiedMemory     uint64        `json:"unified_memory_bytes"`
    MemoryDiff        float64       `json:"memory_difference_percent"`
    LegacyThroughput  float64       `json:"legacy_throughput"`
    UnifiedThroughput float64       `json:"unified_throughput"`
    ThroughputDiff    float64       `json:"throughput_difference_percent"`
}
```

#### Performance Measurement
The tests use the following approach to measure performance:

1. **Warm-up**: Run the function multiple times to allow for JIT compilation and caching.
2. **Time Measurement**: Execute the function multiple times and measure the total time.
3. **Memory Measurement**: Use Go's runtime package to measure memory allocation.
4. **Throughput Measurement**: Calculate requests per second based on execution time.

## Troubleshooting

### Common Issues

1. **Compilation Errors**
   - Ensure all dependencies are properly installed.
   - Check that the Go version is compatible (1.16 or higher).
   - Verify that all import paths are correct.

2. **Test Failures**
   - Check that the database is properly initialized.
   - Ensure that temporary directories can be created.
   - Verify that all handler methods are accessible.

3. **Performance Degradation**
   - Check for resource contention (CPU, memory, disk I/O).
   - Ensure that the system is not running other resource-intensive processes.
   - Verify that the database is properly indexed.

### Debugging Tips

1. **Enable CPU Profiling**: The tests automatically generate a CPU profile (`cpu.prof`) that can be analyzed with:
   ```bash
   go tool pprof cpu.prof
   ```

2. **Run Tests Individually**: Comment out specific test categories to isolate issues.

3. **Increase Logging**: Add additional log statements to track test execution.

## Performance Optimization Recommendations

Based on the test results, consider the following optimizations:

### Database Optimizations
1. **Indexing**: Ensure that frequently queried fields are properly indexed.
2. **Query Optimization**: Review and optimize complex database queries.
3. **Connection Pooling**: Implement or tune database connection pooling.

### API Optimizations
1. **Caching**: Implement caching for frequently accessed data.
2. **Pagination**: Add pagination to large result sets.
3. **Compression**: Enable response compression for large payloads.

### System Optimizations
1. **Memory Management**: Review memory allocation patterns and reduce allocations.
2. **Concurrency**: Optimize concurrent request handling.
3. **Resource Limits**: Adjust system resource limits based on load testing results.

## Continuous Performance Monitoring

To ensure ongoing performance optimization:

1. **CI/CD Integration**: Include performance tests in the CI/CD pipeline.
2. **Performance Baselines**: Establish performance baselines and set thresholds for alerts.
3. **Regular Testing**: Schedule regular performance testing to catch regressions early.
4. **Production Monitoring**: Implement production monitoring to track real-world performance.

## Future Enhancements

Potential improvements to the performance testing suite:

1. **Additional Media Types**: Extend tests to include video and audio files.
2. **Load Testing**: Implement more sophisticated load testing with varying patterns.
3. **Distributed Testing**: Support for distributed performance testing across multiple machines.
4. **Historical Comparison**: Track performance changes over time and compare with historical data.
5. **Integration with Monitoring Tools**: Export results to monitoring and alerting systems.

## Conclusion

The performance testing suite provides comprehensive insights into the performance characteristics of the unified media repository compared to the legacy implementations. By regularly running these tests and analyzing the results, we can ensure that the unified repository maintains or improves upon the performance of the legacy systems while providing the benefits of a unified architecture.