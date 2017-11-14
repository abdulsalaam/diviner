package lmsr

import "testing"

func TestNewEvent(t *testing.T) {
	user := "test"
	title := "title"

	evt, err := NewEvent(user, title, "a", "b")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewEvent(user, title, "a")
	if err == nil {
		t.Fatal("can not new an event with one outcome")
	}

	bytes, err := MarshalEvent(evt)
	if err != nil {
		t.Fatal("marshal failed")
	}

	evt2, err := UnmarshalEvent(bytes)
	if err != nil {
		t.Fatal("unmarshal failed")
	}

	if !CmpEvent(evt, evt2) {
		t.Fatal("data not match")
	}

	evt3, _ := NewEvent(user, title, "a", "b")

	if CmpEvent(evt3, evt) {
		t.Fatal("compare failed")
	}

	evt3.Id = evt.Id
	if CmpEvent(evt3, evt) {
		t.Fatal("compare failed")
	}

	evt4, _ := NewEvent(user, title, "a", "b", "c")
	evt4.Id = evt.Id

	if CmpEvent(evt4, evt) {
		t.Fatal("compare failed")
	}
}
