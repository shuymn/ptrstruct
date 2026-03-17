# SuggestedFix Autofix

## Goal

Enable `ptrstruct -fix` / `-fix -diff` without replacing the existing `singlechecker` driver.

## Constraints

- Preserve current diagnostics, suppression behavior, and analyzer flags.
- Limit edits to the violating type leaf so fixes stay local and mergeable.
- Do not attempt repo-wide refactors such as call-site or return-expression rewrites.

## Core Design

- Store machine-readable violation path steps alongside the human-readable `Violation.Path`.
- Build `analysis.SuggestedFix` from the AST node already associated with each diagnostic site.
- Traverse the diagnostic type expression by violation path and replace only the leaf range with `*` plus the original source text.
- Keep diagnostics when a fix cannot be derived safely from the AST shape.

## Acceptance

- `analysistest.RunWithSuggestedFixes` passes for direct, nested, alias, embedded, and optional-flag container cases.
- CLI replay proves `ptrstruct -fix -diff` prints a patch and `ptrstruct -fix` rewrites files.
- `task check` passes with no regression in existing diagnostics.
