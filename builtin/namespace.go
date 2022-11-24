package builtin

import (
	"reflect"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/util/checking"
)

type BuiltinNamespace interface {
	Name() string
	MemberContainer
	MemberChecker
}

type MemberContainer interface {
	Get(name string) (Member, bool)
}

type MemberChecker interface {
	Check(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info)
}
