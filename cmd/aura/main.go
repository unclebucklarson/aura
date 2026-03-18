// Package main provides the CLI entry point for the Aura toolchain.
//
// Usage:
//
//      aura format <file.aura>        - Parse and reformat an Aura source file
//      aura parse  <file.aura>        - Parse and display the AST (for debugging)
//      aura check  [--json] <file.aura> - Type-check an Aura source file
package main

import (
        "fmt"
        "os"
        "strings"

        "github.com/unclebucklarson/aura/pkg/ast"
        "github.com/unclebucklarson/aura/pkg/checker"
        "github.com/unclebucklarson/aura/pkg/formatter"
        "github.com/unclebucklarson/aura/pkg/lexer"
        "github.com/unclebucklarson/aura/pkg/parser"
)

func main() {
        if len(os.Args) < 3 {
                printUsage()
                os.Exit(1)
        }

        command := os.Args[1]

        switch command {
        case "format":
                filePath := os.Args[2]
                src, err := readFile(filePath)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "error: %v\n", err)
                        os.Exit(1)
                }
                os.Exit(runFormat(src, filePath))

        case "parse":
                filePath := os.Args[2]
                src, err := readFile(filePath)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "error: %v\n", err)
                        os.Exit(1)
                }
                os.Exit(runParse(src, filePath))

        case "check":
                jsonOutput := false
                filePath := ""
                for _, arg := range os.Args[2:] {
                        if arg == "--json" {
                                jsonOutput = true
                        } else {
                                filePath = arg
                        }
                }
                if filePath == "" {
                        fmt.Fprintln(os.Stderr, "error: no file specified")
                        printUsage()
                        os.Exit(1)
                }
                src, err := readFile(filePath)
                if err != nil {
                        fmt.Fprintf(os.Stderr, "error: %v\n", err)
                        os.Exit(1)
                }
                os.Exit(runCheck(src, filePath, jsonOutput))

        default:
                fmt.Fprintf(os.Stderr, "error: unknown command %q\n", command)
                printUsage()
                os.Exit(1)
        }
}

func readFile(path string) (string, error) {
        data, err := os.ReadFile(path)
        if err != nil {
                return "", fmt.Errorf("cannot read file %q: %v", path, err)
        }
        return string(data), nil
}

func printUsage() {
        fmt.Fprintln(os.Stderr, "Usage: aura <command> [options] <file.aura>")
        fmt.Fprintln(os.Stderr, "")
        fmt.Fprintln(os.Stderr, "Commands:")
        fmt.Fprintln(os.Stderr, "  format  Parse and reformat an Aura source file")
        fmt.Fprintln(os.Stderr, "  parse   Parse and display the AST (for debugging)")
        fmt.Fprintln(os.Stderr, "  check   Type-check an Aura source file")
        fmt.Fprintln(os.Stderr, "")
        fmt.Fprintln(os.Stderr, "Options for 'check':")
        fmt.Fprintln(os.Stderr, "  --json  Output diagnostics in JSON format (AI-parseable)")
}

func runFormat(src, file string) int {
        // Lex
        l := lexer.New(src, file)
        tokens, lexErrors := l.Tokenize()
        if len(lexErrors) > 0 {
                for _, e := range lexErrors {
                        fmt.Fprintf(os.Stderr, "%s:%s\n", file, e)
                }
                return 1
        }

        // Parse
        p := parser.New(tokens, file)
        module, parseErrors := p.Parse()
        if len(parseErrors) > 0 {
                for _, e := range parseErrors {
                        fmt.Fprintf(os.Stderr, "%s:%s\n", file, e)
                }
                return 1
        }

        // Format
        f := formatter.New()
        output := f.Format(module)
        fmt.Print(output)
        return 0
}

func runParse(src, file string) int {
        // Lex
        l := lexer.New(src, file)
        tokens, lexErrors := l.Tokenize()
        if len(lexErrors) > 0 {
                for _, e := range lexErrors {
                        fmt.Fprintf(os.Stderr, "%s:%s\n", file, e)
                }
                return 1
        }

        // Print tokens
        fmt.Println("=== Tokens ===")
        for _, tok := range tokens {
                fmt.Printf("  %s\n", tok)
        }

        // Parse
        p := parser.New(tokens, file)
        module, parseErrors := p.Parse()
        if len(parseErrors) > 0 {
                fmt.Println("\n=== Parse Errors ===")
                for _, e := range parseErrors {
                        fmt.Fprintf(os.Stderr, "%s:%s\n", file, e)
                }
                return 1
        }

        // Print AST
        fmt.Println("\n=== AST ===")
        printAST(module, 0)
        return 0
}

func runCheck(src, file string, jsonOutput bool) int {
        // Lex
        l := lexer.New(src, file)
        tokens, lexErrors := l.Tokenize()
        if len(lexErrors) > 0 {
                for _, e := range lexErrors {
                        fmt.Fprintf(os.Stderr, "%s:%s\n", file, e)
                }
                return 1
        }

        // Parse
        p := parser.New(tokens, file)
        module, parseErrors := p.Parse()
        if len(parseErrors) > 0 {
                for _, e := range parseErrors {
                        fmt.Fprintf(os.Stderr, "%s:%s\n", file, e)
                }
                return 1
        }

        // Type-check
        c := checker.New(module)
        errs := c.Check()

        if len(errs) == 0 {
                if jsonOutput {
                        fmt.Println(checker.FormatErrorsJSON(errs))
                } else {
                        fmt.Println("✓ No type errors found.")
                }
                return 0
        }

        // Output errors
        if jsonOutput {
                fmt.Println(checker.FormatErrorsJSON(errs))
        } else {
                fmt.Fprint(os.Stderr, checker.FormatErrors(errs))
        }

        // Return 1 if any error-severity issues, 0 for warnings only
        for _, e := range errs {
                if e.Severity == checker.SeverityError {
                        return 1
                }
        }
        return 0
}

func printAST(node ast.Node, depth int) {
        indent := strings.Repeat("  ", depth)

        switch n := node.(type) {
        case *ast.Module:
                fmt.Printf("%sModule", indent)
                if n.Name != nil {
                        fmt.Printf(" %s", n.Name.String())
                }
                fmt.Println()
                for _, imp := range n.Imports {
                        printAST(imp, depth+1)
                }
                for _, item := range n.Items {
                        printAST(item, depth+1)
                }

        case *ast.ImportNode:
                fmt.Printf("%sImport %s", indent, n.Path.String())
                if n.Alias != "" {
                        fmt.Printf(" as %s", n.Alias)
                }
                if n.Names != nil {
                        fmt.Printf(" names=%v", n.Names)
                }
                fmt.Println()

        case *ast.TypeDef:
                fmt.Printf("%sTypeDef %s%s = ", indent, visStr(n.Visibility), n.Name)
                printTypeExpr(n.Body)
                fmt.Println()

        case *ast.StructDef:
                fmt.Printf("%sStructDef %s%s\n", indent, visStr(n.Visibility), n.Name)
                for _, f := range n.Fields {
                        fmt.Printf("%s  Field %s%s: ", indent, visStr(f.Visibility), f.Name)
                        printTypeExpr(f.TypeExpr)
                        if f.Default != nil {
                                fmt.Printf(" = ...")
                        }
                        fmt.Println()
                }

        case *ast.EnumDef:
                fmt.Printf("%sEnumDef %s%s\n", indent, visStr(n.Visibility), n.Name)
                for _, v := range n.Variants {
                        fmt.Printf("%s  Variant %s", indent, v.Name)
                        if len(v.Fields) > 0 {
                                fmt.Printf("(")
                                for i, f := range v.Fields {
                                        if i > 0 {
                                                fmt.Printf(", ")
                                        }
                                        printTypeExpr(f)
                                }
                                fmt.Printf(")")
                        }
                        fmt.Println()
                }

        case *ast.SpecBlock:
                fmt.Printf("%sSpecBlock %s\n", indent, n.Name)
                if n.Doc != "" {
                        fmt.Printf("%s  doc: %q\n", indent, n.Doc)
                }
                for _, inp := range n.Inputs {
                        fmt.Printf("%s  input: %s: ", indent, inp.Name)
                        printTypeExpr(inp.TypeExpr)
                        fmt.Println()
                }
                for _, g := range n.Guarantees {
                        fmt.Printf("%s  guarantee: %q\n", indent, g.Condition)
                }
                if len(n.Effects) > 0 {
                        fmt.Printf("%s  effects: %v\n", indent, n.Effects)
                }
                for _, e := range n.Errors {
                        fmt.Printf("%s  error: %s\n", indent, e.TypeName)
                }

        case *ast.FnDef:
                fmt.Printf("%sFnDef %s%s(", indent, visStr(n.Visibility), n.Name)
                for i, p := range n.Params {
                        if i > 0 {
                                fmt.Printf(", ")
                        }
                        fmt.Printf("%s", p.Name)
                        if p.TypeExpr != nil {
                                fmt.Printf(": ")
                                printTypeExpr(p.TypeExpr)
                        }
                }
                fmt.Printf(")")
                if n.ReturnType != nil {
                        fmt.Printf(" -> ")
                        printTypeExpr(n.ReturnType)
                }
                if len(n.Effects) > 0 {
                        fmt.Printf(" with %s", strings.Join(n.Effects, ", "))
                }
                if n.Satisfies != "" {
                        fmt.Printf(" satisfies %s", n.Satisfies)
                }
                fmt.Println()
                for _, stmt := range n.Body {
                        printStmt(stmt, depth+1)
                }

        case *ast.ConstDef:
                fmt.Printf("%sConstDef %s%s\n", indent, visStr(n.Visibility), n.Name)

        case *ast.TestBlock:
                fmt.Printf("%sTestBlock %q\n", indent, n.Name)
                for _, stmt := range n.Body {
                        printStmt(stmt, depth+1)
                }

        case *ast.TraitDef:
                fmt.Printf("%sTraitDef %s%s\n", indent, visStr(n.Visibility), n.Name)

        case *ast.ImplBlock:
                fmt.Printf("%sImplBlock", indent)
                if n.TraitName != "" {
                        fmt.Printf(" %s for", n.TraitName)
                }
                fmt.Println()
        }
}

func printStmt(stmt ast.Statement, depth int) {
        indent := strings.Repeat("  ", depth)
        switch s := stmt.(type) {
        case *ast.LetStmt:
                fmt.Printf("%sLet %s\n", indent, s.Name)
        case *ast.ReturnStmt:
                fmt.Printf("%sReturn\n", indent)
        case *ast.IfStmt:
                fmt.Printf("%sIf\n", indent)
        case *ast.MatchStmt:
                fmt.Printf("%sMatch\n", indent)
        case *ast.ForStmt:
                fmt.Printf("%sFor %s\n", indent, s.Variable)
        case *ast.WhileStmt:
                fmt.Printf("%sWhile\n", indent)
        case *ast.ExprStmt:
                fmt.Printf("%sExprStmt\n", indent)
        case *ast.AssertStmt:
                fmt.Printf("%sAssert\n", indent)
        case *ast.BreakStmt:
                fmt.Printf("%sBreak\n", indent)
        case *ast.ContinueStmt:
                fmt.Printf("%sContinue\n", indent)
        case *ast.AssignStmt:
                fmt.Printf("%sAssign\n", indent)
        case *ast.WithStmt:
                fmt.Printf("%sWith\n", indent)
        default:
                fmt.Printf("%s<stmt>\n", indent)
        }
}

func printTypeExpr(te ast.TypeExpr) {
        switch t := te.(type) {
        case *ast.NamedType:
                fmt.Printf("%s", t.Name)
                if len(t.Args) > 0 {
                        fmt.Printf("[")
                        for i, a := range t.Args {
                                if i > 0 {
                                        fmt.Printf(", ")
                                }
                                printTypeExpr(a)
                        }
                        fmt.Printf("]")
                }
        case *ast.QualifiedType:
                fmt.Printf("%s.%s", t.Qualifier, t.Name)
        case *ast.UnionType:
                printTypeExpr(t.Left)
                fmt.Printf(" | ")
                printTypeExpr(t.Right)
        case *ast.ListType:
                fmt.Printf("[")
                printTypeExpr(t.Element)
                fmt.Printf("]")
        case *ast.OptionType:
                printTypeExpr(t.Inner)
                fmt.Printf("?")
        case *ast.RefinementType:
                printTypeExpr(t.Base)
                fmt.Printf(" where ...")
        case *ast.StringLitType:
                fmt.Printf("%q", t.Value)
        default:
                fmt.Printf("<type>")
        }
}

func visStr(v ast.Visibility) string {
        if v == ast.Public {
                return "pub "
        }
        return ""
}
