package market

import (
	"testing"
	"time"
)

func afterOneYear() time.Time {
	return time.Now().AddDate(1, 0, 0)
}

func afterTwoYear() time.Time {
	return time.Now().AddDate(2, 0, 0)
}

func beforeOneYear() time.Time {
	return time.Now().AddDate(-1, 0, 0)
}

func TestNewEvent(t *testing.T) {
	user := "test"
	title := "title"

	evt, err := NewEvent(user, title, afterOneYear(), "a", "b")
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewEvent(user, title, beforeOneYear(), "a", "b")
	if err == nil {
		t.Fatal("can not set an event with wrong time")
	}

	_, err = NewEvent(user, title, afterOneYear(), "a")
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

	if evt.String() != evt2.String() {
		t.Fatal("data not match")
	}
}
