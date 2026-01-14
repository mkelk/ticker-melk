---
description: Start a new specification document with interactive dialogue
argument-hint: [project-name]
---

## Context

- Today's date: !`date +%Y-%m-%d`
- Current directory: !`pwd`
- Current git branch: !`git branch --show-current 2>/dev/null || echo ""`
- Project name argument: $ARGUMENTS

## Your Task

You are starting a new specification document. Follow these steps:

### Step 1: Determine the Naming Pattern

**Option A: Explicit argument provided**

If `$ARGUMENTS` is not empty, use it as the project name. Generate a new naming pattern:
1. Generate a random 4-digit hex code (e.g., `a3f2`, `7b1c`)
2. Form the naming pattern: `<today's date>-<hex>-<project-name>`

**Option B: Infer from worktree**

If `$ARGUMENTS` is empty, check if we're in a worktree created by `/sdd:createworktree`:

1. Check the current directory name for the pattern: `<repo>-YYYY-MM-DD-XXXX-<project-name>`
2. Or check the git branch for the pattern: `feature/YYYY-MM-DD-XXXX-<project-name>`

If a matching pattern is found (e.g., `feature/2026-01-09-a3f2-dark-mode` or directory ending in `-2026-01-09-a3f2-dark-mode`):
- Extract the naming pattern: `2026-01-09-a3f2-dark-mode`
- Use this for the spec directory and file names
- Inform the user: "Detected worktree naming pattern: `<pattern>`. Using this for the spec."

**Option C: No argument and not in a worktree**

If `$ARGUMENTS` is empty AND no worktree pattern is detected, ask:
> What would you like to call this spec? (e.g., "dark-mode", "user-auth", "api-refactor")

Then generate a new naming pattern as in Option A.

Convert the project name to kebab-case if needed.

### Step 2: Gather Project Context

**IMPORTANT: Check for `docs/current-setup/` directory**

Before creating the spec, check if `docs/current-setup/` exists in the project root. This folder contains crucial context about:
- How the project is organized
- How the project is tested
- Build and development workflows
- Any project-specific conventions

If this folder exists:
1. Read all files in `docs/current-setup/`
2. Pay special attention to testing documentation
3. Use this context to inform the spec dialogue
4. Reference relevant conventions in the spec

If this folder does NOT exist:
1. Note this in the spec under "Project Context"
2. Explore README.md, package.json, or other config files for testing info
3. Flag that testing approach may need clarification during spec development

### Step 3: Create the Spec Directory and File

Use the naming pattern determined in Step 1 (either inferred from worktree or newly generated).

1. Create the directory: `/docs/<naming-pattern>/`
2. Create the spec file: `/docs/<naming-pattern>/<naming-pattern>-spec.md`

Example structure:
```
docs/
  2026-01-09-a3f2-dark-mode/
    2026-01-09-a3f2-dark-mode-spec.md      # The main specification
    2026-01-09-a3f2-dark-mode-research.md  # Additional documents use same prefix
    2026-01-09-a3f2-dark-mode-notes.md
```

Create the spec file with this template:

```markdown
# <Project Name> Specification

**Created:** <today's date>
**Status:** Draft

## Braindump / introductory thoughts

[To be defined]

## Overview

[To be defined]

## Project Context

**Current Setup Documentation:** [Found at `docs/current-setup/` | Not found - see notes below]

[If docs/current-setup exists, summarize key points relevant to this spec]
[If not found, note any discovered context about project organization]

## Testing Approach

**Existing Test Infrastructure:**
- Test framework: [e.g., Jest, Go test, pytest - discovered from project]
- Test command: [e.g., `npm test`, `go test ./...`]
- Test location: [e.g., `__tests__/`, `*_test.go`, `tests/`]

**Testing Strategy for This Feature:**
[How will this feature be tested? What types of tests are needed?]

**Testing Clarity Check:**
- [ ] Test infrastructure is understood
- [ ] Test patterns for this codebase are clear
- [ ] Testing approach for this feature is defined

> **BLOCKER RULE:** If any boxes above are unchecked when spec is "Ready for Implementation", the implementation plan MUST include a blocking task to clarify testing before any code changes.

## Open Questions

- [ ]

## References

[To be defined]
```

### Step 4: Start the Dialogue

After creating the file, begin a dialogue on the spec. During the dialogue:

1. **Explore testing early**: If `docs/current-setup/` was found, reference it. If not, ask clarifying questions about how the project is tested.

2. **Don't proceed without testing clarity**: The spec should not be marked "Ready for Implementation" if the testing approach is unclear. This is a blocking concern.

3. **Document testing decisions**: Any decisions about testing approach should be captured in the "Testing Approach" section.

Continue until the spec is ready.

### Step 5: Iterate until satisfactory

Continue the dialogue until the spec is good. Ensure:
- All open questions are resolved
- Testing approach is clearly defined
- Testing Clarity Check boxes are all checked

### Step 6: Remove braindump

Upon spec ready, remove the Braindump section, making sure that everything important from there is now recorded elsewhere in the spec.

### Step 7: Final Testing Check

Before marking spec as "Ready for Implementation", verify:
1. Testing Approach section is complete
2. All Testing Clarity Check boxes are checked
3. If testing clarity is insufficient, keep status as "Draft" and note: "**BLOCKER:** Testing approach must be clarified before implementation can begin."
