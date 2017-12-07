package market

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

// NewEvent ...
func NewEvent(user, title string, end time.Time, outcomes ...string) (*Event, error) {
	if len(outcomes) <= 1 {
		return nil, fmt.Errorf("length of outcomes must larger than 1")
	}

	if end.Before(time.Now()) {
		return nil, fmt.Errorf("end time must be latter than now: %v", end)
	}

	ts, err := ptypes.TimestampProto(end)
	if err != nil {
		return nil, err
	}

	event := &Event{
		Id:       EventID(),
		User:     user,
		Title:    title,
		Approved: false,
		Allowed:  true,
		End:      ts,
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

// FindOutcome ...
func FindOutcome(event *Event, outcome string) int {
	for i, x := range event.Outcomes {
		if x.Id == outcome {
			return i
		}
	}

	return -1
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
