package monad

// MapToFloat64 ...
func MapToFloat64(x interface{}, op func(interface{}) float64) []float64 {
	v, ok := valueOf(x)
	if !ok {
		return nil
	}

	len := v.Len()

	result := make([]float64, len)

	for i := 0; i < len; i++ {
		result[i] = op(v.Index(i).Interface())
	}

	return result

}

// FoldLeftFloat64 ...
func FoldLeftFloat64(x interface{}, z float64, op func(float64, interface{}) float64) (float64, bool) {
	v, ok := valueOf(x)
	if !ok {
		return 0.0, false
	}

	len := v.Len()

	for i := 0; i < len; i++ {
		z = op(z, v.Index(i).Interface())
	}

	return z, true
}
