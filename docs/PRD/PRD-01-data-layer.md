# PRD-01: Data Layer — Schema, Migrations & Repository

## Status

`finished`

## Dependencies

None — this is the foundation. All other PRDs depend on this one being complete.

---

## 1. Objective

Build the data persistence layer for gask using **Ent ORM** and **SQLite**.

This layer is the only part of the codebase that touches the database. When this PRD is complete, every other command PRD (add, list, status, edit) can be built on top of it without touching a single SQL string.

This PRD also establishes the `gask init` command and the `.gask/` project folder — the single source of truth for where all gask data lives.

---

## 2. Scope

**In scope:**

- `gask init` command that creates `.gask/` and initialises the DB
- Ent schema definition for the `Tasks` entity
- SQLite connection setup with WAL mode
- DB file located at `.gask/gask.db` (not `.gask.db` in root)
- Legacy `.gask.db` detection: `gask init` fails if legacy file exists
- Auto-migration on startup
- Repository struct with all CRUD operations returning domain-specific structs
- Input validation before any DB mutation
- Tests for all repository methods

**Out of scope:**

- All other CLI commands (`add`, `list`, `status`, `edit`) — those are separate PRDs
- `priority` field — not part of the core workflow defined in CONTEXT.md
- `metadata` field — no feature currently requires it; add in a future PRD

---

## 3. The `.gask/` Folder

### Why a folder, not a file

Putting the DB directly in the project root as `.gask.db` creates problems as the project grows — config, logs, templates, and other gask internals have nowhere clean to live. A dedicated `.gask/` folder solves this now before it becomes a migration problem.

### Structure after `gask init`

```
your-project/
└── .gask/
    └── gask.db     ← the SQLite database, lives here and only here
```

Future PRDs may add files inside `.gask/` (e.g. config). Nothing outside `.gask/` should be created or modified by gask.

### `.gitignore` recommendation

`gask init` must print a reminder (not auto-modify `.gitignore`) that users should add `.gask/` to their `.gitignore` if they do not want tasks committed to version control. The message should be:

```
Hint: add .gask/ to your .gitignore to keep tasks out of version control.
```

---

## 4. `gask init` Command

### File location

`cmd/init.go`

### Behaviour

1. Check if `.gask.db` (legacy file) exists in the current working directory.
   - If it exists → print `Found legacy .gask.db file. Please move it to .gask/gask.db or delete it before running init.` and exit with error.
2. Check if `.gask/` already exists in the current working directory.
   - If it exists → print `gask is already initialised in this directory.` and exit cleanly (no error, no overwrite).
3. If it does not exist → create `.gask/` directory.
4. Open the DB connection (this triggers auto-migration, creating `gask.db`).
5. Print success message and the `.gitignore` hint.

### Output on success

```
✓ Initialised gask in .gask/
  DB: .gask/gask.db

Hint: add .gask/ to your .gitignore to keep tasks out of version control.
```

### Error cases

| Situation                                | Error message                                   |
| :--------------------------------------- | :---------------------------------------------- |
| No write permission in current directory | `init: cannot create .gask/: permission denied` |
| DB migration fails                       | `init: migration failed: <wrapped error>`       |

---

## 5. DB Initialisation Guard

Every command other than `init` must check that `.gask/gask.db` exists before doing anything. If it does not exist, the command must exit with:

```
gask is not initialised in this directory. Run 'gask init' first.
```

This check lives in `internal/db/sqlite.go` as a helper function `EnsureInitialized(dir string) error`. The `cmd/root.go` persistent pre-run hook calls it so every subcommand gets the check for free — no need to repeat it in each command file.

**Exception:** `gask init` itself skips this check (it is the one that creates the folder).

---

## 6. Data Model: Tasks

### Fields

| Field         | Ent Type | Constraints                        | Default      |
| :------------ | :------- | :--------------------------------- | :----------- |
| `id`          | `int`    | Auto-increment, Primary Key        | auto         |
| `title`       | `string` | Required, max 100 chars            | —            |
| `description` | `string` | Optional (Markdown content)        | `""`         |
| `status`      | `enum`   | `todo`, `doing`, `done`, `blocked` | `todo`       |
| `created_at`  | `time`   | Immutable, set on insert           | `time.Now()` |
| `updated_at`  | `time`   | Updated on every mutation          | `time.Now()` |

### Schema file location

`ent/schema/tasks.go` — the only file in `ent/` you are allowed to edit.
Always run `go generate ./ent/...` after every change to this file.

---

## 7. SQLite Configuration

**Driver:** `github.com/mattn/go-sqlite3` — CGO required.

**File path:** `.gask/gask.db` relative to the current working directory.

**Connection string:**

```
file:.gask/gask.db?_fk=1&_journal_mode=WAL
```

WAL mode is required for concurrent reads during writes. Foreign keys must be enabled explicitly — SQLite does not enable them by default.

---

## 8. Files To Create or Modify

```
cmd/
└── init.go              ← 'gask init' command

internal/db/
├── sqlite.go            ← DB init, driver setup, auto-migration, EnsureInitialized()
├── repository.go        ← ALL CRUD lives here and only here
└── task.go              ← Task domain struct (separate from Ent-generated types)

ent/schema/
└── tasks.go              ← Ent schema — edit this, then go generate
```

No other files should be created or modified by this PRD.

---

## 9. Protected Paths — DO NOT MODIFY

The following paths are **auto-generated** by `go generate ./ent/...`.
Editing them manually will be silently overwritten on the next generate run and may corrupt the build.

```
ent/                     ← everything in here is generated
  (except ent/schema/)   ← this subfolder is the only exception
```

Every generated file starts with this header — if you see it, stop and do not edit:

```go
// Code generated by ent, DO NOT EDIT.
```

**If something in `ent/` looks wrong, the fix is always:**

1. Edit `ent/schema/task.go`
2. Run `go generate ./ent/...`

Never edit generated files directly. Not even to fix a one-line compilation error.

---

## 10. Repository API

The `Repository` struct must expose exactly these methods. No more, no less — additional methods will be added in future PRDs when a command actually needs them.

```go
// CreateTask inserts a new task. Returns validation error if title is empty or too long.
func (r *Repository) CreateTask(ctx context.Context, title string) (*db.Task, error)

// GetAllTasks returns all tasks ordered by created_at descending.
func (r *Repository) GetAllTasks(ctx context.Context) ([]*db.Task, error)

// GetTaskByID returns a single task by its integer ID.
// Returns a sentinel ErrNotFound (not an Ent error) if no task exists with that ID.
func (r *Repository) GetTaskByID(ctx context.Context, id int) (*db.Task, error)

// UpdateStatus changes the status of a task.
// Returns a validation error if the status string is not one of the four valid values.
func (r *Repository) UpdateStatus(ctx context.Context, id int, status string) (*db.Task, error)

// UpdateDescription sets the Markdown description on a task (used by gask edit).
func (r *Repository) UpdateDescription(ctx context.Context, id int, description string) (*db.Task, error)
```

---

## 11. Validation Rules

Use `go-playground/validator` with typed request structs. Raw strings from CLI args must never reach Ent directly.

```go
// Example — CreateTask validates via a struct before any DB call
type createTaskInput struct {
    Title string `validate:"required,min=1,max=100"`
}
```

Validation errors must be returned as `error` values — never panic, never log-and-continue.

---

## 12. Error Handling

- Wrap all errors with context before returning: `fmt.Errorf("repository.GetTaskByID: %w", err)`
- When a task is not found, return the sentinel `ErrNotFound` — do not leak Ent's internal `NotFoundError` type upward into `cmd/`
- Never return `nil, nil` — always return either a valid result or an error

```go
// Sentinel error defined in internal/db/repository.go
var ErrNotFound = errors.New("task not found")
```

---

## 13. Acceptance Criteria

- [x] **AC1** — `go generate ./ent/...` runs without errors and produces a fully typed client.
- [x] **AC2** — `go build -o gask` compiles cleanly.
- [x] **AC3** — `./gask init` creates `.gask/gask.db` and prints the success message.
- [x] **AC4** — Running `./gask init` a second time prints the already-initialised message and exits cleanly without overwriting anything.
- [x] **AC5** — Running any command other than `init` without `.gask/` present exits with the "run gask init first" message.
- [x] **AC6** — Creating a task with an empty title returns a validation error. The DB is never touched.
- [x] **AC7** — Creating a task with a title over 100 characters returns a validation error. The DB is never touched.
- [x] **AC8** — `created_at` and `updated_at` are set automatically by Ent hooks. No `cmd/` code sets timestamps manually.
- [x] **AC9** — `go test ./internal/db/...` passes with no failures.x

---

## 14. Tests

All tests run against an **in-memory SQLite instance** — never against `.gask/gask.db`.

```go
// In-memory connection string for tests only
"file:testdb?mode=memory&cache=shared&_fk=1"
```

### Required test cases

| Test name                       | What it verifies                                                            |
| :------------------------------ | :-------------------------------------------------------------------------- |
| `TestCreateTask_Success`        | Task saved and retrieved by ID with correct title and default status `todo` |
| `TestCreateTask_EmptyTitle`     | Returns validation error, no DB row created                                 |
| `TestCreateTask_TitleTooLong`   | Returns validation error for title > 100 chars                              |
| `TestGetTaskByID_NotFound`      | Returns `ErrNotFound` sentinel for a missing ID                             |
| `TestUpdateStatus_Valid`        | Status updates correctly, `updated_at` is after the original value          |
| `TestUpdateStatus_InvalidEnum`  | Returns validation error for an unknown status string                       |
| `TestGetAllTasks_Order`         | Multiple tasks returned in `created_at` descending order                    |
| `TestEnsureInitialized_Missing` | Returns error when `.gask/` does not exist                                  |
| `TestEnsureInitialized_Present` | Returns nil when `.gask/gask.db` exists                                     |

Use `require` for setup steps that must succeed. Use `assert` for the actual outcome checks.

---

## 15. Build & Verify Checklist

Run these in order after implementation. Do not proceed to the next step if one fails.

```
1. go generate ./ent/...      ← regenerate Ent client from schema
2. go test ./internal/db/...  ← all 9 tests must pass
3. go test ./...              ← no regressions anywhere
4. go build -o gask           ← binary must compile cleanly
5. ./gask init                ← verify .gask/gask.db is created
6. ./gask init                ← verify second run prints already-initialised message
7. rm -rf .gask && ./gask list ← verify "run gask init first" guard works
```
