---
name: suggest-commit
description: Suggest a title for the latest commit on the current branch
---

Suggest one or more titles without providing explanations for the most recent commit on the current branch.

The title should follow the convention below whenever possible:

| Type       | Description                                  |
| ---------- | -------------------------------------------- |
| `feat`     | New functionality                            |
| `fix`      | Bug fix                                      |
| `chore`    | Maintenance (configs, dependencies, tooling) |
| `docs`     | Documentation                                |
| `test`     | Tests                                        |
| `refactor` | Refactoring without behavior change          |
| `perf`     | Performance improvements                     |

**Examples:**

```
feat(lexer): add token consumer for strings
fix(token): correct position tracking
chore: update linter config
docs: add contributing guide
```