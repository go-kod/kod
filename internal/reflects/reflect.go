package reflects

import "reflect"

// TypeFor returns the reflect.Type of T.
func TypeFor[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
