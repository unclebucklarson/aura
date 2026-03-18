# Aura Language Roadmap

> A phased plan for building Aura from a working parser into a fully functional programming language.

---

## Phase Overview

| Phase | Name | Status | Estimated Effort |
|-------|------|--------|------------------|
| 1 | Syntax (Lexer, Parser, Formatter) | ✅ COMPLETE | — |
| 2 | Semantic Analysis | 🔲 Not Started | 8–12 weeks |
| 3 | Code Generation | 🔲 Not Started | 6–10 weeks |
| 4 | Runtime & Standard Library | 🔲 Not Started | 8–12 weeks |
| 5 | Advanced Tooling & Ecosystem | 🔲 Not Started | Ongoing |

---

## Phase 1: Syntax — ✅ COMPLETE

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

## Phase 2: Semantic Analysis

**Goal:** Validate that parsed programs are meaningful — names resolve, types check, effects are tracked, and specs are verified.

**Dependencies:** Phase 1 (complete)

### 2.1 Symbol Table & Name Resolution

**Complexity:** Medium | **Estimate:** 2–3 weeks

- [ ] Build a hierarchical symbol table (module → function → block scopes)
- [ ] Resolve all identifier references to their declarations
- [ ] Handle qualified names (`std.time.Instant`, `TaskError.NotFound`)
- [ ] Detect undefined names, duplicate definitions, and shadowing
- [ ] Resolve `import` statements and track module dependencies
- [ ] Handle visibility (`pub` vs private) access rules

**Package:** `pkg/resolver`

### 2.2 Type Checker

**Complexity:** High | **Estimate:** 3–4 weeks

- [ ] Implement bidirectional type inference (check mode + infer mode)
- [ ] Type-check all expression forms (binary ops, calls, field access, index, etc.)
- [ ] Validate struct field types and default values
- [ ] Check function parameter types, return types, and all return paths
- [ ] Implement structural subtyping rules (width subtyping for structs)
- [ ] Handle generic type parameters and type argument inference
- [ ] Validate union types (`"pending" | "done"`) and literal types
- [ ] Implement `Option[T]` / `Result[T, E]` type checking
- [ ] Handle the `?` propagation operator desugaring
- [ ] Validate pattern exhaustiveness in `match` statements
- [ ] Implement type compatibility checks (`T <: U`)

**Package:** `pkg/typechecker`

### 2.3 Refinement Type Checking

**Complexity:** High | **Estimate:** 2–3 weeks

- [ ] Parse and validate refinement predicates (`where len >= 1 and len <= 64`)
- [ ] Static predicate evaluation for constant expressions
- [ ] Insert runtime assertion stubs for predicates that can't be statically verified
- [ ] Track refinement flow through assignments and function boundaries
- [ ] Support built-in predicates: `len`, `self`, `matches`, `in`

**Package:** `pkg/typechecker/refinement`

### 2.4 Effect System

**Complexity:** Medium | **Estimate:** 2 weeks

- [ ] Parse effect annotations from function signatures (`with db, time`)
- [ ] Build a call graph from the AST
- [ ] Compute transitive effect closure for each function
- [ ] Verify declared effects match actual effects (required ⊆ declared)
- [ ] Enforce explicit effects on `pub` functions
- [ ] Infer effects for private functions without annotations
- [ ] Verify spec ↔ function effect agreement

**Package:** `pkg/effects`

### 2.5 Spec Validation

**Complexity:** Medium | **Estimate:** 1–2 weeks

- [ ] Bind `satisfies` declarations to spec blocks
- [ ] Validate spec input names and types match function parameters
- [ ] Validate spec effects match function effects (exact match)
- [ ] Validate spec error types are covered by the function's return type
- [ ] Emit warnings for possible spec guarantee violations (where detectable)
- [ ] Enforce spec uniqueness (one function per spec, one spec per function)

**Package:** `pkg/speccheck`

### Phase 2 Milestone

At the end of Phase 2, `aura check <file>` should:
- Report all name resolution errors with source locations
- Report type errors with clear messages and suggestions
- Report effect mismatches and missing capabilities
- Report spec violations
- Successfully validate all testdata files

---

## Phase 3: Code Generation

**Goal:** Execute Aura programs, either via interpretation or compilation.

**Dependencies:** Phase 2 (semantic analysis)

### 3.1 Tree-Walk Interpreter (Recommended First Target)

**Complexity:** Medium-High | **Estimate:** 4–6 weeks

- [ ] Implement a value system (AuraInt, AuraString, AuraStruct, AuraEnum, etc.)
- [ ] Evaluate all expression types
- [ ] Execute all statement types (let, assign, return, if, match, for, while)
- [ ] Implement function calls with argument binding and default values
- [ ] Implement pattern matching evaluation
- [ ] Handle struct construction and field access
- [ ] Implement list comprehensions and lambda evaluation
- [ ] Implement string interpolation at runtime
- [ ] Handle `Option` and `Result` types with `?` propagation
- [ ] Implement effect capability injection and checking at runtime

**Package:** `pkg/interpreter`

### 3.2 Bytecode Compiler (Optional, Future)

**Complexity:** Very High | **Estimate:** 6–10 weeks

- [ ] Design a bytecode instruction set for Aura
- [ ] Implement a bytecode compiler from the typed AST
- [ ] Build a stack-based virtual machine
- [ ] Implement garbage collection
- [ ] Add debug information (source maps, breakpoints)

**Package:** `pkg/compiler`, `pkg/vm`

### 3.3 Transpilation to Go (Alternative)

**Complexity:** Medium-High | **Estimate:** 4–6 weeks

- [ ] Generate Go source code from the typed Aura AST
- [ ] Map Aura types to Go types
- [ ] Handle indentation-based blocks → Go brace blocks
- [ ] Implement effect tracking as Go interfaces/context injection
- [ ] Generate Go test files from Aura test blocks

**Package:** `pkg/codegen/golang`

### Phase 3 Milestone

`aura run <file>` should execute a complete Aura program and produce output.

---

## Phase 4: Runtime & Standard Library

**Goal:** Provide the standard library and runtime support needed for real programs.

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
