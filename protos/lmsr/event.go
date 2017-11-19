package lmsr

import (
	proto "github.com/golang/protobuf/proto"
	perrors "github.com/pkg/errors"
)

// NewEvent ...
func NewEvent(user, title string, outcomes ...string) (*Event, error) {
	if len(outcomes) <= 1 {
		return nil, perrors.New("length of outcomes must larger than 1")
	}

	event := &Event{
		Id:       EventID(),
		User:     user,
		Title:    title,
		Approved: false,
	}

	tmp := make([]*Outcome, len(outcomes))

	for i, x := range outcomes {
		o := &Outcome{
			Id:    OutcomeID(event.Id, i),
			Title: x,
		}

		tmp[i] = o
	}

	event.Outcomes = tmp

	return event, nil
}

func FindOutcome(event *Event, outcome string) int {
	for i, x := range event.Outcomes {
		if x.Id == outcome {
			return i
		}
	}

	return -1
}

// CmpEvent ...
func CmpEvent(e1, e2 *Event) bool {
	if e1.Id != e2.Id || e1.User != e2.User || e1.Title != e2.Title || e1.Approved != e2.Approved {
		return false
	}

	if len(e1.Outcomes) != len(e2.Outcomes) {
		return false
	}

	for i := range e1.Outcomes {
		if *(e1.Outcomes[i]) != *(e2.Outcomes[i]) {
			return false
		}
	}

	return true
}

// UnmarshalEvent ...
func UnmarshalEvent(data []byte) (*Event, error) {
	evt := &Event{}
	if err := proto.Unmarshal(data, evt); err != nil {
		return nil, err
	}
	return evt, nil
}

// MarshalEvent ...
func MarshalEvent(e *Event) ([]byte, error) {
	return proto.Marshal(e)
}
