# AI Next Session Context — Phase 4.1 Chunk 4

> **Purpose:** This file provides complete context for any AI agent (or human developer) to pick up exactly where we left off. Read this before starting any work.
>
> **Last updated:** 2026-03-20
> **Last session by:** DeepAgent (Abacus.AI)

---

## Current Status

### What's Complete

| Phase | Description | Status |
|-------|-------------|--------|
| Phase 1 | Lexer, Parser, AST, Formatter, CLI | ✅ Complete |
| Phase 2 | Type Checker, Symbol Table, Semantic Analysis | ✅ Complete |
| Phase 3 | Tree-walk Interpreter | ✅ Complete |
| Phase 3+ | String Interpolation, Pipeline Operator (`|>`), Option Chaining (`?.`) | ✅ Complete |
| Phase 4.1 Chunk 1 | Method Dispatch Infrastructure + String Methods (22 methods) | ✅ Complete |
| Phase 4.1 Chunk 2 | List Methods (27 methods) | ✅ Complete |
| Phase 4.1 Chunk 3 | Map Methods (24 methods) | ✅ Complete |

- **Test count:** 379 test functions passing (259 interpreter, 11 lexer, 16 parser, plus checker/symbols/types/formatter)
- **Language version:** Pre-1.0, Phase 4.1 Chunk 3 complete
- **Repository:** `github.com/unclebucklarson/aura`
- **Branch:** `main`

### Method Registry Summary

**Architecture** (in `pkg/interpreter/methods.go`):
- `MethodFunc` type: `func(receiver Value, args []Value) Value`
- `methodRegistry` map: `map[ValueType]map[string]MethodFunc{}`
- `RegisterMethod(vt, name, fn)` — registers a method
- `LookupMethod(vt, name)` — looks up a method (returns nil if not found)
- `resolveMethod(obj, name)` — returns a `BuiltinFnVal` closure for a method
- Methods are registered via `init()` functions in type-specific files

**String methods (22 methods in `methods_string.go`):**
`len`, `length`, `upper`, `to_upper`, `lower`, `to_lower`, `contains`, `split`, `trim`, `trim_left`, `trim_right`, `starts_with`, `ends_with`, `replace`, `replace_first`, `slice`, `index_of`, `chars`, `join`, `repeat`, `is_empty`, `reverse`, `pad_left`, `pad_right`

**List methods (27 methods in `methods_list.go`):**
`len`, `length`, `append`, `push`, `contains`, `is_empty`, `first`, `last`, `get`, `pop`, `remove`, `reverse`, `slice`, `join`, `index_of`, `map`, `filter`, `reduce`, `for_each`, `flat_map`, `flatten`, `any`, `all`, `count`, `unique`, `sum`, `min`, `max`, `sort`, `zip`, `enumerate`

**Map methods (24 methods in `methods_map.go`):**
- Size/emptiness: `len`, `length`, `size`, `is_empty`
- Key/value accessors: `keys`, `values`, `entries`
- Lookup: `has`, `contains_key`, `contains_value`, `get` (returns Option), `get_or` (with default)
- Mutation: `set`, `remove` (returns Option of removed value), `delete` (returns Bool), `clear`, `merge`
- Higher-order: `filter`, `map`, `for_each`, `reduce`, `any`, `all`, `count` (optional predicate)
- Utility: `to_list`, `find` (returns Option)

**Helper functions (in `methods_list.go`):**
- `callValue(fn Value, args []Value) Value` — invokes any callable (FunctionVal, LambdaVal, BuiltinFnVal)
- `cmpValues(a, b Value) int` — compares two values for ordering

**Helper functions (in `methods_map.go`):**
- `mapFindKey(m *MapVal, key Value) int` — finds index of key in map, or -1

**Existing Option/Result methods:** None (constructors `Some`, `Ok`, `Err` exist as builtins)

---

## What We Just Did (Session Ending 2026-03-20)

1. **Implemented 24 map methods** in `pkg/interpreter/methods_map.go`:
   - Size: `len`, `length`, `size`, `is_empty`
   - Accessors: `keys`, `values`, `entries` (returns list of [key, value] lists)
   - Lookup: `has`, `contains_key`, `contains_value`, `get` (returns Option), `get_or` (with default)
   - Mutation: `set` (add/update), `remove` (returns Option), `delete` (returns Bool), `clear`, `merge` (overwrites existing keys)
   - Higher-order: `filter(fn)`, `map(fn)`, `for_each(fn)`, `reduce(init, fn)`, `any(fn)`, `all(fn)`, `count(fn?)`
   - Utility: `to_list` (alias for entries), `find(fn)` (returns Option of [key, value])

2. **Created `mapFindKey` helper** for DRY key lookup across map methods

3. **Wrote 41 new test functions** in `pkg/interpreter/methods_test.go`:
   - Individual tests for all map methods
   - Lambda and named function tests
   - Method chaining: `filter→map→len`
   - Empty map operations
   - Mutation behavior: `set`, `remove`, `delete`, `clear`, `merge` mutate the receiver
   - Non-mutation: `filter`, `map` return new maps without modifying original
   - Error handling: `merge` with non-map argument
   - Edge cases: empty maps, missing keys, vacuous truth for `all` on empty
   - All 379 tests pass (338 existing + 41 new, zero regressions)

### Important Design Notes

- **Map field access vs method access:** `MapVal` dot access (`m.foo`) first checks for a string key matching the field name, then falls through to the method registry. This means if a map has a key named `"len"`, `m.len` returns the key's value, not the method. Use `m.len()` through the call path which properly resolves via the method registry.
- **Insertion order preserved:** All map methods maintain insertion order since `MapVal` uses parallel `Keys` and `Values` slices.
- **Higher-order map methods pass (key, value):** Unlike list HOFs which pass a single element, map `filter`, `map`, `for_each`, `reduce`, `any`, `all`, `count`, `find` pass both key and value to the callback. `reduce` passes (acc, key, value).

---

## Next Task: Phase 4.1 Chunk 4 — Option/Result Methods

### Goal

Extend the method registry with Option and Result methods to complete Phase 4.1.

### Option Methods to Implement in `pkg/interpreter/methods_option.go`

| Method | Signature | Description | Notes |
|--------|-----------|-------------|-------|
| `is_some()` | `() -> Bool` | Check if Some | |
| `is_none()` | `() -> Bool` | Check if None | |
| `unwrap()` | `() -> T` | Get value or panic | |
| `unwrap_or(default)` | `(T) -> T` | Get value or return default | |
| `unwrap_or_else(fn)` | `(Fn() -> T) -> T` | Get value or call fn | |
| `map(fn)` | `(Fn(T) -> U) -> Option[U]` | Transform inner value | |
| `flat_map(fn)` | `(Fn(T) -> Option[U]) -> Option[U]` | Transform, flatten Option | |
| `filter(fn)` | `(Fn(T) -> Bool) -> Option[T]` | Keep if predicate true | |
| `or_else(fn)` | `(Fn() -> Option[T]) -> Option[T]` | Alternative if None | |
| `and_then(fn)` | `(Fn(T) -> Option[U]) -> Option[U]` | Alias for flat_map | |
| `expect(msg)` | `(String) -> T` | Get value or panic with message | |

### Result Methods to Implement in `pkg/interpreter/methods_result.go` (or same file)

| Method | Signature | Description | Notes |
|--------|-----------|-------------|-------|
| `is_ok()` | `() -> Bool` | Check if Ok | |
| `is_err()` | `() -> Bool` | Check if Err | |
| `unwrap()` | `() -> T` | Get Ok value or panic | |
| `unwrap_err()` | `() -> E` | Get Err value or panic | |
| `unwrap_or(default)` | `(T) -> T` | Get Ok value or default | |
| `unwrap_or_else(fn)` | `(Fn(E) -> T) -> T` | Get Ok or transform Err | |
| `map(fn)` | `(Fn(T) -> U) -> Result[U,E]` | Transform Ok value | |
| `map_err(fn)` | `(Fn(E) -> F) -> Result[T,F]` | Transform Err value | |
| `and_then(fn)` | `(Fn(T) -> Result[U,E]) -> Result[U,E]` | Chain Results | |
| `or_else(fn)` | `(Fn(E) -> Result[T,F]) -> Result[T,F]` | Alternative on Err | |
| `ok()` | `() -> Option[T]` | Convert to Option (drops Err) | |
| `err()` | `() -> Option[E]` | Convert Err to Option | |

### Testing Target

- 25-35 new test functions
- Target total: ~410+ tests

---

## What Comes After (Post Phase 4.1)

After Chunk 4 → Phase 4.1 (Core Runtime Methods) is complete → move to:

1. **Import System & Module Resolution** (Phase 4.2 prerequisite)
2. **Standard Library Foundation** (`std.testing`, `std.json`, `std.io`)
3. **Effect Runtime** (Phase 4.3)

---

## Technical Details

### Key Code Patterns

**Method registration** (established in Chunk 1):
```go
RegisterMethod(TypeMap, "get", func(receiver Value, args []Value) Value {
    m := receiver.(*MapVal)
    // ... implementation
})
```

**Calling Aura functions from Go** (established in Chunk 2):
```go
result := callValue(fn, []Value{key, value})
```

**Map key lookup helper** (established in Chunk 3):
```go
idx := mapFindKey(m, key)
if idx < 0 { /* not found */ }
```

**Returning Options** (common pattern for safe access):
```go
return &OptionVal{IsSome: true, Val: value}  // Some(value)
return &OptionVal{IsSome: false}               // None
```

**Runtime errors** in method implementations:
```go
panic(&RuntimeError{Message: "error message"})
```

### Files Created/Modified

| File | Description |
|------|-------------|
| `pkg/interpreter/methods.go` | Method registry infrastructure |
| `pkg/interpreter/methods_string.go` | 22 string methods |
| `pkg/interpreter/methods_list.go` | 27 list methods + `callValue` + `cmpValues` helpers |
| `pkg/interpreter/methods_map.go` | 24 map methods + `mapFindKey` helper |
| `pkg/interpreter/methods_test.go` | 133 test functions (40 string + 52 list + 41 map) |

### Running Tests

```bash
cd /path/to/aura
go test ./pkg/interpreter/ -v            # interpreter tests only
go test ./...                             # all tests (379 passing)
go test ./pkg/interpreter/ -run TestMap   # map method tests only
go test ./pkg/interpreter/ -run TestOption # option method tests (Chunk 4)
```

---

## References

| Resource | Path |
|----------|------|
| Method registry | `pkg/interpreter/methods.go` |
| String methods | `pkg/interpreter/methods_string.go` |
| List methods | `pkg/interpreter/methods_list.go` |
| Map methods | `pkg/interpreter/methods_map.go` |
| Method tests | `pkg/interpreter/methods_test.go` |
| Interpreter source | `pkg/interpreter/eval.go` (main evaluation loop) |
| Value definitions | `pkg/interpreter/value.go` |
| Existing builtins | `pkg/interpreter/interpreter.go` → `registerBuiltins()` |
| Interpreter tests | `pkg/interpreter/interpreter_test.go` (126 tests) |
| Lexer | `pkg/lexer/lexer.go` |
| Parser | `pkg/parser/parser.go` |
| AST nodes | `pkg/ast/ast.go` |
| ROADMAP | `ROADMAP.md` |
| AI Mission | `AI_MISSION.md` |

---

*This context file was generated at the end of the 2026-03-20 session. It should be updated or replaced at the end of the next session.*
