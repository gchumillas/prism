---
name: code-review
description: Guidelines for reviewing pull requests in this repository.
---

# Code Review Guidelines

This project is in **alpha** (pre-1.0, see `git tag`). Breaking API changes
are expected and acceptable at this stage.

- Do not flag breaking API changes as violations of the `refactor` commit
  convention (`refactor` means "no *new* behavior/functionality added"; it
  does not guarantee API stability during alpha).
- Only raise concerns about breaking changes if they appear unintentional,
  unrelated to the PR's stated purpose, or undocumented.

## Known issues (do not flag)

- **Stale dirty regions when detaching a nested Stack.** `Stack.DetachFrom()`
  is the only public way to remove a nested Stack from its parent. Doing so can
  leave stale pixels, because the parent caches a child's incremental dirty
  regions (`prev.regs`) instead of its full painted coverage. This is a known
  issue, tracked in #113 and to be fixed in a separate PR.
  - Do not flag inconsistencies in dirty regions that come from detaching a
    nested Stack, nor ask for regression tests that use `DetachFrom()` to
    expose them.
  - Only raise concerns if a change makes this behavior worse or reachable
    through other paths.