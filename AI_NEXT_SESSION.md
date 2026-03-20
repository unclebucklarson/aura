# AI Next Session Context — Phase 4.1 Chunk 1

> **Purpose:** This file provides complete context for any AI agent (or human developer) to pick up exactly where we left off. Read this before starting any work.
>
> **Last updated:** 2026-03-19
> **Last session by:** DeepAgent (Abacus.AI)

---

## Current Status

### What's Complete

| Phase | Description | Status |
|-------|-------------|--------|
| Phase 1 | Lexer, Parser, AST, Formatter, CLI | ✅ Complete |
| Phase 2 | Type Checker, Symbol Table, Semantic Analysis | ✅ Complete |
| Phase 3 | Tree-walk Interpreter | ✅ Complete |
| Phase 3+ | String Interpolation, Pipeline Operator (`\|>`), Option Chaining (`?.`) | ✅ Complete |

- **Test count:** 246 test functions passing (126 interpreter, 11 lexer, 16 parser, plus checker/symbols/types/formatter)
- **Language version:** Pre-1.0, Phase 3 complete
- **Repository:** `github.com/unclebucklarson/aura`
- **Branch:** `main`

### What Exists Today for Method Calls

The interpreter already has **basic** method dispatch via `evalListMethod()` and `evalStringMethod()` in `pkg/interpreter/eval.go`. These are implemented as inline switch statements returning `BuiltinFnVal` closures:

**Existing String methods:** `len`/`length`, `upper`, `lower`, `contains`, `split`
**Existing List methods:** `len`/`length`, `append`, `contains`
**Existing Map methods:** None (Maps exist as `MapVal` but have no methods)
**Existing Option/Result methods:** None (constructors `Some`, `Ok`, `Err` exist as builtins)

---

## What We Just Did (Session Ending 2026-03-19)

1. **Implemented Phase 3 deferred features:**
   - String interpolation (`"Hello, {name}!"`)
   - Pipeline operator (`value |> transform |> format`)
   - Option chaining (`user?.address?.city`)

2. **Resolved merge conflicts** on `feat/pipeline-operator` branch and merged to `main`

3. **Analyzed the full ROADMAP.md** against current codebase state and produced a prioritized recommendation document (`/home/ubuntu/aura_next_steps_recommendations.md`)

4. **Decided on next task:** Phase 4.1 Chunk 1 — Foundation + String Methods (Option A)

---

## Next Task: Phase 4.1 Chunk 1 — Foundation + String Methods

### Goal

Refactor the method dispatch system from ad-hoc inline switch statements into a clean, extensible architecture, then use that architecture to implement a comprehensive set of string methods.

### Why This First

- **Immediate usability** — AI-generated Aura code immediately reaches for `.trim()`, `.replace()`, `.starts_with()`, etc. Without these, every AI-generated program hits a wall.
- **No architectural dependencies** — Doesn't require imports, modules, or any deferred Phase 2 work.
- **Foundation for everything** — The stdlib, effect runtime, and AI integration all assume these primitives exist.
- **Low risk** — Well-scoped, incremental, fully testable.

---

## Implementation Plan

### Part 1: Method Dispatch Infrastructure (Day 1)

**Problem:** The current approach (switch statements in `evalListMethod` / `evalStringMethod`) doesn't scale. Adding 40+ methods across 4+ types will create unmaintainable spaghetti.

**Solution:** Create a method registry system.

#### Option A: Registry Map (Recommended)

Create a new file `pkg/interpreter/methods.go` with:

```go
// MethodFunc is the signature for all built-in methods.
// receiver is the value the method is called on, args are the call arguments.
type MethodFunc func(receiver Value, args []Value) Value

// methodRegistry maps type → method name → implementation
var methodRegistry = map[ValueType]map[string]MethodFunc{}

// RegisterMethod registers a built-in method for a value type.
func RegisterMethod(vt ValueType, name string, fn MethodFunc) {
    if methodRegistry[vt] == nil {
        methodRegistry[vt] = map[string]MethodFunc{}
    }
    methodRegistry[vt][name] = fn
}

// LookupMethod returns the method function for a given type and name, or nil.
func LookupMethod(vt ValueType, name string) MethodFunc {
    if methods, ok := methodRegistry[vt]; ok {
        return methods[name]
    }
    return nil
}
```

Then refactor the `FieldAccess` evaluation in `eval.go` (~line 550-660) to use:

```go
if method := LookupMethod(obj.Type(), e.Field); method != nil {
    captured := obj // capture for closure
    return &BuiltinFnVal{
        Name: fmt.Sprintf("%s.%s", valueTypeNames[obj.Type()], e.Field),
        Fn:   func(args []Value) Value { return method(captured, args) },
    }
}
```

#### Files to Create/Modify

| File | Action | Description |
|------|--------|-------------|
| `pkg/interpreter/methods.go` | **CREATE** | Method registry, `RegisterMethod`, `LookupMethod` |
| `pkg/interpreter/methods_string.go` | **CREATE** | All string method registrations |
| `pkg/interpreter/methods_list.go` | **CREATE** | Migrate existing + add new list methods (Chunk 2) |
| `pkg/interpreter/methods_map.go` | **CREATE** | Map methods (Chunk 3) |
| `pkg/interpreter/methods_option.go` | **CREATE** | Option/Result methods (Chunk 4) |
| `pkg/interpreter/eval.go` | **MODIFY** | Replace switch statements with registry lookup |
| `pkg/interpreter/value.go` | **MODIFY** | Ensure `ValueType` constants are exported and complete |

### Part 2: String Methods Implementation (Days 2-3)

Implement these string methods in `pkg/interpreter/methods_string.go`:

| Method | Signature | Go Implementation | Notes |
|--------|-----------|-------------------|-------|
| `len()` | `() -> Int` | `len(s.Val)` | Already exists, migrate |
| `upper()` | `() -> String` | `strings.ToUpper(s.Val)` | Already exists, migrate |
| `lower()` | `() -> String` | `strings.ToLower(s.Val)` | Already exists, migrate |
| `contains(sub)` | `(String) -> Bool` | `strings.Contains(s.Val, sub)` | Already exists, migrate |
| `split(sep)` | `(String) -> List[String]` | `strings.Split(s.Val, sep)` | Already exists, migrate |
| `trim()` | `() -> String` | `strings.TrimSpace(s.Val)` | **NEW** |
| `trim_left()` | `() -> String` | `strings.TrimLeft(s.Val, " \t\n\r")` | **NEW** |
| `trim_right()` | `() -> String` | `strings.TrimRight(s.Val, " \t\n\r")` | **NEW** |
| `starts_with(pre)` | `(String) -> Bool` | `strings.HasPrefix(s.Val, pre)` | **NEW** |
| `ends_with(suf)` | `(String) -> Bool` | `strings.HasSuffix(s.Val, suf)` | **NEW** |
| `replace(old, new)` | `(String, String) -> String` | `strings.ReplaceAll(s.Val, old, new)` | **NEW** |
| `replace_first(old, new)` | `(String, String) -> String` | `strings.Replace(s.Val, old, new, 1)` | **NEW** |
| `slice(start, end?)` | `(Int, Int?) -> String` | `s.Val[start:end]` | **NEW**, bounds checking! |
| `index_of(sub)` | `(String) -> Option[Int]` | `strings.Index(s.Val, sub)` → `Some(i)` or `None` | **NEW**, returns Option |
| `chars()` | `() -> List[String]` | Split into individual chars | **NEW** |
| `join(list)` | `(List[String]) -> String` | `strings.Join(...)` | **NEW** — called as static or on separator |
| `repeat(n)` | `(Int) -> String` | `strings.Repeat(s.Val, n)` | **NEW** |
| `is_empty()` | `() -> Bool` | `len(s.Val) == 0` | **NEW** |
| `reverse()` | `() -> String` | Reverse runes | **NEW** |
| `pad_left(n, char?)` | `(Int, String?) -> String` | Left-pad to length n | **NEW** |
| `pad_right(n, char?)` | `(Int, String?) -> String` | Right-pad to length n | **NEW** |

### Part 3: Testing (Day 3-4)

Add tests to `pkg/interpreter/interpreter_test.go` (or create a new `pkg/interpreter/methods_test.go`).

**Testing approach:**
- One test function per method (e.g., `TestStringTrim`, `TestStringStartsWith`)
- Each test function should cover: basic usage, edge cases (empty string, not found), error cases (wrong arg type)
- Tests use the existing `runFunc` / `runModule` helpers
- Target: **20-25 new test functions** for string methods

**Example test pattern:**

```go
func TestStringTrim(t *testing.T) {
    src := `
fn test_trim() -> String
    let s = "  hello  "
    return s.trim()
`
    result := runFunc(t, src, "test_trim", nil)
    expectString(t, result, "hello")
}
```

### Part 4: Migrate Existing Methods (Day 1, alongside Part 1)

Move the 5 existing string methods and 3 existing list methods from inline switch statements in `eval.go` into the new registry. **All 246 existing tests must continue to pass.**

---

## Technical Details

### Key Code Patterns

**Value types** are defined in `pkg/interpreter/value.go`:
```go
type IntVal struct{ Val int64 }
type FloatVal struct{ Val float64 }
type StringVal struct{ Val string }
type BoolVal struct{ Val bool }
type NoneVal struct{}
type ListVal struct{ Elements []Value }
type MapVal struct{ Entries map[string]Value; Order []string }
```

**BuiltinFnVal** wraps Go functions as Aura callable values:
```go
type BuiltinFnVal struct {
    Name string
    Fn   func(args []Value) Value
}
```

**Method dispatch** currently happens in `pkg/interpreter/eval.go` around line 545-660, inside the `FieldAccess` case of `evalExpr`.

**Runtime panics** use `runtimePanic(span, format, args...)` for error reporting.

**Existing builtins** are registered in `pkg/interpreter/interpreter.go` in `registerBuiltins()` — these are *global functions* like `print`, `len`, `str`, `int`, `float`, `range`, `type_of`, `abs`, `min`, `max`.

### Import Paths

```
github.com/unclebucklarson/aura/pkg/interpreter
github.com/unclebucklarson/aura/pkg/lexer
github.com/unclebucklarson/aura/pkg/parser
github.com/unclebucklarson/aura/pkg/ast
```

### Running Tests

```bash
cd /path/to/aura
go test ./pkg/interpreter/ -v          # interpreter tests only
go test ./...                           # all tests
go test ./pkg/interpreter/ -run TestStringTrim  # single test
```

---

## Success Criteria

Chunk 1 is **complete** when:

- [ ] `pkg/interpreter/methods.go` exists with a clean registry system
- [ ] `pkg/interpreter/methods_string.go` exists with 20+ string methods registered
- [ ] All 5 existing string methods migrated to the new registry
- [ ] All 3 existing list methods migrated to the new registry (kept working, even if not extended yet)
- [ ] All 246 existing tests still pass (zero regressions)
- [ ] 20-25 new test functions for string methods pass
- [ ] Total test count: ~270+
- [ ] `evalStringMethod()` and `evalListMethod()` in `eval.go` are removed (replaced by registry)
- [ ] Code is clean, documented, and follows existing patterns

---

## Estimated Timeline

| Day | Focus | Deliverable |
|-----|-------|-------------|
| Day 1 | Infrastructure + migration | `methods.go` registry, refactored `eval.go`, all 246 tests green |
| Day 2 | New string methods (first half) | `trim`, `starts_with`, `ends_with`, `replace`, `slice`, `index_of`, `chars`, `repeat`, `is_empty` |
| Day 3 | New string methods (second half) + tests | `reverse`, `pad_left`, `pad_right`, `replace_first`, `trim_left`, `trim_right`, `join` + all new tests |
| Day 4 | Polish, edge cases, docs | Full test coverage, error messages, update ROADMAP.md, commit |

---

## What Comes After (Chunks 2-4)

| Chunk | Scope | Depends On |
|-------|-------|------------|
| **Chunk 2:** List Methods | `map`, `filter`, `reduce`, `sort`, `reverse`, `flat_map`, `slice`, `get`, `pop`, `push` + tests | Chunk 1 (uses same registry) |
| **Chunk 3:** Map Methods | `get`, `set`, `delete`, `keys`, `values`, `entries`, `contains_key`, `len` + tests | Chunk 1 |
| **Chunk 4:** Option/Result Methods | `unwrap`, `unwrap_or`, `map`, `flat_map`, `is_some`, `is_none`, `is_ok`, `is_err`, `or_else` + tests | Chunk 1 |

After all 4 chunks → Phase 4.1 complete → move to Import System (Priority 2).

---

## References

| Resource | Path |
|----------|------|
| Interpreter source | `pkg/interpreter/eval.go` (main evaluation loop) |
| Value definitions | `pkg/interpreter/value.go` |
| Existing builtins | `pkg/interpreter/interpreter.go` → `registerBuiltins()` |
| Existing method dispatch | `pkg/interpreter/eval.go` lines ~545-660 |
| Interpreter tests | `pkg/interpreter/interpreter_test.go` (126 tests) |
| Lexer | `pkg/lexer/lexer.go` |
| Parser | `pkg/parser/parser.go` |
| AST nodes | `pkg/ast/ast.go` |
| ROADMAP | `ROADMAP.md` |
| AI Mission | `AI_MISSION.md` |
| Recommendations doc | `/home/ubuntu/aura_next_steps_recommendations.md` (external, not in repo) |

---

*This context file was generated at the end of the 2026-03-19 session. It should be updated or replaced at the end of the next session.*
