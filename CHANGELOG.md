# Changelog

All notable changes to the Aura toolchain are documented here.

Format based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/).

---

## [v0.4.0] — 2026-03-20

### Phase 4.1: Core Runtime Methods — COMPLETE

Major release adding **108+ built-in methods** across 5 core types, delivered in 4 implementation chunks.

### Added

#### Method Dispatch Infrastructure
- Centralized method registry system (`pkg/interpreter/methods.go`) using `RegisterMethod(ValueType, "name", func)` pattern
- `callValue()` helper for invoking Aura lambdas/closures from Go method implementations
- `cmpValues()` helper for type-safe ordering comparisons (Int, Float, String)
- Registry-based method resolution replaces inline switch statements in `eval.go`

#### String Methods (22) — `methods_string.go`
- Core: `len`, `upper`/`to_upper`, `lower`/`to_lower`, `contains`, `split`, `trim`, `trim_start`, `trim_end`
- Search: `starts_with`, `ends_with`, `index_of` (returns Option), `replace`
- Transform: `repeat`, `reverse`, `chars`, `slice` (with bounds checking)
- Aliases: `length` → `len`

#### List Methods (27) — `methods_list.go`
- Core: `len`/`length`, `append`/`push`, `contains`, `is_empty`
- Safe accessors: `first()`, `last()`, `get(index)` — all return Option
- Mutation: `pop()` (returns Option), `remove(index)`
- Transforms: `reverse()`, `slice(start, end?)` (supports negative indices), `join(sep)`, `index_of(item)` (returns Option)
- Higher-order: `map(fn)`, `filter(fn)`, `reduce(init, fn)`, `for_each(fn)`, `flat_map(fn)`, `flatten()`
- Predicates: `any(fn)`, `all(fn)`, `count(fn?)`
- Utilities: `unique()`, `sum()`, `min()`/`max()` (return Option), `sort()`, `zip(other)`, `enumerate()`

#### Map Methods (24) — `methods_map.go`
- Size/emptiness: `len`/`length`/`size`, `is_empty`
- Key/value access: `keys()`, `values()`, `entries()`, `get(key)` (returns Option), `get_or(key, default)`
- Lookup: `has(key)`, `contains_key(key)`, `contains_value(value)`
- Mutation: `set(key, value)`, `remove(key)` (returns Option), `delete(key)` (returns Bool), `clear()`, `merge(other)`
- Higher-order: `filter(fn)`, `map(fn)`, `for_each(fn)`, `reduce(init, fn)`, `any(fn)`, `all(fn)`, `count(fn?)`
- Utilities: `to_list()`, `find(fn)` (returns Option)

#### Option Methods (17) — `methods_option.go`
- Predicates: `is_some()`, `is_none()`
- Extraction: `unwrap()`, `expect(msg)`, `unwrap_or(default)`, `unwrap_or_else(fn)`
- Monadic transforms: `map(fn)`, `flat_map(fn)`, `and_then(fn)`, `filter(fn)`, `flatten()`
- Combinators: `or_else(fn)`, `or(alt)`, `and(other)`, `zip(other)`
- Query: `contains(value)`
- Conversion: `to_result(err_val)`

#### Result Methods (18) — `methods_option.go`
- Predicates: `is_ok()`, `is_err()`
- Extraction: `unwrap()`, `unwrap_err()`, `expect(msg)`, `unwrap_or(default)`, `unwrap_or_else(fn)`
- Monadic transforms: `map(fn)`, `map_err(fn)`, `and_then(fn)`, `or_else(fn)`, `flatten()`
- Combinators: `or(alt)`, `and(other)`
- Query: `contains(value)`, `contains_err(value)`
- Conversion: `ok()`, `err()`, `to_option()`

#### Tests
- 222 new method-specific tests in `methods_test.go`
- Total test count: **468 tests** across all packages (up from 232)
- Comprehensive coverage including: success cases, error/panic conditions, None/Err edge cases, method chaining, monadic composition, Option↔Result round-trip conversions

### Changed
- `eval.go`: Refactored FieldAccess evaluation to use method registry instead of inline switch statements
- Interpreter package now contains 12 source files (up from 6)

---

## [v0.3.1] — 2026-03-19

### Added
- **String interpolation** — `"Hello, {name}!"` with full expression support in lexer, parser, and interpreter
- **Pipeline operator** (`|>`) — Lexer tokenization, parser precedence handling, interpreter evaluation with lambda support
- **Option chaining** (`?.`) — None short-circuiting for `?` postfix operator
- 14 new pipeline operator tests (232+ total tests)

---

## [v0.3.0] — 2026-03-17

### Phase 3: Tree-Walk Interpreter — COMPLETE

### Added
- Tree-walk interpreter (`pkg/interpreter/`) with value system, environment, evaluator, module runner, test runner
- CLI commands: `aura run`, `aura test`, `aura repl`
- Full expression/statement evaluation (arithmetic, comparison, logic, control flow, structs, enums, match, closures, lambdas, list comprehensions)
- 14 built-in functions: `print`, `len`, `str`, `int`, `float`, `range`, `type_of`, `abs`, `min`, `max`, `Ok`, `Err`, `Some`, `None`
- 112 interpreter tests (211 total)

---

## [v0.2.0] — 2026-03-17

### Phase 2: Semantic Analysis — COMPLETE

### Added
- Symbol table and scope management (`pkg/symbols/`) — 9 tests
- Type system representation with subtyping (`pkg/types/`) — 26 tests
- Multi-pass type checker (`pkg/checker/`) — 48 tests
- AI-parseable structured error output with JSON format
- CLI `aura check` command with `--json` flag
- 83 new tests (119 total)

---

## [v0.1.0] — 2026-03-17

### Phase 1: Syntax — COMPLETE

### Added
- Indentation-sensitive lexer (`pkg/lexer/`) — 11 tests
- Recursive descent parser (`pkg/parser/`) — 16 tests
- Complete AST node definitions (`pkg/ast/`)
- Canonical source formatter (`pkg/formatter/`) — 9 tests with round-trip guarantee
- CLI entry point with `format` and `parse` commands
- 36 tests total
