---
name: golang-extreme-performance
description: >-
  Use this skill when auditing, benchmarking, or optimizing Go code for extreme throughput,
  sub-microsecond latencies, zero heap allocations, and optimal memory layout.
---

# Go Extreme Performance & Optimization Skill

This skill provides guidelines and procedures for achieving maximum performance, zero memory allocations, and sub-microsecond latency in Go.

## High-Performance Engineering Principles

1. **Pre-allocate Slice Capacity**:
   - Always initialize slices with known capacity to eliminate re-allocation overhead:
     ```go
     result := make([]T, 0, count)
     ```

2. **O(1) Map Lookups with Read Mutexing**:
   - Use `sync.RWMutex` with `RLock()` / `RUnlock()` for read-heavy workloads to allow concurrent goroutine reads without blocking.

3. **Zero-Allocation Hot Paths**:
   - Avoid `fmt.Sprintf` in hot lookup loops. Use string builders, bytes buffers, or direct slice indexing.

4. **Struct Field Alignment (Cache Line Locality)**:
   - Align struct fields from largest size (16-byte strings/pointers) to smallest size (1-byte bools) to prevent memory padding and maximize CPU cache line efficiency.

5. **Benchmark Verification Protocol**:
   - Measure performance metrics using Go benchmark testing:
     ```bash
     go test -bench=. -benchmem ./test/...
     ```
   - Target metrics: `< 500 ns/op`, `< 64 B/op`, `0-1 allocs/op` on lookup hot paths.

## Verification Workflow
Run benchmarks after any performance optimization:
```bash
go test -bench=. -benchmem ./test
```
