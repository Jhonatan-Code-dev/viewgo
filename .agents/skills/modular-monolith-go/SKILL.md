---
name: modular-monolith-go
description: >-
  Use this skill when organizing Go applications into isolated feature modules (Modular Monolith pattern).
---

# Modular Monolith in Go

Guidelines for maintaining clean boundaries between modules in a single Go repository.

## Rules of Engagement

1. **Isolation**:
   - Each module inside `internal/modules/` is self-contained.
   - Never import implementation details from another module's `infrastructure` or `delivery` packages.

2. **Cross-Module Communication**:
   - Communicate strictly through public interfaces or domain models exposed via a module's `domain` or `application` export.
   - For asynchronous decoupling, use internal events or shared kernel contracts (`internal/shared`).

3. **Package Visibility**:
   - Leverage `internal/` to enforce compiler-level boundaries against outside packages.
