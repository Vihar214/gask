# PRD-01b: Editor Config — Init Prompt & Config File

## Status

`complete`

## Dependencies

- **PRD-01-data-layer.md** must be complete — `.gask/` folder and `gask init` command must exist.

---

## 1. Objective

Extend `gask init` to ask the user which editor they use and save that choice to `.gask/config.json`. This means `gask add` never fails due to a missing `$EDITOR` after a proper init.

---

## 2. Scope

**In scope:**

- [x] Editor prompt during `gask init` (idempotent: skip only if `config.json` exists)
- [x] `exec.LookPath` validation for **all** editor inputs (presets and custom)
- [x] Specific error messaging/hints for failed validation (e.g., VS Code PATH instructions)
- [x] Writing editor choice to `.gask/config.json`
- [x] `internal/config/config.go` — read/write helper for the config file
- [x] `internal/editor/editor.go` — validation logic and editor hints
- [x] Tests for prompt, validation, and config read/write

**Out of scope:**

- Any other config keys — this file only stores `editor` for now
- Changes to `gask add` or editor spawning — that is PRD-03's responsibility
- Changing anything already implemented in PRD-01

---

## 3. `.gask/config.json`

### Structure

```json
{
  "editor": "nvim"
}
```

Simple flat JSON. One key for now. Future PRDs may add keys — the read/write helper must not break if unknown keys are present (use `json.Decoder` not `json.Unmarshal` on a strict struct).

### File location

`.gask/config.json` alongside `.gask/gask.db`.

```
.gask/
├── gask.db         ← created by PRD-01
└── config.json     ← created by this PRD
```

---

## 4. Architecture

### Config Helper (`internal/config/config.go`)

```go
type Config struct {
    Editor string `json:"editor"`
}

// Load reads .gask/config.json from the given dir.
// Returns ErrConfigNotFound if the file does not exist.
func Load(dir string) (*Config, error)

// Save writes the config to .gask/config.json in the given dir.
// Creates the file if it does not exist, overwrites if it does.
func Save(dir string, cfg *Config) error
```

### Editor Domain (`internal/editor/editor.go`)

```go
// Validate checks if the editor exists in PATH.
// Returns the absolute path or the command name if valid.
// Returns a descriptive error with "how-to-fix" hints if invalid.
func Validate(name string) (string, error)
```

---

## 5. Init Prompt

### When to show it

Show the editor prompt during `gask init` if `.gask/config.json` does **not** exist.
If the `.gask/` directory exists but the config is missing, prompt anyway.
If the config exists, skip the prompt.

### Prompt flow

```
Which editor do you use?
  1. nvim
  2. vim
  3. nano
  4. code (VS Code)
  5. Other

Enter choice (1-5):
```

If the user enters `5` (Other):

```
Enter editor command (must be in PATH):
>
```

### Validation rules

Run `exec.LookPath(name)` for **every** input.

If validation fails, show a specific hint if available:

- **VS Code (`code`, `code-insiders`)**: `Note: ensure 'code' is in your PATH. In VS Code: Cmd+Shift+P -> 'Shell Command: Install code command in PATH'`
- **General**: `"[editor]" was not found in PATH. Is it installed?`

### Re-prompt logic

Maximum 3 attempts. After 3 failures, do not block init — complete init without editor config and print:

```
Could not validate editor after 3 attempts.
Run 'gask init' again or set $EDITOR in your shell to fix this later.
```

### Updated init success output

```
✓ Initialised gask in .gask/
  DB:     .gask/gask.db
  Editor: nvim

Hint: add .gask/ to your .gitignore to keep tasks out of version control.
```

---

## 6. Files To Create or Modify

```
internal/config/
└── config.go       ← CREATE — Load, Save, Config struct, ErrConfigNotFound

internal/editor/
└── editor.go       ← MODIFY/CREATE — Validate function and hints

cmd/
└── init.go         ← MODIFY — add editor prompt after DB init step
```

---

## 7. Protected Paths — DO NOT MODIFY

```
ent/    ← fully generated, do not touch
```

---

## 8. Acceptance Criteria

- [x] **AC1** — `gask init` prompts for editor choice if `.gask/config.json` is missing.
- [x] **AC2** — All inputs (presets and custom) are validated via `LookPath`.
- [x] **AC3** — Failed validation shows helpful hints (especially for VS Code).
- [x] **AC4** — Choosing `Other` and entering an invalid binary re-prompts up to 3 times.
- [x] **AC5** — Running `gask init` when config exists skips the prompt.
- [x] **AC6** — After 3 failed attempts, init completes without editor config and prints the fallback message.
- [x] **AC7** — `go test ./internal/config/...` and `./internal/editor/...` pass.

---

## 9. Tests

| Test name                        | What it verifies                                         |
| :------------------------------- | :------------------------------------------------------- |
| `TestConfigSave_And_Load`        | Save a config, load it back, values match                |
| `TestConfigLoad_NotFound`        | Returns `ErrConfigNotFound` when file missing            |
| `TestConfigLoad_UnknownKeys`     | Does not error if config.json has extra unknown keys     |
| `TestValidateEditor_ValidBinary` | LookPath succeeds for a real binary (e.g., `ls` or `sh`) |
| `TestValidateEditor_Invalid`     | Returns descriptive error/hint for invalid binary        |

---

## 10. Build & Verify Checklist

```
1. go test ./internal/config/...   ← [PASSED]
2. go test ./internal/editor/...   ← [PASSED]
3. go build -o gask                ← [PASSED]
4. rm -rf .gask && ./gask init     ← [VERIFIED]
5. cat .gask/config.json           ← [VERIFIED]
6. ./gask init                     ← [VERIFIED]
```
