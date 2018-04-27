package models

type (
	EventUser struct {
		ID               int    `json:"id"`
		EventID          int    `json:"eventId"`
		UserID           int    `json:"userId"`
		MemberRole       string `json:"memberRole"`
		MemberPermission string `json:"memberPermission"`
	}
)
