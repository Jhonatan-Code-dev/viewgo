---
name: golang-clean-architecture
description: >-
  Use this skill when scaffolding, refactoring, or reviewing Go projects using Clean Architecture principles.
  Defines the standard 4-layer structure: Domain, Application, Infrastructure, and Delivery.
---

# Go Clean Architecture Skill

This skill guides the design and implementation of Go features using Clean Architecture.

## Layer Definitions & Responsibilities

### 1. Domain Layer (`internal/modules/<module>/domain`)
- **Entities**: Business models with private fields and validation methods.
- **Value Objects**: Immutable data structures representing domain attributes.
- **Repository/Provider Interfaces**: Abstractions required by the business domain.
- **Domain Errors**: Sentinel errors (`ErrNotFound`, `ErrInvalidCode`).

### 2. Application Layer (`internal/modules/<module>/application`)
- **Use Cases**: Encapsulate specific business workflows or queries.
- **DTOs**: Input and output data structures for boundaries.
- **Services**: Orchestrate operations across multiple domain entities.

### 3. Infrastructure Layer (`internal/modules/<module>/infrastructure`)
- **Providers / Repositories**: Implementations of domain interfaces (e.g. standard library tzdata, Unicode CLDR, SQL/Mongo).

### 4. Delivery Layer (`internal/modules/<module>/delivery`)
- **HTTP / REST / CLI / gRPC**: Transport controllers handling requests, validation, and JSON serialization.

## Verification Workflow
Run tests after any architectural change:
```bash
go test -v ./internal/modules/...
```
