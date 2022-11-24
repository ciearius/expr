package builtin

type BaseFunc struct {
}

func (b *BaseFunc) Callable() bool {
	return true
}

type Function interface {
	Arguments() []Argument
}
