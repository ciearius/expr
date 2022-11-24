package checking

import (
	"reflect"

	"github.com/antonmedv/expr/ast"
)

type Info struct {
	Method bool
}

type VisitFunction func(node ast.Node) (reflect.Type, Info)
type ErrorFunction func(node ast.Node, format string, args ...interface{}) (reflect.Type, Info)
type AddCollectionFunction func(collection reflect.Type)
type PopCollectionFunction func()

type ExternVisitor struct {
	visit         VisitFunction
	error         ErrorFunction
	addCollection AddCollectionFunction
	popCollection PopCollectionFunction
}

func (e *ExternVisitor) Visit(node ast.Node) (reflect.Type, Info) {
	return e.visit(node)
}

func (e *ExternVisitor) Error(node ast.Node, format string, args ...interface{}) (reflect.Type, Info) {
	return e.error(node, format, args)
}

func (e *ExternVisitor) AddCollection(collection reflect.Type) {
	e.addCollection(collection)
}

func (e *ExternVisitor) PopCollection() {
	e.popCollection()
}

func Create(v VisitFunction, err ErrorFunction, add AddCollectionFunction, pop PopCollectionFunction) *ExternVisitor {
	return &ExternVisitor{
		visit:         v,
		error:         err,
		addCollection: add,
		popCollection: pop,
	}
}
