package models

import (
	"time"
)

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

// Permissions
var (
	writeEvent  = 1
	readEvent   = 2
	deleteEvent = 4
	addUser     = 8
	removeUser  = 16
)

// User groups
var administrator = writeEvent | readEvent | deleteEvent | addUser | removeUser
var owner = readEvent | writeEvent | deleteEvent | removeUser
var editor = writeEvent | readEvent
var member = readEvent

func (e *EventUser) GrantMember() int {
	return member
}
func (e *EventUser) GrantEditor() int {
	return editor
}
func (e *EventUser) GrantOwner() int {
	return owner
}

func (e *EventUser) CanEditEvent() bool {
	if e.MemberPermission&editor == editor {
		return true
	}
	return false
}

func (e *EventUser) CanDeleteEvent() bool {
	if e.MemberPermission&owner == owner {
		return true
	}
	return false
}
