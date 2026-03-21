package interpreter

import (
        "testing"
)

// =============================================================================
// Method Registry Tests
// =============================================================================

func TestMethodRegistry(t *testing.T) {
        // Verify string methods are registered
        if LookupMethod(TypeString, "len") == nil {
                t.Fatal("expected String.len to be registered")
        }
        if LookupMethod(TypeString, "trim") == nil {
                t.Fatal("expected String.trim to be registered")
        }
        // Verify list methods are registered
        if LookupMethod(TypeList, "len") == nil {
                t.Fatal("expected List.len to be registered")
        }
        if LookupMethod(TypeList, "append") == nil {
                t.Fatal("expected List.append to be registered")
        }
        // Verify unknown type/method returns nil
        if LookupMethod(TypeInt, "nonexistent") != nil {
                t.Fatal("expected nil for unregistered method")
        }
}

// =============================================================================
// String Method Tests
// =============================================================================

func TestStringLen(t *testing.T) {
        src := `
fn test() -> Int:
    let s = "hello"
    return s.len()
`
        result := runFunc(t, src, "test", nil)
        expectInt(t, result, 5)
}

func TestStringLenEmpty(t *testing.T) {
        src := `
fn test() -> Int:
    let s = ""
    return s.len()
`
        result := runFunc(t, src, "test", nil)
        expectInt(t, result, 0)
}

func TestStringLenUnicode(t *testing.T) {
        src := `
fn test() -> Int:
    let s = "héllo"
    return s.len()
`
        result := runFunc(t, src, "test", nil)
        expectInt(t, result, 5) // rune count, not byte count
}

func TestStringUpper(t *testing.T) {
        src := `
fn test() -> String:
    return "hello world".upper()
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "HELLO WORLD")
}

func TestStringToUpper(t *testing.T) {
        src := `
fn test() -> String:
    return "hello".to_upper()
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "HELLO")
}

func TestStringLower(t *testing.T) {
        src := `
fn test() -> String:
    return "HELLO WORLD".lower()
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "hello world")
}

func TestStringToLower(t *testing.T) {
        src := `
fn test() -> String:
    return "HELLO".to_lower()
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "hello")
}

func TestStringContains(t *testing.T) {
        src := `
fn test_yes() -> Bool:
    return "hello world".contains("world")

fn test_no() -> Bool:
    return "hello world".contains("xyz")
`
        result := runFunc(t, src, "test_yes", nil)
        expectBool(t, result, true)
        result = runFunc(t, src, "test_no", nil)
        expectBool(t, result, false)
}

func TestStringSplit(t *testing.T) {
        src := `
fn test() -> List:
    return "a,b,c".split(",")
`
        result := runFunc(t, src, "test", nil)
        list, ok := result.(*ListVal)
        if !ok {
                t.Fatalf("expected ListVal, got %T", result)
        }
        if len(list.Elements) != 3 {
                t.Fatalf("expected 3 elements, got %d", len(list.Elements))
        }
        expectString(t, list.Elements[0], "a")
        expectString(t, list.Elements[1], "b")
        expectString(t, list.Elements[2], "c")
}

func TestStringSplitDefault(t *testing.T) {
        src := `
fn test() -> List:
    return "hello world foo".split()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 3 {
                t.Fatalf("expected 3 elements, got %d", len(list.Elements))
        }
        expectString(t, list.Elements[0], "hello")
        expectString(t, list.Elements[1], "world")
        expectString(t, list.Elements[2], "foo")
}

func TestStringTrim(t *testing.T) {
        src := `
fn test() -> String:
    let s = "  hello  "
    return s.trim()
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "hello")
}

func TestStringTrimWhitespace(t *testing.T) {
        src := `
fn test() -> String:
    return "   hello   ".trim()
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "hello")
}

func TestStringTrimEmpty(t *testing.T) {
        src := `
fn test() -> String:
    return "".trim()
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "")
}

func TestStringTrimLeft(t *testing.T) {
        src := `
fn test() -> String:
    return "  hello  ".trim_left()
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "hello  ")
}

func TestStringTrimRight(t *testing.T) {
        src := `
fn test() -> String:
    return "  hello  ".trim_right()
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "  hello")
}

func TestStringStartsWith(t *testing.T) {
        src := `
fn test_yes() -> Bool:
    return "hello world".starts_with("hello")

fn test_no() -> Bool:
    return "hello world".starts_with("world")

fn test_empty() -> Bool:
    return "hello".starts_with("")
`
        expectBool(t, runFunc(t, src, "test_yes", nil), true)
        expectBool(t, runFunc(t, src, "test_no", nil), false)
        expectBool(t, runFunc(t, src, "test_empty", nil), true)
}

func TestStringEndsWith(t *testing.T) {
        src := `
fn test_yes() -> Bool:
    return "hello world".ends_with("world")

fn test_no() -> Bool:
    return "hello world".ends_with("hello")

fn test_empty() -> Bool:
    return "hello".ends_with("")
`
        expectBool(t, runFunc(t, src, "test_yes", nil), true)
        expectBool(t, runFunc(t, src, "test_no", nil), false)
        expectBool(t, runFunc(t, src, "test_empty", nil), true)
}

func TestStringReplace(t *testing.T) {
        src := `
fn test() -> String:
    return "hello world world".replace("world", "aura")
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "hello aura aura")
}

func TestStringReplaceFirst(t *testing.T) {
        src := `
fn test() -> String:
    return "hello world world".replace_first("world", "aura")
`
        result := runFunc(t, src, "test", nil)
        expectString(t, result, "hello aura world")
}

func TestStringSlice(t *testing.T) {
        src := `
fn test_basic() -> String:
    return "hello world".slice(0, 5)

fn test_from() -> String:
    return "hello world".slice(6)

fn test_negative() -> String:
    return "hello world".slice(-5)

fn test_neg_end() -> String:
    return "hello world".slice(0, -6)
`
        expectString(t, runFunc(t, src, "test_basic", nil), "hello")
        expectString(t, runFunc(t, src, "test_from", nil), "world")
        expectString(t, runFunc(t, src, "test_negative", nil), "world")
        expectString(t, runFunc(t, src, "test_neg_end", nil), "hello")
}

func TestStringSliceEmpty(t *testing.T) {
        src := `
fn test() -> String:
    return "hello".slice(5, 5)
`
        expectString(t, runFunc(t, src, "test", nil), "")
}

func TestStringSliceBoundsCheck(t *testing.T) {
        src := `
fn test() -> String:
    return "hello".slice(0, 100)
`
        expectString(t, runFunc(t, src, "test", nil), "hello")
}

func TestStringIndexOf(t *testing.T) {
        src := `
fn test_found() -> Int:
    let result = "hello world".index_of("world")
    match result:
        case Some(i):
            return i
        case _:
            return -1

fn test_not_found() -> Bool:
    let result = "hello world".index_of("xyz")
    match result:
        case Some(_):
            return false
        case _:
            return true
`
        expectInt(t, runFunc(t, src, "test_found", nil), 6)
        expectBool(t, runFunc(t, src, "test_not_found", nil), true)
}

func TestStringChars(t *testing.T) {
        src := `
fn test() -> List:
    return "abc".chars()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 3 {
                t.Fatalf("expected 3 elements, got %d", len(list.Elements))
        }
        expectString(t, list.Elements[0], "a")
        expectString(t, list.Elements[1], "b")
        expectString(t, list.Elements[2], "c")
}

func TestStringCharsEmpty(t *testing.T) {
        src := `
fn test() -> List:
    return "".chars()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 0 {
                t.Fatalf("expected 0 elements, got %d", len(list.Elements))
        }
}

func TestStringJoin(t *testing.T) {
        src := `
fn test() -> String:
    let items = ["a", "b", "c"]
    return ", ".join(items)
`
        expectString(t, runFunc(t, src, "test", nil), "a, b, c")
}

func TestStringJoinEmpty(t *testing.T) {
        src := `
fn test() -> String:
    let items = ["x", "y"]
    return "".join(items)
`
        expectString(t, runFunc(t, src, "test", nil), "xy")
}

func TestStringRepeat(t *testing.T) {
        src := `
fn test() -> String:
    return "ab".repeat(3)
`
        expectString(t, runFunc(t, src, "test", nil), "ababab")
}

func TestStringRepeatZero(t *testing.T) {
        src := `
fn test() -> String:
    return "hello".repeat(0)
`
        expectString(t, runFunc(t, src, "test", nil), "")
}

func TestStringIsEmpty(t *testing.T) {
        src := `
fn test_empty() -> Bool:
    return "".is_empty()

fn test_not_empty() -> Bool:
    return "hello".is_empty()
`
        expectBool(t, runFunc(t, src, "test_empty", nil), true)
        expectBool(t, runFunc(t, src, "test_not_empty", nil), false)
}

func TestStringReverse(t *testing.T) {
        src := `
fn test() -> String:
    return "hello".reverse()
`
        expectString(t, runFunc(t, src, "test", nil), "olleh")
}

func TestStringReverseEmpty(t *testing.T) {
        src := `
fn test() -> String:
    return "".reverse()
`
        expectString(t, runFunc(t, src, "test", nil), "")
}

func TestStringPadLeft(t *testing.T) {
        src := `
fn test() -> String:
    return "42".pad_left(5, "0")

fn test_default() -> String:
    return "hi".pad_left(5)
`
        expectString(t, runFunc(t, src, "test", nil), "00042")
        expectString(t, runFunc(t, src, "test_default", nil), "   hi")
}

func TestStringPadLeftNoOp(t *testing.T) {
        src := `
fn test() -> String:
    return "hello".pad_left(3)
`
        expectString(t, runFunc(t, src, "test", nil), "hello")
}

func TestStringPadRight(t *testing.T) {
        src := `
fn test() -> String:
    return "42".pad_right(5, "0")

fn test_default() -> String:
    return "hi".pad_right(5)
`
        expectString(t, runFunc(t, src, "test", nil), "42000")
        expectString(t, runFunc(t, src, "test_default", nil), "hi   ")
}

func TestStringPadRightNoOp(t *testing.T) {
        src := `
fn test() -> String:
    return "hello".pad_right(3)
`
        expectString(t, runFunc(t, src, "test", nil), "hello")
}

func TestStringMethodChaining(t *testing.T) {
        src := `
fn test() -> String:
    return "  Hello World  ".trim().lower()
`
        expectString(t, runFunc(t, src, "test", nil), "hello world")
}

func TestStringMethodChainingComplex(t *testing.T) {
        src := `
fn test() -> String:
    return "  hello world  ".trim().replace("world", "aura").upper()
`
        expectString(t, runFunc(t, src, "test", nil), "HELLO AURA")
}

func TestStringMethodErrorNoMethod(t *testing.T) {
        src := `
fn test() -> String:
    return "hello".nonexistent_method()
`
        expectRuntimeError(t, src, "test", "cannot access field 'nonexistent_method'")
}
