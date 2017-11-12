package monad

import (
	"reflect"
)

func valueOf(x interface{}) (reflect.Value, bool) {
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Slice {
		return v, false
	}
	return v, true
}
