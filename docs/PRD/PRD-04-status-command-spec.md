# Spec: PRD-04 Status Command

Status: `finished`

## My Understanding
Implement the `gask status` command to allow users to update the status of an existing task. The command requires an ID and a new status string. It includes validation, a `done → todo` transition guard (requiring a `--force` flag to bypass), and follows the existing repository and error-handling patterns.

## Files I Will Touch
- `cmd/status.go` (CREATE)
- `cmd/status_test.go` (CREATE)

## Decisions
- I will implement the command using `cobra` to maintain consistency with existing `gask` commands.
- I will use `github.com/go-playground/validator/v10` as specified in the PRD for validation.
- I will structure the `status.go` file to use `RunE` to follow the established CLI pattern for error propagation.

## Assumptions
- The `Repository` interface and `UpdateStatus`, `GetTaskByID` methods in `internal/db/repository.go` have the signatures required for this command as per PRD-01.
- `cmd/root.go` already exists and handles the root command setup.
- The `gask` CLI context (database initialization) is available via the existing `internal/db` or config management.

## Gaps
- None — PRD is unambiguous.

## Implementation Order
1. Create `cmd/status.go` with basic command definition, flag handling (`--force`), and validation logic.
2. Implement the command logic inside `RunE`, fetching the task, applying the `done → todo` guard, and calling `repo.UpdateStatus`.
3. Add success and error output formatting.
4. Create `cmd/status_test.go` and implement the test cases defined in the PRD using the in-memory SQLite repository.
5. Run tests and verify the build/verify checklist.

## AC Mapping
- **AC1** (Valid update): Covered by `TestUpdateStatus_Valid`
- **AC2** (Invalid status): Covered by `TestUpdateStatus_InvalidStatus`
- **AC3** (Invalid ID): Covered by `TestUpdateStatus_InvalidID`
- **AC4** (Not found): Covered by `TestUpdateStatus_NotFound`
- **AC5** (Done-to-Todo guard): Covered by `TestUpdateStatus_DoneToTodo_NoForce`
- **AC6** (Done-to-Todo force): Covered by `TestUpdateStatus_DoneToTodo_Force`
- **AC7** (Done transitions): Covered by `TestUpdateStatus_DoneToDoing`, `TestUpdateStatus_DoneToBlocked`
- **AC8** (Updated at): Implicitly covered by DB interaction, verified in `TestUpdateStatus_Valid`
- **AC9** (Init guard): Implicitly handled by root command/db init
- **AC10** (All tests pass): Final verification step

## Protected Paths
- [x] `ent/` — fully generated, do not touch
