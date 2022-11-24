package lib_std

import (
	"reflect"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/builtin"
	"github.com/antonmedv/expr/util/checking"
	. "github.com/antonmedv/expr/util/typing"
)

func (f *F_len) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	param, _ := v.Visit(node.Arguments[0])

	if IsArray(param) || IsMap(param) || IsString(param) {
		return IntegerType, checking.Info{}
	}
	if IsAny(param) {
		return AnyType, checking.Info{}
	}

	return v.Error(node, "invalid argument for len (type %s)", param)
}
func (f *F_all) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	return visit_reducers(f, v, node)
}

func (f *F_none) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	return visit_reducers(f, v, node)
}

func (f *F_any) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	return visit_reducers(f, v, node)
}

func (f *F_one) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	return visit_reducers(f, v, node)
}

func (f *F_filter) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	collection, _ := v.Visit(node.Arguments[0])
	if !IsArray(collection) && !IsAny(collection) {
		return v.Error(node.Arguments[0], "builtin %s takes only array (got %s)", node, collection)
	}

	v.AddCollection(collection)
	closure, _ := v.Visit(node.Arguments[1])
	v.PopCollection()

	if IsFunc(closure) &&
		closure.NumOut() == 1 &&
		closure.NumIn() == 1 && IsAny(closure.In(0)) {

		if !IsBool(closure.Out(0)) && !IsAny(closure.Out(0)) {
			return v.Error(node.Arguments[1], "closure should return boolean (got %s)", closure.Out(0))
		}
		if IsAny(collection) {
			return ArrayType, checking.Info{}
		}
		return reflect.SliceOf(collection.Elem()), checking.Info{}
	}
	return v.Error(node.Arguments[1], "closure should has one input and one output param")
}

func (f *F_map) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	collection, _ := v.Visit(node.Arguments[0])
	if !IsArray(collection) && !IsAny(collection) {
		return v.Error(node.Arguments[0], "builtin %s takes only array (got %s)", node, collection)
	}

	v.AddCollection(collection)
	closure, _ := v.Visit(node.Arguments[1])
	v.PopCollection()

	if IsFunc(closure) &&
		closure.NumOut() == 1 &&
		closure.NumIn() == 1 && IsAny(closure.In(0)) {

		return reflect.SliceOf(closure.Out(0)), checking.Info{}
	}
	return v.Error(node.Arguments[1], "closure should has one input and one output param")
}

func (f *F_count) Visit(v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	collection, _ := v.Visit(node.Arguments[0])
	if !IsArray(collection) && !IsAny(collection) {
		return v.Error(node.Arguments[0], "builtin %s takes only array (got %s)", node, collection)
	}

	v.AddCollection(collection) // v.collections = append(v.collections, collection)
	closure, _ := v.Visit(node.Arguments[1])
	v.PopCollection() // v.collections = v.collections[:len(v.collections)-1]

	if IsFunc(closure) &&
		closure.NumOut() == 1 &&
		closure.NumIn() == 1 && IsAny(closure.In(0)) {
		if !IsBool(closure.Out(0)) && !IsAny(closure.Out(0)) {
			return v.Error(node.Arguments[1], "closure should return boolean (got %s)", closure.Out(0))
		}

		return IntegerType, checking.Info{}
	}
	return v.Error(node.Arguments[1], "closure should has one input and one output param")
}

func visit_reducers(f builtin.Function, v *checking.ExternVisitor, node *ast.BuiltinNode) (reflect.Type, checking.Info) {
	collection, _ := v.Visit(node.Arguments[0])
	if !IsArray(collection) && !IsAny(collection) {
		return v.Error(node.Arguments[0], "builtin %s takes only array (got %s)", node, collection)
	}

	v.AddCollection(collection)
	closure, _ := v.Visit(node.Arguments[1])
	v.PopCollection()

	if IsFunc(closure) &&
		closure.NumOut() == 1 &&
		closure.NumIn() == 1 && IsAny(closure.In(0)) {

		if !IsBool(closure.Out(0)) && !IsAny(closure.Out(0)) {
			return v.Error(node.Arguments[1], "closure should return boolean (got %s)", closure.Out(0))
		}
		return BoolType, checking.Info{}
	}

	return v.Error(node.Arguments[1], "closure should has one input and one output param")
}
