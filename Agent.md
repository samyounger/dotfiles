# Agent Guide

This repository is a small Go-based dotfiles/bootstrap project. Changes should stay predictable, easy to review, and well covered by unit tests.

## Preferred workflow

Follow this order when making changes:

1. Read the existing code and tests in the package you are about to change.
2. Write or update unit tests first.
3. Run the new tests and confirm they fail for the expected reason.
4. Implement the code change.
5. Re-run the focused tests for the touched package until they pass.
6. Run the full repository checks.
7. Run existing lint commands, if the repository defines them.

In short: **tests first, implementation second, verification last**.

## Repository-specific testing guidance

The codebase already uses a few patterns that make tests straightforward:

- Prefer dependency injection through package-level variables for OS-facing behavior.
  - Examples in the repo include `setupFunc`, `installFunc`, `setupMachineFunc`, `runCommand`, `runCommandOutput`, `execCommand`, and `execLookPath`.
- Prefer `io.Reader` and `io.Writer` parameters for interactive flows so tests can use `strings.NewReader` and `bytes.Buffer`.
- Reset mocks with dedicated helper functions after each test when overriding package-level variables.
- Keep tests deterministic. Avoid relying on real machine state, real shell commands, or the real home directory unless the test is explicitly designed for that.

## Recommended command sequence

When working on a change, prefer this validation flow:

```bash
go test ./...
```

If you are iterating in one package, start with a focused test run and then finish with the full suite:

```bash
go test ./install/...
go test ./...
```

If the repository later adds a lint target, run it after tests pass. Use the lint command that already exists in the repository rather than inventing a new tool or config.

## Implementation guidance

- Keep changes surgical and match existing package boundaries.
- Put command construction in helper/build functions when possible, and keep side-effect execution thin.
- Return or surface contextual errors instead of swallowing them.
- Preserve current output wording unless the change intentionally updates user-facing text.
- Update docs when behavior, workflows, or repository structure changes.

## Working with command execution

Much of the repository configures a machine by building shell commands and then executing them. Prefer these patterns:

- Build commands as data where possible.
- Keep OS calls behind a small seam that tests can replace.
- When a command failure is expected to be non-fatal, make that explicit in the code and in the output.
- When a command is required, return an error with command context.

## Before finishing

Before considering a task complete:

1. Ensure the relevant unit tests exist and pass.
2. Run `go test ./...`.
3. Run any existing lint command if one is available.
4. Check whether `README.md`, `Architecture.md`, or package-level docs should be updated.
5. Avoid reverting unrelated user changes in the worktree.

## Documentation to keep in sync

- `Architecture.md` explains package layout, runtime flow, and coding conventions.
- `README.md` is the short onboarding document for the repo.

