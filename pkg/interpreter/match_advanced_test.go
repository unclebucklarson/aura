package interpreter

import (
        "strings"
        "testing"
)

// ============================================================
// Guard Clause Tests
// ============================================================

func TestGuardClauseBasic(t *testing.T) {
        src := `module test
fn classify(x: Int) -> String:
    let result = match x:
        n when n > 0 -> "positive"
        n when n < 0 -> "negative"
        _ -> "zero"
    return result
`
        result := runFunc(t, src, "classify", []Value{&IntVal{Val: 5}})
        expectString(t, result, "positive")

        result = runFunc(t, src, "classify", []Value{&IntVal{Val: -3}})
        expectString(t, result, "negative")

        result = runFunc(t, src, "classify", []Value{&IntVal{Val: 0}})
        expectString(t, result, "zero")
}

func TestGuardClauseWithTuplePattern(t *testing.T) {
        src := `module test
fn check(x: Int, y: Int) -> String:
    let point = (x, y)
    let result = match point:
        (a, b) when a == b -> "diagonal"
        (a, b) when a > b -> "right"
        _ -> "other"
    return result
`
        result := runFunc(t, src, "check", []Value{&IntVal{Val: 3}, &IntVal{Val: 3}})
        expectString(t, result, "diagonal")

        result = runFunc(t, src, "check", []Value{&IntVal{Val: 5}, &IntVal{Val: 2}})
        expectString(t, result, "right")

        result = runFunc(t, src, "check", []Value{&IntVal{Val: 1}, &IntVal{Val: 9}})
        expectString(t, result, "other")
}

func TestGuardClauseFallthrough(t *testing.T) {
        // When guard fails, should continue to next arm
        src := `module test
fn check(x: Int) -> String:
    let result = match x:
        n when n > 100 -> "huge"
        n when n > 10 -> "big"
        n when n > 0 -> "small"
        _ -> "non-positive"
    return result
`
        result := runFunc(t, src, "check", []Value{&IntVal{Val: 5}})
        expectString(t, result, "small")

        result = runFunc(t, src, "check", []Value{&IntVal{Val: 50}})
        expectString(t, result, "big")

        result = runFunc(t, src, "check", []Value{&IntVal{Val: 200}})
        expectString(t, result, "huge")
}

func TestGuardClauseWithConstructorPattern(t *testing.T) {
        src := `module test
fn check(x: Int) -> String:
    let val = Some(x)
    let result = match val:
        Some(n) when n > 10 -> "big some"
        Some(n) -> "small some"
        None -> "none"
    return result
`
        result := runFunc(t, src, "check", []Value{&IntVal{Val: 20}})
        expectString(t, result, "big some")

        result = runFunc(t, src, "check", []Value{&IntVal{Val: 3}})
        expectString(t, result, "small some")
}

// ============================================================
// Or Pattern Tests
// ============================================================

func TestOrPatternBasicLiterals(t *testing.T) {
        src := `module test
fn check(x: Int) -> String:
    let result = match x:
        1 | 2 | 3 -> "small"
        4 | 5 | 6 -> "medium"
        _ -> "other"
    return result
`
        result := runFunc(t, src, "check", []Value{&IntVal{Val: 2}})
        expectString(t, result, "small")

        result = runFunc(t, src, "check", []Value{&IntVal{Val: 5}})
        expectString(t, result, "medium")

        result = runFunc(t, src, "check", []Value{&IntVal{Val: 9}})
        expectString(t, result, "other")
}

func TestOrPatternStrings(t *testing.T) {
        src := `module test
fn check(s: String) -> String:
    let result = match s:
        "yes" | "ok" | "true" -> "affirmative"
        "no" | "false" -> "negative"
        _ -> "unknown"
    return result
`
        result := runFunc(t, src, "check", []Value{&StringVal{Val: "ok"}})
        expectString(t, result, "affirmative")

        result = runFunc(t, src, "check", []Value{&StringVal{Val: "no"}})
        expectString(t, result, "negative")

        result = runFunc(t, src, "check", []Value{&StringVal{Val: "maybe"}})
        expectString(t, result, "unknown")
}

func TestOrPatternWithWildcard(t *testing.T) {
        src := `module test
fn check(x: Int) -> String:
    let result = match x:
        0 | 1 -> "binary"
        _ -> "not binary"
    return result
`
        result := runFunc(t, src, "check", []Value{&IntVal{Val: 1}})
        expectString(t, result, "binary")

        result = runFunc(t, src, "check", []Value{&IntVal{Val: 7}})
        expectString(t, result, "not binary")
}

func TestOrPatternWithGuard(t *testing.T) {
        src := `module test
fn check(x: Int) -> String:
    let result = match x:
        1 | 2 | 3 when x > 1 -> "small but > 1"
        1 | 2 | 3 -> "small"
        _ -> "other"
    return result
`
        result := runFunc(t, src, "check", []Value{&IntVal{Val: 3}})
        expectString(t, result, "small but > 1")

        result = runFunc(t, src, "check", []Value{&IntVal{Val: 1}})
        expectString(t, result, "small")
}

func TestOrPatternWithBinding(t *testing.T) {
        src := `module test
fn check(x: Int) -> Int:
    let result = match x:
        n | n -> n * 2
        _ -> 0
    return result
`
        // n binds to the value in whichever alternative matches
        result := runFunc(t, src, "check", []Value{&IntVal{Val: 5}})
        expectInt(t, result, 10)
}

// ============================================================
// Function Parameter Pattern Tests
// ============================================================

func TestFuncParamTuplePattern(t *testing.T) {
        src := `module test
fn add_point((x, y)) -> Int:
    return x + y
`
        result := runFunc(t, src, "add_point", []Value{
                &TupleVal{Elements: []Value{&IntVal{Val: 3}, &IntVal{Val: 4}}},
        })
        expectInt(t, result, 7)
}

func TestFuncParamTuplePatternMultipleArgs(t *testing.T) {
        src := `module test
fn distance((x1, y1), (x2, y2)) -> Int:
    return (x2 - x1) + (y2 - y1)
`
        result := runFunc(t, src, "distance", []Value{
                &TupleVal{Elements: []Value{&IntVal{Val: 1}, &IntVal{Val: 2}}},
                &TupleVal{Elements: []Value{&IntVal{Val: 4}, &IntVal{Val: 6}}},
        })
        expectInt(t, result, 7)
}

func TestFuncParamListPattern(t *testing.T) {
        src := `module test
fn first_elem([x, ...rest]) -> Int:
    return x
`
        result := runFunc(t, src, "first_elem", []Value{
                &ListVal{Elements: []Value{&IntVal{Val: 10}, &IntVal{Val: 20}, &IntVal{Val: 30}}},
        })
        expectInt(t, result, 10)
}

func TestFuncParamPatternMismatch(t *testing.T) {
        src := `module test
fn add_point((x, y)) -> Int:
    return x + y
`
        interp := runModule(t, src)
        _, err := interp.RunFunction("add_point", []Value{&IntVal{Val: 42}})
        if err == nil {
                t.Fatal("expected error on pattern mismatch")
        }
        if !strings.Contains(err.Error(), "does not match pattern") {
                t.Fatalf("unexpected error: %s", err.Error())
        }
}

// ============================================================
// Let Pattern Destructuring Tests
// ============================================================

func TestLetPatternListSpread(t *testing.T) {
        src := `module test
fn check() -> Int:
    let list = [10, 20, 30, 40, 50]
    let [first, ...rest] = list
    return first
`
        result := runFunc(t, src, "check", nil)
        expectInt(t, result, 10)
}

func TestLetPatternListSpreadRest(t *testing.T) {
        src := `module test
fn check() -> Int:
    let list = [10, 20, 30, 40, 50]
    let [first, ...rest] = list
    return rest.len()
`
        result := runFunc(t, src, "check", nil)
        expectInt(t, result, 4)
}

func TestLetPatternConstructorSome(t *testing.T) {
        src := `module test
fn check() -> Int:
    let val = Some(42)
    let Some(x) = val
    return x
`
        result := runFunc(t, src, "check", nil)
        expectInt(t, result, 42)
}

func TestLetPatternConstructorOk(t *testing.T) {
        src := `module test
fn check() -> String:
    let val = Ok("success")
    let Ok(msg) = val
    return msg
`
        result := runFunc(t, src, "check", nil)
        expectString(t, result, "success")
}

func TestLetPatternDestructureFail(t *testing.T) {
        src := `module test
fn check() -> Int:
    let val = None
    let Some(x) = val
    return x
`
        interp := runModule(t, src)
        _, err := interp.RunFunction("check", nil)
        if err == nil {
                t.Fatal("expected error on failed pattern destructure")
        }
        if !strings.Contains(err.Error(), "pattern destructure failed") {
                t.Fatalf("unexpected error: %s", err.Error())
        }
}

func TestLetPatternListExact(t *testing.T) {
        src := `module test
fn check() -> Int:
    let list = [1, 2, 3]
    let [a, b, c] = list
    return a + b + c
`
        result := runFunc(t, src, "check", nil)
        expectInt(t, result, 6)
}

// ============================================================
// Combined / Integration Tests
// ============================================================

func TestGuardWithOrPatternAndBinding(t *testing.T) {
        src := `module test
fn describe(x: Int) -> String:
    let result = match x:
        1 | 2 | 3 -> "tiny"
        n when n < 0 -> "negative"
        n when n > 100 -> "huge"
        _ -> "normal"
    return result
`
        result := runFunc(t, src, "describe", []Value{&IntVal{Val: 2}})
        expectString(t, result, "tiny")

        result = runFunc(t, src, "describe", []Value{&IntVal{Val: -5}})
        expectString(t, result, "negative")

        result = runFunc(t, src, "describe", []Value{&IntVal{Val: 200}})
        expectString(t, result, "huge")

        result = runFunc(t, src, "describe", []Value{&IntVal{Val: 50}})
        expectString(t, result, "normal")
}

func TestLetPatternWithSpreadMiddle(t *testing.T) {
        src := `module test
fn check() -> Int:
    let list = [1, 2, 3, 4, 5]
    let [first, ...middle, last] = list
    return first + last
`
        result := runFunc(t, src, "check", nil)
        expectInt(t, result, 6)
}

func TestOrPatternBooleans(t *testing.T) {
        src := `module test
fn check(x: Bool) -> String:
    let result = match x:
        true | false -> "boolean"
    return result
`
        result := runFunc(t, src, "check", []Value{&BoolVal{Val: true}})
        expectString(t, result, "boolean")

        result = runFunc(t, src, "check", []Value{&BoolVal{Val: false}})
        expectString(t, result, "boolean")
}

func TestFuncParamTupleWithWildcard(t *testing.T) {
        src := `module test
fn second((_, y)) -> Int:
    return y
`
        result := runFunc(t, src, "second", []Value{
                &TupleVal{Elements: []Value{&IntVal{Val: 99}, &IntVal{Val: 42}}},
        })
        expectInt(t, result, 42)
}

func TestLetPatternListSpreadEmpty(t *testing.T) {
        // When spread captures zero elements
        src := `module test
fn check() -> Int:
    let list = [1]
    let [first, ...rest] = list
    return rest.len()
`
        result := runFunc(t, src, "check", nil)
        expectInt(t, result, 0)
}
