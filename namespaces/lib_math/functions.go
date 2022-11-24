package lib_math

import "github.com/antonmedv/expr/builtin"

type F_abs struct {
	builtin.BaseFunc
}

func (f *F_abs) Name() string {
	return "abs"
}

func (f *F_abs) Arguments() []builtin.Argument {
	return []builtin.Argument{
		{ParserType: builtin.Expression},
	}
}

type C_pi struct {
	builtin.BaseFunc
}

func (f *C_pi) Name() string {
	return "pi"
}
