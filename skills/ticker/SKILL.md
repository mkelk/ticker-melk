---
name: ticker
description: Work with Ticks issue tracker and Ticker AI agent runner. Use when managing tasks/issues with `tk` commands, running AI agents on epics, creating ticks from a SPEC.md, or working in a repo with a `.tick/` directory. Triggers on phrases like "create ticks", "tk", "run ticker", "epic", "close the task", "plan this", "break this down".
---

# Ticker & Ticks Workflow

Ticker runs AI agents (Claude Code, Codex, Gemini CLI) in continuous loops to complete coding tasks from the Ticks issue tracker.

## Skill Workflow

When invoked, follow this workflow:

### Step 0: Check Prerequisites

**1. Git repository:**
```bash
git status 2>/dev/null || git init
```

**2. GitHub remote (optional but recommended):**
```bash
gh repo view 2>/dev/null || gh repo create <name> --private --source=. --push
```
Ask user for repo name if creating new. Skip if they prefer local-only.

**3. Tools installed:**
```bash
which tk && which ticker
```

If not installed:
```bash
# Install ticks (tk CLI)
curl -fsSL https://raw.githubusercontent.com/pengelbrecht/ticks/main/scripts/install.sh | sh

# Install ticker
curl -fsSL https://raw.githubusercontent.com/pengelbrecht/ticker/main/scripts/install.sh | sh
```

**4. Ticks initialized:**
```bash
ls .tick/ 2>/dev/null || tk init
```

**5. Detect project code:**

Detect the project code from the environment. Check in this order (first match wins):

1. **Git branch:** Extract from branch name pattern
   ```bash
   git branch --show-current 2>/dev/null
   ```
   Look for pattern: `(\d{4}-\d{2}-\d{2}-[\w-]+)` (YYYY-MM-DD followed by project name)
   Example: `feature/2026-01-14-project-dim` → `2026-01-14-project-dim`

2. **Project directory:** Check if working in a project directory under `docs/projects/`
   ```bash
   pwd | grep -oE 'docs/projects/[^/]+' | sed 's|docs/projects/||'
   ```
   Example: Working in `docs/projects/2026-01-14-project-dim/` → `2026-01-14-project-dim`

   Also check for subdirectories:
   ```bash
   # If in repo root, look for project directories
   ls -d docs/projects/????-??-??-*/ 2>/dev/null | head -1 | xargs basename
   ```
   Example: `docs/projects/2026-01-14-project-dim/` → `2026-01-14-project-dim`

**Pattern regex:** `(\d{4}-\d{2}-\d{2}-[\w-]+)`

**If project detected:** Notify the user and store for use when creating ticks:
```
Detected project: 2026-01-14-project-dim
(from branch: feature/2026-01-14-project-dim)
```

**If no project detected:** Continue without a project code. When creating ticks (Step 3), ask the user if they want to:
1. Generate a new project code (recommended)
2. Continue without a project code
3. Enter a project code manually

### Step 1: Check for SPEC.md and Project Context

**1a. Check for `docs/current-setup/` directory**

Before creating specs or ticks, check if `docs/current-setup/` exists. This folder contains crucial grounding on:
- How the project is organized
- How the project is tested (test framework, commands, patterns)
- Build and development workflows
- Any project-specific conventions

```bash
ls docs/current-setup/ 2>/dev/null
```

If this folder exists, ensure you understand the files in it before proceeding. Pay special attention to testing documentation - understanding how the project is tested is critical for creating good ticks.

If this folder does NOT exist, you must explore the codebase to understand:
- Test framework (Jest, Go test, pytest, etc.)
- Test command (`npm test`, `go test ./...`, etc.)
- Test file locations (`__tests__/`, `*_test.go`, `tests/`)
- Any CI configuration that runs tests

**1b. Look for spec file**

Look for a SPEC.md (or similar spec file) in the repo root or in the project dir.

**If no spec exists:** → Go to Step 2a (Create Spec) → then Step 2c (Challenge Spec)
**If spec exists but incomplete:** → Go to Step 2b (Complete Spec) → then Step 2c (Challenge Spec)
**If spec is complete:** → Go to Step 2c (Challenge Spec) → then Step 3 (Create Ticks)

### Step 2a: Create Spec Through Conversation

Have a natural conversation with the user to understand their idea:

1. **Let them describe it** — Don't interrupt, let them explain the full vision
2. **Ask clarifying questions** — Dig into unclear areas through back-and-forth dialogue
3. **Optionally use AskUserQuestion** — For quick multiple-choice decisions
4. **Write SPEC.md** — Once you have enough detail, generate the spec

**Conversation topics to explore:**
- What problem does this solve? Who's it for?
- Core features vs nice-to-haves
- Technical constraints or preferences
- How will users interact with it?
- What does "done" look like?

**Example flow:**
```
User: "I want to build a CLI tool for managing dotfiles"

Claude: "Interesting! A few questions to understand the scope:
- Should it handle syncing between machines, or just organization?
- Are you thinking symlinks, copies, or templating?
- Any existing dotfile managers you like or want to improve on?"

User: "Syncing via git, symlinks, and I want it simpler than chezmoi"

Claude: "Got it - simple git-based syncing with symlinks. What about:
- Multi-machine configs (work vs personal)?
- Secret handling (API keys, tokens)?
- Bootstrap for new machines?"
...
```

Continue until you can write a complete SPEC.md with clear features and acceptance criteria.

### Step 2b: Complete Existing Spec

If SPEC.md exists but has gaps:

1. **Read the spec** — Identify what's missing or unclear
2. **Ask targeted questions** — Focus on the gaps, don't re-ask obvious things
3. **Update SPEC.md** — Fill in the missing details

Use AskUserQuestion for quick decisions, conversation for complex topics.

### Step 2c: Challenge the Spec

Before creating ticks, validate the spec against the codebase to catch issues early. This prevents wasted effort on specs that don't fit the existing architecture.

**1. Technical Consistency**

Check that the spec accurately references the codebase:

```bash
# Verify referenced files exist
# Check that APIs/interfaces described are accurate
# Confirm naming conventions match the codebase
```

- Do referenced files, modules, or directories actually exist?
- Are the described APIs, functions, or interfaces accurate?
- Does the spec use correct naming conventions matching the codebase?

**2. Technology Alignment**

Verify the spec fits the existing tech stack:

- Does the spec propose technologies already in use, or new ones?
- If new technologies, are there conflicts with existing choices?
- Are there existing patterns in the codebase the spec should follow?

**3. Architecture Fit**

Check the proposed design against existing architecture:

- Does the proposed design fit the existing architecture? Check against docs/current, if it exists.
- Are there existing utilities or components that should be reused?
- Does it follow established patterns in the codebase?

**4. Gap Analysis**

Identify gaps that could cause implementation uncertainty:

- Undefined edge cases or error handling
- Missing details about data flow or state management
- Unclear integration points with existing code
- Ambiguous requirements that need clarification

**5. Report Findings**

Present issues to the user in priority order:

| Priority | Action |
|----------|--------|
| **Critical** | Must fix before creating ticks — blocks implementation |
| **Warning** | Should address — could cause problems during implementation |
| **Suggestion** | Nice to have — would strengthen the spec |

**If issues found:**
1. Explain each issue clearly with specific references to code/spec
2. Propose fixes or ask clarifying questions
3. Update the spec with the user's input
4. Re-validate until clean

**If spec is valid:**
```
✓ Spec validated against codebase
  - Technical references verified
  - Technology alignment confirmed
  - Architecture fit checked
  - No blocking gaps found

Ready to create ticks.
```

Proceed to Step 3 only when the spec passes validation.

## Test-Driven Development (Critical)

**AI agents work best with test-driven tasks.** Tests provide:
- Clear acceptance criteria the agent can verify
- Immediate feedback on correctness
- Guard rails against regressions

When creating ticks, structure them for TDD:

1. **Write test first** — Each feature tick should specify expected test cases
2. **Include test commands** — Tell the agent how to run tests (`go test`, `npm test`, etc.)
3. **Define success criteria** — "Tests pass" is unambiguous; "looks good" is not

**Good tick (test-driven):**
```bash
tk create "Add email validation to registration" \
  -d "Implement email validation with test cases:
- valid@example.com → valid
- invalid@ → invalid
- @nodomain.com → invalid
- empty string → invalid

Run: go test ./internal/validation/..." \
  -acceptance "All validation tests pass, no regressions" \
  -parent <epic-id>
```

**Bad tick (no tests):**
```bash
tk create "Add email validation" -d "Make sure emails are valid"
# No acceptance criteria, no test cases - agent will guess
```

See `references/tick-patterns.md` for more TDD patterns.

### Step 3: Create Ticks from Spec

Transform the spec into ticks organized by epic.

**Apply detected project to all ticks:**

If a project code was detected in Step 0.5, apply it to all created ticks using the `-project` flag:

```bash
# Create epics with project
tk create "Authentication" -t epic -project "$DETECTED_PROJECT"

# Tasks under epic will inherit project from parent
tk create "Add JWT token generation" -parent <auth-epic>

# Or explicitly set project on tasks
tk create "Add login endpoint" -project "$DETECTED_PROJECT"
```

The project code links all related ticks together for filtering and budget tracking.

**If no project was detected in Step 0.5:**

Before creating ticks, prompt the user with `AskUserQuestion`:

```
No project code detected. Would you like to:
1. Generate a new project code (recommended)
2. Continue without a project code
3. Enter a project code manually
```

- **Option 1 (Generate new):** Create a project code using format `YYYY-MM-DD-project-name` where the date is today and project-name is derived from the spec/repo name. Apply to all ticks.
- **Option 2 (Continue without):** Create ticks without the `-project` flag. Not recommended for larger projects.
- **Option 3 (Enter manually):** Use the user-provided project code for all ticks.

**CRITICAL: First Task Must Be Environment Validation**

The very first task in ANY epic must be running existing tests and validating the environment. This ensures we start from a healthy state before making any code changes.

```bash
# Always create this as the first task, blocking all others
tk create "Run existing tests and validate environment" \
  -d "Before any code changes, verify the project is in a healthy state:
1. Run the build: [build command from docs/current-setup/ or discovered]
2. Run all tests: [test command from docs/current-setup/ or discovered]
3. Verify dev environment setup (dependencies, tools, etc.)
4. Check docs/current-setup/ for any documented known issues

BLOCKING CONDITIONS - Do not proceed if:
- Build fails
- Tests fail (unless documented as known issues in docs/current-setup/)
- Environment is misconfigured (missing deps, wrong versions, etc.)
- Any unexpected inconsistencies discovered

If problems are found, this task should EJECT with details so they can be fixed before implementation begins." \
  -acceptance "Build passes, all tests pass (or match documented known issues), environment healthy" \
  -parent <epic-id> \
  -p 0
```

All other implementation tasks should be blocked by this validation task.

**If validation fails:** The agent should signal `<promise>EJECT: [describe the problem]</promise>` so the issue can be resolved before any code changes are made. Never proceed with implementation on an unhealthy codebase.

**CRITICAL: Testing Understanding Blocker**

If you don't have sufficient understanding of how this project is tested (no `docs/current-setup/`, unclear test patterns, or unfamiliar test framework), you MUST create a blocking task:

```bash
tk create "Document testing approach and patterns" --manual \
  -d "Testing understanding is insufficient to proceed confidently.

Need to clarify:
- Test framework and how to run tests
- Test file organization and naming conventions
- How to write tests for this codebase
- Any mocking/stubbing patterns used

This blocks implementation until testing approach is clear." \
  -acceptance "Testing approach documented, patterns understood" \
  -parent <epic-id> \
  -p 1
```

This is marked `--manual` because it may require human input. It should block any tasks that involve writing code.

**For phased specs:** Focus on creating ticks for the current/next phase only. Don't create ticks for future phases—they may change based on learnings from earlier phases.

**Use AskUserQuestion** if questions arise while creating ticks:
- Unclear requirements or edge cases
- Missing acceptance criteria
- Ambiguous priorities or dependencies
- Implementation approach decisions

**Epic organization:**
1. Group related tasks into logical epics (auth, API, UI, etc.)
2. Create a **"Manual Tasks"** epic for anything requiring human intervention
3. Set up dependencies between tasks using `-blocked-by`

```bash
# Create epics (with detected project if available)
tk create "Authentication" -t epic -project "$DETECTED_PROJECT"
tk create "API Endpoints" -t epic -project "$DETECTED_PROJECT"

# Create tasks with acceptance criteria (inherit project from parent)
tk create "Add JWT token generation" \
  -d "Implement JWT signing and verification" \
  -acceptance "JWT tests pass, tokens validate correctly" \
  -parent <auth-epic>

tk create "Add login endpoint" \
  -d "POST /api/login with email/password" \
  -acceptance "Login endpoint tests pass, returns valid JWT" \
  -parent <api-epic> \
  -blocked-by <jwt-task>

# Manual tasks - use -manual flag (skipped by tk next)
tk create "Set up production database" -manual \
  -d "Create RDS instance and configure access" \
  -acceptance "Database accessible, migrations run" \
  -project "$DETECTED_PROJECT"

tk create "Create Stripe API keys" -manual \
  -d "Set up Stripe account and get API credentials" \
  -project "$DETECTED_PROJECT"
```

**Manual tasks** (use `-manual` flag):
- Setting up external services (databases, auth providers)
- Creating accounts or API keys
- Design decisions needing human judgment
- Anything requiring credentials or secrets

Manual tasks are skipped by `tk next` and ticker automation. They appear in `tk list -manual`.

### Step 3b: Guide User Through Blocking Manual Tasks

**Critical:** If manual tasks block automated tasks, guide the user through them before running ticker.

```bash
# Check for blocking manual tasks
tk list -manual
tk blocked  # See what's waiting on manual tasks
```

**When manual tasks block automation:**

1. **Identify blocking manual tasks** — Find manual tasks that other tasks depend on
2. **Guide user step-by-step** — Walk them through each manual task
3. **Verify completion** — Confirm the task is done before closing
4. **Close and unblock** — `tk close <id> "reason"` to unblock dependent tasks

**Example guidance flow:**

```
I see 2 manual tasks that block automated work:

1. **Set up PostgreSQL database** (blocks: API endpoints epic)
   - Create database instance (RDS, Supabase, or local)
   - Note the connection string
   - Run: `tk close abc "Created RDS instance, connection string in .env"`

2. **Create Stripe API keys** (blocks: payment tasks)
   - Go to dashboard.stripe.com
   - Create test API keys
   - Add to .env: STRIPE_SECRET_KEY=sk_test_...
   - Run: `tk close def "Stripe keys configured in .env"`

Once these are done, I can run the automated epics.
```

Always resolve blocking manual tasks before starting ticker, otherwise automation will stall.

### Step 4: Optimize for Parallelization

Review each epic and consider splitting if:
- Epic has many independent tasks (no dependencies between them)
- Tasks could run in parallel but are grouped together

**Split large epics:**
```
Before: "Build Dashboard" (8 independent tasks)
After:  "Build Dashboard (1/2)" (4 tasks)
        "Build Dashboard (2/2)" (4 tasks)
```

This allows ticker to run both epic halves in parallel.

**Guidelines:**
- Aim for 3-5 tasks per epic for optimal parallelization
- Keep dependent task chains in the same epic
- Independent tasks can be split across epics

### Step 5: Run Ticker

Ask the user how they want to run:

```
How would you like to run these epics?

1. Headless (I'll run ticker for you)
   - Runs in background, I'll report results

2. Interactive TUI (you run it)
   - You get real-time visibility and control
   - Command: ticker run <epic-ids...>
```

**If headless:**
```bash
ticker run <epic1> <epic2> --headless --parallel <n>
```

**If TUI:**
Provide the command for the user to run:
```bash
ticker run <epic1> <epic2>
```

## Quick Reference

### Ticks CLI (`tk`)

```bash
# Create ticks
tk create "Title" -d "Description" -acceptance "Tests pass"  # Task with acceptance criteria
tk create "Title" -t epic                                    # Create epic
tk create "Title" -parent <epic-id>                          # Task under epic
tk create "Title" -blocked-by <task-id>                      # Blocked task
tk create "Title" -manual                                    # Manual task (skipped by automation)
tk create "Title" -project "2026-01-14-project-name"            # Task with project code

# List and query
tk list                                      # All open ticks
tk list -t epic                              # Epics only
tk list -parent <epic-id>                    # Tasks in epic
tk list --project "2026-01-14-project-name"     # Filter by project
tk ready                                     # Unblocked tasks
tk ready --project "2026-01-14-project-name"    # Ready tasks for project
tk next <epic-id>                            # Next task to work on
tk next --project "2026-01-14-project-name"     # Next task for project

# Manage
tk show <id>                                 # Show details
tk close <id> "reason"                       # Close tick
tk note <id> "note text"                     # Add note
```

See `references/tk-commands.md` for full reference.

### Running Ticker

```bash
# Interactive TUI
ticker run <epic-id>                         # Single epic
ticker run <epic1> <epic2>                   # Multiple epics

# Headless
ticker run <epic-id> --headless              # Single epic
ticker run <epic1> <epic2> --headless        # Parallel epics
ticker run --auto --parallel 3               # Auto-select epics
```

## Signal Protocol

When working on a tick, signal completion with XML tags:

| Signal | Tag | When to Use |
|--------|-----|-------------|
| COMPLETE | `<promise>COMPLETE</promise>` | All work done, tests pass |
| EJECT | `<promise>EJECT: reason</promise>` | Need human help |
| BLOCKED | `<promise>BLOCKED: reason</promise>` | Missing credentials |

## Creating Good Ticks

See `references/tick-patterns.md` for detailed patterns.

**Key principles:**
1. **Atomic** — One clear deliverable per tick
2. **Testable** — Clear acceptance criteria
3. **Independent** — Minimize dependencies
4. **AI-friendly** — Include enough context for autonomous completion

**Bad tick:**
```
Title: Build the feature
```

**Good tick:**
```
Title: Add email validation to registration form
Description:
- Validate email format on blur
- Show error message below input
- Prevent form submission if invalid
- Add unit tests for validation
```
