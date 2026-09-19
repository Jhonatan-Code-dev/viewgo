# Go Clean Architecture & Modular Monolith Guidelines

This project adheres strictly to **Clean Architecture** principles embedded within a **Modular Monolith** pattern in Go, inspired by the core software engineering philosophy of Go's creators (Rob Pike, Ken Thompson, Robert Griesemer, Russ Cox).

## Core Principles & Uncompromising Truth

1. **Absolute Technical Honesty & Zero Sycophancy**:
   - Always state objective technical truth based on empirical evidence.
   - Never flatter the user or agree blindly ("no me de la razón"). If code has defects or anti-patterns, report it explicitly.
   - Never invent or assume file paths, package schemas, or test results ("no invente"). Always verify using file viewing and command execution.
   - Never skip audit steps ("nunca salte").

2. **Core Go Engineering Philosophy**:
   - **Clear is better than clever**: Code must be explicit, maintainable, and self-documenting. Avoid unnecessary reflection, magic, or deep inheritance hierarchies.
   - **Accept interfaces, return concrete structs**: Define small consumer-centric interfaces (1-3 methods) and return concrete structs from constructors (`NewProvider()`).
   - **Errors are values**: Treat errors as first-class domain values. Wrap errors using `fmt.Errorf("context: %w", err)` and handle them explicitly.
   - **Share memory by communicating**: Use channels for goroutine signaling and `sync.RWMutex` for protecting shared state.

## Architectural Guidelines

1. **Modular Monolith Structure**:
   - Feature modules live isolated under `internal/modules/<module_name>/`.
   - Public SDK packages exported for external consumers live under `pkg/<module_name>/`.
   - Executables live under `cmd/<binary_name>/main.go`.
   - Each module MUST contain isolated sub-packages representing Clean Architecture layers:
     - `domain`: Pure business entities, value objects, domain events, and repository/provider interfaces.
     - `application`: Use cases, query handlers, command handlers, and DTOs.
     - `infrastructure`: Concrete implementations of data providers, repositories, external APIs, and persistence.
     - `delivery`: Transport handlers (HTTP/REST controllers, CLI handlers, gRPC, event listeners).
   - Cross-module communication MUST occur through explicit domain interfaces or shared kernel contracts (`internal/shared/`).

2. **Clean Code & Best Practices in Go**:
   - **No Hardcoded Static Datasets**: Global reference data (countries, timezones, currencies) MUST be derived from official standard libraries (`time/tzdata`, `golang.org/x/text`) or dynamic Unicode CLDR registries.
   - **Explicit Error Handling**: Always wrap errors using `fmt.Errorf("context: %w", err)` or custom domain errors.
   - **Dependency Injection**: Use constructor injection for dependencies (`NewUseCase(repo)`).
   - **Concurrency Safety**: Domain repositories and providers must be safe for concurrent access (`sync.RWMutex` where needed).
   - **Testing**: Every module must have unit tests covering domain logic and infrastructure data extraction (`go test -v ./...`).

3. **Strict Professional Documentation Standard (No Emojis)**:
   - All documentation files (`README.md`, `doc.go`, markdown files, release notes) and log outputs MUST BE strictly professional and corporate.
   - **DO NOT USE EMOJIS OR GRAPHICAL ICONS** anywhere in documentation, README titles, headers, terminal logs, or code comments. Use clean standard markdown typography and structured ASCII formatting.
