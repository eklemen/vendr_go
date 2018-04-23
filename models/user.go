package models

import "github.com/satori/go.uuid"

type (
	User struct {
		ID            int       `json:"-"`
		Uuid          uuid.UUID `json:"uuid"`
		Token         string    `json:"token"`
		Email         string    `json:"email"`
		Phone         string    `json:"phone"`
		IgID          string    `json:"igId"`
		IgUsername    string    `json:"igUsername"`
		IgFullName    string    `json:"igFullName"`
		IgToken       string    `json:"igToken"`
		IgPic         string    `json:"igPic"`
		CreatedEvents []Event   `json:"createdEvents" gorm:"foreignkey:CreatorID"`
		//DeletedAt     string    `json:"-"`
	}
)

// How to convert this to User
type (
	IgUser struct {
		Email       string `json:"email"`
		UserId      string `json:"igId"`
		NickName    string `json:"igUsername"`
		AccessToken string `json:"igToken"`
		AvatarURL   string `json:"igPic"`
	}
)

func NewUser() *User {
	return &User{
		CreatedEvents: []Event{},
	}
}
