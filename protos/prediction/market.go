package prediction

import (
	fmt "fmt"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
)

// NewMarketCreateRequest ...
func NewMarketCreateRequest(user, event string, num float64, fund bool, start, end time.Time) (*MarketCreateRequest, error) {
	if user = strings.TrimSpace(user); user == "" {
		return nil, fmt.Errorf("user is empty")
	}

	if event = strings.TrimSpace(event); event == "" {
		return nil, fmt.Errorf("event is empty")
	}

	if num <= 0 {
		return nil, fmt.Errorf("num must be larger than 0, but %v", num)
	}

	if start.After(end) || end.Before(time.Now()) {
		return nil, fmt.Errorf("start and end time error: %v ~ %v", start, end)
	}

	s, err := ptypes.TimestampProto(start)
	if err != nil {
		return nil, err
	}

	e, err := ptypes.TimestampProto(end)
	if err != nil {
		return nil, err
	}

	return &MarketCreateRequest{
		User:   user,
		Event:  event,
		Num:    num,
		IsFund: fund,
		Start:  s,
		End:    e,
	}, nil

}

// CheckMarketCreateRequest ...
func CheckMarketCreateRequest(req *MarketCreateRequest) (bool, error) {
	if strings.TrimSpace(req.User) == "" {
		return false, fmt.Errorf("user is empty")
	}

	if strings.TrimSpace(req.Event) == "" {
		return false, fmt.Errorf("event is empty")
	}

	if req.Num <= 0 {
		return false, fmt.Errorf("num must be larger than 0, but %v", req.Num)
	}

	if req.Start.Seconds >= req.End.Seconds || req.End.Seconds <= ptypes.TimestampNow().Seconds {
		return false, fmt.Errorf("start and end time error: %v ~ %v", req.Start, req.End)
	}

	return true, nil
}


