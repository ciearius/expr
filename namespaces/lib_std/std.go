package lib_std

import (
	"github.com/antonmedv/expr/builtin"
)

func NewBuiltinStandard() builtin.BuiltinNamespace {
	return &BuiltinStandard{
		builtin.ContainerWith(
			&F_len{},
			&F_all{},
			&F_none{},
			&F_any{},
			&F_one{},
			&F_filter{},
			&F_map{},
			&F_count{},
		),
	}
}

type BuiltinStandard struct {
	builtin.BaseNamespace
}

func (b *BuiltinStandard) Name() string {
	return ""
}
