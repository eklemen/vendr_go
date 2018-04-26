package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type (
	User struct {
		ID              int        `json:"-"`
		Uuid            uuid.UUID  `json:"uuid"`
		Email           string     `json:"email"`
		Phone           string     `json:"phone"`
		IgID            string     `json:"igId"`
		IgUsername      string     `json:"igUsername"`
		IgFullName      string     `json:"igFullName"`
		IgToken         string     `json:"igToken"`
		IgPic           string     `json:"igPic"`
		DeletedAt       *time.Time `json:"-"`
		CreatedEvents   []Event    `json:"createdEvents" gorm:"foreignkey:CreatorID"`
		EventsAttending []Event    `json:"eventsAttending"`
	}
)

func NewUser() *User {
	return &User{
		CreatedEvents: []Event{},
	}
}
