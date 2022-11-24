package namespaces

import (
	"github.com/antonmedv/expr/builtin"
	"github.com/antonmedv/expr/namespaces/lib_math"
)

type namespaceMap map[string]builtin.BuiltinNamespace

func (n namespaceMap) add(b builtin.BuiltinNamespace) {
	n[b.Name()] = b
}

var mapped namespaceMap = namespaceMap{}

func init() {
	mapped.add(Stdlib)
	mapped.add(lib_math.NewBuiltinMath())
}
