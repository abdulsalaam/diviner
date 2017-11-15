package lmsr

import "testing"

func TestAsset(t *testing.T) {
	if _, err := NewAsset("a", "b", -1); err == nil {
		t.Error("can not create a asset with -1 volume")
	}

	a1, err := NewAsset("a", "b", 100.0)
	if err != nil {
		t.Errorf("create a asset failed: %v", err)
	}

	if a1.Id != AssetID("a", "b") {
		t.Errorf("id failed: %s, %s", a1.Id, AssetID("a", "b"))
	}

	if a1.Volume != 100.0 {
		t.Errorf("volume failed")
	}

	bytes, err := MarshalAsset(a1)
	if err != nil {
		t.Errorf("marshal failed: %v", err)
	}

	a2, err := UnmarshalAsset(bytes)
	if err != nil {
		t.Errorf("unmarshal failed: %v", err)
	}

	if *a1 != *a2 {
		t.Errorf("data not match: %v, %v", a1, a2)
	}
}
