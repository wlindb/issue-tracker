## Testing Conventions

Test names follow `Test_<UnitOfWork>_<StateUnderTest>_<ExpectedBehaviour>`:
- `UnitOfWork`: method or unit being tested (e.g. `Create`, `RowToProject`)
- `StateUnderTest`: inputs or conditions (e.g. `NilDescription`, `DuplicateID`)
- `ExpectedBehaviour`: the result (e.g. `SuccessfulProjectCreation`, `ReturnsError`)

Example: `Test_Create_ValidDescription_SuccessfulProjectCreation`

Variable naming in assertions:
- `expected` — the value you set up / expect
- `actual` — the value returned by the code under test

## Variable Naming

Use full, unabbreviated names for all variables and fields. Never use abbreviations:
- `repository` not `repo`
- `service` not `svc`
- `query` not `q`
- `description` not `desc`
- `queries` not `q` (for sqlc Queries struct fields)

This applies to local variables, function parameters, and struct fields. Go receiver names (single character or short abbreviation) follow standard Go conventions and are exempt.
