# Aura Development Guide

This document covers everything you need to know to contribute to the Aura toolchain.

> 🤖 **Aura is an AI-first language.** All development decisions should be evaluated through the lens of AI-human collaboration. See [AI_MISSION.md](AI_MISSION.md) for the full mission statement.

---

## Designing for AI Developers

Aura's primary audience is **AI agents generating and reasoning about code**, with human developers as reviewers and collaborators. Every feature, API, and error message should be designed with this in mind.

### Design Decision Framework

When faced with trade-offs, apply this priority order:

1. **AI flow** — Does this make AI code generation faster and more accurate?
2. **Compiler verifiability** — Can the compiler check this automatically?
3. **Human readability** — Is this clear for human review?
4. **Brevity** — Is this concise? (Lowest priority — clarity always wins over conciseness)

### Code Review Checklist: AI-First Design

When reviewing PRs or designing features, ask these questions:

- [ ] **Does this feature help AI understand intent?** — Can an AI read the syntax/output and know exactly what to do without surrounding context?
- [ ] **Is the representation structured?** — Prefer structured data (spec blocks, typed annotations) over freeform text (comments, naming conventions).
- [ ] **Are error messages machine-parseable?** — Error output should include structured fields (error code, location, expected vs actual) that AI agents can parse and act on automatically.
- [ ] **Does this integrate with specs?** — Every new feature should consider how it interacts with the specification system. Can specs reference it? Can the compiler validate it?
- [ ] **Are effects explicit?** — If a feature introduces side effects, are they tracked in the effect system?
- [ ] **Is it deterministic?** — Given the same input, does the feature always produce the same output? AI agents depend on deterministic behavior.

### Testing: AI Code Generation Scenarios

When writing tests for new features, include scenarios that validate AI-relevant use cases:

- **Spec-to-implementation validation** — Test that code satisfying a spec actually passes all spec checks.
- **Round-trip stability** — AI-generated code, when formatted, should be identical to human-written canonical form.
- **Error message quality** — Test that error messages include enough information for an AI to fix the issue automatically (error code, location, suggestion).
- **Effect tracking accuracy** — Test that the effect system correctly identifies all effects, especially for complex call graphs that AI might generate.
- **Edge cases from AI generation** — AI may produce valid but unusual code patterns. Test that these are handled correctly (e.g., deeply nested expressions, max-length identifiers, unusual but valid type combinations).

---

## Architecture Overview

### Pipeline

Aura source code flows through the toolchain in stages:

```
Source (.aura)
    │
    ▼
┌──────────┐     ┌──────────┐     ┌────────────┐     ┌───────────┐
│  Lexer   │────▶│  Parser  │────▶│  Semantic  │────▶│  CodeGen  │
│ (tokens) │     │  (AST)   │     │  Analysis  │     │ (output)  │
└──────────┘     └──────────┘     └────────────┘     └───────────┘
    │                 │                 │                   │
    ▼                 ▼                 ▼                   ▼
  Token stream    Raw AST         Typed AST           Executable
                      │
                      ▼
                ┌────────────┐
                │ Formatter  │
                │ (source)   │
                └────────────┘
```

### Package Layout

```
aura-toolchain/
├── cmd/
│   └── aura/
│       └── main.go              # CLI entry point
├── pkg/
│   ├── token/
│   │   └── token.go             # Token types, positions, spans
│   ├── lexer/
│   │   ├── lexer.go             # Indentation-sensitive lexer
│   │   └── lexer_test.go        # 11 tests
│   ├── ast/
│   │   └── ast.go               # Complete AST node definitions
│   ├── parser/
│   │   ├── parser.go            # Recursive descent parser
│   │   └── parser_test.go       # 16 tests
│   ├── formatter/
│   │   ├── formatter.go         # AST → canonical source
│   │   └── formatter_test.go    # 9 tests (incl. round-trip)
│   ├── resolver/                # [Phase 2] Name resolution
│   ├── typechecker/             # [Phase 2] Type checking
│   ├── effects/                 # [Phase 2] Effect system
│   ├── speccheck/               # [Phase 2] Spec validation
│   └── interpreter/             # [Phase 3] Tree-walk interpreter
├── testdata/                    # Sample .aura files
├── user_docs/                   # User-facing documentation
├── ROADMAP.md                   # Development roadmap
├── DEVELOPMENT.md               # This file
└── README.md                    # Project overview
```

### Key Design Decisions

1. **Go implementation** — Chosen for fast compilation, easy cross-compilation, and strong tooling support.
2. **Indentation-sensitive lexer** — The lexer emits INDENT/DEDENT tokens, so the parser never deals with whitespace directly.
3. **Paren-depth tracking** — Inside `()`, `[]`, `{}`, newlines and indentation are suppressed. This allows multi-line expressions without explicit line continuation.
4. **Recursive descent parser** — Simple, predictable, and easy to extend. Operator precedence climbing handles expression parsing.
5. **Round-trip formatting** — The formatter produces deterministic output from any valid AST, ensuring `parse → format → parse → format` is stable.

---

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git

### Clone & Build

```bash
git clone https://github.com/unclebucklarson/aura.git
cd aura
go build -o aura ./cmd/aura
```

### Run Tests

```bash
# Run all tests
go test ./... -v

# Run tests for a specific package
go test ./pkg/lexer -v
go test ./pkg/parser -v
go test ./pkg/formatter -v

# Run with race detection
go test ./... -race

# Run with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### CLI Usage

```bash
# Format a file (print to stdout)
./aura format testdata/models.aura

# Format in-place
./aura format -w testdata/models.aura

# Parse and dump tokens + AST
./aura parse testdata/specs.aura
```

---

## Implementation Checklists

### Phase 2: Semantic Analysis

#### 2.1 Name Resolution (`pkg/resolver`)

```
[ ] Define Scope and Symbol types
[ ] Implement hierarchical scope chain (module → function → block)
[ ] Walk the AST to register all declarations
[ ] Walk the AST to resolve all references
[ ] Handle qualified name resolution (e.g., `TaskError.NotFound`)
[ ] Handle import resolution
[ ] Report errors: undefined names, duplicate declarations, visibility violations
[ ] Write tests for each error case
```

#### 2.2 Type Checker (`pkg/typechecker`)

```
[ ] Define internal Type representations (PrimitiveType, StructType, EnumType, etc.)
[ ] Implement type environment mapping names → types
[ ] Implement bidirectional type checking:
    [ ] Infer mode: expression → type
    [ ] Check mode: expression against expected type
[ ] Type-check literals (Int, Float, String, Bool, None)
[ ] Type-check binary operators (arithmetic, comparison, logical)
[ ] Type-check unary operators (negation, not)
[ ] Type-check function calls (argument count, types, named args, defaults)
[ ] Type-check field access on structs
[ ] Type-check index expressions on lists and maps
[ ] Type-check struct construction
[ ] Type-check pattern matching (exhaustiveness, type compatibility)
[ ] Type-check list comprehensions
[ ] Type-check lambda expressions (infer param types from context)
[ ] Type-check Option[T] and Result[T, E] operations
[ ] Type-check ? propagation operator
[ ] Implement structural subtyping for structs
[ ] Implement generic type parameter instantiation
[ ] Handle union types and literal string types
[ ] Report type errors with source locations and suggestions
[ ] Write comprehensive tests
```

#### 2.3 Refinement Types (`pkg/typechecker/refinement`)

```
[ ] Represent refinement predicates as a data structure
[ ] Evaluate constant predicates at compile time
[ ] Generate runtime assertion code for non-constant predicates
[ ] Track refinement information through variable assignments
[ ] Support built-in predicate vocabulary (len, self, matches, in)
[ ] Write tests for static and dynamic predicate checking
```

#### 2.4 Effect System (`pkg/effects`)

```
[ ] Extract effect annotations from function signatures
[ ] Build call graph from resolved AST
[ ] Compute transitive effect closure per function
[ ] Verify declared effects ⊇ required effects
[ ] Enforce explicit annotations on pub functions
[ ] Infer effects for private functions
[ ] Report effect mismatch errors with call chain context
[ ] Write tests for effect propagation and violation scenarios
```

#### 2.5 Spec Validation (`pkg/speccheck`)

```
[ ] Resolve satisfies references to spec blocks
[ ] Validate input names and types match
[ ] Validate effect lists are identical
[ ] Validate error types are covered
[ ] Detect possible guarantee violations where feasible
[ ] Enforce uniqueness constraints
[ ] Write tests for each validation rule
```

### Phase 3: Code Generation

#### 3.1 Tree-Walk Interpreter (`pkg/interpreter`)

```
[ ] Define Value types (IntVal, FloatVal, StringVal, BoolVal, NoneVal, etc.)
[ ] Define Environment (variable bindings per scope)
[ ] Implement expression evaluation:
    [ ] Literals
    [ ] Identifiers (lookup in environment)
    [ ] Binary and unary operators
    [ ] Function calls
    [ ] Field access
    [ ] Index access
    [ ] Struct construction
    [ ] List literals and list comprehensions
    [ ] Map literals
    [ ] String interpolation
    [ ] Lambda expressions
    [ ] If expressions
    [ ] Pipeline operator
    [ ] Option chaining (?)
[ ] Implement statement execution:
    [ ] Let bindings
    [ ] Assignment
    [ ] Return
    [ ] If/elif/else
    [ ] Match/case
    [ ] For loops
    [ ] While loops
    [ ] Break/continue
    [ ] Assert
    [ ] Expression statements
[ ] Implement function definition and calling convention
[ ] Implement effect capability injection
[ ] Implement test block runner
[ ] Add REPL support (read-eval-print loop)
[ ] Write tests for each evaluation form
```

---

## Testing Strategy

### Test Categories

1. **Unit tests** — Each package has `_test.go` files testing individual functions and components in isolation.

2. **Round-trip tests** — The formatter tests verify `source → parse → format → parse → format` stability. This catches bugs in both the parser and formatter.

3. **Integration tests** — Test the full pipeline from source to output. Located alongside the package tests or in a dedicated `integration/` directory.

4. **Testdata files** — The `testdata/` directory contains representative `.aura` files covering all language constructs:
   - `simple.aura` — Minimal: type, struct, enum, function
   - `models.aura` — Structs with refinement types, enums
   - `specs.aura` — Spec blocks with all sections
   - `service.aura` — Functions with effects, satisfies, complex bodies
   - `control_flow.aura` — if/elif/else, match, for, while
   - `expressions.aura` — Pipelines, comprehensions, lambdas, operators
   - `comments.aura` — Comment handling edge cases
   - `empty.aura` — Empty module edge case

### Writing Tests

- Use table-driven tests where appropriate
- Test both success cases and error cases
- Include source location validation in error tests
- Use `testdata/` files for integration-level tests
- Aim for >80% code coverage on new packages

### Running Tests

```bash
# All tests, verbose
go test ./... -v

# Specific package
go test ./pkg/parser -v -run TestMatchStatement

# With coverage report
go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out
```

---

## Code Organization Guidelines

### Package Responsibilities

Each package should have a **single, clear responsibility**:

- `token` — Token type definitions only. No logic.
- `lexer` — Source text → token stream. No AST knowledge.
- `ast` — AST node type definitions. Minimal logic (just constructors and accessors).
- `parser` — Token stream → AST. No type checking or validation.
- `formatter` — AST → canonical source text. No parsing.
- `resolver` — AST → AST with resolved names. No type checking.
- `typechecker` — Resolved AST → typed AST. Depends on resolver output.
- `effects` — Typed AST → effect-checked AST. Depends on typechecker and call graph.
- `interpreter` — Typed AST → execution. Depends on all analysis phases.

### Naming Conventions

- Go standard conventions: `CamelCase` for exports, `camelCase` for unexported
- AST node types: `PascalCase` matching the grammar (e.g., `StructDef`, `LetStmt`, `BinaryOp`)
- Test functions: `TestComponentName_scenario` (e.g., `TestParser_MatchStatement`)
- Error types: `ErrorCategory` prefix (e.g., `ResolveError`, `TypeError`, `EffectError`)

### Error Reporting

All errors should include:
- Source file path
- Line and column number
- Error code (e.g., `E101`, `E201`, `W301`)
- Clear message describing the problem
- Suggestion for fixing (where possible)

Use the `token.Span` from AST nodes to generate precise error locations.

---

## How to Contribute

### Workflow

1. **Check the roadmap** — Pick a task from [ROADMAP.md](ROADMAP.md)
2. **Create a branch** — `git checkout -b feature/your-feature`
3. **Implement** — Write code following the guidelines above
4. **Test** — Add tests and ensure all existing tests pass (`go test ./...`)
5. **Format** — Run `gofmt -w .` to format Go code
6. **Commit** — Use conventional commit messages:
   - `feat: Add name resolution for imports`
   - `fix: Handle empty match bodies in parser`
   - `test: Add type checker tests for union types`
   - `docs: Update roadmap with Phase 2 progress`
7. **Push & PR** — Push your branch and open a pull request

### Commit Message Format

```
<type>: <short description>

<optional longer description>

<optional references>
```

Types: `feat`, `fix`, `test`, `docs`, `refactor`, `chore`

### Before Submitting a PR

- [ ] All tests pass: `go test ./...`
- [ ] No race conditions: `go test ./... -race`
- [ ] Code is formatted: `gofmt -l .` (should output nothing)
- [ ] New code has tests
- [ ] Commit messages follow conventions
- [ ] ROADMAP.md is updated if a task is completed
