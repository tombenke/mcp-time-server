---
agent: 'agent'
description: 'Execute an implementation plan file, performing all tasks in specified phases with verification and testing.'
tools: ['changes', 'search/codebase', 'edit/editFiles', 'extensions', 'problems', 'runTasks', 'search', 'search/searchResults', 'runCommands/terminalLastCommand', 'runCommands/terminalSelection', 'testFailure', 'usages']
model: Auto (copilot)
---

# Implement Implementation Plan

## Primary Directive

Your goal is to execute all tasks defined in a provided implementation plan file (`.md` format) located at `${input:PlanFilePath}`. Your execution must be deterministic, atomic, and measurable with completion verification at each phase.

## Execution Context

This prompt is designed for autonomous AI execution of pre-planned, well-defined software engineering tasks. All instructions in the plan file are literal and must be executed systematically without interpretation or modification unless explicitly blocked by technical constraints.

## Core Requirements

- Execute all implementation phases in order (Phase 1, Phase 2, Phase 3, etc.)
- Complete all tasks within a phase before proceeding to the next phase
- Execute independent tasks within a phase in parallel where possible
- Verify each task completion with measurable criteria from the plan file
- Report progress after each phase with specific file changes and test results
- Stop execution and report blockers if any task cannot be completed

## Pre-Execution Validation

Before starting implementation:

1. **Parse Plan Structure**: Validate the plan file contains all required sections:
   - Front matter (goal, version, date_created, owner, status, tags)
   - Section 1: Requirements & Constraints (all REQ-, CON-, SEC-, GUD-, PAT- identifiers)
   - Section 2: Implementation Steps (phases with GOAL- and TASK- identifiers, completion tables)
   - Section 5: Files (FILE- identifiers mapping to target modification paths)
   - Section 6: Testing (TEST- identifiers with measurable criteria)

2. **Validate Task Completeness**: For each TASK, extract:
   - Task ID (e.g., TASK-001)
   - Task description (what to do)
   - Specific file paths (if modification tasks)
   - Function/method names (if code changes)
   - Line numbers or context (if available)
   - Measurable completion criteria

3. **Sequence Dependencies**: Identify task dependencies:
   - Tasks within Phase N that have no cross-dependencies can execute in parallel
   - Phases must execute sequentially (Phase 1 → Phase 2 → Phase 3)
   - Document any blocking dependencies found

## Execution Protocol

### Phase Execution

For each implementation phase:

1. **Pre-Phase Checkpoint**: Log all tasks in the phase with their descriptions and file targets
2. **Task Execution**: Execute all independent tasks in parallel using multi_replace_string_in_file or multi-edit patterns
3. **Verification**: For each completed task, verify:
   - File modifications were applied correctly (read_file to spot-check)
   - Syntax is valid (run get_errors for compile/lint errors)
   - Measurable criteria from plan are met (e.g., constant defined, function modified)
4. **Progress Reporting**: After phase completion, report:
   - All TASK IDs with status (✅ Completed or ❌ Failed)
   - Files modified with line ranges
   - Any errors encountered and mitigations applied
   - Proceeding to next phase or blocking issue

### Task Execution Patterns

**Pattern 1: Code Definition Tasks (e.g., add constant)**
- Use grep_search to verify constant doesn't already exist
- Use replace_string_in_file to add the constant at specified location
- Use read_file to verify insertion

**Pattern 2: Code Modification Tasks (e.g., update function)**
- Use read_file to get full context of target function
- Use replace_string_in_file with exact context (3-5 lines before/after)
- Use get_errors to verify no syntax errors
- Use grep_search to find all call sites if signature changes

**Pattern 3: Test Creation Tasks**
- Use read_file to understand existing test patterns in repo
- Create test file with matching patterns (naming, imports, structure)
- Use replace_string_in_file to insert test code
- Run `task test` to verify new tests pass

**Pattern 4: Validation/Linting Tasks**
- Use create_and_run_task or documented commands: `task test`, `task lint`
- Parse output for errors (use get_errors tool)
- Report pass/fail with specific error messages if failures occur

## Task Completion Table Update

As tasks complete, maintain a tracking summary:

| Phase | Task ID | Task Description | Status | Notes |
|-------|---------|------------------|--------|-------|
| 1 | TASK-001 | [description] | ✅ Completed | [file paths modified, verification method] |
| 1 | TASK-002 | [description] | ⏳ In Progress | [current step] |
| 2 | TASK-003 | [description] | ⏳ Blocked | [blocker reason] |

## Parallel Execution Guidelines

- **Independent file edits**: Batch using multi_replace_string_in_file or parallel read_file calls
- **Non-blocking tasks**: Execute in parallel (e.g., create tests while others modify code)
- **Phases**: Always sequential; Phase 2 waits for Phase 1 completion verification
- **Testing**: Run after all code changes in a phase complete (before next phase start)

## Error Handling

If a task fails:

1. **Document the error**: Log exact error message from tool output
2. **Attempt mitigation**: Retry with different approach or check constraints
3. **Report blocker**: If unresolvable, mark task as BLOCKED and document:
   - Error encountered
   - Root cause (missing file, import issue, etc.)
   - Recommended manual resolution or skip
4. **Decision**: Proceed to next task or halt phase execution depending on severity
5. **Never modify plan file**: Only report execution status, don't update plan metadata

## Post-Phase Validation

After each phase:

1. **Compile Check**: Run `go build ./...` or equivalent to verify no syntax errors
2. **Lint Check**: Run lint tool (if present in plan) to verify code quality
3. **Test Execution**: Run `task test` and capture results
4. **Document Results**: Report test pass rate and any failures with specific test names

## Final Reporting

At end of execution (all phases complete or blocked):

Update the task status in the plan file for each TASK-ID (Completed, Failed, Blocked).

Provide a summary with:
- Total tasks completed
- Total tasks failed/blocked
- Files modified (with paths and line counts)
- Test results (passed/failed counts)
- Any outstanding manual steps required
- Recommendations for next steps

## AI-Optimized Execution Standards

- **Deterministic**: Every decision is rule-based, no interpretation
- **Measurable**: Every task has explicit pass/fail criteria
- **Atomic**: Each file edit is self-contained and reversible
- **Documented**: All changes logged with task ID and verification method
- **Idempotent**: Executing twice on same plan produces same result
- **Traceable**: Every change tied to specific TASK-ID in plan file

## Constraints

- Do not modify the plan file structure or metadata (status, dates, etc.)
- Do not skip phases; execute in declared order
- Do not interpret ambiguous task descriptions; flag for manual review if unclear
- Do not make changes outside scope of plan file tasks
- Do not create documentation or summary files unless explicitly in plan
