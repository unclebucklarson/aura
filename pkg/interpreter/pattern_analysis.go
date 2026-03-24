package interpreter

import (
        "fmt"
        "strings"

        "github.com/unclebucklarson/aura/pkg/ast"
        "github.com/unclebucklarson/aura/pkg/token"
)

// --- Warning System ---

// WarningKind classifies the type of warning.
type WarningKind int

const (
        WarnNonExhaustive WarningKind = iota
        WarnUnreachable
        WarnRedundant
)

func (k WarningKind) String() string {
        switch k {
        case WarnNonExhaustive:
                return "Non-exhaustive match"
        case WarnUnreachable:
                return "Unreachable pattern"
        case WarnRedundant:
                return "Redundant pattern"
        default:
                return "Warning"
        }
}

// PatternWarning represents a warning generated during pattern analysis.
type PatternWarning struct {
        Kind       WarningKind
        Span       token.Span
        Message    string
        Suggestion string
}

// String formats the warning for display.
func (w *PatternWarning) String() string {
        var b strings.Builder
        b.WriteString(fmt.Sprintf("Warning: [line %d] %s", w.Span.Start.Line, w.Kind))
        if w.Message != "" {
                b.WriteString("\n  " + w.Message)
        }
        if w.Suggestion != "" {
                b.WriteString("\n  Suggestion: " + w.Suggestion)
        }
        return b.String()
}

// WarningCollector accumulates warnings during analysis.
type WarningCollector struct {
        Warnings []*PatternWarning
}

// NewWarningCollector creates a new empty collector.
func NewWarningCollector() *WarningCollector {
        return &WarningCollector{}
}

// Add appends a warning.
func (wc *WarningCollector) Add(w *PatternWarning) {
        wc.Warnings = append(wc.Warnings, w)
}

// HasWarnings returns true if any warnings were collected.
func (wc *WarningCollector) HasWarnings() bool {
        return len(wc.Warnings) > 0
}

// Clear removes all warnings.
func (wc *WarningCollector) Clear() {
        wc.Warnings = nil
}

// FormatAll returns all warnings as a single string.
func (wc *WarningCollector) FormatAll() string {
        if !wc.HasWarnings() {
                return ""
        }
        var parts []string
        for _, w := range wc.Warnings {
                parts = append(parts, w.String())
        }
        return strings.Join(parts, "\n\n")
}

// --- Pattern Category (for exhaustiveness) ---

// patternCategory describes what a pattern covers.
type patternCategory int

const (
        catCatchAll    patternCategory = iota // wildcard or binding — matches everything
        catLiteral                            // matches a specific literal
        catConstructor                        // matches a specific constructor (Some, None, Ok, Err, enum variant)
        catComplex                            // list/tuple/or pattern — too complex for simple analysis
)

// --- Exhaustiveness Checking ---

// analyzeMatchExprExhaustiveness checks a MatchExpr for exhaustiveness and unreachable patterns.
func analyzeMatchExprExhaustiveness(m *ast.MatchExpr) []*PatternWarning {
        patterns := make([]ast.Pattern, len(m.Arms))
        spans := make([]token.Span, len(m.Arms))
        hasGuards := make([]bool, len(m.Arms))
        for i, arm := range m.Arms {
                patterns[i] = arm.Pattern
                spans[i] = arm.Pattern.GetSpan()
                hasGuards[i] = arm.Guard != nil
        }
        return analyzePatterns(patterns, spans, hasGuards, m.Span)
}

// analyzeMatchStmtExhaustiveness checks a MatchStmt for exhaustiveness and unreachable patterns.
func analyzeMatchStmtExhaustiveness(m *ast.MatchStmt) []*PatternWarning {
        patterns := make([]ast.Pattern, len(m.Cases))
        spans := make([]token.Span, len(m.Cases))
        hasGuards := make([]bool, len(m.Cases))
        for i, c := range m.Cases {
                patterns[i] = c.Pattern
                spans[i] = c.Pattern.GetSpan()
                hasGuards[i] = c.Guard != nil
        }
        return analyzePatterns(patterns, spans, hasGuards, m.Span)
}

// analyzePatterns is the core analysis function.
func analyzePatterns(patterns []ast.Pattern, spans []token.Span, hasGuards []bool, matchSpan token.Span) []*PatternWarning {
        var warnings []*PatternWarning

        if len(patterns) == 0 {
                return warnings
        }

        // Track which pattern indices are reachable
        catchAllIdx := -1 // index of first catch-all pattern (wildcard/binding) without guard

        for i, p := range patterns {
                cat := categorizePattern(p)

                // Check if this pattern is after a catch-all (unreachable)
                if catchAllIdx >= 0 && !hasGuards[catchAllIdx] {
                        warnings = append(warnings, &PatternWarning{
                                Kind:       WarnUnreachable,
                                Span:       spans[i],
                                Message:    fmt.Sprintf("Pattern at line %d is unreachable — shadowed by catch-all pattern at line %d", spans[i].Start.Line, spans[catchAllIdx].Start.Line),
                                Suggestion: "Remove this unreachable pattern or reorder patterns so specific ones come before catch-all",
                        })
                        continue
                }

                // Check for duplicate literal patterns
                if cat == catLiteral {
                        for j := 0; j < i; j++ {
                                if !hasGuards[j] && patternsEquivalent(patterns[j], p) {
                                        warnings = append(warnings, &PatternWarning{
                                                Kind:       WarnRedundant,
                                                Span:       spans[i],
                                                Message:    fmt.Sprintf("Pattern at line %d is redundant — already covered by pattern at line %d", spans[i].Start.Line, spans[j].Start.Line),
                                                Suggestion: "Remove the redundant pattern",
                                        })
                                        break
                                }
                        }
                }

                // Check for duplicate constructor patterns (without fields / with catch-all fields)
                if cat == catConstructor {
                        for j := 0; j < i; j++ {
                                if !hasGuards[j] && patternsEquivalent(patterns[j], p) {
                                        warnings = append(warnings, &PatternWarning{
                                                Kind:       WarnRedundant,
                                                Span:       spans[i],
                                                Message:    fmt.Sprintf("Pattern at line %d is redundant — already covered by pattern at line %d", spans[i].Start.Line, spans[j].Start.Line),
                                                Suggestion: "Remove the redundant pattern",
                                        })
                                        break
                                }
                        }
                }

                // Record the first catch-all
                if cat == catCatchAll && catchAllIdx < 0 {
                        catchAllIdx = i
                }
        }

        // Exhaustiveness checking
        // If there's a catch-all (without guard), the match is exhaustive
        hasCatchAll := false
        for i, p := range patterns {
                if !hasGuards[i] && categorizePattern(p) == catCatchAll {
                        hasCatchAll = true
                        break
                }
        }

        if !hasCatchAll {
                missing := checkExhaustiveness(patterns, hasGuards)
                if len(missing) > 0 {
                        warnings = append(warnings, &PatternWarning{
                                Kind:       WarnNonExhaustive,
                                Span:       matchSpan,
                                Message:    "Missing patterns:\n    - " + strings.Join(missing, "\n    - "),
                                Suggestion: "Add the missing patterns or a catch-all pattern (_) to handle all cases",
                        })
                }
        }

        return warnings
}

// categorizePattern classifies a pattern for analysis.
func categorizePattern(p ast.Pattern) patternCategory {
        switch pat := p.(type) {
        case *ast.WildcardPattern:
                return catCatchAll
        case *ast.BindingPattern:
                return catCatchAll
        case *ast.LiteralPattern:
                return catLiteral
        case *ast.ConstructorPattern:
                return catConstructor
        case *ast.OrPattern:
                // If any alternative is a catch-all, the whole or-pattern is catch-all
                for _, sub := range pat.Patterns {
                        if categorizePattern(sub) == catCatchAll {
                                return catCatchAll
                        }
                }
                return catComplex
        default:
                return catComplex
        }
}

// patternsEquivalent checks if two patterns match the same set of values.
func patternsEquivalent(a, b ast.Pattern) bool {
        switch pa := a.(type) {
        case *ast.LiteralPattern:
                pb, ok := b.(*ast.LiteralPattern)
                if !ok {
                        return false
                }
                return pa.Kind == pb.Kind && pa.Value == pb.Value
        case *ast.ConstructorPattern:
                pb, ok := b.(*ast.ConstructorPattern)
                if !ok {
                        return false
                }
                if pa.TypeName != pb.TypeName {
                        return false
                }
                if len(pa.Fields) != len(pb.Fields) {
                        return false
                }
                for i := range pa.Fields {
                        if !patternsEquivalent(pa.Fields[i], pb.Fields[i]) {
                                return false
                        }
                }
                return true
        case *ast.WildcardPattern:
                _, ok := b.(*ast.WildcardPattern)
                return ok
        case *ast.BindingPattern:
                _, ok := b.(*ast.BindingPattern)
                return ok
        default:
                return false
        }
}

// checkExhaustiveness determines missing patterns based on what constructors appear.
func checkExhaustiveness(patterns []ast.Pattern, hasGuards []bool) []string {
        // Collect all constructor names used (without guards)
        constructors := make(map[string]bool)
        boolValues := make(map[string]bool)
        hasLiterals := false

        for i, p := range patterns {
                if hasGuards[i] {
                        continue // guarded patterns don't count for exhaustiveness
                }
                collectPatternCoverage(p, constructors, boolValues, &hasLiterals)
        }

        var missing []string

        // Check boolean exhaustiveness
        if len(boolValues) > 0 {
                if !boolValues["true"] {
                        missing = append(missing, "true")
                }
                if !boolValues["false"] {
                        missing = append(missing, "false")
                }
        }

        // Check Option exhaustiveness (Some/None)
        hasSome := constructors["Some"]
        hasNone := constructors["None"]
        if hasSome || hasNone {
                if !hasSome {
                        missing = append(missing, "Some(_)")
                }
                if !hasNone {
                        missing = append(missing, "None")
                }
        }

        // Check Result exhaustiveness (Ok/Err)
        hasOk := constructors["Ok"]
        hasErr := constructors["Err"]
        if hasOk || hasErr {
                if !hasOk {
                        missing = append(missing, "Ok(_)")
                }
                if !hasErr {
                        missing = append(missing, "Err(_)")
                }
        }

        // For literal-only patterns with no bool/option/result, suggest wildcard
        if hasLiterals && len(boolValues) == 0 && !hasSome && !hasNone && !hasOk && !hasErr && len(constructors) == 0 {
                missing = append(missing, "_ (wildcard)")
        }

        // If only enum constructors are present (non-standard), suggest wildcard
        hasNonStandard := false
        for k := range constructors {
                if k != "Some" && k != "None" && k != "Ok" && k != "Err" {
                        hasNonStandard = true
                        break
                }
        }
        if hasNonStandard && len(missing) == 0 {
                // We can't verify enum exhaustiveness without type info, so suggest a wildcard
                missing = append(missing, "_ (wildcard for unhandled variants)")
        }

        return missing
}

// collectPatternCoverage gathers coverage info from a single pattern.
func collectPatternCoverage(p ast.Pattern, constructors map[string]bool, boolValues map[string]bool, hasLiterals *bool) {
        switch pat := p.(type) {
        case *ast.LiteralPattern:
                *hasLiterals = true
                if pat.Kind == token.BOOL_LIT {
                        boolValues[pat.Value] = true
                }
        case *ast.ConstructorPattern:
                // Extract variant name (handle dotted names)
                name := pat.TypeName
                if idx := strings.LastIndex(name, "."); idx >= 0 {
                        name = name[idx+1:]
                }
                constructors[name] = true
        case *ast.OrPattern:
                for _, sub := range pat.Patterns {
                        collectPatternCoverage(sub, constructors, boolValues, hasLiterals)
                }
        }
}

// --- Package-level active warning collector ---

// activeWarnings is set by the Interpreter before execution to collect warnings
// from standalone eval functions (evalMatchExpr, execMatchStmt).
var activeWarnings *WarningCollector

// SetActiveWarnings sets the package-level warning collector.
// Called by Interpreter.Run() before execution begins.
func SetActiveWarnings(wc *WarningCollector) {
        activeWarnings = wc
}

// emitWarnings sends warnings to the active collector if one is set.
func emitWarnings(warnings []*PatternWarning) {
        if activeWarnings == nil || len(warnings) == 0 {
                return
        }
        for _, w := range warnings {
                activeWarnings.Add(w)
        }
}

// --- Integration helpers ---

// AnalyzeMatchExpr performs static analysis on a match expression and returns warnings.
func AnalyzeMatchExpr(m *ast.MatchExpr) []*PatternWarning {
        return analyzeMatchExprExhaustiveness(m)
}

// AnalyzeMatchStmt performs static analysis on a match statement and returns warnings.
func AnalyzeMatchStmt(m *ast.MatchStmt) []*PatternWarning {
        return analyzeMatchStmtExhaustiveness(m)
}
