package lib_std

import "github.com/antonmedv/expr/builtin"

var arg_expression []builtin.Argument = []builtin.Argument{
	{ParserType: builtin.Expression},
}

var arg_expression_and_closure []builtin.Argument = []builtin.Argument{
	{ParserType: builtin.Expression},
	{ParserType: builtin.Closure},
}

type ExpressionClosureAccepting struct {
	builtin.BaseFunc
}

type F_len struct {
	builtin.BaseFunc
}

func (e *ExpressionClosureAccepting) Arguments() []builtin.Argument {
	return arg_expression_and_closure
}

func (f *F_len) Name() string {
	return "len"
}

func (f *F_len) Arguments() []builtin.Argument {
	return arg_expression
}

type F_all struct {
	ExpressionClosureAccepting
}

func (f *F_all) Name() string {
	return "all"
}

type F_none struct {
	ExpressionClosureAccepting
}

func (f *F_none) Name() string {
	return "none"
}

type F_any struct {
	ExpressionClosureAccepting
}

func (f *F_any) Name() string {
	return "any"
}

type F_one struct {
	ExpressionClosureAccepting
}

func (f *F_one) Name() string {
	return "one"
}

type F_filter struct {
	ExpressionClosureAccepting
}

func (f *F_filter) Name() string {
	return "filter"
}

type F_map struct {
	ExpressionClosureAccepting
}

func (f *F_map) Name() string {
	return "map"
}

type F_count struct {
	ExpressionClosureAccepting
}

func (f *F_count) Name() string {
	return "count"
}
