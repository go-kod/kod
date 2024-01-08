package kod

import (
	"reflect"
)

func rtype[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
