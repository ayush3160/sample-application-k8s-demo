# Load Testing Tool

This directory contains a standalone load testing tool for the E-commerce API.

## Build

```bash
go build -o load_test main.go
```

## Run

```bash
./load_test
```

## Configuration

Edit `main.go` and modify these constants:

```go
const (
    baseURL            = "http://localhost:8080"  // API endpoint
    totalRequests      = 1000                     // Total requests to make
    concurrentRequests = 50                       // Concurrent workers
)
```

## Features

- Tests all API endpoints randomly
- Generates realistic test data
- Concurrent request execution
- Performance metrics reporting:
  - Success/failure rates
  - Average latency
  - Requests per second
  - Total duration

## From Root Directory

You can also use Make commands from the project root:

```bash
# From /home/ayush/sample-application
make load-test        # Build and run
make load-test-build  # Just build
```
