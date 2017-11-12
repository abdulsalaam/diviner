package lmsr

import "testing"

func TestNewEvent(t *testing.T) {
	user := "test"
	title := "title"

	_, err := NewEvent(user, title, "a", "b")
	if err != nil {
		t.Fatal(err)
	}
}
