package eval

import (
	"fmt"
	"monkey/object"
)

// todo: len array support, first, last, push, rest

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument type given to `len` not supported, got=%s", args[0].Type())
			}
		},
	},

	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[0]
			default:
				return newError("argument type given to `first` not supported, got=%s", args[0].Type())
			}
		},
	},

	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newError("argument type given to `last` not supported, got=%s", args[0].Type())
			}
		},
	},

	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, expected=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) > 0 {
					elems := make([]object.Object, len(arg.Elements)-1)
					copy(elems, arg.Elements[1:])
					return &object.Array{Elements: arg.Elements[1:]}
				}
				return NULL
			default:
				return newError("argument type given to `rest` not supported, got=%s", args[0].Type())
			}
		},
	},

	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, expected=2", len(args))
			}

			switch first := args[0].(type) {
			case *object.Array:
				elems := make([]object.Object, len(first.Elements)+1)
				copy(elems, first.Elements)
				elems[len(first.Elements)] = args[1]
				return &object.Array{Elements: elems}
			default:
				return newError("argument type given to `push` not supported, got=%s", args[0].Type())
			}
		},
	},

	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
