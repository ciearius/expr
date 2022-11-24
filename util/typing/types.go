package typing

import (
	"reflect"
	"time"

	"github.com/antonmedv/expr/conf"
)

var (
	NilType      = reflect.TypeOf(nil)
	BoolType     = reflect.TypeOf(true)
	IntegerType  = reflect.TypeOf(0)
	FloatType    = reflect.TypeOf(float64(0))
	StringType   = reflect.TypeOf("")
	ArrayType    = reflect.TypeOf([]interface{}{})
	MapType      = reflect.TypeOf(map[string]interface{}{})
	AnyType      = reflect.TypeOf(new(interface{})).Elem()
	TimeType     = reflect.TypeOf(time.Time{})
	DurationType = reflect.TypeOf(time.Duration(0))
	ErrorType    = reflect.TypeOf((*error)(nil)).Elem()
)

type TypeMatcher func(reflect.Type) bool

func MatchesAny(t reflect.Type, fns ...TypeMatcher) bool {
	for _, fn := range fns {
		if fn(t) {
			return true
		}
	}
	return false
}

func DualAnyOf(l, r reflect.Type, fns ...TypeMatcher) bool {
	if IsAny(l) && IsAny(r) {
		return true
	}
	if IsAny(l) && MatchesAny(r, fns...) {
		return true
	}
	if IsAny(r) && MatchesAny(l, fns...) {
		return true
	}
	return false
}

func IsAny(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Interface:
			return true
		}
	}
	return false
}

func IsTime(t reflect.Type) bool {
	if t != nil {
		switch t {
		case TimeType:
			return true
		}
	}
	return IsAny(t)
}

func IsDuration(t reflect.Type) bool {
	if t != nil {
		switch t {
		case DurationType:
			return true
		}
	}
	return false
}

func IsBool(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Bool:
			return true
		}
	}
	return false
}

func IsString(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.String:
			return true
		}
	}
	return false
}

func IsArray(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Ptr:
			return IsArray(t.Elem())
		case reflect.Slice, reflect.Array:
			return true
		}
	}
	return false
}

func IsMap(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Ptr:
			return IsMap(t.Elem())
		case reflect.Map:
			return true
		}
	}
	return false
}

func IsStruct(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Ptr:
			return IsStruct(t.Elem())
		case reflect.Struct:
			return true
		}
	}
	return false
}

func IsFunc(t reflect.Type) bool {
	if t != nil {
		switch t.Kind() {
		case reflect.Ptr:
			return IsFunc(t.Elem())
		case reflect.Func:
			return true
		}
	}
	return false
}

func FetchField(t reflect.Type, name string) (reflect.StructField, bool) {
	if t != nil {
		// First check all structs fields.
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			// Search all fields, even embedded structs.
			if conf.FieldName(field) == name {
				return field, true
			}
		}

		// Second check fields of embedded structs.
		for i := 0; i < t.NumField(); i++ {
			anon := t.Field(i)
			if anon.Anonymous {
				if field, ok := FetchField(anon.Type, name); ok {
					field.Index = append(anon.Index, field.Index...)
					return field, true
				}
			}
		}
	}
	return reflect.StructField{}, false
}

func Deref(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() == reflect.Interface {
		return t, true
	}
	found := false
	for t != nil && t.Kind() == reflect.Ptr {
		e := t.Elem()
		switch e.Kind() {
		case reflect.Struct, reflect.Map, reflect.Array, reflect.Slice:
			return t, false
		default:
			found = true
			t = e
		}
	}
	return t, found
}

func Name(a interface{}) string {
	return reflect.TypeOf(a).Name()
}
