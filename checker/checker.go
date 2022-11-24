package checker

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/antonmedv/expr/ast"
	"github.com/antonmedv/expr/conf"
	"github.com/antonmedv/expr/file"
	"github.com/antonmedv/expr/namespaces"
	"github.com/antonmedv/expr/parser"
	"github.com/antonmedv/expr/util/checking"
	. "github.com/antonmedv/expr/util/checking"
	. "github.com/antonmedv/expr/util/typing"
	"github.com/antonmedv/expr/vm"
)

func Check(tree *parser.Tree, config *conf.Config) (t reflect.Type, err error) {
	if config == nil {
		config = conf.New(nil)
	}

	v := &CheckVisitor{
		config:      config,
		collections: make([]reflect.Type, 0),
		parents:     make([]ast.Node, 0),
	}

	v.ex = checking.Create(
		func(node ast.Node) (reflect.Type, checking.Info) {
			return v.visit(node)
		},
		func(node ast.Node, format string, args ...interface{}) (reflect.Type, checking.Info) {
			return v.error(node, format, args...)
		},
		func(collection reflect.Type) {
			// v.collections = append(v.collections, collection)
			v.collections = append(v.collections, collection)
		}, func() {
			// v.collections = v.collections[:len(v.collections)-1]
			v.collections = v.collections[:len(v.collections)-1]
		},
	)

	t, _ = v.visit(tree.Node)

	if v.err != nil {
		return t, v.err.Bind(tree.Source)
	}

	if v.config.Expect != reflect.Invalid {
		switch v.config.Expect {
		case reflect.Int, reflect.Int64, reflect.Float64:
			if !IsNumber(t) {
				return nil, fmt.Errorf("expected %v, but got %v", v.config.Expect, t)
			}
		default:
			if t == nil || t.Kind() != v.config.Expect {
				return nil, fmt.Errorf("expected %v, but got %v", v.config.Expect, t)
			}
		}
	}

	return t, nil
}

type CheckVisitor struct {
	config      *conf.Config
	collections []reflect.Type
	parents     []ast.Node
	err         *file.Error
	ex          *ExternVisitor
}

func (v *CheckVisitor) visit(node ast.Node) (reflect.Type, Info) {
	var t reflect.Type
	var i Info
	v.parents = append(v.parents, node)
	switch n := node.(type) {
	case *ast.NilNode:
		t, i = v.NilNode(n)
	case *ast.IdentifierNode:
		t, i = v.IdentifierNode(n)
	case *ast.IntegerNode:
		t, i = v.IntegerNode(n)
	case *ast.FloatNode:
		t, i = v.FloatNode(n)
	case *ast.BoolNode:
		t, i = v.BoolNode(n)
	case *ast.StringNode:
		t, i = v.StringNode(n)
	case *ast.ConstantNode:
		t, i = v.ConstantNode(n)
	case *ast.UnaryNode:
		t, i = v.UnaryNode(n)
	case *ast.BinaryNode:
		t, i = v.BinaryNode(n)
	case *ast.ChainNode:
		t, i = v.ChainNode(n)
	case *ast.MemberNode:
		t, i = v.MemberNode(n)
	case *ast.SliceNode:
		t, i = v.SliceNode(n)
	case *ast.CallNode:
		t, i = v.CallNode(n)
	case *ast.BuiltinNode:
		t, i = v.BuiltinNode(n)
	case *ast.ClosureNode:
		t, i = v.ClosureNode(n)
	case *ast.PointerNode:
		t, i = v.PointerNode(n)
	case *ast.ConditionalNode:
		t, i = v.ConditionalNode(n)
	case *ast.ArrayNode:
		t, i = v.ArrayNode(n)
	case *ast.MapNode:
		t, i = v.MapNode(n)
	case *ast.PairNode:
		t, i = v.PairNode(n)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
	v.parents = v.parents[:len(v.parents)-1]
	node.SetType(t)
	return t, i
}

func (v *CheckVisitor) error(node ast.Node, format string, args ...interface{}) (reflect.Type, Info) {
	if v.err == nil { // show first error
		v.err = &file.Error{
			Location: node.Location(),
			Message:  fmt.Sprintf(format, args...),
		}
	}
	return AnyType, Info{} // interface represent undefined type
}

func (v *CheckVisitor) NilNode(*ast.NilNode) (reflect.Type, Info) {
	return NilType, Info{}
}

func (v *CheckVisitor) IdentifierNode(node *ast.IdentifierNode) (reflect.Type, Info) {
	if v.config.Types == nil {
		node.Deref = true
		return AnyType, Info{}
	}
	if t, ok := v.config.Types[node.Value]; ok {
		if t.Ambiguous {
			return v.error(node, "ambiguous identifier %v", node.Value)
		}
		d, c := Deref(t.Type)
		node.Deref = c
		node.Method = t.Method
		node.MethodIndex = t.MethodIndex
		node.FieldIndex = t.FieldIndex
		return d, Info{Method: t.Method}
	}
	if !v.config.Strict {
		if v.config.DefaultType != nil {
			return v.config.DefaultType, Info{}
		}
		return AnyType, Info{}
	}
	return v.error(node, "unknown name %v", node.Value)
}

func (v *CheckVisitor) IntegerNode(*ast.IntegerNode) (reflect.Type, Info) {
	return IntegerType, Info{}
}

func (v *CheckVisitor) FloatNode(*ast.FloatNode) (reflect.Type, Info) {
	return FloatType, Info{}
}

func (v *CheckVisitor) BoolNode(*ast.BoolNode) (reflect.Type, Info) {
	return BoolType, Info{}
}

func (v *CheckVisitor) StringNode(*ast.StringNode) (reflect.Type, Info) {
	return StringType, Info{}
}

func (v *CheckVisitor) ConstantNode(node *ast.ConstantNode) (reflect.Type, Info) {
	return reflect.TypeOf(node.Value), Info{}
}

func (v *CheckVisitor) UnaryNode(node *ast.UnaryNode) (reflect.Type, Info) {
	t, _ := v.visit(node.Node)

	switch node.Operator {

	case "!", "not":
		if IsBool(t) {
			return BoolType, Info{}
		}
		if IsAny(t) {
			return BoolType, Info{}
		}

	case "+", "-":
		if IsNumber(t) {
			return t, Info{}
		}
		if IsAny(t) {
			return AnyType, Info{}
		}

	default:
		return v.error(node, "unknown operator (%v)", node.Operator)
	}

	return v.error(node, `invalid operation: %v (mismatched type %v)`, node.Operator, t)
}

func (v *CheckVisitor) BinaryNode(node *ast.BinaryNode) (reflect.Type, Info) {
	l, _ := v.visit(node.Left)
	r, _ := v.visit(node.Right)

	// check operator overloading
	if fns, ok := v.config.Operators[node.Operator]; ok {
		t, _, ok := conf.FindSuitableOperatorOverload(fns, v.config.Types, l, r)
		if ok {
			return t, Info{}
		}
	}

	switch node.Operator {
	case "==", "!=":
		if IsNumber(l) && IsNumber(r) {
			return BoolType, Info{}
		}
		if l == nil || r == nil { // It is possible to compare with nil.
			return BoolType, Info{}
		}
		if l.Kind() == r.Kind() {
			return BoolType, Info{}
		}
		if IsAny(l) || IsAny(r) {
			return BoolType, Info{}
		}

	case "or", "||", "and", "&&":
		if IsBool(l) && IsBool(r) {
			return BoolType, Info{}
		}
		if DualAnyOf(l, r, IsBool) {
			return BoolType, Info{}
		}

	case "<", ">", ">=", "<=":
		if IsNumber(l) && IsNumber(r) {
			return BoolType, Info{}
		}
		if IsString(l) && IsString(r) {
			return BoolType, Info{}
		}
		if IsTime(l) && IsTime(r) {
			return BoolType, Info{}
		}
		if DualAnyOf(l, r, IsNumber, IsString, IsTime) {
			return BoolType, Info{}
		}

	case "-":
		if IsNumber(l) && IsNumber(r) {
			return ToNumbertype(l, r), Info{}
		}
		if IsTime(l) && IsTime(r) {
			return DurationType, Info{}
		}
		if DualAnyOf(l, r, IsNumber, IsTime) {
			return AnyType, Info{}
		}

	case "/", "*":
		if IsNumber(l) && IsNumber(r) {
			return ToNumbertype(l, r), Info{}
		}
		if DualAnyOf(l, r, IsNumber) {
			return AnyType, Info{}
		}

	case "**", "^":
		if IsNumber(l) && IsNumber(r) {
			return FloatType, Info{}
		}
		if DualAnyOf(l, r, IsNumber) {
			return FloatType, Info{}
		}

	case "%":
		if IsInteger(l) && IsInteger(r) {
			return ToNumbertype(l, r), Info{}
		}
		if DualAnyOf(l, r, IsInteger) {
			return AnyType, Info{}
		}

	case "+":
		if IsNumber(l) && IsNumber(r) {
			return ToNumbertype(l, r), Info{}
		}
		if IsString(l) && IsString(r) {
			return StringType, Info{}
		}
		if IsTime(l) && IsDuration(r) {
			return TimeType, Info{}
		}
		if IsDuration(l) && IsTime(r) {
			return TimeType, Info{}
		}
		if DualAnyOf(l, r, IsNumber, IsString, IsTime, IsDuration) {
			return AnyType, Info{}
		}

	case "in":
		if (IsString(l) || IsAny(l)) && IsStruct(r) {
			return BoolType, Info{}
		}
		if IsMap(r) {
			return BoolType, Info{}
		}
		if IsArray(r) {
			return BoolType, Info{}
		}
		if IsAny(l) && MatchesAny(r, IsString, IsArray, IsMap) {
			return BoolType, Info{}
		}
		if IsAny(l) && IsAny(r) {
			return BoolType, Info{}
		}

	case "matches":
		if s, ok := node.Right.(*ast.StringNode); ok {
			r, err := regexp.Compile(s.Value)
			if err != nil {
				return v.error(node, err.Error())
			}
			node.Regexp = r
		}
		if IsString(l) && IsString(r) {
			return BoolType, Info{}
		}
		if DualAnyOf(l, r, IsString) {
			return BoolType, Info{}
		}

	case "contains", "startsWith", "endsWith":
		if IsString(l) && IsString(r) {
			return BoolType, Info{}
		}
		if DualAnyOf(l, r, IsString) {
			return BoolType, Info{}
		}

	case "..":
		ret := reflect.SliceOf(IntegerType)
		if IsInteger(l) && IsInteger(r) {
			return ret, Info{}
		}
		if DualAnyOf(l, r, IsInteger) {
			return ret, Info{}
		}

	default:
		return v.error(node, "unknown operator (%v)", node.Operator)

	}

	return v.error(node, `invalid operation: %v (mismatched types %v and %v)`, node.Operator, l, r)
}

func (v *CheckVisitor) ChainNode(node *ast.ChainNode) (reflect.Type, Info) {
	return v.visit(node.Node)
}

func (v *CheckVisitor) MemberNode(node *ast.MemberNode) (reflect.Type, Info) {
	base, _ := v.visit(node.Node)
	prop, _ := v.visit(node.Property)

	if name, ok := node.Property.(*ast.StringNode); ok {
		if base == nil {
			return v.error(node, "type %v has no field %v", base, name.Value)
		}
		// First, check methods defined on base type itself,
		// independent of which type it is. Without dereferencing.
		if m, ok := base.MethodByName(name.Value); ok {
			node.Method = true
			node.MethodIndex = m.Index
			node.Name = name.Value
			if base.Kind() == reflect.Interface {
				// In case of interface type method will not have a receiver,
				// and to prevent checker decreasing numbers of in arguments
				// return method type as not method (second argument is false).
				return m.Type, Info{}
			} else {
				return m.Type, Info{Method: true}
			}
		}
	}

	if base.Kind() == reflect.Ptr {
		base = base.Elem()
	}

	switch base.Kind() {
	case reflect.Interface:
		node.Deref = true
		return AnyType, Info{}

	case reflect.Map:
		if !prop.AssignableTo(base.Key()) {
			return v.error(node.Property, "cannot use %v to get an element from %v", prop, base)
		}
		t, c := Deref(base.Elem())
		node.Deref = c
		return t, Info{}

	case reflect.Array, reflect.Slice:
		if !IsInteger(prop) && !IsAny(prop) {
			return v.error(node.Property, "array elements can only be selected using an integer (got %v)", prop)
		}
		t, c := Deref(base.Elem())
		node.Deref = c
		return t, Info{}

	case reflect.Struct:
		if name, ok := node.Property.(*ast.StringNode); ok {
			propertyName := name.Value
			if field, ok := FetchField(base, propertyName); ok {
				t, c := Deref(field.Type)
				node.Deref = c
				node.FieldIndex = field.Index
				node.Name = propertyName
				return t, Info{}
			}
			if len(v.parents) > 1 {
				if _, ok := v.parents[len(v.parents)-2].(*ast.CallNode); ok {
					return v.error(node, "type %v has no method %v", base, propertyName)
				}
			}
			return v.error(node, "type %v has no field %v", base, propertyName)
		}
	}

	return v.error(node, "type %v[%v] is undefined", base, prop)
}

func (v *CheckVisitor) SliceNode(node *ast.SliceNode) (reflect.Type, Info) {
	t, _ := v.visit(node.Node)

	switch t.Kind() {
	case reflect.Interface:
		// ok
	case reflect.String, reflect.Array, reflect.Slice:
		// ok
	default:
		return v.error(node, "cannot slice %v", t)
	}

	if node.From != nil {
		from, _ := v.visit(node.From)
		if !IsInteger(from) && !IsAny(from) {
			return v.error(node.From, "non-integer slice index %v", from)
		}
	}
	if node.To != nil {
		to, _ := v.visit(node.To)
		if !IsInteger(to) && !IsAny(to) {
			return v.error(node.To, "non-integer slice index %v", to)
		}
	}
	return t, Info{}
}

func (v *CheckVisitor) CallNode(node *ast.CallNode) (reflect.Type, Info) {
	fn, fnInfo := v.visit(node.Callee)

	fnName := "function"
	if identifier, ok := node.Callee.(*ast.IdentifierNode); ok {
		fnName = identifier.Value
	}
	if member, ok := node.Callee.(*ast.MemberNode); ok {
		if name, ok := member.Property.(*ast.StringNode); ok {
			fnName = name.Value
		}
	}

	switch fn.Kind() {
	case reflect.Interface:
		return AnyType, Info{}
	case reflect.Func:
		inputParamsCount := 1 // for functions
		if fnInfo.Method {
			inputParamsCount = 2 // for methods
		}

		if !IsAny(fn) &&
			fn.IsVariadic() &&
			fn.NumIn() == inputParamsCount &&
			((fn.NumOut() == 1 && // Function with one return value
				fn.Out(0).Kind() == reflect.Interface) ||
				(fn.NumOut() == 2 && // Function with one return value and an error
					fn.Out(0).Kind() == reflect.Interface &&
					fn.Out(1) == ErrorType)) {
			rest := fn.In(fn.NumIn() - 1) // function has only one param for functions and two for methods
			if rest.Kind() == reflect.Slice && rest.Elem().Kind() == reflect.Interface {
				node.Fast = true
			}
		}

		return v.checkFunc(fn, fnInfo.Method, node, fnName, node.Arguments)
	}
	return v.error(node, "%v is not callable", fn)
}

// checkFunc checks func arguments and returns "return type" of func or method.
func (v *CheckVisitor) checkFunc(fn reflect.Type, method bool, node *ast.CallNode, name string, arguments []ast.Node) (reflect.Type, Info) {
	if IsAny(fn) {
		return AnyType, Info{}
	}

	if fn.NumOut() == 0 {
		return v.error(node, "func %v doesn't return value", name)
	}
	if numOut := fn.NumOut(); numOut > 2 {
		return v.error(node, "func %v returns more then two values", name)
	}

	numIn := fn.NumIn()

	// If func is method on an env, first argument should be a receiver,
	// and actual arguments less than numIn by one.
	if method {
		numIn--
	}

	if fn.IsVariadic() {
		if len(arguments) < numIn-1 {
			return v.error(node, "not enough arguments to call %v", name)
		}
	} else {
		if len(arguments) > numIn {
			return v.error(node, "too many arguments to call %v", name)
		}
		if len(arguments) < numIn {
			return v.error(node, "not enough arguments to call %v", name)
		}
	}

	offset := 0

	// Skip first argument in case of the receiver.
	if method {
		offset = 1
	}

	for i, arg := range arguments {
		t, _ := v.visit(arg)

		var in reflect.Type
		if fn.IsVariadic() && i >= numIn-1 {
			// For variadic arguments fn(xs ...int), go replaces type of xs (int) with ([]int).
			// As we compare arguments one by one, we need underling type.
			in = fn.In(fn.NumIn() - 1).Elem()
		} else {
			in = fn.In(i + offset)
		}

		if IsIntegerOrArithmeticOperation(arg) {
			t = in
			SetTypeForIntegers(arg, t)
		}

		if t == nil {
			continue
		}

		if !t.AssignableTo(in) && t.Kind() != reflect.Interface {
			return v.error(arg, "cannot use %v as argument (type %v) to call %v ", t, in, name)
		}
	}

	if !fn.IsVariadic() {
	funcTypes:
		for i := range vm.FuncTypes {
			if i == 0 {
				continue
			}
			typed := reflect.ValueOf(vm.FuncTypes[i]).Elem().Type()
			if typed.Kind() != reflect.Func {
				continue
			}
			if typed.NumOut() != fn.NumOut() {
				continue
			}
			for j := 0; j < typed.NumOut(); j++ {
				if typed.Out(j) != fn.Out(j) {
					continue funcTypes
				}
			}
			if typed.NumIn() != len(arguments) {
				continue
			}
			for j, arg := range arguments {
				if typed.In(j) != arg.Type() {
					continue funcTypes
				}
			}
			node.Typed = i
		}
	}

	return fn.Out(0), Info{}
}

func (v *CheckVisitor) BuiltinNode(node *ast.BuiltinNode) (reflect.Type, Info) {
	space, ok := namespaces.Get(node.Namespace)

	if !ok {
		return v.error(node, "there is no builtin namespace %s", node.Namespace)
	}

	return space.Check(v.ex, node)
}

func (v *CheckVisitor) ClosureNode(node *ast.ClosureNode) (reflect.Type, Info) {
	t, _ := v.visit(node.Node)
	return reflect.FuncOf([]reflect.Type{AnyType}, []reflect.Type{t}, false), Info{}
}

func (v *CheckVisitor) PointerNode(node *ast.PointerNode) (reflect.Type, Info) {
	if len(v.collections) == 0 {
		return v.error(node, "cannot use pointer accessor outside closure")
	}

	collection := v.collections[len(v.collections)-1]
	switch collection.Kind() {
	case reflect.Interface:
		return AnyType, Info{}
	case reflect.Array, reflect.Slice:
		return collection.Elem(), Info{}
	}
	return v.error(node, "cannot use %v as array", collection)
}

func (v *CheckVisitor) ConditionalNode(node *ast.ConditionalNode) (reflect.Type, Info) {
	c, _ := v.visit(node.Cond)
	if !IsBool(c) && !IsAny(c) {
		return v.error(node.Cond, "non-bool expression (type %v) used as condition", c)
	}

	t1, _ := v.visit(node.Exp1)
	t2, _ := v.visit(node.Exp2)

	if t1 == nil && t2 != nil {
		return t2, Info{}
	}
	if t1 != nil && t2 == nil {
		return t1, Info{}
	}
	if t1 == nil && t2 == nil {
		return NilType, Info{}
	}
	if t1.AssignableTo(t2) {
		return t1, Info{}
	}
	return AnyType, Info{}
}

func (v *CheckVisitor) ArrayNode(node *ast.ArrayNode) (reflect.Type, Info) {
	for _, node := range node.Nodes {
		v.visit(node)
	}
	return ArrayType, Info{}
}

func (v *CheckVisitor) MapNode(node *ast.MapNode) (reflect.Type, Info) {
	for _, pair := range node.Pairs {
		v.visit(pair)
	}
	return MapType, Info{}
}

func (v *CheckVisitor) PairNode(node *ast.PairNode) (reflect.Type, Info) {
	v.visit(node.Key)
	v.visit(node.Value)
	return NilType, Info{}
}
