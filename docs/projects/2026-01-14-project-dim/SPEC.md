# Project Dimension for Ticks

**Created:** 2026-01-14
**Status:** Draft
**Project Code:** 2026-01-14-6453-project-dim

## Overview

Add a first-class `project` field to ticks that links tasks and epics to a project code. This enables:
- Grouping related work across epics
- Filtering and querying by project
- Budget tracking per project
- Auto-detection from git worktrees/branches

The project code format follows the SDD naming convention: `YYYY-MM-DD-XXXX-name` (e.g., `2026-01-09-a3f2-dark-mode`).

## Background

The SDD (Spec-Driven Development) workflow uses a naming pattern to tie together:
- Spec directories (`docs/2026-01-09-a3f2-dark-mode/`)
- Git branches (`feature/2026-01-09-a3f2-dark-mode`)
- Git worktrees (`repo-2026-01-09-a3f2-dark-mode`)

Currently, ticks have no way to link to this project context. This spec adds that connection.

## Components Affected

1. **Ticks CLI (`tk`)** - Add project field to schema and commands
2. **Ticker Go Program** - Add project filtering and budget tracking
3. **Ticker Skill** - Auto-detect and apply project codes

---

## Part 1: Ticks CLI Changes

### 1.1 Schema Changes

Add `project` field to tick JSON schema:

```json
{
  "id": "abc123",
  "title": "Add JWT validation",
  "type": "task",
  "status": "open",
  "project": "2026-01-09-a3f2-dark-mode",
  ...
}
```

**Field properties:**
- Optional string field
- No format validation (any string allowed)
- Recommended format: `YYYY-MM-DD-XXXX-name`

### 1.2 Command Changes

#### `tk create`

Add `-project` flag:

```bash
tk create "Title" -project "2026-01-14-6453-project-dim"
```

| Flag | Description |
|------|-------------|
| `-project` | Project code to associate with this tick |

#### `tk list`

Add `--project` filter:

```bash
tk list --project "2026-01-14-6453-project-dim"
tk list -t epic --project "..."
tk list -s open --project "..."
```

#### `tk ready`

Add `--project` filter:

```bash
tk ready --project "2026-01-14-6453-project-dim"
```

#### `tk next`

Add `--project` filter:

```bash
tk next --project "2026-01-14-6453-project-dim"
tk next <epic-id>  # Works as before, project from epic context
```

#### `tk update`

Add `--project` flag:

```bash
tk update <id> --project "2026-01-14-6453-project-dim"
tk update <id> --project ""  # Clear project
```

#### `tk show`

Display project in output:

```
ID:       abc123
Title:    Add JWT validation
Type:     task
Status:   open
Project:  2026-01-14-6453-project-dim
...
```

### 1.3 Inheritance Behavior

When creating a task with `-parent`:
- If `-project` is specified, use it
- If `-project` is not specified, inherit from parent epic
- If parent has no project, task has no project

```bash
# Epic with project
tk create "Auth System" -t epic -project "2026-01-14-6453-auth"

# Task inherits project from parent
tk create "Add JWT" -parent <epic-id>
# Result: task.project = "2026-01-14-6453-auth"

# Task overrides project
tk create "Add JWT" -parent <epic-id> -project "other-project"
# Result: task.project = "other-project"
```

### 1.4 JSON Output

Include project in `--json` output:

```bash
tk list --json
```

```json
{
  "ticks": [
    {
      "id": "abc123",
      "title": "Add JWT",
      "project": "2026-01-14-6453-auth",
      ...
    }
  ]
}
```

---

## Part 2: Ticker Go Program Changes

### 2.1 CLI Flag

Add `--project` flag to `ticker run`:

```bash
ticker run --project "2026-01-14-6453-project-dim"
ticker run --auto --project "..."
ticker run <epic-id> --project "..."  # Validates epic matches
```

**Behavior:**
- Filters epics/tasks to only those matching project
- With `--auto`, only considers epics with matching project
- With explicit epic ID, validates epic has matching project (warn if mismatch)

### 2.2 Budget Tracking Per Project

Track and report costs aggregated by project:

```go
type ProjectBudget struct {
    Project    string
    Iterations int
    Tokens     int64
    Cost       float64
}
```

**Storage:** Add to checkpoint data or separate file (`.ticker/projects.json`).

**Reporting:**
```
Project: 2026-01-14-6453-auth
  Iterations: 47
  Tokens: 1.2M
  Cost: $12.34
```

### 2.3 TUI Enhancements

**Header/Status bar:**
- Show current project if `--project` flag used
- Show project of selected epic

**Task list:**
- Optionally show project column
- Filter indicator when project filter active

### 2.4 Headless Output

Include project in headless output:

```
[2026-01-14 10:30:00] Project: 2026-01-14-6453-auth
[2026-01-14 10:30:00] Starting epic: abc123 (Auth System)
...
[2026-01-14 10:45:00] Project Summary:
  Iterations: 12
  Cost: $3.45
```

---

## Part 3: Ticker Skill Changes

### 3.1 Auto-Detection

When skill is invoked, detect project code from environment:

**Detection order:**
1. Git branch: `feature/YYYY-MM-DD-XXXX-name` → extract `YYYY-MM-DD-XXXX-name`
2. Directory name: `repo-YYYY-MM-DD-XXXX-name` → extract `YYYY-MM-DD-XXXX-name`
3. Spec directory: `docs/YYYY-MM-DD-XXXX-name/` → extract from path

**Pattern regex:** `(\d{4}-\d{2}-\d{2}-[0-9a-f]{4}-[\w-]+)`

**User notification:**
```
Detected project: 2026-01-14-6453-project-dim
(from branch: feature/2026-01-14-6453-project-dim)
```

### 3.2 Apply to Created Ticks

When creating ticks (Step 3 of skill workflow):

```bash
# Apply detected project to all created ticks
tk create "Auth System" -t epic -project "$DETECTED_PROJECT"
tk create "Add JWT" -parent <epic> -project "$DETECTED_PROJECT"
```

If no project detected, ask user:
```
No project code detected. Would you like to:
1. Generate a new project code (recommended)
2. Continue without a project code
3. Enter a project code manually
```

### 3.3 Skill Reference Updates

Update `references/tk-commands.md` to document `-project` flag.

Update `SKILL.md` Step 3 to include project handling.

---

## Testing Approach

### Ticks CLI Tests

```bash
# Test create with project
tk create "Test" -project "test-project"
tk show <id> | grep "Project: test-project"

# Test list filter
tk create "A" -project "proj-a"
tk create "B" -project "proj-b"
tk list --project "proj-a"  # Should only show A

# Test inheritance
tk create "Epic" -t epic -project "proj-x"
tk create "Task" -parent <epic-id>
tk show <task-id> | grep "Project: proj-x"

# Test ready/next with project
tk ready --project "proj-a"
tk next --project "proj-a"
```

### Ticker Go Program Tests

```go
// Test project filter
func TestRunWithProjectFilter(t *testing.T) {
    // Create epics with different projects
    // Run with --project flag
    // Verify only matching epics are processed
}

// Test budget tracking
func TestProjectBudgetTracking(t *testing.T) {
    // Run tasks with project
    // Verify budget aggregated by project
}
```

### Ticker Skill Tests

Manual testing:
1. Create worktree with project naming
2. Invoke skill, verify project detected
3. Create ticks, verify project applied
4. Run ticker with --project filter

---

## Migration

**Existing ticks without project field:**
- Treat as `project: null` (no project)
- No migration script needed
- `--project` filter excludes ticks without project

---

## Open Questions

- [x] Project code format: Allow any string (no validation)
- [x] Inheritance: Tasks inherit from parent epic if not specified
- [x] Migration: No migration needed, null is valid

---

## Acceptance Criteria

1. [ ] `tk create -project` adds project to tick
2. [ ] `tk list --project` filters by project
3. [ ] `tk ready --project` filters ready tasks by project
4. [ ] `tk next --project` filters next task by project
5. [ ] `tk update --project` updates tick's project
6. [ ] `tk show` displays project field
7. [ ] Tasks inherit project from parent epic
8. [ ] `ticker run --project` filters by project
9. [ ] Ticker tracks budget per project
10. [ ] Ticker skill auto-detects project from environment
11. [ ] Ticker skill applies project to created ticks

---

## Implementation Order

1. **Ticks CLI** - Foundation (must be first)
   - Schema change
   - Create/update flags
   - List/ready/next filters
   - Inheritance logic

2. **Ticker Go Program** - Depends on ticks CLI
   - `--project` flag
   - Budget tracking
   - TUI display

3. **Ticker Skill** - Depends on both above
   - Auto-detection logic
   - Workflow updates
   - Reference docs
