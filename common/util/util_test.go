package util

import "testing"

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
