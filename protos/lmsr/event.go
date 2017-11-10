package lmsr

import (
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
