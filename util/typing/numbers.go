package typing

import (
	"reflect"

	"github.com/antonmedv/expr/ast"
)

func ToNumbertype(a, b reflect.Type) reflect.Type {
	if a.Kind() == b.Kind() {
		return a
	}
	if IsFloat(a) || IsFloat(b) {
		return FloatType
	}
	return IntegerType
}

func IsInteger(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			fallthrough
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		}
	}
	return false
}

func IsFloat(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Float32, reflect.Float64:
			return true
		}
	}
	return false
}

func IsNumber(t reflect.Type) bool {
	return IsInteger(t) || IsFloat(t)
}

func IsIntegerOrArithmeticOperation(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.IntegerNode:
		return true
	case *ast.UnaryNode:
		switch n.Operator {
		case "+", "-":
			return true
		}
	case *ast.BinaryNode:
		switch n.Operator {
		case "+", "/", "-", "*":
			return true
		}
	}
	return false
}

func SetTypeForIntegers(node ast.Node, t reflect.Type) {
	switch n := node.(type) {
	case *ast.IntegerNode:
		n.SetType(t)
	case *ast.UnaryNode:
		switch n.Operator {
		case "+", "-":
			SetTypeForIntegers(n.Node, t)
		}
	case *ast.BinaryNode:
		switch n.Operator {
		case "+", "/", "-", "*":
			SetTypeForIntegers(n.Left, t)
			SetTypeForIntegers(n.Right, t)
		}
	}
}
