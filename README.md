# Aura Toolchain

A complete toolchain for the **Aura programming language** — a Python-inspired, statically typed language with specification-driven development, algebraic types, and effect tracking.

Built in Go. Implements lexing, parsing, AST construction, and canonical source formatting.

## Project Structure

```
aura-toolchain/
├── cmd/aura/main.go           # CLI entry point (format, parse commands)
├── pkg/
│   ├── token/token.go         # Token types, positions, spans
│   ├── lexer/lexer.go         # Indentation-sensitive lexer (INDENT/DEDENT)
│   ├── ast/ast.go             # Complete AST node definitions
│   ├── parser/parser.go       # Recursive descent parser
│   └── formatter/formatter.go # AST → canonical source formatter
├── testdata/                  # Sample .aura files
│   ├── models.aura            # AuraTask models (struct, enum, type aliases)
│   ├── specs.aura             # Specification blocks
│   ├── service.aura           # Functions with effects & satisfies clauses
│   ├── control_flow.aura      # if/elif/else, match, for loops
│   ├── expressions.aura       # Pipelines, list comprehensions, lambdas
│   ├── comments.aura          # Comment handling
│   ├── simple.aura            # Minimal example
│   └── empty.aura             # Edge case: empty module
└── README.md
```

## Quick Start

### Prerequisites

- Go 1.22+

### Build

```bash
go build -o aura ./cmd/aura
```

### Usage

**Format** an Aura source file (prints canonical formatting to stdout):

```bash
./aura format testdata/models.aura
```

**Parse** an Aura source file (dumps tokens and AST):

```bash
./aura parse testdata/specs.aura
```

**Format in-place** with the `-w` flag:

```bash
./aura format -w testdata/service.aura
```

## Language Features Supported

### Types & Definitions
- **Type aliases** with refinement types: `type TaskId = String where len >= 1`
- **Union types**: `type Status = "pending" | "done" | "cancelled"`
- **Structs** with default values and optional fields: `pub struct Task:`
- **Enums** with data variants: `pub enum TaskError: NotFound(TaskId)`
- **Traits** and **impl blocks**

### Functions
- Named & default parameters, return types
- **Effect tracking**: `fn save() -> Result with db, time`
- **Satisfies clauses**: `fn create_task(...) satisfies CreateNewTask`
- **Guard clauses** with `where`

### Spec Blocks
- `doc`, `inputs`, `guarantees`, `effects`, `errors` sections
- Typed inputs with descriptions
- Guarantee strings and error variant descriptions

### Control Flow
- `if` / `elif` / `else`
- `match` with patterns (literals, wildcards, constructors, destructuring)
- `for ... in` loops
- `while` loops
- `return`, `break`, `continue`

### Expressions
- Binary and unary operators
- Pipeline operator: `data |> transform |> format`
- List comprehensions: `[x * 2 for x in items if x > 0]`
- Lambda expressions: `|x| x + 1`
- Optional chaining: `task.completed_at?`
- Unwrap: `maybe_task!`
- String interpolation: `"Hello, {name}!"`
- Struct construction with named fields

### Indentation
- Python-style significant whitespace
- INDENT/DEDENT token generation
- 4-space canonical indentation (enforced by formatter)

## Testing

Run the full test suite:

```bash
go test ./... -v
```

**Test breakdown:**
- `pkg/lexer/` — 11 tests covering tokenization, indentation, comments, edge cases
- `pkg/parser/` — 16 tests covering all language constructs
- `pkg/formatter/` — 9 tests including round-trip verification (parse → format → parse = same AST)

### Round-Trip Guarantee

The formatter produces deterministic output. Formatting source code, then parsing and formatting again, always produces identical output:

```
source → parse → AST₁ → format → source₂ → parse → AST₂ → format → source₃
                                  source₂ == source₃  ✓
```

## Architecture

### Lexer (`pkg/lexer`)
Scans Aura source into tokens with full position tracking. Key features:
- **Indentation tracking** via an indent stack — emits `INDENT` and `DEDENT` tokens
- **Paren depth tracking** — suppresses NEWLINE/INDENT/DEDENT inside `()`, `[]`, `{}`
- **Comment handling** — `#` line comments and `##` doc comments
- **Blank line handling** — properly emits DEDENTs across blank-line gaps between blocks

### Parser (`pkg/parser`)
Recursive descent parser that builds a complete AST. Features:
- Operator precedence climbing for expressions
- Pattern matching for `match` cases (wildcards, literals, constructors, lists, tuples)
- Type expression parsing (generics, optionals, result types, maps, tuples)
- Spec block parsing with all section types

### AST (`pkg/ast`)
Complete node definitions for the Aura language including:
- Module-level declarations (types, structs, enums, functions, specs, traits, impls, tests)
- Statements (let, assignment, return, if, match, for, while, expression statements)
- Expressions (literals, binary/unary ops, calls, field access, index, pipe, comprehensions, lambdas)
- Patterns and type expressions

### Formatter (`pkg/formatter`)
Converts AST back to canonical Aura source with:
- Consistent 4-space indentation
- Deterministic output ordering
- Blank line separation between top-level declarations
- Proper handling of all expression precedence (parenthesization where needed)

## Example

Input (`testdata/models.aura`):

```aura
module auratask.models

import std.time as time

type TaskId = String where len >= 1 and len <= 64

pub struct Task:
    pub id: TaskId
    pub title: String where len >= 1 and len <= 200
    pub status: TaskStatus = "pending"
    pub priority: Priority = 3
    pub created_at: time.Instant
    pub tags: [String] = []
```

The formatter preserves this exact canonical form. The parser produces a full AST that can be inspected, transformed, or used for code generation.

## License

MIT
