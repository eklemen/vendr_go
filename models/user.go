package models

import (
	"github.com/satori/go.uuid"
	"time"
)

type (
	User struct {
		ID              int          `json:"-"`
		Uuid            uuid.UUID    `json:"uuid"`
		Email           string       `json:"email"`
		Phone           string       `json:"phone"`
		IgID            string       `json:"igId"`
		IgUsername      string       `json:"igUsername"`
		IgFullName      string       `json:"igFullName"`
		IgToken         string       `json:"igToken"`
		IgPic           string       `json:"igPic"`
		DeletedAt       *time.Time   `json:"-"`
		CreatedEvents   []*Event     `json:"createdEvents,omitempty" gorm:"foreignkey:CreatorID"`
		EventsAttending []*EventUser `json:"eventsAttending,omitempty"`
		ContactList     []*User      `json:"contacts,omitempty" gorm:"many2many:contacts;association_jointable_foreignkey:contact_id"`
	}
)

func NewUser() *User {
	return &User{
		CreatedEvents:   []*Event{},
		EventsAttending: []*EventUser{},
		ContactList:     []*User{},
	}
}
