package util

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
