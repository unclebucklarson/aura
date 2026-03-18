package interpreter

import (
        "fmt"

        "github.com/unclebucklarson/aura/pkg/ast"
        "github.com/unclebucklarson/aura/pkg/token"
)

// Interpreter executes an Aura module.
type Interpreter struct {
        module *ast.Module
        env    *Environment
}

// New creates a new interpreter for the given module.
func New(module *ast.Module) *Interpreter {
        interp := &Interpreter{
                module: module,
                env:    NewEnvironment(),
        }
        interp.registerBuiltins()
        return interp
}

// registerBuiltins adds built-in functions and constructors to the environment.
func (interp *Interpreter) registerBuiltins() {
        env := interp.env

        // Ok constructor
        env.DefineConst("Ok", &BuiltinFnVal{
                Name: "Ok",
                Fn: func(args []Value) Value {
                        if len(args) != 1 {
                                panic(&RuntimeError{Message: "Ok() requires exactly one argument"})
                        }
                        return &ResultVal{IsOk: true, Val: args[0]}
                },
        })

        // Err constructor
        env.DefineConst("Err", &BuiltinFnVal{
                Name: "Err",
                Fn: func(args []Value) Value {
                        if len(args) != 1 {
                                panic(&RuntimeError{Message: "Err() requires exactly one argument"})
                        }
                        return &ResultVal{IsOk: false, Val: args[0]}
                },
        })

        // Some constructor
        env.DefineConst("Some", &BuiltinFnVal{
                Name: "Some",
                Fn: func(args []Value) Value {
                        if len(args) != 1 {
                                panic(&RuntimeError{Message: "Some() requires exactly one argument"})
                        }
                        return &OptionVal{IsSome: true, Val: args[0]}
                },
        })

        // None value
        env.DefineConst("None", &OptionVal{IsSome: false})

        // print function
        env.DefineConst("print", &BuiltinFnVal{
                Name: "print",
                Fn: func(args []Value) Value {
                        parts := make([]string, len(args))
                        for i, a := range args {
                                parts[i] = a.String()
                        }
                        fmt.Println(joinStrings(parts, " "))
                        return &NoneVal{}
                },
        })

        // len function
        env.DefineConst("len", &BuiltinFnVal{
                Name: "len",
                Fn: func(args []Value) Value {
                        if len(args) != 1 {
                                panic(&RuntimeError{Message: "len() requires exactly one argument"})
                        }
                        switch v := args[0].(type) {
                        case *ListVal:
                                return &IntVal{Val: int64(len(v.Elements))}
                        case *MapVal:
                                return &IntVal{Val: int64(len(v.Keys))}
                        case *SetVal:
                                return &IntVal{Val: int64(len(v.Elements))}
                        case *StringVal:
                                return &IntVal{Val: int64(len(v.Val))}
                        case *TupleVal:
                                return &IntVal{Val: int64(len(v.Elements))}
                        default:
                                panic(&RuntimeError{Message: fmt.Sprintf("len() not supported for %s", valueTypeNames[args[0].Type()])})
                        }
                },
        })

        // str function
        env.DefineConst("str", &BuiltinFnVal{
                Name: "str",
                Fn: func(args []Value) Value {
                        if len(args) != 1 {
                                panic(&RuntimeError{Message: "str() requires exactly one argument"})
                        }
                        return &StringVal{Val: args[0].String()}
                },
        })

        // int function
        env.DefineConst("int", &BuiltinFnVal{
                Name: "int",
                Fn: func(args []Value) Value {
                        if len(args) != 1 {
                                panic(&RuntimeError{Message: "int() requires exactly one argument"})
                        }
                        switch v := args[0].(type) {
                        case *IntVal:
                                return v
                        case *FloatVal:
                                return &IntVal{Val: int64(v.Val)}
                        case *StringVal:
                                n, err := fmt.Sscanf(v.Val, "%d", new(int64))
                                if err != nil || n != 1 {
                                        panic(&RuntimeError{Message: fmt.Sprintf("cannot convert '%s' to int", v.Val)})
                                }
                                var i int64
                                fmt.Sscanf(v.Val, "%d", &i)
                                return &IntVal{Val: i}
                        case *BoolVal:
                                if v.Val {
                                        return &IntVal{Val: 1}
                                }
                                return &IntVal{Val: 0}
                        default:
                                panic(&RuntimeError{Message: fmt.Sprintf("cannot convert %s to int", valueTypeNames[args[0].Type()])})
                        }
                },
        })

        // float function
        env.DefineConst("float", &BuiltinFnVal{
                Name: "float",
                Fn: func(args []Value) Value {
                        if len(args) != 1 {
                                panic(&RuntimeError{Message: "float() requires exactly one argument"})
                        }
                        switch v := args[0].(type) {
                        case *FloatVal:
                                return v
                        case *IntVal:
                                return &FloatVal{Val: float64(v.Val)}
                        default:
                                panic(&RuntimeError{Message: fmt.Sprintf("cannot convert %s to float", valueTypeNames[args[0].Type()])})
                        }
                },
        })

        // range function for for loops
        env.DefineConst("range", &BuiltinFnVal{
                Name: "range",
                Fn: func(args []Value) Value {
                        var start, end, step int64
                        switch len(args) {
                        case 1:
                                e, ok := args[0].(*IntVal)
                                if !ok {
                                        panic(&RuntimeError{Message: "range() requires integer arguments"})
                                }
                                start, end, step = 0, e.Val, 1
                        case 2:
                                s, ok1 := args[0].(*IntVal)
                                e, ok2 := args[1].(*IntVal)
                                if !ok1 || !ok2 {
                                        panic(&RuntimeError{Message: "range() requires integer arguments"})
                                }
                                start, end, step = s.Val, e.Val, 1
                        case 3:
                                s, ok1 := args[0].(*IntVal)
                                e, ok2 := args[1].(*IntVal)
                                st, ok3 := args[2].(*IntVal)
                                if !ok1 || !ok2 || !ok3 {
                                        panic(&RuntimeError{Message: "range() requires integer arguments"})
                                }
                                start, end, step = s.Val, e.Val, st.Val
                        default:
                                panic(&RuntimeError{Message: "range() requires 1-3 arguments"})
                        }
                        if step == 0 {
                                panic(&RuntimeError{Message: "range() step cannot be zero"})
                        }
                        elems := make([]Value, 0)
                        if step > 0 {
                                for i := start; i < end; i += step {
                                        elems = append(elems, &IntVal{Val: i})
                                }
                        } else {
                                for i := start; i > end; i += step {
                                        elems = append(elems, &IntVal{Val: i})
                                }
                        }
                        return &ListVal{Elements: elems}
                },
        })

        // type function
        env.DefineConst("type_of", &BuiltinFnVal{
                Name: "type_of",
                Fn: func(args []Value) Value {
                        if len(args) != 1 {
                                panic(&RuntimeError{Message: "type_of() requires exactly one argument"})
                        }
                        return &StringVal{Val: valueTypeNames[args[0].Type()]}
                },
        })

        // abs function
        env.DefineConst("abs", &BuiltinFnVal{
                Name: "abs",
                Fn: func(args []Value) Value {
                        if len(args) != 1 {
                                panic(&RuntimeError{Message: "abs() requires exactly one argument"})
                        }
                        switch v := args[0].(type) {
                        case *IntVal:
                                if v.Val < 0 {
                                        return &IntVal{Val: -v.Val}
                                }
                                return v
                        case *FloatVal:
                                if v.Val < 0 {
                                        return &FloatVal{Val: -v.Val}
                                }
                                return v
                        default:
                                panic(&RuntimeError{Message: fmt.Sprintf("abs() not supported for %s", valueTypeNames[args[0].Type()])})
                        }
                },
        })

        // min / max
        env.DefineConst("min", &BuiltinFnVal{
                Name: "min",
                Fn: func(args []Value) Value {
                        if len(args) < 2 {
                                panic(&RuntimeError{Message: "min() requires at least 2 arguments"})
                        }
                        result := args[0]
                        for _, a := range args[1:] {
                                if compareRaw(a, result) < 0 {
                                        result = a
                                }
                        }
                        return result
                },
        })

        env.DefineConst("max", &BuiltinFnVal{
                Name: "max",
                Fn: func(args []Value) Value {
                        if len(args) < 2 {
                                panic(&RuntimeError{Message: "max() requires at least 2 arguments"})
                        }
                        result := args[0]
                        for _, a := range args[1:] {
                                if compareRaw(a, result) > 0 {
                                        result = a
                                }
                        }
                        return result
                },
        })
}

func joinStrings(parts []string, sep string) string {
        result := ""
        for i, p := range parts {
                if i > 0 {
                        result += sep
                }
                result += p
        }
        return result
}

func compareRaw(a, b Value) int {
        switch av := a.(type) {
        case *IntVal:
                switch bv := b.(type) {
                case *IntVal:
                        if av.Val < bv.Val {
                                return -1
                        }
                        if av.Val > bv.Val {
                                return 1
                        }
                        return 0
                case *FloatVal:
                        af := float64(av.Val)
                        if af < bv.Val {
                                return -1
                        }
                        if af > bv.Val {
                                return 1
                        }
                        return 0
                }
        case *FloatVal:
                var bf float64
                switch bv := b.(type) {
                case *FloatVal:
                        bf = bv.Val
                case *IntVal:
                        bf = float64(bv.Val)
                }
                if av.Val < bf {
                        return -1
                }
                if av.Val > bf {
                        return 1
                }
                return 0
        }
        return 0
}

// Run executes the module's top-level items.
func (interp *Interpreter) Run() (result Value, err error) {
        defer func() {
                if r := recover(); r != nil {
                        switch e := r.(type) {
                        case *RuntimeError:
                                err = e
                        case returnSignal:
                                result = e.val
                        default:
                                err = fmt.Errorf("runtime panic: %v", r)
                        }
                }
        }()

        result = &NoneVal{}

        // First pass: register all type definitions, functions, constants, enums
        for _, item := range interp.module.Items {
                switch it := item.(type) {
                case *ast.StructDef:
                        fields := make([]string, len(it.Fields))
                        for i, f := range it.Fields {
                                fields[i] = f.Name
                        }
                        interp.env.DefineStruct(it.Name, fields)

                case *ast.EnumDef:
                        variants := make(map[string]int, len(it.Variants))
                        for _, v := range it.Variants {
                                variants[v.Name] = len(v.Fields)
                        }
                        interp.env.DefineEnum(it.Name, variants)

                case *ast.FnDef:
                        fn := &FunctionVal{
                                Name:   it.Name,
                                Params: it.Params,
                                Body:   it.Body,
                                Env:    interp.env,
                        }
                        interp.env.DefineConst(it.Name, fn)

                case *ast.ConstDef:
                        val := EvalExpr(it.Value, interp.env)
                        interp.env.DefineConst(it.Name, val)

                case *ast.TypeDef:
                        // Type aliases: no runtime effect
                case *ast.TraitDef:
                        // Traits: no runtime effect for now
                case *ast.ImplBlock:
                        // Impl blocks: no runtime effect for now
                case *ast.SpecBlock:
                        // Specs: no runtime effect
                case *ast.TestBlock:
                        // Tests: handled separately
                }
        }

        return result, nil
}

// RunFunction calls a named function with the given arguments.
func (interp *Interpreter) RunFunction(name string, args []Value) (result Value, err error) {
        defer func() {
                if r := recover(); r != nil {
                        switch e := r.(type) {
                        case *RuntimeError:
                                err = e
                        case returnSignal:
                                result = e.val
                        default:
                                err = fmt.Errorf("runtime panic: %v", r)
                        }
                }
        }()

        fnVal, ok := interp.env.Get(name)
        if !ok {
                return nil, fmt.Errorf("function '%s' not found", name)
        }

        fn, ok := fnVal.(*FunctionVal)
        if !ok {
                return nil, fmt.Errorf("'%s' is not a function", name)
        }

        var span ast.Node
        if len(fn.Params) > 0 {
                span = fn.Params[0]
        }
        var spanVal token.Span
        if span != nil {
                spanVal = span.GetSpan()
        }
        result = callUserFn(spanVal, fn, args)
        return result, nil
}

// Env returns the interpreter's environment (for testing/REPL).
func (interp *Interpreter) Env() *Environment {
        return interp.env
}
