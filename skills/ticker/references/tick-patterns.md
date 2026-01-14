# Tick Patterns

Patterns for creating effective ticks that AI agents can complete autonomously.

## Before Creating Ticks: Check Project Context

**Always check for `docs/current-setup/` first.** This folder contains crucial grounding on:

- How the project is organized
- How the project is tested (framework, commands, patterns)
- Build and development workflows
- Project-specific conventions

```bash
ls docs/current-setup/ 2>/dev/null
```

If this folder exists, read all files before creating ticks. The testing documentation is especially important for writing good acceptance criteria.

If this folder doesn't exist, explore the codebase to understand:
- Test framework (Jest, Go test, pytest, etc.)
- Test command (`npm test`, `go test ./...`, etc.)
- Test file locations (`__tests__/`, `*_test.go`, `tests/`)

## The Ideal Tick

A well-formed tick has:

1. **Clear title** — Action verb + specific target
2. **Context** — What exists, what's needed
3. **Acceptance criteria** — How to verify done
4. **Bounded scope** — Completable in 1-3 iterations

## Tick Sizing

### Too Small
```
Title: Add semicolon to line 42
```
Waste of overhead. Fix inline.

### Too Large
```
Title: Build complete user management system
```
Break into epic with tasks.

### Just Right
```
Title: Add email validation to registration form
Description:
- Validate email format on blur
- Show error message below input
- Prevent form submission if invalid
- Add unit tests for validation
```

## Pattern: Environment Validation (Required First Task)

**Every epic MUST start with this task.** It ensures we begin from a healthy baseline.

```bash
tk create "Run existing tests and validate environment" \
  -d "Before any code changes, verify the project is in a healthy state:

1. Run the build: [build command]
2. Run all tests: [test command]
3. Verify dev environment setup (dependencies, tools, etc.)
4. Check docs/current-setup/ for documented known issues

BLOCKING CONDITIONS - Do not proceed if:
- Build fails
- Tests fail (unless documented as known issues in docs/current-setup/)
- Environment is misconfigured (missing deps, wrong versions, etc.)
- Any unexpected inconsistencies discovered

If problems are found, EJECT with details so they can be fixed first." \
  -acceptance "Build passes, all tests pass (or match documented known issues), environment healthy" \
  -parent <epic-id> \
  -p 0
```

**Why this matters:**
- Establishes a known-good baseline before changes
- Surfaces environment issues early (not mid-implementation)
- Prevents implementing on a broken codebase (which wastes effort)
- Ensures all dependencies and tools are properly configured

**If validation fails:** Signal `<promise>EJECT: [describe the problem]</promise>`. Never proceed with implementation on an unhealthy codebase. Problems must be fixed before any code changes.

All other implementation tasks should use `-blocked-by <this-task-id>`.

## Pattern: Testing Blocker

When testing approach is unclear, create a blocking manual task to prevent premature implementation.

```bash
tk create "Document testing approach and patterns" --manual \
  -d "Testing understanding is insufficient to proceed confidently.

Need to clarify:
- Test framework and how to run tests
- Test file organization and naming conventions
- How to write tests for this codebase
- Any mocking/stubbing patterns used
- Integration vs unit test patterns

Check docs/current-setup/ if it exists.
This blocks implementation until testing approach is clear." \
  -acceptance "Testing approach documented, patterns understood" \
  -parent <epic-id> \
  -p 1
```

**When to use:**
- No `docs/current-setup/` folder exists
- Unfamiliar test framework
- Unclear test patterns in existing code
- Complex testing setup (mocking, fixtures, etc.)

This is marked `--manual` because it may require human input to clarify.

## Pattern: Bug Fix

```
Title: Fix [specific symptom]

Description:
Current behavior: [what happens now]
Expected behavior: [what should happen]
Reproduction: [steps to reproduce]

Files likely involved: [paths if known]
```

## Pattern: Feature Addition

```
Title: Add [feature name] to [component]

Description:
User story: As a [user], I want [action] so that [benefit]

Requirements:
- [Requirement 1]
- [Requirement 2]

Acceptance criteria:
- [ ] [Testable criterion 1]
- [ ] [Testable criterion 2]
```

## Pattern: Refactor

```
Title: Refactor [component] to [goal]

Description:
Current state: [what's wrong/suboptimal]
Target state: [desired architecture]

Constraints:
- Must maintain backward compatibility
- No behavior changes
- Tests must pass
```

## Pattern: Test Addition

```
Title: Add tests for [component/function]

Description:
Test cases needed:
- [ ] Happy path: [scenario]
- [ ] Edge case: [scenario]
- [ ] Error case: [scenario]

Coverage target: [percentage or specific paths]
```

## Test-Driven Development (TDD)

**Critical for AI agent success.** Tests give agents:
- Unambiguous success criteria
- Immediate feedback loop
- Regression protection

### TDD Tick Pattern

```bash
tk create "Add [feature]" \
  -d "Implement [feature] with test cases:
- Input: [x] → Expected: [y]
- Input: [a] → Expected: [b]
- Edge case: [condition] → Expected: [behavior]
- Error case: [bad input] → Expected: [error handling]

Run: [test command]" \
  -acceptance "All tests pass, no regressions"
```

### TDD Feature Example

```bash
tk create "Add password strength validator" \
  -d "Implement password validation with scoring:

Test cases:
- \"abc\" → score 1 (weak), reasons: [\"too short\", \"no numbers\"]
- \"abc12345\" → score 2 (medium), reasons: [\"no special chars\"]
- \"Abc123!@#\" → score 3 (strong), reasons: []
- \"\" → error: \"password required\"

Run: go test ./internal/auth/... -v" \
  -acceptance "All password tests pass, validator integrated in registration"
```

### TDD Bug Fix Example

```bash
tk create "Fix email parsing for plus addresses" \
  -d "Plus addresses (user+tag@domain.com) are rejected incorrectly.

Test cases to add:
- \"user+newsletter@gmail.com\" → valid
- \"user+shop@example.org\" → valid
- \"user++double@test.com\" → valid

Current failing: Returns \"invalid email format\"
Expected: All plus addresses should validate

Run: npm test -- --grep \"email\"" \
  -acceptance "New plus-address tests pass, existing email tests pass"
```

### Why TDD Matters for Agents

1. **Clear completion signal** — "Tests pass" vs "looks right"
2. **Prevents scope creep** — Agent knows exactly what to implement
3. **Catches regressions** — Agent can verify it didn't break other code
4. **Self-documenting** — Tests show intended behavior

### Anti-Pattern: No Test Criteria

Bad:
```
Title: Add input validation
Description: Validate user inputs appropriately
```

The agent has no way to verify "appropriately" — it will guess and may be wrong.

## Anti-Patterns

### Vague Titles
- Bad: "Improve performance"
- Good: "Add database index for user lookup query"

### Missing Context
- Bad: "Fix the bug in auth"
- Good: "Fix OAuth callback failing when user has no email"

### Unbounded Scope
- Bad: "Make the app better"
- Good: "Add loading spinner to dashboard data fetch"

### Implicit Dependencies
- Bad: Create tasks without explicit blockers
- Good: Use `-blocked-by` to make order clear

## Epic Structure

Group related tasks under an epic:

```bash
# Create epic
tk create "Search Feature" -t epic -d "Full-text search for documents"

# Create tasks with dependencies
tk create "Add search index schema" -parent <epic>
tk create "Implement indexing service" -parent <epic> -blocked-by <schema-task>
tk create "Add search API endpoint" -parent <epic> -blocked-by <indexing-task>
tk create "Add search UI component" -parent <epic> -blocked-by <api-task>
```

## Priority Guidelines

| Priority | Use For |
|----------|---------|
| P0 Critical | Production down, security issues |
| P1 High | Blocking other work, user-facing bugs |
| P2 Medium | Normal feature work (default) |
| P3 Low | Nice-to-have, minor improvements |
| P4 Backlog | Future consideration |
