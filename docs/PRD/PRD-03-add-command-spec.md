# SPEC: PRD-03: CLI Command — `gask add`

## Status

`finished`

## PRD Reference

`docs/PRD/PRD-03-add-command.md`

---

## My Understanding — What I Am Building

I am building the `gask add` command, which provides a zero-friction task capture flow. The command will accept a task title (validated for emptiness and length), resolve an editor (preferring `$EDITOR` over the saved config), spawn the editor to capture an optional description, and finally use existing repository methods (`CreateTask` and `UpdateDescription`) to persist the data.

## Files I Will Touch

| File                             | Action | Reason                                                     |
| :------------------------------- | :----- | :--------------------------------------------------------- |
| `cmd/add.go`                     | CREATE | Define the `add` command, validation logic, and flow.      |
| `internal/editor/editor.go`      | MODIFY | Add `ResolveEditor`, `Open`, and `ErrNoEditor` sentinel.   |
| `internal/editor/editor_test.go` | MODIFY | Add tests for `ResolveEditor` logic and mock spawning.     |

**I will not open, read, or modify any file not listed in this table.**

## Decisions I Am Making

- [x] **Temp file naming**: Use `os.CreateTemp("", "gask-*.md")` for consistent naming and automatic cleanup via `defer`.
- [x] **Command parsing**: The command will accept arguments either quoted or unquoted, joining them into a single string for the `title`.
- [x] **Editor spawning**: Use `os/exec` to spawn the editor, directly piping `os.Stdin`, `os.Stdout`, and `os.Stderr`.
- [x] **Test Mocking**: Implement `var EditorRunner = exec.Command` in `internal/editor/editor.go` to inject a mockable runner for tests.
- [x] **Validation Timing**: Title is validated *before* the editor is opened, ensuring an invalid title aborts the flow immediately.

## Assumptions I Am Filling In

- [x] **Validation**: Title is invalid if empty (after trimming) or > 100 characters.
- [x] **Process Lifecycle**: The command will correctly handle the editor process exit code, ensuring the task is not created if the user aborts (non-zero exit).

## Gaps & Questions

- None — PRD is unambiguous.

## Implementation Order

- [x] Step 1 — Modify `internal/editor/editor.go` to add `ResolveEditor`, `Open` method, and the mockable runner variable.
- [x] Step 2 — Implement tests in `internal/editor/editor_test.go` using the mock runner.
- [x] Step 3 — Create `cmd/add.go`, implementing the orchestration logic (resolve -> validate -> open -> save).
- [x] Step 4 — Verify with build and manual CLI checks, ensuring the init guard works correctly.

## Protected Paths Acknowledgement

- [x] I have read the Protected Paths section of the PRD.
      I will not modify any file that begins with a generated-code header (specifically the `ent/` directory).
      If I encounter such a file that needs to change, I will stop,
      fix the source (e.g. the schema), and re-run the generator instead.

## AC Mapping

| AC   | Satisfied by step |
| :--- | :---------------- |
| AC1  | Step 3            |
| AC2  | Step 3            |
| AC3  | Step 3            |
| AC4  | Step 3            |
| AC5  | Step 3            |
| AC6  | Step 3            |
| AC7  | Step 3            |
| AC8  | Step 3            |
| AC9  | Step 3            |
| AC10 | Step 3            |
| AC11 | Step 4            |
