## [0.1.0] - 2026-03-24

### 🚀 Features

- Implement ptrstruct analyzer
- Add stdlib/third-party exemption options
- Add autofix for struct-by-value violations
- *(cmd)* Migrate from singlechecker to checker
- Add type-sig checks with opposite defaults
- Add structpolicy-tune skill

### 🐛 Bug Fixes

- *(nolint)* Pass matcher by pointer

### 💼 Other

- Add valuestruct build target

### 🚜 Refactor

- *(walk)* Extract walker for cycle detection
- *(config)* Default to receiver-only checks
- Reorganize into internal/pkg packages
- Replace *T returns with (T, bool)
- Unexport internal package identifiers

### 📚 Documentation

- Update README for ptrstruct
- Document -fix flags and limitations
- Add ADR and plan for autofix
- Drop autofix documentation
- Document multi-analyzer structure
- Add ADR-005 for performance defaults
- Update default profile documentation
- Add skills reference to README
- Clarify guidelines for struct return types and field pointers in SKILL.md
- Enhance guidelines for pointer semantics and slice element types in SKILL.md

### 🧪 Testing

- *(testdata)* Move want comment to correct line
- Expand allChecksAnalyzer to cover all flags
- Add typesigflags test cases

### ⚙️ Miscellaneous Tasks

- Initialize from template
- Remove SuggestedFix autofix implementation
