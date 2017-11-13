package lmsr

import (
	"github.com/gogo/protobuf/proto"
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

// UnmarshalEvent ...
func UnmarshalEvent(data []byte) (*Event, error) {
	evt := &Event{}
	err := proto.Unmarshal(data, evt)
	return evt, err
}

// MarshalEvent ...
func MarshalEvent(e *Event) ([]byte, error) {
	return proto.Marshal(e)
}
