---
name: suggest-issue-title
description: Suggest a GitHub issue title from a failing test or bug description
---

Given a failing test, code snippet, or bug description provided by the user, suggest one or more GitHub issue titles without providing explanations.

Titles should be:

- Written in English.
- Concise and descriptive of the observed (wrong) behavior, not the solution.
- Short: ideally under 70 characters, never more than one line.
- Free of test names or issue numbers.

**Examples:**

```
Shared Button only clickable at last rendered position
HStack.Remove panics when the widget is not comparable
```
