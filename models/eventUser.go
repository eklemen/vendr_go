package models

type (
	EventUser struct {
		ID               int    `json:"-"`
		EventID          int    `json:"-"`
		Event            *Event `json:"event,omitempty"`
		UserID           int    `json:"-"`
		User             *User  `json:"user,omitempty"`
		MemberRole       string `json:"memberRole"`
		MemberPermission string `json:"memberPermission"`
	}
)

//
//func (*EventUser) TableName() string {
//	return "event_users"
//}
