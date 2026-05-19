# PRD-04: CLI Command — `gask status`

## Status

`finished`

## Dependencies

- **PRD-01-data-layer.md** must be complete — `Repository`, `UpdateStatus`, and `.gask/gask.db` must exist.
- **PRD-02-list-command.md** — not a hard dependency but `gask list` is used in the verify checklist.

---

## 1. Objective

Allow users to transition a task between lifecycle states from the terminal. One command, one task, one status — fast and explicit.

---

## 2. Scope

**In scope:**

- `gask status <id> <status>` command
- ID and status validation before any DB call
- `done → todo` guard (see Section 4)
- Success and error output
- Tests for all validation and transition rules

**Out of scope:**

- Interactive TUI or hotkey-based status updates — future PRD
- Interactive selection menu when status argument is omitted — future PRD
- Any changes to `internal/db/repository.go` — `UpdateStatus` already exists from PRD-01

---

## 3. Command Behaviour

### Usage

```
gask status <id> <status> [flags]
```

### Flags

- `-f`, `--force` : Force transition back to todo from done.

### Full flow

1. Parse `<id>` as integer.
2. Parse `<status>` as string.
3. Validate ID and status using `github.com/go-playground/validator/v10` with `updateStatusInput` struct.
4. Check for `.gask` initialization.
5. Fetch task by ID via `repo.GetTaskByID` — if not found, exit with error.
6. Check `done → todo` guard: If current status is `done`, requested is `todo`, and `--force` is not set, exit with guard message.
7. Call `repo.UpdateStatus(ctx, id, status)`.
8. Print confirmation and exit.

### Success output

```
✓ Task #<id> status updated: <old_status> → <new_status>
```

### Argument errors

| Situation                            | Error message                                                 |
| :----------------------------------- | :------------------------------------------------------------ |
| `<id>` is not an integer             | `"invalid task ID: must be a number"`                         |
| `<status>` is not a valid enum value | `"invalid status: must be one of todo, doing, blocked, done"` |
| Task ID does not exist               | `"task #<id> not found"`                                      |
| `done → todo` blocked by guard       | `"Task #<id> is already done. Are you sure you want to move it back to todo?\nRun: gask status <id> todo --force"` |

---

## 4. The `done → todo` Guard

CONTEXT.md states: a `done` task cannot go back to `todo` without explicit intent.

When the current status is `done` and the requested status is `todo`, the command must exit with:

```
Task #<id> is already done. Are you sure you want to move it back to todo?
Run: gask status <id> todo --force
```

The `--force` flag bypasses the guard and allows the transition. Without `--force` the command exits cleanly with no DB change.

All other transitions from `done` (`done → doing`, `done → blocked`) are allowed without a force flag — only `done → todo` requires it.

---

## 5. Validation and Error Handling

- **Validation:** Use `github.com/go-playground/validator/v10` against the `updateStatusInput` struct.
- **Error Handling:** Follow the established CLI pattern used in `cmd/add.go` and `cmd/list.go`: implement the command logic within `RunE` and return formatted errors which the CLI prints automatically.

```go
type updateStatusInput struct {
    ID     int    `validate:"required,min=1"`
    Status string `validate:"required,oneof=todo doing blocked done"`
}
```

---

## 6. Files To Create or Modify

```
cmd/
└── status.go       ← CREATE — command definition, validation, guard, output
```

`internal/db/repository.go` — **do not modify**. `UpdateStatus` and `GetTaskByID` already exist.

No other files should be created or modified by this PRD.

---

## 7. Protected Paths — DO NOT MODIFY

```
ent/    ← fully generated, do not touch
```

---

## 8. Acceptance Criteria

- [ ] **AC1** — `gask status 1 done` updates the task status in the DB and prints the confirmation.
- [ ] **AC2** — `gask status 1 complete` exits with the invalid status error message.
- [ ] **AC3** — `gask status abc done` exits with the invalid ID error message.
- [ ] **AC4** — `gask status 999 done` exits with the not found error message.
- [ ] **AC5** — `gask status 1 todo` on a `done` task exits with the force prompt message. DB is not changed.
- [ ] **AC6** — `gask status 1 todo --force` on a `done` task updates the DB successfully.
- [ ] **AC7** — `done → doing` and `done → blocked` transitions work without `--force`.
- [ ] **AC8** — `updated_at` is updated on a successful status change (handled by Ent, verified in test).
- [ ] **AC9** — Running `gask status` without `.gask/` present exits with the init guard message.
- [ ] **AC10** — `go test ./...` passes with no failures.

---

## 9. Tests

| Test name                             | What it verifies                                 |
| :------------------------------------ | :----------------------------------------------- |
| `TestUpdateStatus_Valid`              | Valid transition updates status and `updated_at` |
| `TestUpdateStatus_InvalidStatus`      | Unknown status string returns validation error   |
| `TestUpdateStatus_InvalidID`          | Non-integer ID returns parse error               |
| `TestUpdateStatus_NotFound`           | Missing task ID returns `ErrNotFound`            |
| `TestUpdateStatus_DoneToTodo_NoForce` | Returns guard message, DB unchanged              |
| `TestUpdateStatus_DoneToTodo_Force`   | With `--force`, transition succeeds              |
| `TestUpdateStatus_DoneToDoing`        | `done → doing` succeeds without `--force`        |
| `TestUpdateStatus_DoneToBlocked`      | `done → blocked` succeeds without `--force`      |

All tests use the in-memory SQLite instance from PRD-01.

---

## 10. Build & Verify Checklist

```
1. go test ./cmd/...              ← status command tests
2. go test ./...                  ← no regressions
3. go build -o gask               ← must compile
4. ./gask status 1 doing          ← verify success output
5. ./gask status 1 complete       ← verify invalid status error
6. ./gask status abc done         ← verify invalid ID error
7. ./gask status 999 done         ← verify not found error
8. ./gask status 1 todo           ← verify done→todo guard (after marking done)
9. ./gask status 1 todo --force   ← verify force flag bypasses guard
10. ./gask list                   ← verify status change is reflected
```
