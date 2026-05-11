# SPEC: PRD-01b: Editor Config — Init Prompt & Config File

## Status

`finished`

## PRD Reference

`docs/PRD/PRD-01b-editor-config.md`

---

## My Understanding — What I Am Building

[x] I am extending the `gask init` command to include an interactive editor selection process. This involves creating a configuration layer to persist the user's choice in `.gask/config.json`, an editor domain layer to validate the selected editor using `exec.LookPath` (with specific hints for common editors like VS Code), and updating the `init` command to handle the interactive loop and idempotency.

---

## Files I Will Touch

| File                             | Action | Reason                                              |
| :------------------------------- | :----- | :-------------------------------------------------- |
| `internal/config/config.go`      | CREATE | Handle loading and saving `.gask/config.json`       |
| `internal/config/config_test.go` | CREATE | Unit tests for configuration persistence            |
| `internal/editor/editor.go`      | MODIFY | Implement editor validation logic and fix-it hints  |
| `internal/editor/editor_test.go` | CREATE | Unit tests for editor validation logic              |
| `cmd/init.go`                    | MODIFY | Add the interactive prompt and 3-attempt retry loop |

**I will not open, read, or modify any file not listed in this table.**

---

## Decisions I Am Making

- [x] **Indented JSON**: I will use `json.NewEncoder(file).SetIndent("", "  ")` in `Save` to keep the config file human-readable.
- [x] **Prompt Input**: I will use `fmt.Scanln` or a `bufio.Scanner` to capture user input during the `init` command for the numbered menu and custom command string.
- [x] **Retry Messaging**: The 3-attempt counter will be managed within a loop in `cmd/init.go`, displaying the remaining attempts if validation fails.

---

## Assumptions I Am Filling In

- [x] **Validation scope**: `exec.LookPath` is considered sufficient for "found in PATH" validation on Linux/macOS environments where this CLI is likely to run.

---

## Gaps & Questions

None — PRD is unambiguous.

---

## Implementation Order

- [x] **Step 1**: Create `internal/config/config.go` with `Config` struct, `Load`, `Save`, and `ErrConfigNotFound` sentinel.
- [x] **Step 2**: Create `internal/config/config_test.go` to verify `Load`, `Save`, and `ErrConfigNotFound` behavior.
- [x] **Step 3**: Implement `Validate(name string) (string, error)` in `internal/editor/editor.go` with specific hints for VS Code.
- [x] **Step 4**: Create `internal/editor/editor_test.go` to verify validation for presets and "Other" commands.
- [x] **Step 5**: Modify `cmd/init.go` to add the interactive prompt logic after DB migration, ensuring idempotency (checking for `config.json` existence).
- [x] **Step 6**: Manual verification: run `gask init`, verify file creation, verify re-run behavior, and verify retry limit.

---

## Protected Paths Acknowledgement

- [x] I have read the Protected Paths section of the PRD.
      I will not modify any file that begins with a generated-code header (specifically the `ent/` directory).
      If I encounter such a file that needs to change, I will stop,
      fix the source (e.g. the schema), and re-run the generator instead.

---

## AC Mapping

| AC                                      | Satisfied by step |
| :-------------------------------------- | :---------------- |
| AC1 — Prompt if config missing          | Step 5            |
| AC2 — All inputs validated via LookPath | Step 3, 5         |
| AC3 — Validation shows helpful hints    | Step 3, 5         |
| AC4 — "Other" re-prompts up to 3 times  | Step 5            |
| AC5 — Run skips if config exists        | Step 5            |
| AC6 — Fail gracefully after 3 attempts  | Step 5            |
| AC7 — Unit tests pass                   | Step 2, 4         |

---

## Ready to Implement?

- [x] All gaps above are resolved by the human
- [x] Human has changed Status to `approved`
- [x] Every AC maps to at least one implementation step

**Do not begin implementation until all three boxes are checked.**
