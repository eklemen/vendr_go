package models

import "time"

type (
	EventUser struct {
		ID               int        `json:"-"`
		EventID          int        `json:"-"`
		Event            *Event     `json:"event,omitempty"`
		UserID           int        `json:"-"`
		User             *User      `json:"user,omitempty"`
		MemberRole       string     `json:"memberRole"`
		MemberPermission int        `json:"memberPermission"`
		Reports          int        `json:"-";gorm:"default:0"`
		DeletedAt        *time.Time `json:"-"`
	}
)

//
//func (*EventUser) TableName() string {
//	return "event_users"
//}
