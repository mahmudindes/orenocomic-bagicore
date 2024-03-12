package utila

import "reflect"

func NilData(val any) bool {
	if val == nil {
		return true
	}
	v := reflect.ValueOf(val)
	switch v.Kind() {
	case reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Pointer,
		reflect.UnsafePointer,
		reflect.Slice:
		return v.IsNil()
	}
	return false
}
