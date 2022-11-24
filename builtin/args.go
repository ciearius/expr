package builtin

type ParserArgType int

const (
	Expression ParserArgType = iota
	Closure
)

type Argument struct {
	ParserType ParserArgType
}

func ParserArguments(m Function) []ParserArgType {
	args := m.Arguments()
	res := make([]ParserArgType, len(args))

	for i, v := range args {
		res[i] = v.ParserType
	}

	return res
}
