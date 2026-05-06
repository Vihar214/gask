# gask — Context

## What This Is

`gask` is a developer-centric CLI task manager that lives inside your project folder.
It is built around a "Raw-to-Polished" capture workflow: tasks are born as quick titles
in the terminal, then deepened into full Markdown documents inside Neovim.

The core insight: most task managers fight your terminal workflow. gask lives inside it.
One DB per project. No sync. No accounts. No internet.

---

## The Raw-to-Polished Workflow

This is the central concept of gask. Every task moves through two states of completeness:

**Raw** — A task captured in under 5 seconds. Title only, no description.
The goal is zero friction at capture time. Don't think, just log it.

**Polished** — A task with a full Markdown description written in Neovim.
This is where you think: reproduce steps, acceptance criteria, notes, links.

The workflow in practice:

```
1. gask add "fix the auth timeout bug"     → Raw task created (todo, no description)
2. gask list                               → shows ID 4 for that task
3. gask edit 4                             → Neovim opens, you write the full description
4. (task is now Polished)
5. gask status 4 doing                     → you start working on it
6. gask status 4 done                      → task complete, moves to archive view
```

A Raw task is valid and useful. Polishing is optional, not required.

---

## Core Domain Concepts

**Task**
The atomic unit of gask. Contains:

- `id` — auto-incremented integer, the only way to reference a task in mutations
- `title` — short, human-readable label (required at creation)
- `description` — full Markdown body (optional, written in Neovim)
- `status` — current lifecycle state (see Status below)
- `created_at`, `updated_at` — timestamps, managed automatically

**Status**
The lifecycle state of a Task. Four valid values:

- `todo` — default state at creation. Not started.
- `doing` — actively being worked on.
- `done` — completed.
- `blocked` — cannot progress, waiting on something external.

Status transitions: any state can move to any other state except — a `done` task
cannot go back to `todo` without explicit intent (ask the user before doing this).

**The Buffer**
The temporary Markdown file opened by Neovim during `gask add` or `gask edit`.
It is a scratch space — written to disk as a temp file, read after Neovim exits,
then immediately deleted. The Buffer never persists. Only its content (saved to DB) survives.

**The Registry**
The collection of all tasks stored in `.gask.db`. What `gask list` renders.
"Check the Registry" means query the DB. "Add to the Registry" means insert a task.

**Local DB**
The `.gask.db` SQLite file in the current working directory.
It is project-scoped by design — tasks created inside `/my-project/` live there, not globally.
There is no global task list unless the user explicitly runs gask from their home directory.

**The Editor**
Always refers to the program set in `$EDITOR` (expected: `nvim`).
gask does not provide a text editing experience — it delegates entirely to the user's editor.

---

## Key Relationships

- A Task is always identified by its **integer ID** in commands — never by title
- A Task's description lives in the DB, not in the filesystem (the Buffer is temporary)
- The Local DB is tied to the **directory gask is run from**, not where the binary lives
- Status and ID are the only two fields that change after a task is created
- `done` tasks stay in the DB permanently — they are never deleted, only filtered in list view

---

## Naming Conventions

| Term in conversation | What it maps to in code                         |
| -------------------- | ----------------------------------------------- |
| task                 | `ent/schema/task.go` → `Task` struct            |
| registry             | the SQLite DB / `internal/db/repository.go`     |
| buffer               | the temp `*.md` file in `editor.go`             |
| polish / polishing   | calling `gask edit <id>` to write a description |
| status update        | calling `gask status <id> <status>`             |
| local DB             | `.gask.db` in the working directory             |

Always use these terms consistently in code, comments, error messages, and variable names.
A variable holding a task description should be `description`, not `body`, `content`, or `notes`.

---

## What This Is NOT

- **Not a project management suite** — no Kanban boards, sprints, priorities, or tags
- **Not a cloud app** — no sync, no accounts, no API keys, no internet required ever
- **Not a global to-do list** — tasks are scoped to the directory, not the whole machine
- **Not a text editor** — gask opens Neovim and gets out of the way; it never renders Markdown
- **Not a team tool** — single user, single machine, no sharing or permissions model
