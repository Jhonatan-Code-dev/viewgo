---
name: golang-extreme-performance
description: >-
  Use this skill when auditing, benchmarking, or optimizing Go code for extreme throughput,
  sub-microsecond latencies, zero heap allocations, and optimal memory layout.
---

# Go Extreme Performance & Optimization Skill

This skill provides guidelines and procedures for achieving maximum performance, zero memory allocations, and sub-10 nanosecond latency in Go.

## High-Performance Engineering Principles

1. **Compact `uint16` Direct Array Indexing (`O(1)` Sub-10ns CPU Latency)**:
   - For fixed alphabet domains (e.g. ISO 3166-1 alpha-2, ISO 4217 alpha-3), encode characters directly to an integer index `(c0-'A')*26 + (c1-'A')`.
   - Store 16-bit indices (`uint16`) pointing to contiguous slices instead of 64-bit pointers (`*T`), reducing lookup table RAM footprint by 75%:
     ```go
     // 676 uint16 values = 1.35 KB RAM (vs 5.4 KB for pointers)
     alpha2Table [676]uint16
     ```

2. **Pre-allocate Slice Capacity**:
   - Always initialize slices with known capacity to eliminate dynamic array reallocations:
     ```go
     result := make([]T, 0, count)
     ```

3. **Zero Lock Contention on Read-Only Tables**:
   - populated immutable lookup tables (like `alpha2Table` and `codeTable`) require no mutex locks during concurrent reads, allowing CPU L1/L2 cache line hits at ~8-10 ns/op.

4. **Zero-Allocation Hot Paths**:
   - Avoid `fmt.Sprintf` or dynamic string string concats in hot lookup loops. Use stack allocation and direct array indexing.

5. **Struct Field Alignment (Cache Line Locality)**:
   - Align struct fields from largest size (16-byte strings/pointers) to smallest size (1-byte bools) to prevent memory padding and maximize CPU cache line efficiency.

6. **Benchmark Verification Protocol**:
   - Measure performance metrics using Go benchmark testing:
     ```bash
     go test -bench=. -benchmem ./test/...
     ```
   - Target metrics: `< 15 ns/op`, `0 B/op`, `0 allocs/op` on lookup hot paths.

## Verification Workflow
Run benchmarks after any performance optimization:
```bash
go test -bench=. -benchmem ./test
```
