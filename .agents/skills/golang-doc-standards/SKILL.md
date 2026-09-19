---
name: golang-doc-standards
description: >-
  Use this skill when auditing or formatting Go documentation, ensuring staticcheck ST1000 compliance,
  package comments, exported symbol comments, and clean ASCII typography.
---

# Go Documentation & Staticcheck ST1000 Skill

This skill enforces official Go documentation standards and staticcheck ST1000 rules.

## Core Rules

1. **Package Comments (`ST1000`)**:
   - Every package must have a top-level package doc comment.
   - Format: `// Package <name> <description>.`

2. **Exported Symbol Comments**:
   - Every exported type, struct, function, interface, variable, or constant must be preceded by a doc comment starting with the symbol's name.

3. **No Emojis**:
   - Package comments must be strictly professional without emojis or decorative symbols.

## Verification Workflow
Run staticcheck or go vet to verify package comments:
```bash
go vet ./...
```
