package interpreter

import "fmt"

// TestRegistration holds a registered test from std.testing.
type TestRegistration struct {
	Name string
	Fn   Value
}

// Global test registry for std.testing
var registeredTests []TestRegistration

// createStdTestingExports creates exports for the std.testing module.
func createStdTestingExports() map[string]Value {
	exports := make(map[string]Value)

	// assert(condition, message?) - Assertion function
	exports["assert"] = &BuiltinFnVal{
		Name: "testing.assert",
		Fn: func(args []Value) Value {
			if len(args) < 1 || len(args) > 2 {
				panic(&RuntimeError{Message: "testing.assert() requires 1-2 arguments (condition, message?)"})
			}
			if !IsTruthy(args[0]) {
				msg := "assertion failed"
				if len(args) >= 2 {
					if s, ok := args[1].(*StringVal); ok {
						msg = s.Val
					} else {
						msg = "assertion failed: " + args[1].String()
					}
				}
				panic(&RuntimeError{Message: msg})
			}
			return &BoolVal{Val: true}
		},
	}

	// assert_eq(actual, expected) - Equality assertion
	exports["assert_eq"] = &BuiltinFnVal{
		Name: "testing.assert_eq",
		Fn: func(args []Value) Value {
			if len(args) < 2 || len(args) > 3 {
				panic(&RuntimeError{Message: "testing.assert_eq() requires 2-3 arguments (actual, expected, message?)"})
			}
			actual := args[0]
			expected := args[1]
			if !Equal(actual, expected) {
				msg := fmt.Sprintf("assertion failed: expected %s, got %s", valueRepr(expected), valueRepr(actual))
				if len(args) >= 3 {
					if s, ok := args[2].(*StringVal); ok {
						msg = fmt.Sprintf("%s: expected %s, got %s", s.Val, valueRepr(expected), valueRepr(actual))
					}
				}
				panic(&RuntimeError{Message: msg})
			}
			return &BoolVal{Val: true}
		},
	}

	// assert_ne(actual, expected) - Inequality assertion
	exports["assert_ne"] = &BuiltinFnVal{
		Name: "testing.assert_ne",
		Fn: func(args []Value) Value {
			if len(args) < 2 || len(args) > 3 {
				panic(&RuntimeError{Message: "testing.assert_ne() requires 2-3 arguments (actual, expected, message?)"})
			}
			actual := args[0]
			unexpected := args[1]
			if Equal(actual, unexpected) {
				msg := fmt.Sprintf("assertion failed: values should not be equal, got %s", valueRepr(actual))
				if len(args) >= 3 {
					if s, ok := args[2].(*StringVal); ok {
						msg = fmt.Sprintf("%s: values should not be equal, got %s", s.Val, valueRepr(actual))
					}
				}
				panic(&RuntimeError{Message: msg})
			}
			return &BoolVal{Val: true}
		},
	}

	// assert_true(value) - Assert value is truthy
	exports["assert_true"] = &BuiltinFnVal{
		Name: "testing.assert_true",
		Fn: func(args []Value) Value {
			if len(args) < 1 || len(args) > 2 {
				panic(&RuntimeError{Message: "testing.assert_true() requires 1-2 arguments"})
			}
			if !IsTruthy(args[0]) {
				msg := fmt.Sprintf("expected truthy value, got %s", valueRepr(args[0]))
				if len(args) >= 2 {
					if s, ok := args[1].(*StringVal); ok {
						msg = s.Val
					}
				}
				panic(&RuntimeError{Message: msg})
			}
			return &BoolVal{Val: true}
		},
	}

	// assert_false(value) - Assert value is falsy
	exports["assert_false"] = &BuiltinFnVal{
		Name: "testing.assert_false",
		Fn: func(args []Value) Value {
			if len(args) < 1 || len(args) > 2 {
				panic(&RuntimeError{Message: "testing.assert_false() requires 1-2 arguments"})
			}
			if IsTruthy(args[0]) {
				msg := fmt.Sprintf("expected falsy value, got %s", valueRepr(args[0]))
				if len(args) >= 2 {
					if s, ok := args[1].(*StringVal); ok {
						msg = s.Val
					}
				}
				panic(&RuntimeError{Message: msg})
			}
			return &BoolVal{Val: true}
		},
	}

	// assert_none(value) - Assert value is None
	exports["assert_none"] = &BuiltinFnVal{
		Name: "testing.assert_none",
		Fn: func(args []Value) Value {
			if len(args) != 1 {
				panic(&RuntimeError{Message: "testing.assert_none() requires exactly 1 argument"})
			}
			isNone := false
			switch v := args[0].(type) {
			case *NoneVal:
				isNone = true
			case *OptionVal:
				isNone = !v.IsSome
			}
			if !isNone {
				panic(&RuntimeError{Message: fmt.Sprintf("expected None, got %s", valueRepr(args[0]))})
			}
			return &BoolVal{Val: true}
		},
	}

	// assert_some(value) - Assert value is Some
	exports["assert_some"] = &BuiltinFnVal{
		Name: "testing.assert_some",
		Fn: func(args []Value) Value {
			if len(args) != 1 {
				panic(&RuntimeError{Message: "testing.assert_some() requires exactly 1 argument"})
			}
			opt, ok := args[0].(*OptionVal)
			if !ok || !opt.IsSome {
				panic(&RuntimeError{Message: fmt.Sprintf("expected Some, got %s", valueRepr(args[0]))})
			}
			return opt.Val
		},
	}

	// assert_ok(value) - Assert value is Ok result
	exports["assert_ok"] = &BuiltinFnVal{
		Name: "testing.assert_ok",
		Fn: func(args []Value) Value {
			if len(args) != 1 {
				panic(&RuntimeError{Message: "testing.assert_ok() requires exactly 1 argument"})
			}
			res, ok := args[0].(*ResultVal)
			if !ok || !res.IsOk {
				panic(&RuntimeError{Message: fmt.Sprintf("expected Ok, got %s", valueRepr(args[0]))})
			}
			return res.Val
		},
	}

	// assert_err(value) - Assert value is Err result
	exports["assert_err"] = &BuiltinFnVal{
		Name: "testing.assert_err",
		Fn: func(args []Value) Value {
			if len(args) != 1 {
				panic(&RuntimeError{Message: "testing.assert_err() requires exactly 1 argument"})
			}
			res, ok := args[0].(*ResultVal)
			if !ok || res.IsOk {
				panic(&RuntimeError{Message: fmt.Sprintf("expected Err, got %s", valueRepr(args[0]))})
			}
			return res.Val
		},
	}

	// test(name, fn) - Test registration
	exports["test"] = &BuiltinFnVal{
		Name: "testing.test",
		Fn: func(args []Value) Value {
			if len(args) != 2 {
				panic(&RuntimeError{Message: "testing.test() requires 2 arguments (name, fn)"})
			}
			name, ok := args[0].(*StringVal)
			if !ok {
				panic(&RuntimeError{Message: "testing.test() first argument must be a string"})
			}
			registeredTests = append(registeredTests, TestRegistration{
				Name: name.Val,
				Fn:   args[1],
			})
			return &NoneVal{}
		},
	}

	// run_tests() - Run all registered tests and return results
	exports["run_tests"] = &BuiltinFnVal{
		Name: "testing.run_tests",
		Fn: func(args []Value) Value {
			results := make([]Value, 0, len(registeredTests))
			for _, test := range registeredTests {
				passed := true
				errMsg := ""

				func() {
					defer func() {
						if r := recover(); r != nil {
							passed = false
							switch e := r.(type) {
							case *RuntimeError:
								errMsg = e.Message
							default:
								errMsg = fmt.Sprintf("%v", r)
							}
						}
					}()
					// Call the test function
					callValue(test.Fn, nil)
				}()

				result := &MapVal{
					Keys: []Value{
						&StringVal{Val: "name"},
						&StringVal{Val: "passed"},
						&StringVal{Val: "error"},
					},
					Values: []Value{
						&StringVal{Val: test.Name},
						&BoolVal{Val: passed},
						&StringVal{Val: errMsg},
					},
				}
				results = append(results, result)
			}
			// Clear registered tests after running
			registeredTests = nil
			return &ListVal{Elements: results}
		},
	}

	return exports
}
