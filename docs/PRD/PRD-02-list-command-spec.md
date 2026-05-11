# SPEC: PRD-02: CLI Command — `gask list`

## Status

`finished`

<!-- Human changes this to `approved` before implementation starts.
     Do not change this yourself under any circumstance. -->

## PRD Reference

`docs/PRD/PRD-02-list-command.md`

---

## My Understanding — What I Am Building

I am building the `gask list` command which displays all tasks in a color-coded table. This involves adding a new retrieval method to the repository with custom sorting logic and implementing a terminal-responsive UI using Lip Gloss.

---

## Files I Will Touch

| File                             | Action | Reason                                                         |
| :------------------------------- | :----- | :------------------------------------------------------------- |
| `internal/db/repository.go`      | MODIFY | Add `ListTasks` method with custom status sorting.             |
| `cmd/list.go`                    | CREATE | Define the `list` command, flags, and table rendering logic.   |
| `cmd/age.go`                     | CREATE | Utility for human-readable relative time formatting.           |
| `cmd/age_test.go`                | CREATE | Unit tests for the age formatter.                              |
| `internal/db/repository_test.go` | MODIFY | Add integration tests for `ListTasks` sorting and empty state. |
| `go.mod`                         | MODIFY | Add `github.com/charmbracelet/lipgloss` dependency.            |

**I will not open, read, or modify any file not listed in this table.**

---

## Decisions I Am Making

- [x] **In-memory sorting:** I will perform the custom status sorting in Go within the `ListTasks` method rather than using complex SQL `CASE` statements, as it is more maintainable for this scale.
- [x] **Table Library:** I will use `github.com/charmbracelet/lipgloss/table` for rendering if available in the latest Lip Gloss version, otherwise I will manualy construct the table using Lip Gloss primitives.
- [x] **Word Boundary Truncation:** Truncation will look for the last space before the limit to avoid cutting words, appending `…` only if the string was actually shortened.

---

## Assumptions I Am Filling In

- [x] **Age Rounding:** "4w ago" will be used for any duration >= 7 days, rounding down to the nearest week (e.g., 13 days is "1w ago").
- [x] **Terminal Width:** If `lipgloss` cannot detect the width (e.g. in some CI environments), I will fallback to 100 as specified.

---

## Gaps & Questions

- [x] **Dependency:** `github.com/charmbracelet/lipgloss` is missing from `go.mod`. Should I add it?
- [x] **Title Truncation:** The PRD says "Truncated at word boundary if over 50 chars". Should it be _maximum_ 50 chars including the `…` or 50 chars of text plus `…`? I will assume 50 total characters.
- [x] **Empty Description:** For `--long`, if description is empty, show `—`. Should this em-dash also be centered or left-aligned? I will assume left-aligned to match the description text.

---

## Implementation Order

- [x] Step 1 — Add `github.com/charmbracelet/lipgloss` to `go.mod`.
- [x] Step 2 — Implement `ListTasks` in `internal/db/repository.go` and add tests in `internal/db/repository_test.go`.
- [x] Step 3 — Create `cmd/age.go` and `cmd/age_test.go` with the relative time logic.
- [x] Step 4 — Create `cmd/list.go`, implement the `list` command and table rendering.
- [x] Step 5 — Wire `list` command into `cmd/root.go`.
- [x] Step 6 — Manual verification of all ACs and running `go test ./...`.

---

## Protected Paths Acknowledgement

- [x] I have read the Protected Paths section of the PRD.
      I will not modify any file that begins with a generated-code header.
      If I encounter such a file that needs to change, I will stop,
      fix the source (e.g. the schema), and re-run the generator instead.

---

## AC Mapping

| AC   | Satisfied by step         |
| :--- | :------------------------ |
| AC1  | Step 4                    |
| AC2  | Step 2, 4                 |
| AC3  | Step 4                    |
| AC4  | Step 2                    |
| AC5  | Step 4                    |
| AC6  | Step 4                    |
| AC7  | Step 4                    |
| AC8  | Step 4                    |
| AC9  | Step 4, 5 (Pre-run guard) |
| AC10 | Step 6                    |

---

## Ready to Implement?

- [x] All gaps above are resolved by the human
- [x] Human has changed Status to `approved`
- [x] Every AC maps to at least one implementation step

**Do not begin implementation until all three boxes are checked.**
