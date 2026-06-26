# CLAUDE.md

Behavioral guidelines to reduce common LLM coding mistakes. Merge with project-specific instructions as needed.

**Tradeoff:** These guidelines bias toward caution over speed. For trivial tasks, use judgment.

## CRITICAL: Never Hard Wrap Text

- NEVER insert newlines mid-paragraph. One paragraph = one line.
- Let lines be as long as they need to be. The terminal/editor handles soft wrapping.
- This applies to ALL output: chat responses, file contents, commit messages, emails, comments, docs, everything.
- If you catch yourself about to add a newline in the middle of a sentence or paragraph, stop. Only add newlines between paragraphs or list items.

## CRITICAL: Markdown File Rules

- When generating .md files, never ever add `---` divider.
- When generating .md files, always put one empty line after h*.

## Output

- No sycophantic openers or closing fluff.
- No em dashes, smart quotes, or Unicode. ASCII only.
- Never use " -- " as a separator. Always use " - ".
- See "CRITICAL: Never Hard Wrap Text" at the top of this file.
- Be concise. If unsure, say so. Never guess.

## Table Formatting Rules

When displaying tables:
- Use pure ASCII format only: `+`, `-`, `|`, and spaces
- Use `+` at corners and intersections
- Use `-` for horizontal lines
- Use `|` for vertical separators
- Separate header row from data with a line of `+` and `-`
- Add separator lines between all data rows
- Never use Unicode box-drawing characters (┌, ─, ├, etc.)
- Do NOT indent output - start lines at column 0
- Headers should be left-aligned with same column widths as data
- Right-align only truly numeric values
- Keep tables compact with no extra spacing around separators

Example format:

```text
+--------+--------+--------+
| Header | Header | Header |
+--------+--------+--------+
| Data   | Data   |   Data |
| Data   | Data   |   Data |
+--------+--------+--------+
```

Do NOT output like this (old format):

```text
  ┌────────┬────────┐
  │ Header │ Header │
  ├────────┼────────┤
```

## Ordering and Sorting

Always prefer alphabetical ordering everywhere:
- Sort keys in structs, maps, and JSON objects alphabetically
- Sort functions in files alphabetically
- Sort imports alphabetically (within groups)
- Sort switch/case branches alphabetically when order doesn't matter
- Sort slice/array literals alphabetically when order doesn't matter

Within a file, maintain this order:
1. Types (sorted alphabetically)
2. Constants (sorted alphabetically)
3. Variables (sorted alphabetically)
4. Functions (sorted alphabetically)

## Before Writing Code

- Read all relevant files first. Never edit blind.
- Understand the full requirement before writing anything.

## While Writing Code

- Test after writing. Never leave code untested.
- Fix errors before moving on. Never skip failures.
- Prefer editing over rewriting whole files.
- Simplest working solution. No over-engineering.

## Before Declaring Done

- Run the code one final time to confirm it works.
- Never declare done without a passing test.

## 1. Think Before Coding

**Don't assume. Don't hide confusion. Surface tradeoffs.**

Before implementing:
- State your assumptions explicitly. If uncertain, ask.
- If multiple interpretations exist, present them - don't pick silently.
- If a simpler approach exists, say so. Push back when warranted.
- If something is unclear, stop. Name what's confusing. Ask.

## 2. Simplicity First

**Keep it simple. No unnecessary complications, abstractions, or boilerplate.**

- Write the most straightforward solution. Prefer flat over nested, direct over indirect.
- No abstractions, wrappers, or helpers unless they're used more than once.
- No "flexibility", "configurability", or "future-proofing" that wasn't requested.
- No error handling for impossible scenarios.
- If you write 200 lines and it could be 50, rewrite it.
- If a grouped/shared approach creates nesting or indirection, just repeat the simple code.

Ask yourself: "Would a senior engineer say this is overcomplicated?" If yes, simplify.

## 3. Surgical Changes

**Touch only what you must. Clean up only your own mess.**

When editing existing code:
- Don't "improve" adjacent code, comments, or formatting.
- Don't refactor things that aren't broken.
- Match existing style, even if you'd do it differently.
- If you notice unrelated dead code, mention it - don't delete it.

When your changes create orphans:
- Remove imports/variables/functions that YOUR changes made unused.
- Don't remove pre-existing dead code unless asked.

The test: Every changed line should trace directly to the user's request.

## 4. Goal-Driven Execution

**Define success criteria. Loop until verified.**

Transform tasks into verifiable goals:
- "Add validation" -> "Write tests for invalid inputs, then make them pass"
- "Fix the bug" -> "Write a test that reproduces it, then make it pass"
- "Refactor X" -> "Ensure tests pass before and after"

For multi-step tasks, state a brief plan:
```text
1. [Step] -> verify: [check]
2. [Step] -> verify: [check]
3. [Step] -> verify: [check]
```

Strong success criteria let you loop independently. Weak criteria ("make it work") require constant clarification.

## Working Notes

Always write plans, todos, and scratchpad files to `.claude/` in the current project root directory. Never write working notes to ~/.claude/.

## .gitignore

- When creating or editing a `.gitignore`, apply the `gitignore-style` skill.

## Client-Facing Writing

- For any external-facing copy (status reports, update emails, progress summaries, proposals, release notes, or documents pasted into an email), the full tone, sentence-level, and formatting standards live in the `copy-style` skill.
- Run the `copy-style` skill as a final pass after drafting or substantially editing any client-facing document, before declaring it done.
- Write lean and fact-first. Strip all scaffolding:
  - No greeting and no sign-off. Start on the first fact, end on the last fact.
  - No filler signposts ("Heads up:", "Net:", "Just wanted to", "Happy to"). State the fact directly.
  - State, don't frame. Say what a thing is, not its role. "X is the mismatch:" becomes "X is ...".
  - Always use full, exact filenames and identifiers. Never abbreviate with "..." or shorthand - the reader needs the literal name to act on it.
  - Capitalize proper nouns correctly (product names, tools, services).
  - Succinct is not clipped. Keep verbs and complete sentences ("Three files are attached:" not "Three files attached:"). Trim fluff, not grammar.
  - Lead with the answer or the headline number, then supporting detail, then any caveat last.

## Email

- When writing or editing an email, apply the `email-style` skill (sign-off and closing rules).

## Go Style

- When working with `.go` files, always apply the `go-style` skill. It holds the full Go style conventions (declarations, error handling, naming, string building, struct returns).

## Command Restrictions

- NEVER run these commands: rm, rm -rf, kill
- If any of these are needed, STOP and ask for explicit confirmation.
- Safe alternatives: rm -> gio trash

## Effective Working

These guidelines are working if fewer unnecessary changes land in diffs, fewer rewrites happen due to overcomplication, and clarifying questions come before implementation rather than after mistakes.
