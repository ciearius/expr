package checker

import (
	"reflect"

	. "github.com/antonmedv/expr/ast"
)

var builtin_standard = namespace{
	name: "",
	check: func(v *visitor, node *BuiltinNode) (reflect.Type, info) {

		switch node.Name {

		case "len":
			param, _ := v.visit(node.Arguments[0])
			if isArray(param) || isMap(param) || isString(param) {
				return integerType, info{}
			}
			if isAny(param) {
				return anyType, info{}
			}
			return v.error(node, "invalid argument for len (type %v)", param)

		case "all", "none", "any", "one":
			collection, _ := v.visit(node.Arguments[0])
			if !isArray(collection) && !isAny(collection) {
				return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection)
			}

			v.collections = append(v.collections, collection)
			closure, _ := v.visit(node.Arguments[1])
			v.collections = v.collections[:len(v.collections)-1]

			if isFunc(closure) &&
				closure.NumOut() == 1 &&
				closure.NumIn() == 1 && isAny(closure.In(0)) {

				if !isBool(closure.Out(0)) && !isAny(closure.Out(0)) {
					return v.error(node.Arguments[1], "closure should return boolean (got %v)", closure.Out(0).String())
				}
				return boolType, info{}
			}
			return v.error(node.Arguments[1], "closure should has one input and one output param")

		case "filter":
			collection, _ := v.visit(node.Arguments[0])
			if !isArray(collection) && !isAny(collection) {
				return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection)
			}

			v.collections = append(v.collections, collection)
			closure, _ := v.visit(node.Arguments[1])
			v.collections = v.collections[:len(v.collections)-1]

			if isFunc(closure) &&
				closure.NumOut() == 1 &&
				closure.NumIn() == 1 && isAny(closure.In(0)) {

				if !isBool(closure.Out(0)) && !isAny(closure.Out(0)) {
					return v.error(node.Arguments[1], "closure should return boolean (got %v)", closure.Out(0).String())
				}
				if isAny(collection) {
					return arrayType, info{}
				}
				return reflect.SliceOf(collection.Elem()), info{}
			}
			return v.error(node.Arguments[1], "closure should has one input and one output param")

		case "map":
			collection, _ := v.visit(node.Arguments[0])
			if !isArray(collection) && !isAny(collection) {
				return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection)
			}

			v.collections = append(v.collections, collection)
			closure, _ := v.visit(node.Arguments[1])
			v.collections = v.collections[:len(v.collections)-1]

			if isFunc(closure) &&
				closure.NumOut() == 1 &&
				closure.NumIn() == 1 && isAny(closure.In(0)) {

				return reflect.SliceOf(closure.Out(0)), info{}
			}
			return v.error(node.Arguments[1], "closure should has one input and one output param")

		case "count":
			collection, _ := v.visit(node.Arguments[0])
			if !isArray(collection) && !isAny(collection) {
				return v.error(node.Arguments[0], "builtin %v takes only array (got %v)", node.Name, collection)
			}

			v.collections = append(v.collections, collection)
			closure, _ := v.visit(node.Arguments[1])
			v.collections = v.collections[:len(v.collections)-1]

			if isFunc(closure) &&
				closure.NumOut() == 1 &&
				closure.NumIn() == 1 && isAny(closure.In(0)) {
				if !isBool(closure.Out(0)) && !isAny(closure.Out(0)) {
					return v.error(node.Arguments[1], "closure should return boolean (got %v)", closure.Out(0).String())
				}

				return integerType, info{}
			}
			return v.error(node.Arguments[1], "closure should has one input and one output param")

		default:
			return v.error(node, "unknown builtin %v", node.Name)
		}
	},
}
