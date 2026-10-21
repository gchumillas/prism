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
