package util

import (
	"strconv"
	"testing"
)

func TestCmpBytes(t *testing.T) {
	t1 := []byte("ABC")
	t2 := []byte("ABC")
	t3 := []byte("ABCDE")
	t4 := []byte("DEF")

	if CmpByteArray(t1, t2) != 0 {
		t.Fatal("compare failed")
	}

	if CmpByteArray(t1, t3) != -1 {
		t.Fatal("compare failed")
	}

	if CmpByteArray(t1, t4) != -1 {
		t.Fatal("compare failed")
	}

	if CmpByteArray(t3, t1) != 1 {
		t.Fatal("compare failed")
	}

	if CmpByteArray(t4, t1) != 1 {
		t.Fatal("compare failed")
	}
}

func TestMap(t *testing.T) {
	data := []string{"1.0", "2.0", "3.0"}

	data1 := MapToFloat64(data, func(arg1 interface{}) (float64, bool) {
		tmp, err := strconv.ParseFloat(arg1.(string), 64)
		if err != nil {
			return 0.0, false
		}
		return tmp, true
	})

	if data1 == nil {
		t.Fatal("map failed")
	}

	for i, x := range data {
		if tmp, _ := strconv.ParseFloat(x, 64); tmp != data1[i] {
			t.Fatal("map data error")
		}
	}

	data2 := MapToFloat64("abcdef", func(arg1 interface{}) (float64, bool) {
		return 0.0, true
	})

	if data2 != nil {
		t.Fatal("can not handle non-slice data")
	}

	data3 := MapToFloat64(data, func(arg1 interface{}) (float64, bool) {
		return 0.0, false
	})

	if data3 != nil {
		t.Fatal("can not handle when op returns false")
	}
}

func TestFold(t *testing.T) {
	data := []string{"1.0", "2.0", "3.0"}

	sum1, ok := FoldLeftFloat64(data, 0.0, func(a float64, b interface{}) (float64, bool) {
		tmp, err := strconv.ParseFloat(b.(string), 64)
		if err != nil {
			return 0.0, false
		}
		return a + tmp, true
	})
	if !ok {
		t.Fatal("fold failed")
	}

	sum2 := 0.0
	for _, x := range data {
		tmp, _ := strconv.ParseFloat(x, 64)
		sum2 += tmp
	}

	if sum2 != sum1 {
		t.Fatal("fold sum failed")
	}

	_, ok = FoldLeftFloat64("123", 0.0, func(a float64, b interface{}) (float64, bool) {
		return 0.0, true
	})

	if ok {
		t.Fatal("can not handle non-slice")
	}

	_, ok = FoldLeftFloat64(data, 0.0, func(a float64, b interface{}) (float64, bool) {
		return 0.0, false
	})

	if ok {
		t.Fatal("can not handle when op returns false")
	}
}
