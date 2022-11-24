package namespaces

import (
	"github.com/antonmedv/expr/builtin"
	"github.com/antonmedv/expr/namespaces/lib_std"
)

var Stdlib builtin.BuiltinNamespace = lib_std.NewBuiltinStandard()

func Get(name string) (builtin.BuiltinNamespace, bool) {
	b, ok := mapped[name]
	return b, ok
}
