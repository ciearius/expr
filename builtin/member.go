package builtin

import (
	"reflect"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/util/checking"
)

// TODO:
// - parsing
// - checking
// - compiling
// - execution / runtime

type Member interface {
	Name() string
	Callable() bool
	Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info)
}

type BaseNamespace struct {
	Members map[string]Member
}

func (b *BaseNamespace) Get(name string) (Member, bool) {
	m, ok := b.Members[name]
	return m, ok
}

func (b *BaseNamespace) Check(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	m, ok := b.Get(node.Name)

	if !ok {
		v.Error(node, "%s does not exist in %s", node.Name, node.Namespace)
	}

	return m.Visit(v, node)
}

func ContainerWith(ms ...Member) BaseNamespace {
	container := BaseNamespace{make(map[string]Member)}

	for _, m := range ms {
		container.Members[m.Name()] = m
	}

	return container
}
