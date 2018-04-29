package models

type (
	EventUser struct {
		ID               int    `json:"id"`
		EventID          int    `json:"eventId"`
		Event            *Event `json:"event"`
		UserID           int    `json:"userId"`
		User             *User  `json:"user"`
		MemberRole       string `json:"memberRole"`
		MemberPermission string `json:"memberPermission"`
	}
)

//
//func (*EventUser) TableName() string {
//	return "event_users"
//}
