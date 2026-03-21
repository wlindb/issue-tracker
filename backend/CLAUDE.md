## Testing Conventions

Test names follow `Test_<UnitOfWork>_<StateUnderTest>_<ExpectedBehaviour>`:
- `UnitOfWork`: method or unit being tested (e.g. `Create`, `RowToProject`)
- `StateUnderTest`: inputs or conditions (e.g. `NilDescription`, `DuplicateID`)
- `ExpectedBehaviour`: the result (e.g. `SuccessfulProjectCreation`, `ReturnsError`)

Example: `Test_Create_ValidDescription_SuccessfulProjectCreation`

Variable naming in assertions:
- `expected` — the value you set up / expect
- `actual` — the value returned by the code under test
