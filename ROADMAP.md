# Aura Language Roadmap

> A phased plan for building Aura from a working parser into a fully functional programming language.

> 🤖 **Aura is an AI-first language.** Every phase and feature in this roadmap is evaluated against a core question: **Does this make AI code generation and AI-human collaboration better?** See [AI_MISSION.md](AI_MISSION.md) for the full mission statement.

---

## AI-First Design Principles

These principles guide every phase of development. When evaluating features, trade-offs, or priorities, apply them in order:

1. **AI parseability first** — Can an AI agent read this feature's output and know exactly what to do? Structured, unambiguous representations always win.
2. **Machine-checkable contracts** — Every constraint, effect, and requirement should be verifiable by the compiler, not dependent on human review alone.
3. **Explicit over implicit** — If information exists (types, effects, error cases, constraints), it must be in the syntax. Hidden conventions are the enemy of AI code generation.
4. **Specs as the interface** — Specs are how humans communicate intent to AI. Every feature should consider: how does this interact with the spec system?
5. **Vibe coding flow** — The human writes *what* (specs), the AI writes *how* (implementation), the compiler validates *correctness*. Features should reinforce this loop.

**Feature evaluation checklist:**
- [ ] Does this feature help AI generate correct code faster?
- [ ] Is the feature's syntax unambiguous and machine-parseable?
- [ ] Does it reduce the need for AI to read surrounding context?
- [ ] Can the compiler validate it automatically?
- [ ] Does it integrate with specs and effects?

---

## Phase Overview

| Phase | Name | Status | Estimated Effort |
|-------|------|--------|------------------|
| 1 | Syntax (Lexer, Parser, Formatter) | ✅ COMPLETE | — |
| 2 | Semantic Analysis | ✅ COMPLETE | — |
| 3 | Code Generation (Interpreter) | 🚧 **UP NEXT** | 4–6 weeks |
| 4 | Runtime & Standard Library | 🔲 Not Started | 8–12 weeks |
| 5 | Advanced Tooling & Ecosystem | 🔲 Not Started | Ongoing |

---

## Phase 1: Syntax — ✅ COMPLETE

> 🤖 **AI optimization:** The parser and AST produce structured, unambiguous representations that AI agents can consume directly. The formatter ensures canonical output — AI-generated code always looks the same as human-written code.

The foundation of the Aura toolchain is fully implemented and tested.

### Deliverables

- [x] **Lexer** (`pkg/lexer`) — Indentation-sensitive tokenizer with INDENT/DEDENT, paren-depth tracking, comment handling
- [x] **Parser** (`pkg/parser`) — Recursive descent parser with operator precedence climbing
- [x] **AST** (`pkg/ast`) — Complete node definitions covering all language constructs
- [x] **Formatter** (`pkg/formatter`) — Canonical source formatter with round-trip guarantee
- [x] **CLI** (`cmd/aura`) — `format` and `parse` commands
- [x] **Test suite** — 36 tests across lexer (11), parser (16), formatter (9)

### Key Properties Verified

- Round-trip guarantee: `parse → format → parse → format` produces identical output
- All language constructs parse correctly: structs, enums, traits, impls, specs, functions, control flow, expressions
- Edge cases handled: empty files, blank lines between blocks, nested indentation

---

## Phase 2: Semantic Analysis — ✅ COMPLETE

**Goal:** Validate that parsed programs are meaningful — names resolve, types check, effects are tracked, and specs are verified.

> 🤖 **AI optimization:** This phase is critical for AI code generation. Type checking, effect validation, and spec verification give AI agents **immediate, automated feedback** on whether generated code is correct. Every error message is structured and JSON-serializable for AI to parse and fix automatically.

**Completed:** 2026-03-17

### 2.1 Symbol Table & Scope Management

- [x] Hierarchical symbol table with scope kinds (Module, Function, Block, Loop, Test)
- [x] Symbol definition with duplicate detection
- [x] Hierarchical symbol lookup (walks parent scopes)
- [x] Local-only lookup for shadowing semantics
- [x] Loop context detection (`IsInsideLoop`) for break/continue validation
- [x] Enclosing function resolution for return type checking

**Package:** `pkg/symbols` — 9 tests ✅

### 2.2 Type System

- [x] Complete type representation (Primitive, Struct, Enum, Union, Function, Option, Result, Refinement, TypeParam, Never, Any, None, Alias, Intersection, Tuple, List, Map, Set, StringLiteral)
- [x] Built-in primitive singletons (Int, Float, String, Bool, None, Never, Any)
- [x] Type equality checking (`Equal()`)
- [x] Subtyping/assignability rules (`IsAssignableTo()`) including:
  - Never as bottom type, Any as top type
  - None/T to Option[T], refinement to base type
  - String literal to String/Union, Int to Float widening
  - Struct width subtyping
- [x] Alias and refinement unwrapping (`Underlying()`)
- [x] Type registry with built-in pre-population

**Package:** `pkg/types` — 26 tests ✅

### 2.3 Type Checker

- [x] Multi-pass architecture (types → specs → functions → constants → bodies → spec validation → tests)
- [x] Bidirectional type inference for all expression forms
- [x] Struct construction validation (missing/unknown fields)
- [x] Function call type checking with effect propagation
- [x] Pattern matching with enum exhaustiveness checking
- [x] Effect tracking and validation (declared vs. required effects)
- [x] Spec contract validation (inputs, effects)
- [x] Control flow validation (break/continue in loops, return in functions)
- [x] Mutability enforcement (immutable by default)
- [x] AI-parseable structured error output (JSON format with error codes, expected/got, fix suggestions)
- [x] CLI `aura check` command with `--json` flag for AI agents

**Package:** `pkg/checker` — 48 tests ✅

### Deferred to Future Phases

- Refinement predicate static evaluation (Phase 4 — runtime assertions)
- Import resolution and cross-module type checking (Phase 4/5)
- Visibility (`pub`) enforcement across modules (Phase 4/5)
- Generic type argument inference (Phase 3/4 as needed)
- Transitive effect inference for private functions (Phase 4)

### Phase 2 Milestone — ✅ Achieved

`aura check <file>` now:
- ✅ Reports name resolution errors with source locations
- ✅ Reports type errors with clear messages, expected/got info, and fix suggestions
- ✅ Reports effect mismatches and missing capabilities
- ✅ Reports spec contract violations
- ✅ Outputs structured JSON for AI agent consumption (`--json` flag)
- ✅ **83 new tests** across symbols (9), types (26), checker (48) — all passing

---

## Phase 3: Code Generation — 🚧 UP NEXT

**Goal:** Execute Aura programs via a tree-walk interpreter, completing the end-to-end vibe coding loop.

> 🤖 **AI optimization:** Code generation outputs should be deterministic and predictable so AI agents can reason about the compilation process. The interpreter provides structured error output (JSON-friendly) that AI agents can parse for debugging. A "dry-run" mode validates without executing — useful for AI testing loops.

**Dependencies:** Phase 2 ✅ (semantic analysis complete)

### Why Phase 3 Now? — Rationale

The tree-walk interpreter is the **highest-impact next step** for the AI-first mission:

1. **Closes the vibe coding feedback loop** — Without execution, the workflow is: spec → generate → check. With the interpreter: spec → generate → check → **run → see output → iterate**. This is the complete AI development cycle.
2. **Builds directly on Phase 2** — The typed AST from the checker is the interpreter's input. All the type information, scope resolution, and validation work is already done.
3. **Enables AI test-driven development** — AI agents can generate code, run `aura test` blocks, observe failures, and self-correct. This is the killer feature for vibe coding.
4. **Lower complexity than alternatives** — A tree-walk interpreter (~4–6 weeks) is faster to build than bytecode compilation (~6–10 weeks) and more self-contained than Go transpilation (which requires a Go runtime dependency).
5. **Validates the language design** — Running real programs will surface language design issues early, before investing in optimization.

### 3.1 Tree-Walk Interpreter (Primary Target)

**Complexity:** Medium-High | **Estimate:** 4–6 weeks

#### 3.1.1 Value System
- [ ] Implement Aura value types (AuraInt, AuraFloat, AuraString, AuraBool, AuraNone)
- [ ] Implement composite values (AuraList, AuraMap, AuraSet, AuraTuple)
- [ ] Implement AuraStruct with field access and construction
- [ ] Implement AuraEnum with variant matching
- [ ] Implement AuraOption (Some/None) and AuraResult (Ok/Err)
- [ ] Implement AuraFunction (closures with captured environment)

#### 3.1.2 Expression Evaluation
- [ ] Evaluate literals (int, float, string, bool, none)
- [ ] Evaluate binary and unary operations with type-appropriate semantics
- [ ] Evaluate function calls with argument binding and default values
- [ ] Evaluate field access and index operations
- [ ] Evaluate string interpolation at runtime
- [ ] Evaluate list comprehensions and lambda expressions
- [ ] Evaluate `?` propagation for Option/Result types
- [ ] Evaluate pipeline operator (`|>`)

#### 3.1.3 Statement Execution
- [ ] Execute `let` bindings (immutable and mutable)
- [ ] Execute assignments (with mutability checking)
- [ ] Execute `return`, `break`, `continue` (as control flow signals)
- [ ] Execute `if`/`elif`/`else` chains
- [ ] Execute `match` with pattern matching evaluation
- [ ] Execute `for ... in` loops with iterator protocol
- [ ] Execute `while` loops

#### 3.1.4 Effect & Runtime Infrastructure
- [ ] Environment/scope management for interpreter state
- [ ] Effect capability injection via `with` blocks
- [ ] Built-in print/assert functions
- [ ] Structured runtime error output (JSON for AI agents)
- [ ] Test block runner (`aura test <file>`)

**Package:** `pkg/interpreter`

#### 3.1.5 CLI Integration
- [ ] `aura run <file>` — execute an Aura program
- [ ] `aura test <file>` — run test blocks with pass/fail reporting
- [ ] `--json` flag for structured output (AI agent consumption)
- [ ] `--dry-run` flag for validation without execution

### 3.2 Bytecode Compiler (Future — Phase 5+)

**Complexity:** Very High | **Estimate:** 6–10 weeks

*Deferred. A tree-walk interpreter is sufficient for the AI-first use case where correctness and rapid feedback matter more than raw performance.*

- [ ] Design a bytecode instruction set for Aura
- [ ] Implement a bytecode compiler from the typed AST
- [ ] Build a stack-based virtual machine
- [ ] Implement garbage collection
- [ ] Add debug information (source maps, breakpoints)

**Package:** `pkg/compiler`, `pkg/vm`

### 3.3 Transpilation to Go (Future — Phase 5+)

**Complexity:** Medium-High | **Estimate:** 4–6 weeks

*Deferred. May revisit after the interpreter proves out the language semantics.*

- [ ] Generate Go source code from the typed Aura AST
- [ ] Map Aura types to Go types
- [ ] Handle indentation-based blocks → Go brace blocks
- [ ] Implement effect tracking as Go interfaces/context injection
- [ ] Generate Go test files from Aura test blocks

**Package:** `pkg/codegen/golang`

### Phase 3 Milestone

`aura run <file>` should execute a complete Aura program and produce output. `aura test <file>` should run test blocks and report structured results that AI agents can parse.

---

## Phase 4: Runtime & Standard Library

**Goal:** Provide the standard library and runtime support needed for real programs.

> 🤖 **AI optimization:** Standard library APIs should have spec blocks defining their contracts, making them instantly understandable by AI agents. Effect providers should be mockable by default, enabling AI to generate testable code without external dependencies. Library functions should follow consistent patterns so AI can predict APIs for unfamiliar modules.

**Dependencies:** Phase 3 (code generation)

### 4.1 Core Runtime

**Complexity:** Medium | **Estimate:** 2–3 weeks

- [ ] String operations (len, slice, contains, split, join, trim, replace)
- [ ] List operations (push, pop, get, slice, sort, filter, map, reduce)
- [ ] Map operations (get, set, delete, keys, values, entries)
- [ ] Option/Result utility methods (unwrap, map, flat_map, or_else, is_ok, is_err)
- [ ] Numeric operations and math functions
- [ ] Equality and comparison for all types

### 4.2 Standard Library Modules

**Complexity:** Medium | **Estimate:** 3–4 weeks

- [ ] `std.time` — Instant, Duration, formatting, parsing
- [ ] `std.uuid` — UUID v4 generation
- [ ] `std.io` — Print, read (console I/O)
- [ ] `std.collections` — Extended List, Map, Set operations
- [ ] `std.testing` — Assert utilities, mock framework, test runner
- [ ] `std.json` — JSON serialization/deserialization
- [ ] `std.string` — Extended string utilities and regex

### 4.3 Effect Runtime

**Complexity:** Medium-High | **Estimate:** 2–3 weeks

- [ ] Implement effect capability providers (db, net, fs, time, random, auth, log)
- [ ] Effect mocking framework for tests
- [ ] `with` block for capability injection in test contexts
- [ ] Runtime effect tracking and violation detection

### Phase 4 Milestone

The AuraTask example from the spec should run end-to-end with mocked effects in tests.

---

## Phase 5: Advanced Tooling & Ecosystem

**Goal:** Build the developer experience and ecosystem around Aura.

> 🤖 **AI optimization:** This phase is where AI-first design pays off most. The LSP should expose spec/effect/type information as structured data for AI agents. The spec-to-implementation pipeline (§5.3) is the flagship AI feature — it's the full realization of the vibe coding workflow. Package metadata should be machine-readable so AI can discover and use libraries without human guidance.

**Dependencies:** Phases 2–4 (progressive)

### 5.1 Language Server Protocol (LSP)

**Complexity:** High | **Estimate:** 4–6 weeks

- [ ] Implement an LSP server for Aura
- [ ] Go-to-definition, find-references, rename
- [ ] Hover information (types, docs, effects)
- [ ] Diagnostics (errors and warnings from semantic analysis)
- [ ] Auto-completion for identifiers, types, and keywords
- [ ] Signature help for function calls
- [ ] Code actions (quick fixes for common errors)

**Package:** `cmd/aura-lsp`

### 5.2 Package Manager

**Complexity:** Medium | **Estimate:** 3–4 weeks

- [ ] Module resolution and dependency management
- [ ] Package manifest file format (`aura.toml`)
- [ ] Version resolution and lock files
- [ ] Registry or Git-based package fetching

### 5.3 AI Integration

**Complexity:** Medium | **Estimate:** 2–3 weeks

- [ ] Spec-to-implementation generation pipeline
- [ ] AST-aware code generation prompts
- [ ] Automatic spec validation for AI-generated code
- [ ] Structured output format for AI consumption

### 5.4 Documentation Generator

**Complexity:** Low-Medium | **Estimate:** 1–2 weeks

- [ ] Extract doc comments (`##`) from source files
- [ ] Generate HTML/Markdown documentation from AST
- [ ] Include type signatures, effects, and spec information
- [ ] Cross-reference linking between modules

### 5.5 REPL

**Complexity:** Medium | **Estimate:** 2 weeks

- [ ] Interactive Aura evaluation loop
- [ ] Expression evaluation and pretty-printing
- [ ] History, auto-completion, and multi-line input
- [ ] `:type` and `:effects` introspection commands

---

## Contributing

See [DEVELOPMENT.md](DEVELOPMENT.md) for setup instructions, architecture overview, and contribution guidelines.

---

## Version History

| Date | Version | Notes |
|------|---------|-------|
| 2026-03-17 | v0.1 | Phase 1 complete; roadmap published |
| 2026-03-17 | v0.2 | Phase 2 complete (type checker, 83 tests); Phase 3 (interpreter) selected as next |
