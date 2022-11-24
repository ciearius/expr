package lib_math

import (
	"github.com/antonmedv/expr/builtin"
)

type BuiltinMath struct {
	builtin.BaseNamespace
}

func NewBuiltinMath() builtin.BuiltinNamespace {
	return &BuiltinMath{
		builtin.ContainerWith(
			&F_abs{},
			&C_pi{},
		),
	}
}

func (b *BuiltinMath) Name() string {
	return "math"
}
