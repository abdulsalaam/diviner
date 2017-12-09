package prediction

import (
	fmt "fmt"
	"strings"
	"time"

	proto "github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
)

// NewEventCreateRequest ...
func NewEventCreateRequest(user, title string, end time.Time, outcomes ...string) (*EventCreateRequest, error) {
	if user = strings.TrimSpace(user); user == "" {
		return nil, fmt.Errorf("user is empty")
	}

	if title = strings.TrimSpace(title); title == "" {
		return nil, fmt.Errorf("title is empty")
	}

	if end.Before(time.Now()) {
		return nil, fmt.Errorf("end time is expired: %v", end)
	}

	pbEnd, err := ptypes.TimestampProto(end)
	if err != nil {
		return nil, err
	}

	if len(outcomes) < 2 {
		return nil, fmt.Errorf("size of outcomes must be larger than 1, but %d", len(outcomes))
	}

	for i, x := range outcomes {
		if outcomes[i] = strings.TrimSpace(x); outcomes[i] == "" {
			return nil, fmt.Errorf("%d outcome is empty", i)
		}
	}

	return &EventCreateRequest{
		User:     user,
		Title:    title,
		Outcomes: outcomes,
		End:      pbEnd,
	}, nil
}

// CheckEventCreateRequest ...
func CheckEventCreateRequest(in *EventCreateRequest) (bool, error) {
	if strings.TrimSpace(in.User) == "" {
		return false, fmt.Errorf("user is empty")
	}

	if strings.TrimSpace(in.Title) == "" {
		return false, fmt.Errorf("title is empty")
	}

	if in.End.Seconds <= ptypes.TimestampNow().Seconds {
		return false, fmt.Errorf("end time is expired: %v", in.End)
	}

	if len(in.Outcomes) < 2 {
		return false, fmt.Errorf("size of outcomes must be larger than 1, but %d", len(in.Outcomes))
	}

	for i, x := range in.Outcomes {
		if strings.TrimSpace(x) == "" {
			return false, fmt.Errorf("%d outcome is empty", i)
		}
	}

	return true, nil
}

// NewEvent ...
func NewEvent(req *EventCreateRequest) (*Event, error) {
	if ok, err := CheckEventCreateRequest(req); err != nil {
		return nil, err
	} else if !ok {
		return nil, fmt.Errorf("data error")
	}

	event := &Event{
		Id:       EventID(),
		User:     req.User,
		Title:    req.Title,
		Approved: false,
		Allowed:  true,
		End:      req.End,
	}

	tmp := make([]*Outcome, len(req.Outcomes))

	for i, x := range req.Outcomes {
		o := &Outcome{
			Id:    OutcomeID(event.Id, i),
			Title: x,
		}
		tmp[i] = o
	}

	event.Outcomes = tmp

	return event, nil
}

// UnmarshalEventCreateRequest ...
func UnmarshalEventCreateRequest(data []byte) (*EventCreateRequest, error) {
	ret := new(EventCreateRequest)
	if err := proto.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// UnmarshalEvent ...
func UnmarshalEvent(data []byte) (*Event, error) {
	ret := new(Event)
	if err := proto.Unmarshal(data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

// MarshalEventCreateRequest ...
func MarshalEventCreateRequest(in *EventCreateRequest) ([]byte, error) {
	return proto.Marshal(in)
}

// MarshalEvent ...
func MarshalEvent(in *Event) ([]byte, error) {
	return proto.Marshal(in)
}
