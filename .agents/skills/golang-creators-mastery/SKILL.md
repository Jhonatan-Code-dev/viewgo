---
name: golang-creators-mastery
description: >-
  Elite engineering guidelines distilled from Go's creators (Rob Pike, Ken Thompson, Robert Griesemer, Russ Cox)
  and top Go core maintainers (Dave Cheney, William Kennedy). Use when designing, refactoring, or optimizing Go code.
---

# Go Creators & Core Experts Engineering Mastery

This skill embodies the core software engineering philosophy of Go's original creators and lead architects.

## Core Engineering Proverbs (Rob Pike, Russ Cox, Dave Cheney)

1. **Clear is better than clever**:
   - Write code that is straightforward, readable, and explicit. Avoid unnecessary abstractions, reflection, or magic.

2. **Don't communicate by sharing memory; share memory by communicating**:
   - Use channels for orchestration and signal passing across goroutines. Use `sync.Mutex` for protecting isolated shared state.

3. **Errors are values**:
   - Treat errors as first-class domain values. Inspect, wrap (`fmt.Errorf("context: %w", err)`), and handle errors explicitly at the call site.

4. **Accept interfaces, return concrete structs**:
   - Define interfaces in the consumer package where they are used, keeping them small (1 to 3 methods). Return concrete pointers or structs from constructor functions.

5. **Make the zero value useful**:
   - Design types so their zero value (`sync.Mutex`, `bytes.Buffer`, `context.Background()`) is immediately usable without elaborate setup.

6. **Concurrency is not parallelism**:
   - Structure program logic into independent concurrent units. Let the Go runtime handle CPU parallel execution.

## Memory Efficiency & Concurrency Safety

- **Struct Alignment**: Order fields in structs from largest to smallest memory size (8 bytes -> 4 bytes -> 2 bytes -> 1 byte) to minimize memory padding.
- **Race Condition Prevention**: Always test concurrent packages using `go test -race ./...`.
- **Context Cancellation**: Ensure every goroutine listening to external events or channels respects `ctx.Done()`.

## Verification Checklist

1. Run formatting and lint checks: `go fmt ./...` and `golangci-lint run`.
2. Run data race analysis: `go test -v -race ./...`.
3. Verify zero allocations for hot loops where applicable.
