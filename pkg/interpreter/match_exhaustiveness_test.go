package interpreter

import (
	"strings"
	"testing"
)

// =============================================================================
// Exhaustiveness Tests
// =============================================================================

// --- Boolean Exhaustiveness ---

func TestExhaustivenessBoolean_MissingFalse(t *testing.T) {
	src := `module test
fn check(b: Bool) -> String:
    let result = match b:
        true -> "yes"
    return result
`
	interp := runModule(t, src)
	// Should still work (match panics if no match, but we're testing warnings)
	result, err := interp.RunFunction("check", []Value{&BoolVal{Val: true}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectString(t, result, "yes")

	// Check warnings
	warnings := interp.Warnings()
	if !warnings.HasWarnings() {
		t.Fatal("expected non-exhaustive warning for missing 'false'")
	}
	found := false
	for _, w := range warnings.Warnings {
		if w.Kind == WarnNonExhaustive && strings.Contains(w.Message, "false") {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected warning about missing 'false', got: %s", warnings.FormatAll())
	}
}

func TestExhaustivenessBoolean_MissingTrue(t *testing.T) {
	src := `module test
fn check(b: Bool) -> String:
    let result = match b:
        false -> "no"
    return result
`
	interp := runModule(t, src)
	result, err := interp.RunFunction("check", []Value{&BoolVal{Val: false}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectString(t, result, "no")

	warnings := interp.Warnings()
	found := false
	for _, w := range warnings.Warnings {
		if w.Kind == WarnNonExhaustive && strings.Contains(w.Message, "true") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected warning about missing 'true', got: %s", warnings.FormatAll())
	}
}

func TestExhaustivenessBoolean_Complete(t *testing.T) {
	src := `module test
fn check(b: Bool) -> String:
    let result = match b:
        true -> "yes"
        false -> "no"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&BoolVal{Val: true}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// No non-exhaustive warning expected
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnNonExhaustive {
			t.Fatalf("unexpected non-exhaustive warning: %s", w.String())
		}
	}
}

// --- Option Exhaustiveness ---

func TestExhaustivenessOption_MissingNone(t *testing.T) {
	src := `module test
fn check(x: Option) -> String:
    let result = match x:
        Some(v) -> "has value"
    return result
`
	interp := runModule(t, src)
	result, err := interp.RunFunction("check", []Value{&OptionVal{IsSome: true, Val: &IntVal{Val: 42}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectString(t, result, "has value")

	found := false
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnNonExhaustive && strings.Contains(w.Message, "None") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected warning about missing 'None'")
	}
}

func TestExhaustivenessOption_MissingSome(t *testing.T) {
	src := `module test
fn check(x: Option) -> String:
    let result = match x:
        None -> "empty"
    return result
`
	interp := runModule(t, src)
	result, err := interp.RunFunction("check", []Value{&OptionVal{IsSome: false}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectString(t, result, "empty")

	found := false
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnNonExhaustive && strings.Contains(w.Message, "Some") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected warning about missing 'Some(_)'")
	}
}

func TestExhaustivenessOption_Complete(t *testing.T) {
	src := `module test
fn check(x: Option) -> String:
    let result = match x:
        Some(v) -> "has"
        None -> "empty"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&OptionVal{IsSome: true, Val: &IntVal{Val: 1}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnNonExhaustive {
			t.Fatalf("unexpected non-exhaustive warning: %s", w.String())
		}
	}
}

// --- Result Exhaustiveness ---

func TestExhaustivenessResult_MissingErr(t *testing.T) {
	src := `module test
fn check(x: Result) -> String:
    let result = match x:
        Ok(v) -> "ok"
    return result
`
	interp := runModule(t, src)
	result, err := interp.RunFunction("check", []Value{&ResultVal{IsOk: true, Val: &IntVal{Val: 1}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectString(t, result, "ok")

	found := false
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnNonExhaustive && strings.Contains(w.Message, "Err") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected warning about missing 'Err(_)'")
	}
}

func TestExhaustivenessResult_MissingOk(t *testing.T) {
	src := `module test
fn check(x: Result) -> String:
    let result = match x:
        Err(e) -> "error"
    return result
`
	interp := runModule(t, src)
	result, err := interp.RunFunction("check", []Value{&ResultVal{IsOk: false, Val: &StringVal{Val: "oops"}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectString(t, result, "error")

	found := false
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnNonExhaustive && strings.Contains(w.Message, "Ok") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected warning about missing 'Ok(_)'")
	}
}

func TestExhaustivenessResult_Complete(t *testing.T) {
	src := `module test
fn check(x: Result) -> String:
    let result = match x:
        Ok(v) -> "ok"
        Err(e) -> "err"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&ResultVal{IsOk: true, Val: &IntVal{Val: 1}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnNonExhaustive {
			t.Fatalf("unexpected non-exhaustive warning: %s", w.String())
		}
	}
}

// --- Literal Exhaustiveness ---

func TestExhaustivenessLiteral_NoCatchAll(t *testing.T) {
	src := `module test
fn check(x: Int) -> String:
    let result = match x:
        1 -> "one"
        2 -> "two"
    return result
`
	interp := runModule(t, src)
	result, err := interp.RunFunction("check", []Value{&IntVal{Val: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectString(t, result, "one")

	found := false
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnNonExhaustive && strings.Contains(w.Message, "wildcard") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected warning about missing wildcard pattern")
	}
}

func TestExhaustivenessLiteral_WithCatchAll(t *testing.T) {
	src := `module test
fn check(x: Int) -> String:
    let result = match x:
        1 -> "one"
        _ -> "other"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&IntVal{Val: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnNonExhaustive {
			t.Fatalf("unexpected non-exhaustive warning: %s", w.String())
		}
	}
}

// =============================================================================
// Unreachable Pattern Tests
// =============================================================================

func TestUnreachable_WildcardBeforeSpecific(t *testing.T) {
	src := `module test
fn check(x: Int) -> String:
    let result = match x:
        _ -> "catch-all"
        1 -> "one"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&IntVal{Val: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnUnreachable {
			found = true
		}
	}
	if !found {
		t.Fatal("expected unreachable pattern warning")
	}
}

func TestUnreachable_BindingBeforeLiteral(t *testing.T) {
	src := `module test
fn check(x: Int) -> String:
    let result = match x:
        y -> "bound"
        42 -> "forty-two"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&IntVal{Val: 42}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnUnreachable {
			found = true
		}
	}
	if !found {
		t.Fatal("expected unreachable pattern warning for literal after binding")
	}
}

func TestUnreachable_DuplicatePatterns(t *testing.T) {
	src := `module test
fn check(x: Int) -> String:
    let result = match x:
        1 -> "first"
        1 -> "second"
        _ -> "other"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&IntVal{Val: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnRedundant {
			found = true
		}
	}
	if !found {
		t.Fatal("expected redundant pattern warning for duplicate literal")
	}
}

func TestUnreachable_MultipleAfterCatchAll(t *testing.T) {
	src := `module test
fn check(x: Int) -> String:
    let result = match x:
        _ -> "catch-all"
        1 -> "one"
        2 -> "two"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&IntVal{Val: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	unreachableCount := 0
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnUnreachable {
			unreachableCount++
		}
	}
	if unreachableCount < 2 {
		t.Fatalf("expected at least 2 unreachable warnings, got %d", unreachableCount)
	}
}

// =============================================================================
// Warning System Tests
// =============================================================================

func TestWarningCollector_Empty(t *testing.T) {
	wc := NewWarningCollector()
	if wc.HasWarnings() {
		t.Fatal("expected no warnings")
	}
	if wc.FormatAll() != "" {
		t.Fatal("expected empty format")
	}
}

func TestWarningCollector_AddAndFormat(t *testing.T) {
	wc := NewWarningCollector()
	wc.Add(&PatternWarning{
		Kind:       WarnNonExhaustive,
		Message:    "missing patterns",
		Suggestion: "add catch-all",
	})
	if !wc.HasWarnings() {
		t.Fatal("expected warnings")
	}
	if len(wc.Warnings) != 1 {
		t.Fatalf("expected 1 warning, got %d", len(wc.Warnings))
	}
	formatted := wc.FormatAll()
	if !strings.Contains(formatted, "Non-exhaustive match") {
		t.Fatalf("expected 'Non-exhaustive match' in formatted output, got: %s", formatted)
	}
	if !strings.Contains(formatted, "missing patterns") {
		t.Fatalf("expected message in formatted output")
	}
}

func TestWarningCollector_Clear(t *testing.T) {
	wc := NewWarningCollector()
	wc.Add(&PatternWarning{Kind: WarnRedundant, Message: "test"})
	wc.Clear()
	if wc.HasWarnings() {
		t.Fatal("expected no warnings after clear")
	}
}

func TestWarningCollector_MultipleWarnings(t *testing.T) {
	wc := NewWarningCollector()
	wc.Add(&PatternWarning{Kind: WarnNonExhaustive, Message: "first"})
	wc.Add(&PatternWarning{Kind: WarnUnreachable, Message: "second"})
	wc.Add(&PatternWarning{Kind: WarnRedundant, Message: "third"})
	if len(wc.Warnings) != 3 {
		t.Fatalf("expected 3 warnings, got %d", len(wc.Warnings))
	}
	formatted := wc.FormatAll()
	if !strings.Contains(formatted, "first") || !strings.Contains(formatted, "second") || !strings.Contains(formatted, "third") {
		t.Fatalf("expected all messages in output: %s", formatted)
	}
}

// =============================================================================
// Pattern Analysis Unit Tests
// =============================================================================

func TestPatternAnalysis_NoWarningsForGoodMatch(t *testing.T) {
	src := `module test
fn check(x: Int) -> String:
    let result = match x:
        1 -> "one"
        2 -> "two"
        _ -> "other"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&IntVal{Val: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should have no unreachable/redundant warnings, and no exhaustiveness warning (has wildcard)
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnUnreachable || w.Kind == WarnRedundant || w.Kind == WarnNonExhaustive {
			t.Fatalf("unexpected warning: %s", w.String())
		}
	}
}

func TestPatternAnalysis_GuardedCatchAllDoesNotShadow(t *testing.T) {
	// A guarded catch-all should NOT prevent later patterns from being reachable
	src := `module test
fn check(x: Int) -> String:
    let result = match x:
        y when y > 10 -> "big"
        1 -> "one"
        _ -> "other"
    return result
`
	interp := runModule(t, src)
	result, err := interp.RunFunction("check", []Value{&IntVal{Val: 1}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectString(t, result, "one")

	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnUnreachable {
			t.Fatalf("unexpected unreachable warning — guarded catch-all should not shadow: %s", w.String())
		}
	}
}

func TestPatternAnalysis_WarningKindStrings(t *testing.T) {
	if WarnNonExhaustive.String() != "Non-exhaustive match" {
		t.Fatalf("unexpected string for WarnNonExhaustive: %s", WarnNonExhaustive.String())
	}
	if WarnUnreachable.String() != "Unreachable pattern" {
		t.Fatalf("unexpected string for WarnUnreachable: %s", WarnUnreachable.String())
	}
	if WarnRedundant.String() != "Redundant pattern" {
		t.Fatalf("unexpected string for WarnRedundant: %s", WarnRedundant.String())
	}
}

func TestPatternAnalysis_DuplicateConstructorWarning(t *testing.T) {
	src := `module test
fn check(x: Option) -> String:
    let result = match x:
        Some(v) -> "first"
        Some(w) -> "second"
        None -> "empty"
    return result
`
	interp := runModule(t, src)
	_, err := interp.RunFunction("check", []Value{&OptionVal{IsSome: true, Val: &IntVal{Val: 1}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	found := false
	for _, w := range interp.Warnings().Warnings {
		if w.Kind == WarnRedundant {
			found = true
		}
	}
	if !found {
		t.Fatal("expected redundant warning for duplicate Some(_) pattern")
	}
}
