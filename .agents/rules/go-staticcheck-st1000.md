# Rule: Go Package Documentation Standards (staticcheck ST1000)

## Rules & Constraints

1. **Package Comments Required (ST1000 Compliance)**:
   - At least one source file in every Go package MUST contain a package-level comment starting with `// Package <pkgname> ...`.
   - The package comment must clearly explain the purpose, scope, and responsibilities of the package.
   - Example:
     ```go
     // Package infrastructure provides concrete IANA timezone loading and resolution.
     package infrastructure
     ```

2. **No Emojis in Package Comments**:
   - Package comments must be written in professional, concise, corporate technical prose without emojis or graphical icons.
