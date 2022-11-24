package lib_math

import (
	"reflect"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/util/checking"
	. "github.com/antonmedv/expr/util/typing"
)

func (f *F_abs) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	if len(node.Arguments) != 1 {
		return v.Error(node, "math.abs expects one number as input")
	}

	param, _ := v.Visit(node.Arguments[0])

	if IsFloat(param) {
		return FloatType, checking.Info{}
	} else if IsInteger(param) {
		return IntegerType, checking.Info{}
	}

	return v.Error(node, "math.abs expects a number as input - got: %v", param)
}

func (f *C_pi) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	if node.Arguments != nil {
		return v.Error(node, "math.pi is a constant - it cannot be invoked")
	}

	return FloatType, checking.Info{}
}
