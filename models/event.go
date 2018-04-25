package models

import "github.com/satori/go.uuid"

type (
	Event struct {
		ID        int       `json:"-"`
		Uuid      uuid.UUID `json:"uuid"`
		Venue     string    `json:"venue"`
		EventDate string    `json:"eventDate"`
		Title     string    `json:"title"`
		Creator   User      `json:"creator"`
		CreatorID int       `json:"creatorId"`
	}
)

func NewEvent() *Event {
	return &Event{}
	// TODO: add the array of members later
	//Creator: []Event{},
	//}
}
