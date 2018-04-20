package models

type (
	Event struct {
		ID        int    `json:"-"`
		Uuid      string `json:"uuid"`
		Venue     string `json:"venue"`
		EventDate string `json:"eventDate"`
		Title     string `json:"title"`
		Creator   User   `json:"creator"`
		CreatorID int    `json:"creatorId"`
		DeletedAt string `json:"-"`
	}
)

func NewEvent() *Event {
	return &Event{}
	// TODO: add the array of members later
	//Creator: []Event{},
	//}
}
