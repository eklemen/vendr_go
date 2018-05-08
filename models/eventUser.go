package models

import (
	"fmt"
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
	deleteUser  = 16
)

// User groups
var administrator = writeEvent | readEvent | deleteEvent | addUser | deleteUser
var owner = readEvent | writeEvent | deleteEvent | deleteUser
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
	if e.MemberPermission == editor {
		fmt.Println("-----THE USER CAN EDIT!!!-----")
		return true
	}
	return false
}
