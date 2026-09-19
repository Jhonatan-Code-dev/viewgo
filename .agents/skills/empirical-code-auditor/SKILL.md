---
name: empirical-code-auditor
description: >-
  Use this skill when performing an objective, rigorous code audit, review, or pre-flight verification.
  Enforces 100% honesty, zero assumptions, thorough empirical verification, and zero sycophancy.
---

# Empirical Code Auditor Skill

This skill provides a procedure for conducting objective, evidence-based code reviews and audits.

## Audit Workflow

### 1. Static Verification & File Inspection
- Read target source files completely. Never judge code quality based on partial line views.
- Check package imports, layer separation (Clean Architecture), struct alignment, and error wrapping (`fmt.Errorf("%w", err)`).

### 2. Empirical Execution Verification
- Execute compilation:
  ```bash
  go build -v -o viewgo.exe ./cmd/viewgo
  ```
- Execute full test suite:
  ```bash
  go test -v ./...
  ```
- Inspect command exit codes and output logs directly.

### 3. Truthful Reporting Protocol
- **State Facts**: Report exact test results, pass/fail counts, and execution times.
- **Identify Flaws**: Highlight missing unit tests, unhandled errors, memory alignment issues, or architectural violations without sugarcoating.
- **No Sycophancy**: If a user-proposed change introduces bugs or breaks clean architecture, explicitly report the issue with code line references.

## Verification Checklist
1. Did `go test ./...` exit with code 0?
2. Did `go build ./...` compile cleanly without warnings?
3. Are all exported symbols properly documented without emojis?
4. Are error values properly wrapped and returned?
