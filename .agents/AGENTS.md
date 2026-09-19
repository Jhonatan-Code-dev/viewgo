# Go Clean Architecture & Modular Monolith Development Guidelines

This project adheres strictly to **Clean Architecture** principles embedded within a **Modular Monolith** pattern in Go.

## Architectural Guidelines

1. **Modular Monolith Structure**:
   - Feature modules live isolated under `internal/modules/<module_name>/`.
   - Each module MUST contain isolated sub-packages representing the Clean Architecture layers:
     - `domain`: Pure business entities, value objects, domain events, and repository/provider interfaces.
     - `application`: Use cases, query handlers, command handlers, and DTOs.
     - `infrastructure`: Concrete implementations of data providers, repositories, external APIs, and persistence.
     - `delivery`: Transport handlers (HTTP/REST controllers, CLI handlers, gRPC, event listeners).
   - Cross-module communication MUST occur through explicit domain interfaces or shared kernel contracts (`internal/shared/`).

2. **Clean Code & Best Practices in Go**:
   - **No Hardcoded Static Datasets**: Global reference data (countries, timezones, currencies) MUST be derived from official standard libraries (`time/tzdata`, `golang.org/x/text`) or dynamic Unicode CLDR registries.
   - **Explicit Error Handling**: Always wrap errors using `fmt.Errorf("context: %w", err)` or custom domain errors.
   - **Dependency Injection**: Use constructor injection for dependencies (`NewUseCase(repo)`). Accept interfaces, return concrete structs.
   - **Concurrency Safety**: Domain repositories and providers must be safe for concurrent access (`sync.RWMutex` where needed).
   - **Testing**: Every module must have unit tests covering domain logic and infrastructure data extraction.

3. **Strict Professional Documentation Standard (No Emojis)**:
   - All documentation files (`README.md`, `doc.go`, markdown files, release notes) and log outputs MUST BE strictly professional and corporate.
   - **DO NOT USE EMOJIS OR GRAPHICAL ICONS** anywhere in documentation, README titles, headers, terminal logs, or code comments. Use clean standard markdown typography and structured ASCII formatting.
