# PRD-02: CLI Command — `gask list`

## Status

`finished`

## Dependencies

- **PRD-01-data-layer.md** must be complete before this PRD is executed.
- `EnsureInitialised` guard, `Repository` struct, and `.gask/gask.db` must all exist.

---

## 1. Objective

Give the user a clear, scannable overview of all tasks in the current project directory. Output is a Lip Gloss rendered table in the terminal. No pagination, no global view, no priorities.

---

## 2. Scope

**In scope:**

- `gask list` command with default and `--long` flag behaviour
- `ListTasks` repository method
- Age formatter utility
- Empty state message
- Terminal width adaptation

**Out of scope:**

- Pagination — display all tasks at once in v1
- Global task view across directories
- `priority` field — not in the schema, not in this command
- Sorting by anything other than `status` and `created_at`

---

## 3. Data Retrieval

The `list` command fetches all tasks via the repository. Sorting is handled **inside the repository method**, not in the command layer.

### Sort order

1. By `status` in this fixed order: `doing` → `todo` → `blocked` → `done`
2. Within each status group, by `created_at` ascending (oldest first)

`done` tasks appear last because they are archive state.

### Repository method signature

Add to `internal/db/repository.go`:

```go
// ListTasks returns all tasks sorted by status group then created_at ascending.
// Sort order: doing → todo → blocked → done.
func (r *Repository) ListTasks(ctx context.Context) ([]*Task, error)
```

---

## 4. Terminal Rendering

Use **Lip Gloss** for all rendering. All rendering logic lives in `cmd/list.go` — none of it goes into `internal/`.

### Table columns

| Column   | Source field     | Notes                                                   |
| :------- | :--------------- | :------------------------------------------------------ |
| `ID`     | `task.ID`        | Integer, right-aligned                                  |
| `Status` | `task.Status`    | Color-coded label (see below)                           |
| `Title`  | `task.Title`     | Truncated at word boundary if over 50 chars, append `…` |
| `Age`    | `task.CreatedAt` | Human-readable relative time (see Section 5)            |

### Status colors

| Status    | Color  |
| :-------- | :----- |
| `todo`    | Blue   |
| `doing`   | Yellow |
| `done`    | Green  |
| `blocked` | Red    |

### Table style

- Header row: bold, underlined
- No alternating row backgrounds in v1 — keep it simple
- Table width: read terminal width via `lipgloss.Width` and fit to it
- If terminal width cannot be determined, default to 100 characters

### Empty state

If `ListTasks` returns zero tasks, do not render a table. Print:

```
No tasks found in this directory. Run 'gask add' to create one.
```

---

## 5. Age Formatter

Create `cmd/age.go` — a small helper used only by the list command.

```go
// FormatAge returns a human-readable string for how long ago t was.
// Examples: "just now", "5m ago", "2h ago", "3d ago", "4w ago"
func FormatAge(t time.Time) string
```

### Formatting rules

| Duration       | Format                   |
| :------------- | :----------------------- |
| Under 1 minute | `just now`               |
| Under 1 hour   | `Nm ago` (e.g. `5m ago`) |
| Under 24 hours | `Nh ago` (e.g. `2h ago`) |
| Under 7 days   | `Nd ago` (e.g. `3d ago`) |
| 7 days or more | `Nw ago` (e.g. `4w ago`) |

No months or years in v1. Anything over 7 days shows in weeks.

---

## 6. `--long` Flag

Add a `-l` / `--long` boolean flag to `gask list`.

When `--long` is set, add a `Description` column after `Title`:

- Show the first 80 characters of `task.Description`
- Truncate at word boundary, append `…` if truncated
- If `task.Description` is empty, show `—` (em dash)

When `--long` is not set, the `Description` column is not rendered at all.

---

## 7. Init Guard

The `list` command must call `db.EnsureInitialised` before any DB access.
This is already wired into `cmd/root.go`'s `PersistentPreRunE` from PRD-01 — no additional work needed here. This is listed explicitly so the agent does not add a duplicate check inside `cmd/list.go`.

---

## 8. Files To Create or Modify

```
cmd/
├── list.go     ← CREATE — command definition, table rendering, flag wiring
└── age.go      ← CREATE — FormatAge helper, used only by list command

internal/db/
└── repository.go  ← MODIFY — add ListTasks method
```

No other files should be created or modified by this PRD.

---

## 9. Protected Paths — DO NOT MODIFY

```
ent/    ← fully generated, do not touch
```

If a change seems needed in `ent/`, it means the schema in `ent/schema/task.go` needs updating — but that is out of scope for this PRD. Raise it as a gap in the spec instead.

---

## 10. Acceptance Criteria

- [x] **AC1** — `gask list` displays a table of all tasks from `.gask/gask.db`.
- [x] **AC2** — The `ID` column matches the SQLite primary key exactly.
- [x] **AC3** — Status labels are color-coded correctly for all four states.
- [x] **AC4** — Tasks are sorted: `doing` first, `done` last.
- [x] **AC5** — Table width adapts to terminal width and does not overflow.
- [x] **AC6** — Titles over 50 characters are truncated at a word boundary with `…`.
- [x] **AC7** — Running `gask list` with no tasks prints the empty state message and exits cleanly.
- [x] **AC8** — `gask list --long` shows the description column; without the flag it is hidden.
- [x] **AC9** — Running `gask list` without `.gask/` present exits with the init guard message.
- [x] **AC10** — `go test ./...` passes with no failures.

---

## 11. Tests

### Repository test

| Test name                 | What it verifies                                                           |
| :------------------------ | :------------------------------------------------------------------------- |
| `TestListTasks_SortOrder` | Returns tasks in correct status group order: doing → todo → blocked → done |
| `TestListTasks_Empty`     | Returns empty slice (not error) when no tasks exist                        |

### Formatter test

| Test name               | What it verifies                  |
| :---------------------- | :-------------------------------- |
| `TestFormatAge_JustNow` | Time 30s ago returns `"just now"` |
| `TestFormatAge_Minutes` | Time 5m ago returns `"5m ago"`    |
| `TestFormatAge_Hours`   | Time 3h ago returns `"3h ago"`    |
| `TestFormatAge_Days`    | Time 2d ago returns `"2d ago"`    |
| `TestFormatAge_Weeks`   | Time 14d ago returns `"2w ago"`   |

All repository tests use the in-memory SQLite instance from PRD-01.

---

## 12. Build & Verify Checklist

```
1. [x] go test ./internal/db/...   ← ListTasks tests must pass
2. [x] go test ./cmd/...           ← FormatAge tests must pass
3. [x] go test ./...               ← no regressions
4. [x] go build -o gask            ← must compile
5. [x] ./gask list                 ← verify empty state message
6. [x] (add a task manually via sqlite) → ./gask list  ← verify table renders
7. [x] ./gask list --long          ← verify description column appears
```
