package interpreter

func init() {
	registerListMethods()
}

func registerListMethods() {
	// len() -> Int — Get list length
	RegisterMethod(TypeList, "len", func(receiver Value, args []Value) Value {
		list := receiver.(*ListVal)
		return &IntVal{Val: int64(len(list.Elements))}
	})

	// length() -> Int — Alias for len
	RegisterMethod(TypeList, "length", func(receiver Value, args []Value) Value {
		list := receiver.(*ListVal)
		return &IntVal{Val: int64(len(list.Elements))}
	})

	// append(item) -> None — Append item to list (mutates)
	RegisterMethod(TypeList, "append", func(receiver Value, args []Value) Value {
		list := receiver.(*ListVal)
		if len(args) < 1 {
			panic(&RuntimeError{Message: "List.append requires at least one argument"})
		}
		list.Elements = append(list.Elements, args[0])
		return &NoneVal{}
	})

	// contains(item) -> Bool — Check if list contains item
	RegisterMethod(TypeList, "contains", func(receiver Value, args []Value) Value {
		list := receiver.(*ListVal)
		if len(args) < 1 {
			return &BoolVal{Val: false}
		}
		for _, elem := range list.Elements {
			if Equal(elem, args[0]) {
				return &BoolVal{Val: true}
			}
		}
		return &BoolVal{Val: false}
	})
}
