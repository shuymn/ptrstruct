# ADR-005: SuggestedFix-based autofix

## Status

Accepted

## Context

`ptrstruct` already ships through `golang.org/x/tools/go/analysis/singlechecker`, which exposes `-fix` and `-diff` driver flags. The missing piece was analyzer-provided `SuggestedFix` data.

Autofix must preserve existing diagnostics and avoid guessing about downstream semantic changes such as call-site rewrites or return-value adaptation.

## Decision

Implement autofix by attaching `analysis.SuggestedFix` values to existing diagnostics and keep the current `singlechecker` entrypoint.

Represent nested violation locations with machine-readable path steps and use them to locate the exact AST leaf to rewrite. Generate edits by prefixing `*` to the original source slice for that leaf, so aliases, selectors, and anonymous structs preserve their spelling.

## Rejected Alternatives

- Replace `singlechecker` with a custom driver and repo-wide refactor engine.
  This expands scope from local type rewrites to whole-program migration logic.
- Reconstruct edits from diagnostic message text.
  Message parsing is brittle and loses syntax fidelity for aliases and anonymous structs.

## Consequence

- `ptrstruct -fix` and `ptrstruct -fix -diff` become usable immediately through existing driver flags.
- Autofix stays local to the reported type expression, so some fixes may still require manual follow-up.
- Future repo-wide refactoring can be added later without changing the analyzer contract introduced here.

## Revisit trigger

Revisit this decision if users need automatic downstream rewrites beyond the reported type expression, or if `singlechecker` fix application cannot safely compose ptrstruct edits.
