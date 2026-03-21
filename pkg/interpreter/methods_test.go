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


// =============================================================================
// List Method Tests
// =============================================================================

func TestListLen(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [1, 2, 3]
    return xs.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestListLenEmpty(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = []
    return xs.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 0)
}

func TestListIsEmpty(t *testing.T) {
        src := `
fn test_empty() -> Bool:
    return [].is_empty()

fn test_not_empty() -> Bool:
    return [1].is_empty()
`
        expectBool(t, runFunc(t, src, "test_empty", nil), true)
        expectBool(t, runFunc(t, src, "test_not_empty", nil), false)
}

func TestListContains(t *testing.T) {
        src := `
fn test_yes() -> Bool:
    return [1, 2, 3].contains(2)

fn test_no() -> Bool:
    return [1, 2, 3].contains(5)
`
        expectBool(t, runFunc(t, src, "test_yes", nil), true)
        expectBool(t, runFunc(t, src, "test_no", nil), false)
}

func TestListFirst(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [10, 20, 30]
    match xs.first():
        case Some(v):
            return v
        case _:
            return -1

fn test_empty() -> Int:
    let xs = []
    match xs.first():
        case Some(v):
            return v
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test", nil), 10)
        expectInt(t, runFunc(t, src, "test_empty", nil), -1)
}

func TestListLast(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [10, 20, 30]
    match xs.last():
        case Some(v):
            return v
        case _:
            return -1

fn test_empty() -> Int:
    match [].last():
        case Some(v):
            return v
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test", nil), 30)
        expectInt(t, runFunc(t, src, "test_empty", nil), -1)
}

func TestListGet(t *testing.T) {
        src := `
fn test_valid() -> Int:
    let xs = [10, 20, 30]
    match xs.get(1):
        case Some(v):
            return v
        case _:
            return -1

fn test_negative() -> Int:
    let xs = [10, 20, 30]
    match xs.get(-1):
        case Some(v):
            return v
        case _:
            return -1

fn test_oob() -> Int:
    let xs = [10, 20, 30]
    match xs.get(5):
        case Some(v):
            return v
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test_valid", nil), 20)
        expectInt(t, runFunc(t, src, "test_negative", nil), 30)
        expectInt(t, runFunc(t, src, "test_oob", nil), -1)
}

func TestListPush(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [1, 2]
    xs.push(3)
    return xs.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestListPop(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [1, 2, 3]
    match xs.pop():
        case Some(v):
            return v
        case _:
            return -1

fn test_empty() -> Int:
    let xs = []
    match xs.pop():
        case Some(v):
            return v
        case _:
            return -1

fn test_mutates() -> Int:
    let xs = [10, 20, 30]
    xs.pop()
    return xs.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
        expectInt(t, runFunc(t, src, "test_empty", nil), -1)
        expectInt(t, runFunc(t, src, "test_mutates", nil), 2)
}

func TestListRemove(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [10, 20, 30]
    let removed = xs.remove(1)
    return removed

fn test_length() -> Int:
    let xs = [10, 20, 30]
    xs.remove(0)
    return xs.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 20)
        expectInt(t, runFunc(t, src, "test_length", nil), 2)
}

func TestListReverse(t *testing.T) {
        src := `
fn test() -> List:
    return [1, 2, 3].reverse()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 3 {
                t.Fatalf("expected 3 elements, got %d", len(list.Elements))
        }
        expectInt(t, list.Elements[0], 3)
        expectInt(t, list.Elements[1], 2)
        expectInt(t, list.Elements[2], 1)
}

func TestListReverseEmpty(t *testing.T) {
        src := `
fn test() -> List:
    return [].reverse()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 0 {
                t.Fatalf("expected 0 elements, got %d", len(list.Elements))
        }
}

func TestListSlice(t *testing.T) {
        src := `
fn test_basic() -> Int:
    let xs = [10, 20, 30, 40, 50]
    let s = xs.slice(1, 4)
    return s.len()

fn test_first() -> Int:
    let xs = [10, 20, 30, 40, 50]
    let s = xs.slice(1, 4)
    let opt = s.get(0)
    match opt:
        case Some(v):
            return v
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test_basic", nil), 3)
        expectInt(t, runFunc(t, src, "test_first", nil), 20)
}

func TestListSliceNegative(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [10, 20, 30, 40, 50]
    let s = xs.slice(-2)
    return s.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 2)
}

func TestListSliceNoEnd(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [10, 20, 30, 40, 50]
    let s = xs.slice(2)
    return s.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestListJoin(t *testing.T) {
        src := `
fn test() -> String:
    return ["a", "b", "c"].join(", ")

fn test_empty() -> String:
    return [].join(", ")

fn test_single() -> String:
    return ["hello"].join(", ")
`
        expectString(t, runFunc(t, src, "test", nil), "a, b, c")
        expectString(t, runFunc(t, src, "test_empty", nil), "")
        expectString(t, runFunc(t, src, "test_single", nil), "hello")
}

func TestListJoinInts(t *testing.T) {
        src := `
fn test() -> String:
    return [1, 2, 3].join("-")
`
        expectString(t, runFunc(t, src, "test", nil), "1-2-3")
}

func TestListIndexOf(t *testing.T) {
        src := `
fn test_found() -> Int:
    let xs = [10, 20, 30]
    match xs.index_of(20):
        case Some(i):
            return i
        case _:
            return -1

fn test_not_found() -> Int:
    match [10, 20, 30].index_of(99):
        case Some(i):
            return i
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test_found", nil), 1)
        expectInt(t, runFunc(t, src, "test_not_found", nil), -1)
}

func TestListMap(t *testing.T) {
        src := `
fn double(x: Int) -> Int:
    return x * 2

fn test() -> List:
    return [1, 2, 3].map(double)
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 3 {
                t.Fatalf("expected 3 elements, got %d", len(list.Elements))
        }
        expectInt(t, list.Elements[0], 2)
        expectInt(t, list.Elements[1], 4)
        expectInt(t, list.Elements[2], 6)
}

func TestListMapLambda(t *testing.T) {
        src := `
fn test() -> List:
    return [1, 2, 3].map(|x| -> x * 10)
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        expectInt(t, list.Elements[0], 10)
        expectInt(t, list.Elements[1], 20)
        expectInt(t, list.Elements[2], 30)
}

func TestListMapEmpty(t *testing.T) {
        src := `
fn test() -> Int:
    let result = [].map(|x| -> x * 2)
    return result.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 0)
}

func TestListFilter(t *testing.T) {
        src := `
fn is_even(x: Int) -> Bool:
    return x % 2 == 0

fn test() -> List:
    return [1, 2, 3, 4, 5, 6].filter(is_even)
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 3 {
                t.Fatalf("expected 3 elements, got %d", len(list.Elements))
        }
        expectInt(t, list.Elements[0], 2)
        expectInt(t, list.Elements[1], 4)
        expectInt(t, list.Elements[2], 6)
}

func TestListFilterLambda(t *testing.T) {
        src := `
fn test() -> List:
    return [1, 2, 3, 4, 5].filter(|x| -> x > 3)
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 2 {
                t.Fatalf("expected 2 elements, got %d", len(list.Elements))
        }
        expectInt(t, list.Elements[0], 4)
        expectInt(t, list.Elements[1], 5)
}

func TestListFilterEmpty(t *testing.T) {
        src := `
fn test() -> Int:
    let result = [1, 2, 3].filter(|x| -> x > 100)
    return result.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 0)
}

func TestListReduce(t *testing.T) {
        src := `
fn test() -> Int:
    return [1, 2, 3, 4].reduce(0, |acc, x| -> acc + x)
`
        expectInt(t, runFunc(t, src, "test", nil), 10)
}

func TestListReduceProduct(t *testing.T) {
        src := `
fn test() -> Int:
    return [1, 2, 3, 4].reduce(1, |acc, x| -> acc * x)
`
        expectInt(t, runFunc(t, src, "test", nil), 24)
}

func TestListReduceEmpty(t *testing.T) {
        src := `
fn test() -> Int:
    return [].reduce(42, |acc, x| -> acc + x)
`
        expectInt(t, runFunc(t, src, "test", nil), 42)
}

func TestListForEach(t *testing.T) {
        // for_each returns None; we test it doesn't crash
        src := `
fn test() -> Int:
    let xs = [1, 2, 3]
    xs.for_each(|x| -> x * 2)
    return xs.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestListFlatMap(t *testing.T) {
        src := `
fn test() -> List:
    return [1, 2, 3].flat_map(|x| -> [x, x * 10])
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 6 {
                t.Fatalf("expected 6 elements, got %d", len(list.Elements))
        }
        expectInt(t, list.Elements[0], 1)
        expectInt(t, list.Elements[1], 10)
        expectInt(t, list.Elements[2], 2)
        expectInt(t, list.Elements[3], 20)
        expectInt(t, list.Elements[4], 3)
        expectInt(t, list.Elements[5], 30)
}

func TestListFlatten(t *testing.T) {
        src := `
fn test() -> List:
    return [[1, 2], [3, 4], [5]].flatten()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 5 {
                t.Fatalf("expected 5 elements, got %d", len(list.Elements))
        }
        expectInt(t, list.Elements[0], 1)
        expectInt(t, list.Elements[4], 5)
}

func TestListFlattenMixed(t *testing.T) {
        src := `
fn test() -> List:
    return [[1, 2], 3, [4]].flatten()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 4 {
                t.Fatalf("expected 4 elements, got %d", len(list.Elements))
        }
}

func TestListFlattenEmpty(t *testing.T) {
        src := `
fn test() -> Int:
    return [].flatten().len()
`
        expectInt(t, runFunc(t, src, "test", nil), 0)
}

func TestListAny(t *testing.T) {
        src := `
fn test_yes() -> Bool:
    return [1, 2, 3].any(|x| -> x > 2)

fn test_no() -> Bool:
    return [1, 2, 3].any(|x| -> x > 10)

fn test_empty() -> Bool:
    return [].any(|x| -> x > 0)
`
        expectBool(t, runFunc(t, src, "test_yes", nil), true)
        expectBool(t, runFunc(t, src, "test_no", nil), false)
        expectBool(t, runFunc(t, src, "test_empty", nil), false)
}

func TestListAll(t *testing.T) {
        src := `
fn test_yes() -> Bool:
    return [2, 4, 6].all(|x| -> x % 2 == 0)

fn test_no() -> Bool:
    return [2, 3, 6].all(|x| -> x % 2 == 0)

fn test_empty() -> Bool:
    return [].all(|x| -> x > 0)
`
        expectBool(t, runFunc(t, src, "test_yes", nil), true)
        expectBool(t, runFunc(t, src, "test_no", nil), false)
        expectBool(t, runFunc(t, src, "test_empty", nil), true)
}

func TestListCount(t *testing.T) {
        src := `
fn test_all() -> Int:
    return [1, 2, 3].count()

fn test_pred() -> Int:
    return [1, 2, 3, 4, 5].count(|x| -> x > 3)
`
        expectInt(t, runFunc(t, src, "test_all", nil), 3)
        expectInt(t, runFunc(t, src, "test_pred", nil), 2)
}

func TestListUnique(t *testing.T) {
        src := `
fn test() -> List:
    return [1, 2, 2, 3, 1, 3].unique()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        if len(list.Elements) != 3 {
                t.Fatalf("expected 3 elements, got %d", len(list.Elements))
        }
        expectInt(t, list.Elements[0], 1)
        expectInt(t, list.Elements[1], 2)
        expectInt(t, list.Elements[2], 3)
}

func TestListUniqueEmpty(t *testing.T) {
        src := `
fn test() -> Int:
    return [].unique().len()
`
        expectInt(t, runFunc(t, src, "test", nil), 0)
}

func TestListSum(t *testing.T) {
        src := `
fn test_int() -> Int:
    return [1, 2, 3, 4].sum()

fn test_empty() -> Int:
    return [].sum()
`
        expectInt(t, runFunc(t, src, "test_int", nil), 10)
        expectInt(t, runFunc(t, src, "test_empty", nil), 0)
}

func TestListSumFloat(t *testing.T) {
        src := `
fn test() -> Float:
    return [1.5, 2.5, 3.0].sum()
`
        result := runFunc(t, src, "test", nil)
        fv, ok := result.(*FloatVal)
        if !ok {
                t.Fatalf("expected FloatVal, got %T", result)
        }
        if fv.Val != 7.0 {
                t.Fatalf("expected 7.0, got %g", fv.Val)
        }
}

func TestListMin(t *testing.T) {
        src := `
fn test() -> Int:
    match [3, 1, 2].min():
        case Some(v):
            return v
        case _:
            return -1

fn test_empty() -> Int:
    match [].min():
        case Some(v):
            return v
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test", nil), 1)
        expectInt(t, runFunc(t, src, "test_empty", nil), -1)
}

func TestListMax(t *testing.T) {
        src := `
fn test() -> Int:
    match [3, 1, 2].max():
        case Some(v):
            return v
        case _:
            return -1

fn test_empty() -> Int:
    match [].max():
        case Some(v):
            return v
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
        expectInt(t, runFunc(t, src, "test_empty", nil), -1)
}

func TestListSort(t *testing.T) {
        src := `
fn test() -> List:
    return [3, 1, 4, 1, 5, 9, 2].sort()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        expected := []int64{1, 1, 2, 3, 4, 5, 9}
        if len(list.Elements) != len(expected) {
                t.Fatalf("expected %d elements, got %d", len(expected), len(list.Elements))
        }
        for i, e := range expected {
                expectInt(t, list.Elements[i], e)
        }
}

func TestListSortStrings(t *testing.T) {
        src := `
fn test() -> List:
    return ["banana", "apple", "cherry"].sort()
`
        result := runFunc(t, src, "test", nil)
        list := result.(*ListVal)
        expectString(t, list.Elements[0], "apple")
        expectString(t, list.Elements[1], "banana")
        expectString(t, list.Elements[2], "cherry")
}

func TestListSortEmpty(t *testing.T) {
        src := `
fn test() -> Int:
    return [].sort().len()
`
        expectInt(t, runFunc(t, src, "test", nil), 0)
}

func TestListZip(t *testing.T) {
        src := `
fn test() -> Int:
    let zipped = [1, 2, 3].zip(["a", "b", "c"])
    return zipped.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestListZipUnequal(t *testing.T) {
        src := `
fn test() -> Int:
    let zipped = [1, 2, 3].zip(["a", "b"])
    return zipped.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 2)
}

func TestListEnumerate(t *testing.T) {
        src := `
fn test() -> Int:
    let items = ["a", "b", "c"].enumerate()
    return items.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestListMethodChaining(t *testing.T) {
        src := `
fn test() -> Int:
    return [1, 2, 3, 4, 5, 6].filter(|x| -> x % 2 == 0).map(|x| -> x * 10).sum()
`
        expectInt(t, runFunc(t, src, "test", nil), 120)
}

func TestListMethodChainingComplex(t *testing.T) {
        src := `
fn test() -> String:
    return [3, 1, 2].sort().map(|x| -> x * 2).join(", ")
`
        expectString(t, runFunc(t, src, "test", nil), "2, 4, 6")
}

func TestListMapFilterReduce(t *testing.T) {
        src := `
fn test() -> Int:
    return [1, 2, 3, 4, 5].filter(|x| -> x > 2).map(|x| -> x * x).reduce(0, |a, b| -> a + b)
`
        // filter: [3, 4, 5], map: [9, 16, 25], reduce: 50
        expectInt(t, runFunc(t, src, "test", nil), 50)
}

func TestListReverseDoesNotMutate(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [1, 2, 3]
    let rev = xs.reverse()
    let opt = xs.get(0)
    match opt:
        case Some(v):
            return v
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test", nil), 1)
}

func TestListSortDoesNotMutate(t *testing.T) {
        src := `
fn test() -> Int:
    let xs = [3, 1, 2]
    let sorted = xs.sort()
    let opt = xs.get(0)
    match opt:
        case Some(v):
            return v
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}



// =============================================================================
// Map Method Tests
// =============================================================================

func TestMapLen(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2, "c": 3}
    return m.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestMapLenEmpty(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {}
    return m.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 0)
}

func TestMapLengthAlias(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"x": 10}
    return m.length()
`
        expectInt(t, runFunc(t, src, "test", nil), 1)
}

func TestMapSizeAlias(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2}
    return m.size()
`
        expectInt(t, runFunc(t, src, "test", nil), 2)
}

func TestMapIsEmpty(t *testing.T) {
        src := `
fn test_empty() -> Bool:
    let m = {}
    return m.is_empty()

fn test_not_empty() -> Bool:
    let m = {"a": 1}
    return m.is_empty()
`
        expectBool(t, runFunc(t, src, "test_empty", nil), true)
        expectBool(t, runFunc(t, src, "test_not_empty", nil), false)
}

func TestMapKeys(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2, "c": 3}
    let ks = m.keys()
    return ks.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestMapValues(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 10, "b": 20}
    let vs = m.values()
    return vs.sum()
`
        // Values should be 10 and 20 (insertion order)
        expectInt(t, runFunc(t, src, "test", nil), 30)
}

func TestMapEntries(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 10, "b": 20}
    let es = m.entries()
    return es.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 2)
}

func TestMapEntriesContent(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"x": 42}
    let es = m.entries()
    match es.get(0):
        case Some(pair):
            return pair[1]
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test", nil), 42)
}

func TestMapHas(t *testing.T) {
        src := `
fn test_found() -> Bool:
    let m = {"a": 1, "b": 2}
    return m.has("a")

fn test_not_found() -> Bool:
    let m = {"a": 1, "b": 2}
    return m.has("z")
`
        expectBool(t, runFunc(t, src, "test_found", nil), true)
        expectBool(t, runFunc(t, src, "test_not_found", nil), false)
}

func TestMapContainsKey(t *testing.T) {
        src := `
fn test() -> Bool:
    let m = {"a": 1}
    return m.contains_key("a")
`
        expectBool(t, runFunc(t, src, "test", nil), true)
}

func TestMapContainsValue(t *testing.T) {
        src := `
fn test_found() -> Bool:
    let m = {"a": 1, "b": 2}
    return m.contains_value(2)

fn test_not_found() -> Bool:
    let m = {"a": 1, "b": 2}
    return m.contains_value(99)
`
        expectBool(t, runFunc(t, src, "test_found", nil), true)
        expectBool(t, runFunc(t, src, "test_not_found", nil), false)
}

func TestMapGet(t *testing.T) {
        src := `
fn test_found() -> Int:
    let m = {"a": 42, "b": 99}
    match m.get("a"):
        case Some(v):
            return v
        case _:
            return -1

fn test_missing() -> Int:
    let m = {"a": 42}
    match m.get("z"):
        case Some(v):
            return v
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test_found", nil), 42)
        expectInt(t, runFunc(t, src, "test_missing", nil), -1)
}

func TestMapGetOr(t *testing.T) {
        src := `
fn test_found() -> Int:
    let m = {"a": 42}
    return m.get_or("a", 0)

fn test_default() -> Int:
    let m = {"a": 42}
    return m.get_or("z", 999)
`
        expectInt(t, runFunc(t, src, "test_found", nil), 42)
        expectInt(t, runFunc(t, src, "test_default", nil), 999)
}

func TestMapSet(t *testing.T) {
        src := `
fn test_add() -> Int:
    let m = {"a": 1}
    m.set("b", 2)
    return m.len()

fn test_update() -> Int:
    let m = {"a": 1}
    m.set("a", 99)
    return m["a"]
`
        expectInt(t, runFunc(t, src, "test_add", nil), 2)
        expectInt(t, runFunc(t, src, "test_update", nil), 99)
}

func TestMapRemove(t *testing.T) {
        src := `
fn test_remove_existing() -> Int:
    let m = {"a": 10, "b": 20}
    match m.remove("a"):
        case Some(v):
            return v
        case _:
            return -1

fn test_remove_missing() -> Int:
    let m = {"a": 10}
    match m.remove("z"):
        case Some(v):
            return v
        case _:
            return -1

fn test_remove_shrinks() -> Int:
    let m = {"a": 10, "b": 20, "c": 30}
    m.remove("b")
    return m.len()
`
        expectInt(t, runFunc(t, src, "test_remove_existing", nil), 10)
        expectInt(t, runFunc(t, src, "test_remove_missing", nil), -1)
        expectInt(t, runFunc(t, src, "test_remove_shrinks", nil), 2)
}

func TestMapDelete(t *testing.T) {
        src := `
fn test_exists() -> Bool:
    let m = {"a": 1, "b": 2}
    return m.delete("a")

fn test_not_exists() -> Bool:
    let m = {"a": 1}
    return m.delete("z")
`
        expectBool(t, runFunc(t, src, "test_exists", nil), true)
        expectBool(t, runFunc(t, src, "test_not_exists", nil), false)
}

func TestMapClear(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2, "c": 3}
    m.clear()
    return m.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 0)
}

func TestMapMerge(t *testing.T) {
        src := `
fn test_size() -> Int:
    let m1 = {"a": 1, "b": 2}
    let m2 = {"c": 3, "d": 4}
    m1.merge(m2)
    return m1.len()

fn test_overwrite() -> Int:
    let m1 = {"a": 1, "b": 2}
    let m2 = {"b": 99, "c": 3}
    m1.merge(m2)
    return m1["b"]
`
        expectInt(t, runFunc(t, src, "test_size", nil), 4)
        expectInt(t, runFunc(t, src, "test_overwrite", nil), 99)
}

func TestMapFilter(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2, "c": 3, "d": 4}
    let filtered = m.filter(|k, v| -> v > 2)
    return filtered.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 2)
}

func TestMapFilterDoesNotMutate(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2, "c": 3}
    let _ = m.filter(|k, v| -> v > 1)
    return m.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestMapMap(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2}
    let doubled = m.map(|k, v| -> v * 2)
    return doubled["a"]
`
        expectInt(t, runFunc(t, src, "test", nil), 2)
}

func TestMapMapDoesNotMutate(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2}
    let _ = m.map(|k, v| -> v * 10)
    return m["a"]
`
        expectInt(t, runFunc(t, src, "test", nil), 1)
}

func TestMapForEach(t *testing.T) {
        // for_each returns None; we test it doesn't crash and iterates all entries
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2, "c": 3}
    m.for_each(|k, v| -> v)
    return m.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestMapReduce(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2, "c": 3}
    return m.reduce(0, |acc, k, v| -> acc + v)
`
        expectInt(t, runFunc(t, src, "test", nil), 6)
}

func TestMapAny(t *testing.T) {
        src := `
fn test_true() -> Bool:
    let m = {"a": 1, "b": 5, "c": 3}
    return m.any(|k, v| -> v > 4)

fn test_false() -> Bool:
    let m = {"a": 1, "b": 2}
    return m.any(|k, v| -> v > 10)
`
        expectBool(t, runFunc(t, src, "test_true", nil), true)
        expectBool(t, runFunc(t, src, "test_false", nil), false)
}

func TestMapAll(t *testing.T) {
        src := `
fn test_true() -> Bool:
    let m = {"a": 1, "b": 2, "c": 3}
    return m.all(|k, v| -> v > 0)

fn test_false() -> Bool:
    let m = {"a": 1, "b": 0, "c": 3}
    return m.all(|k, v| -> v > 0)
`
        expectBool(t, runFunc(t, src, "test_true", nil), true)
        expectBool(t, runFunc(t, src, "test_false", nil), false)
}

func TestMapCount(t *testing.T) {
        src := `
fn test_no_predicate() -> Int:
    let m = {"a": 1, "b": 2, "c": 3}
    return m.count()

fn test_with_predicate() -> Int:
    let m = {"a": 1, "b": 2, "c": 3, "d": 4}
    return m.count(|k, v| -> v > 2)
`
        expectInt(t, runFunc(t, src, "test_no_predicate", nil), 3)
        expectInt(t, runFunc(t, src, "test_with_predicate", nil), 2)
}

func TestMapToList(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2}
    let pairs = m.to_list()
    return pairs.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 2)
}

func TestMapFind(t *testing.T) {
        src := `
fn test_found() -> Int:
    let m = {"a": 1, "b": 42, "c": 3}
    match m.find(|k, v| -> v == 42):
        case Some(pair):
            return pair[1]
        case _:
            return -1

fn test_not_found() -> Int:
    let m = {"a": 1, "b": 2}
    match m.find(|k, v| -> v == 99):
        case Some(pair):
            return pair[1]
        case _:
            return -1
`
        expectInt(t, runFunc(t, src, "test_found", nil), 42)
        expectInt(t, runFunc(t, src, "test_not_found", nil), -1)
}

func TestMapMethodChaining(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2, "c": 3, "d": 4}
    return m.filter(|k, v| -> v > 1).map(|k, v| -> v * 10).len()
`
        expectInt(t, runFunc(t, src, "test", nil), 3)
}

func TestMapEmptyOperations(t *testing.T) {
        src := `
fn test_keys() -> Int:
    let m = {}
    return m.keys().len()

fn test_values() -> Int:
    let m = {}
    return m.values().len()

fn test_entries() -> Int:
    let m = {}
    return m.entries().len()

fn test_filter() -> Int:
    let m = {}
    return m.filter(|k, v| -> true).len()
`
        expectInt(t, runFunc(t, src, "test_keys", nil), 0)
        expectInt(t, runFunc(t, src, "test_values", nil), 0)
        expectInt(t, runFunc(t, src, "test_entries", nil), 0)
        expectInt(t, runFunc(t, src, "test_filter", nil), 0)
}

func TestMapSetAndGet(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {}
    m.set("x", 100)
    m.set("y", 200)
    return m.get_or("x", 0) + m.get_or("y", 0)
`
        expectInt(t, runFunc(t, src, "test", nil), 300)
}

func TestMapMergeDoesNotMutateSource(t *testing.T) {
        src := `
fn test() -> Int:
    let m1 = {"a": 1}
    let m2 = {"b": 2}
    m1.merge(m2)
    return m2.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 1)
}

func TestMapWithNamedFunction(t *testing.T) {
        src := `
fn big_value(k: String, v: Int) -> Bool:
    return v >= 10

fn test() -> Int:
    let m = {"a": 5, "b": 15, "c": 20}
    return m.filter(big_value).len()
`
        expectInt(t, runFunc(t, src, "test", nil), 2)
}

func TestMapRegistryExists(t *testing.T) {
        if LookupMethod(TypeMap, "len") == nil {
                t.Fatal("expected Map.len to be registered")
        }
        if LookupMethod(TypeMap, "keys") == nil {
                t.Fatal("expected Map.keys to be registered")
        }
        if LookupMethod(TypeMap, "filter") == nil {
                t.Fatal("expected Map.filter to be registered")
        }
        if LookupMethod(TypeMap, "merge") == nil {
                t.Fatal("expected Map.merge to be registered")
        }
}

func TestMapMergeErrorOnNonMap(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1}
    m.merge([1, 2, 3])
    return 0
`
        expectRuntimeError(t, src, "test", "Map.merge argument must be a Map")
}

func TestMapReduceStringConcat(t *testing.T) {
        src := `
fn test() -> String:
    let m = {"a": 1, "b": 2}
    return m.reduce("", |acc, k, v| -> acc + k)
`
        // Keys are "a" and "b" in insertion order
        expectString(t, runFunc(t, src, "test", nil), "ab")
}

func TestMapClearThenSet(t *testing.T) {
        src := `
fn test() -> Int:
    let m = {"a": 1, "b": 2}
    m.clear()
    m.set("x", 42)
    return m.len()
`
        expectInt(t, runFunc(t, src, "test", nil), 1)
}

func TestMapAnyOnEmpty(t *testing.T) {
        src := `
fn test() -> Bool:
    let m = {}
    return m.any(|k, v| -> true)
`
        expectBool(t, runFunc(t, src, "test", nil), false)
}

func TestMapAllOnEmpty(t *testing.T) {
        src := `
fn test() -> Bool:
    let m = {}
    return m.all(|k, v| -> false)
`
        // all on empty is vacuously true
        expectBool(t, runFunc(t, src, "test", nil), true)
}
