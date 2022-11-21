package checker

import (
	"reflect"

	. "github.com/antonmedv/expr/ast"
)

var namespaces = map[string]namespace{}

func init() {
	namespaces[builtin_standard.name] = builtin_standard
	namespaces[builtin_math.name] = builtin_math
}

type namespace struct {
	name    string
	check   Checker
	members map[string]Checker
}

type Checker func(v *visitor, node *BuiltinNode) (reflect.Type, info)
