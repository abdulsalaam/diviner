package util

import "reflect"

// CmpByteArray ...
func CmpByteArray(a, b []byte) int {
	lenA := len(a)
	lenB := len(b)

	if lenA > lenB {
		return 1
	} else if lenA < lenB {
		return -1
	}

	for i := 0; i < lenA; i++ {
		if a[i] > b[i] {
			return 1
		} else if a[i] < b[i] {
			return -1
		}
	}

	return 0
}

func valueOf(x interface{}) (reflect.Value, bool) {
	v := reflect.ValueOf(x)
	if v.Kind() != reflect.Slice {
		return v, false
	}
	return v, true
}

// MapToFloat64 ...
func MapToFloat64(x interface{}, op func(interface{}) (float64, bool)) []float64 {
	v, ok := valueOf(x)
	if !ok {
		return nil
	}

	len := v.Len()

	result := make([]float64, len)

	for i := 0; i < len; i++ {
		tmp, ok := op(v.Index(i).Interface())
		if ok {
			result[i] = tmp
		} else {
			return nil
		}
	}

	return result

}

// FoldLeftFloat64 ...
func FoldLeftFloat64(x interface{}, z float64, op func(float64, interface{}) (float64, bool)) (float64, bool) {
	v, ok := valueOf(x)
	if !ok {
		return 0.0, false
	}

	len := v.Len()

	for i := 0; i < len; i++ {
		z, ok = op(z, v.Index(i).Interface())
		if !ok {
			return 0.0, false
		}
	}

	return z, true
}
