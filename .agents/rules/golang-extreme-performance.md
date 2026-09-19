# Rule: Go Extreme Performance & Zero Allocation Standards

## Rules & Directives

1. **Pre-allocated Slice Capacity**:
   - Every slice construction where the target size is known or bounded MUST use `make([]T, 0, capacity)` to eliminate dynamic array re-allocations.

2. **O(1) Hot Path Execution**:
   - Lookups (`GetCountryByCode`, `GetTimezoneByName`, `GetCurrencyByCode`) MUST achieve O(1) time complexity using internal hash maps with RLock concurrency protection.

3. **Benchmarking Mandatory**:
   - Benchmark tests (`Benchmark*`) MUST be maintained in `test/` to track nanoseconds per operation (`ns/op`), bytes per operation (`B/op`), and allocations per operation (`allocs/op`).

4. **Zero Allocation Targets**:
   - Hot path lookups should aim for 0 to 1 heap allocations per call.
