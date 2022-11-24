package builtin

type BaseConstant struct {
}

func (b *BaseConstant) Invokable() bool {
	return false
}

type Constant interface {
	GetValue() interface{}
}
