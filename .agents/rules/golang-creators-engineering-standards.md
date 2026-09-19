# Rule: Go Creators & Core Engineering Standards

## Principles & Rules

1. **Simplicity and Explicitness**:
   - Code must be explicit and self-documenting. No hidden control flow, implicit type casting, or unhandled panics.

2. **Error Wrapping & Inspection**:
   - Always contextualize errors using `fmt.Errorf("operation failed: %w", err)`.
   - Never swallow errors silently or return empty default fallback values when a call fails.

3. **Small Interfaces (Consumer Defined)**:
   - Interfaces belong to the code that consumes them, not the code that implements them.
   - Keep interfaces small (ideally 1 to 2 methods like `io.Reader`, `io.Writer`, `CountryProvider`).

4. **Goroutine Lifecycle & Context Propagation**:
   - Never leak goroutines. Every spawned goroutine must have a deterministic termination path via `context.Context` or channel closure.

5. **No Static Hardcoded Reference Datasets**:
   - Global reference data (countries, timezones, currencies) MUST be loaded dynamically from official Go standard libraries (`time/tzdata`, `golang.org/x/text`) or dynamic Unicode CLDR registries.

6. **Strict Professional Documentation (No Emojis)**:
   - Documentation (`README.md`, `doc.go`, code comments) and console output MUST BE strictly professional and corporate without emojis or decorative icons.
