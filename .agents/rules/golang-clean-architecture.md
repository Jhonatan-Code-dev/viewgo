# Rule: Go Clean Architecture & Quality Standards

## Rules & Constraints

1. **Dependency Rule**:
   - Source code dependencies can only point inward.
   - Domain layer imports NOTHING from application, infrastructure, or delivery.
   - Application layer imports ONLY domain layer.
   - Infrastructure and Delivery layers import application and domain layers.

2. **Official Datasets & No Hardcoded Magic Data**:
   - Reference data for countries, timezones, and currencies MUST NEVER be hardcoded as static Go slices or maps in business logic.
   - Use official runtime registry providers (`golang.org/x/text/language`, `golang.org/x/text/currency`, `time/tzdata`).

3. **Interface Segregation & Small Interfaces**:
   - Define small, targeted interfaces in the domain or consumer package.
   - Example: `type CountryProvider interface { ListCountries(ctx context.Context) ([]Country, error) }`

4. **Context & Cancellation**:
   - All use cases, repositories, and network handlers MUST take `context.Context` as their first parameter.

5. **Constructor Functions**:
   - Instantiate structs using `New<StructName>` functions returning a pointer or interface.
