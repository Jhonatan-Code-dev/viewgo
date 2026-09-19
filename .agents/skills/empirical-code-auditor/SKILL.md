---
name: empirical-code-auditor
description: >-
  Use this skill when performing an objective, rigorous code audit, review, or pre-flight verification.
  Enforces 100% honesty, zero assumptions, exhaustive recursive directory scanning, and zero sycophancy.
---

# Empirical Code Auditor Skill

This skill provides an uncompromised procedure for conducting objective, evidence-based code reviews across all workspace directories and subdirectories.

## Exhaustive Audit Workflow ("Revisa Todo")

### 1. Recursive Directory & Subdirectory Scanning
- Traverse every folder in the workspace recursively (`cmd/`, `internal/`, `pkg/`, `.agents/`, root files).
- Open and inspect every Go source file, Markdown documentation, and configuration file.
- Verify zero hardcoded static datasets, zero static fallback maps/switches, zero emojis/icons, and proper struct alignment.

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
- **Identify Flaws**: Highlight missing unit tests, unhandled errors, memory alignment issues, or static data fallbacks without sugarcoating.
- **No Sycophancy**: If code contains defects, report them explicitly with file paths and line numbers.

## Verification Checklist
1. Were all directories and subdirectories scanned completely?
2. Did `go test ./...` exit with code 0?
3. Did `go build ./...` compile cleanly without warnings?
4. Are all exported symbols properly documented without emojis?
5. Are error values properly wrapped and returned?
