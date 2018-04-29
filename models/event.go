package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type (
	Event struct {
		ID        int          `json:"-"`
		Uuid      uuid.UUID    `json:"uuid"`
		Venue     string       `json:"venue"`
		EventDate string       `json:"eventDate"`
		Title     string       `json:"title"`
		Creator   User         `json:"creator"`
		Attendees []*EventUser `json:"attendees"`
		DeletedAt *time.Time   `json:"-"`
		CreatorID int          `json:"-"`
	}
)

func NewEvent() *Event {
	return &Event{}
	// TODO: add the array of members later
	//Creator: []Event{},
	//}
}
