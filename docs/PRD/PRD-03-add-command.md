# PRD-03: CLI Command — `gask add`

## Status

`finished`

## Dependencies

- **PRD-01-data-layer.md** must be complete — `Repository`, `CreateTask`, and `.gask/gask.db` must exist.
- **PRD-01b-editor-config.md** must be complete — `internal/config/config.go` and `.gask/config.json` must exist.
- **PRD-02-list-command.md** — not a hard dependency but `gask list` is used in the verify checklist.

---

## 1. Objective

Provide a zero-friction task capture experience. The user types a title in the terminal, the editor opens blank for an optional description, and the task is saved.

---

## 2. Scope

**In scope:**

- `gask add "<title>"` (or unquoted `gask add title ...`) command
- Title validation before the editor opens
- Two-step editor resolution (`$EDITOR` → `.gask/config.json`)
- Editor spawning via resolved editor binary
- Temp file lifecycle — create, read, delete
- Saving title + description to the DB via existing repository methods

**Out of scope:**

- `priority` field — not in the schema, not in this command
- Any interactive prompts asking the user about descriptions
- Adding new repository methods — reuse what PRD-01 already provides
- Editor configuration — that is PRD-01b's responsibility

---

## 3. Command Behaviour

### Usage

```
gask add "Implement auth middleware"
gask add Implement auth middleware
```

### Full flow

1. Resolve editor (see Section 4) — if no editor found, exit with error.
2. Join all arguments to form the `title`.
3. Validate the `title` (see Section 5) — if invalid, exit with validation error. Editor never opens.
4. Create a temp file in the system temp folder (see Section 6).
5. Spawn the resolved editor with the temp file path — connect Stdin, Stdout, Stderr.
6. Wait for the editor process to exit.
7. If editor exit code != 0, print error to Stderr and exit (no task created).
8. Read the temp file contents as `description`.
9. Delete the temp file (via `defer`).
10. Call `repo.CreateTask(ctx, title)`.
11. If `CreateTask` succeeds, call `repo.UpdateDescription(ctx, id, description)`.
    - If `UpdateDescription` fails, print warning `! Warning: Description could not be saved: <error>`.
12. Print confirmation: `✓ Task #<id> created: "<title>"`

---

## 4. Editor Resolution

The editor is resolved in this exact order — stop at the first match:

```
1. $EDITOR env var          ← always wins if set
2. .gask/config.json        ← set during gask init (PRD-01b)
3. error                    ← never a silent fallback
```

### Implementation

```go
// ResolveEditor returns the editor binary to use.
// Checks $EDITOR first, then .gask/config.json.
// Returns ErrNoEditor if neither is set.
func ResolveEditor(dir string) (string, error) {
    if e := os.Getenv("EDITOR"); e != "" {
        return e, nil
    }
    cfg, err := config.Load(dir)
    if err == nil && cfg.Editor != "" {
        return cfg.Editor, nil
    }
    return "", ErrNoEditor
}
```

### Error sentinel

```go
var ErrNoEditor = errors.New(
    "$EDITOR is not set and no editor is configured. " +
    "Run 'gask init' to configure one, or export EDITOR=nvim in your shell.",
)
```

This function lives in `internal/editor/editor.go` alongside the existing `Open` function.

---

### Initialisation Guard

The command must first verify that the application is initialized by checking for the existence of the `.gask` configuration directory. If the directory does not exist, the command must exit with a message: `Error: gask not initialized. Run 'gask init' first.`

---

## 5. Validation Rules

Title is validated before the editor opens. The implementation must enforce the following rules:

```go
type addTaskInput struct {
    Title string
}
```

| Rule                | Implementation     | Error message                          |
| :------------------ | :----------------- | :------------------------------------- |
| Empty/Missing title | `len(trim) == 0`   | `"title is required"`                  |
| Over 100 characters | `len(title) > 100` | `"title cannot exceed 100 characters"` |

---

## 6. Temp File Management

### Location

System temp folder via `os.CreateTemp("", "gask-*.md")`.

### Cleanup rule

```go
f, err := os.CreateTemp("", "gask-*.md")
if err != nil { return err }
defer os.Remove(f.Name())
```

---

## 7. Repository Usage

Do **not** add new repository methods. Use the two methods that already exist from PRD-01:

```go
task, err := repo.CreateTask(ctx, title)
err = repo.UpdateDescription(ctx, task.ID, description)
```

---

## 8. Files To Create or Modify

```
cmd/
└── add.go                  ← CREATE — command definition, validation, flow orchestration

internal/editor/
└── editor.go               ← MODIFY — add ResolveEditor(), Open(), ErrNoEditor sentinel
```

`internal/db/repository.go` — **do not modify**.
`internal/config/config.go` — **do not modify**.

---

## 9. Protected Paths — DO NOT MODIFY

```
ent/    ← fully generated, do not touch
```

---

## 10. Acceptance Criteria

- [x] **AC1** — `gask add "Fix auth"` opens a blank editor.
- [x] **AC2** — If `$EDITOR` is set, it takes priority.
- [x] **AC3** — If `$EDITOR` is unset but `.gask/config.json` has an editor, that editor is used.
- [x] **AC4** — If neither is set, command exits with `ErrNoEditor` message.
- [x] **AC5** — An empty title exits with a validation error.
- [x] **AC6** — A title over 100 characters exits with a validation error.
- [x] **AC7** — Exiting editor with save (exit 0) saves file content as description.
- [x] **AC8** — Closing/canceling editor (exit non-zero) aborts, deletes temp file, no task created.
- [x] **AC9** — The newly created task appears in `gask list`.
- [x] **AC10** — Running `gask add` without `.gask/` present exits with the init guard message.
- [x] **AC11** — `go test ./...` passes.

---

## 11. Tests

### Editor resolution tests

| Test name                             | What it verifies                                      |
| :------------------------------------ | :---------------------------------------------------- |
| `TestResolveEditor_EnvVarWins`        | `$EDITOR` set → returns env var value, ignores config |
| `TestResolveEditor_FallsBackToConfig` | `$EDITOR` unset → returns editor from config.json     |
| `TestResolveEditor_NeitherSet`        | Both unset → returns `ErrNoEditor`                    |

### Mock editor pattern

```go
// In tests, inject a function that writes to the file instead of spawning a real editor
editor.Runner = func(name string, arg ...string) error {
    return os.WriteFile(arg[0], []byte("mock description"), 0644)
}
```

---

## 12. Build & Verify Checklist

```
1. go test ./internal/editor/...
2. go test ./cmd/...
3. go test ./...
4. go build -o gask
5. EDITOR=nvim ./gask add test task
6. unset EDITOR && ./gask add "test"
7. ./gask add "ab" (verify min len 1)
```
