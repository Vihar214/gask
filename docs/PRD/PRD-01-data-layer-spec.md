# SPEC: PRD-01: Data Layer — Schema, Migrations & Repository

## Status

`finished`

<!-- Human changes this to `approved` before implementation starts.
     Do not change this yourself under any circumstance. -->

## PRD Reference

`docs/PRD/PRD-01-data-layer.md`

---

## My Understanding — What I Am Building

I am building the foundational data persistence layer for `gask`. This involves setting up Ent ORM with a SQLite backend (WAL mode), defining the `Tasks` schema, implementing a repository for CRUD operations, and creating the `gask init` command to bootstrap the project's `.gask/` directory.

---

## Files I Will Touch

| File                             | Action | Reason                                                                                    |
| :------------------------------- | :----- | :---------------------------------------------------------------------------------------- |
| `cmd/root.go`                    | MODIFY | Implement a scalable Cobra root command structure with a dedicated execution entry point. |
| `cmd/init.go`                    | CREATE | Implement the `gask init` command.                                                        |
| `ent/schema/tasks.go`            | MODIFY | Ensure schema matches PRD requirements (renaming if necessary to match entity naming).    |
| `internal/db/sqlite.go`          | MODIFY | Implement SQLite connection logic, WAL mode, and `EnsureInitialized` helper.              |
| `internal/db/repository.go`      | CREATE | Implement the `Repository` struct and CRUD methods.                                       |
| `internal/db/task.go`            | MODIFY | Define the domain-specific `Task` struct.                                                 |
| `internal/db/repository_test.go` | CREATE | Implement tests for the repository and initialization guard.                              |

**I will not open, read, or modify any file not listed in this table.**

---

## Decisions I Am Making

- [x] Decision — I will use `tasks.go` as the schema file name as it's already present, but ensure the entity is correctly pluralized/singularized by Ent.
- [x] Decision — In `cmd/root.go`, I will check the command being executed to skip `EnsureInitialized` if it's the `init` command or `help`.
- [x] Decision — `Repository` will take the Ent client as a dependency for better testability.
- [x] Decision — I will use a scalable Cobra structure (e.g., dedicated `Execute` function and package-level variables for commands).

---

## Assumptions I Am Filling In

- [x] Assumption — The existing empty files in `cmd/` and `internal/db/` are safe to overwrite or fill.
- [x] Assumption — `go generate ./ent/...` will be executed via `run_shell_command` during development.
- [x] Assumption — The user has Go installed and environment set up to run `go test` and `go build`.

---

## Gaps & Questions

- [x] Gap — The PRD has a minor inconsistency: Section 6 says `ent/schema/tasks.go` but Section 8 says `ent/schema/task.go`. I will stick with `tasks.go`.
- [x] Gap — PRD says "Raw strings from CLI args must never reach Ent directly." I will implement validation in the Repository methods using `go-playground/validator`.

---

## Implementation Order

- [x] Step 1 — Initialize `cmd/root.go` with a proper, scalable Cobra structure.
- [x] Step 2 — Refine `ent/schema/tasks.go` and run `go generate ./ent/...`.
- [x] Step 3 — Implement domain struct in `internal/db/task.go`.
- [x] Step 4 — Implement connection and guard logic in `internal/db/sqlite.go`.
- [x] Step 5 — Implement `Repository` and CRUD methods in `internal/db/repository.go`.
- [x] Step 6 — Implement `gask init` in `cmd/init.go`.
- [x] Step 7 — Implement tests in `internal/db/repository_test.go`.
- [x] Step 8 — Update `cmd/root.go` with `PersistentPreRun` hook for the DB guard.
- [x] Step 9 — Verify with build and manual CLI checks.

---

## Protected Paths Acknowledgement

- [x] I have read the Protected Paths section of the PRD.
      I will not modify any file that begins with a generated-code header.
      If I encounter such a file that needs to change, I will stop,
      fix the source (e.g. the schema), and re-run the generator instead.

---

## AC Mapping

| AC  | Satisfied by step |
| :-- | :---------------- |
| AC1 | Step 2            |
| AC2 | Step 9            |
| AC3 | Step 6, 9         |
| AC4 | Step 6, 9         |
| AC5 | Step 8, 9         |
| AC6 | Step 5, 7         |
| AC7 | Step 5, 7         |
| AC8 | Step 2            |
| AC9 | Step 7            |

---

## Ready to Implement?

- [x] All gaps above are resolved by the human
- [x] Human has changed Status to `approved`
- [x] Every AC maps to at least one implementation step

**Do not begin implementation until all three boxes are checked.**
