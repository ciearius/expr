package checker

import (
	"reflect"

	. "github.com/antonmedv/expr/ast"
)

var builtin_math = namespace{
	name: "math",
	members: map[string]Checker{
		"abs": func(v *visitor, node *BuiltinNode) (reflect.Type, info) {
			if len(node.Arguments) != 1 {
				return v.error(node, "math.abs expects one number as input")
			}

			param, _ := v.visit(node.Arguments[0])

			if isFloat(param) {
				return floatType, info{}
			} else if isInteger(param) {
				return integerType, info{}
			}

			return v.error(node, "math.abs expects a number as input - got: %v", param)
		},
	},
}
