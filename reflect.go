package kod

import (
	"reflect"
)

// rtype returns the reflect.Type of T.
func rtype[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
