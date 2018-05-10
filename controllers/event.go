package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/eklemen/vendr/models"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"net/http"
)

type memberRoleAndPerm struct {
	Role string `json:"role"`
}

func ListEvents(c echo.Context) error {
	var e []models.Event
	err := DB.Preload("Creator").
		Preload("Attendees.User").
		Find(&e).Error
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, e)
}

func GetEvent(c echo.Context) error {
	e := new(models.Event)
	uid, _ := uuid.FromString(c.Param("uuid"))
	r := DB.Preload("Creator").
		Preload("Attendees.User").
		Where(&models.Event{Uuid: uid}).
		First(&e)
	if r.RecordNotFound() {
		return c.JSON(http.StatusNotFound, "Record not found")
	}
	return c.JSON(http.StatusOK, r.Value)
}

func CreateEvent(c echo.Context) error {
	//e := new(models.Event)
	//eu := new(models.EventUser)
	var (
		e  = models.Event{}
		eu = models.EventUser{}
	)

	// get user id of token bearer
	userId := c.Get("userId").(int)
	if err := c.Bind(&e); err != nil {
		return err
	}

	// set the creator
	e.CreatorID = userId
	e.Uuid = uuid.NewV4()
	e.Attendees = []*models.EventUser{
		{
			UserID:           userId,
			MemberPermission: eu.GrantOwner(),
			MemberRole:       "vendor",
		},
	}

	// create the event
	DB.Set("gorm:association_autoupdate", false).Save(&e)

	err := DB.Preload("Attendees.User").
		Preload("Creator").
		First(&e, e.ID).Error
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, e)
}

func UpdateEvent(c echo.Context) error {
	uid, _ := uuid.FromString(c.Param("uuid"))
	e := &models.Event{Uuid: uid}

	// get the event users details and check the permission
	i := c.Get("myEvent").(*models.EventUser)
	if !i.CanEditEvent() {
		return c.JSON(http.StatusNotFound, "Not found")
	}

	// else bind the incoming data and update
	if err := c.Bind(e); err != nil {
		return err
	}
	DB.Model(&e).Updates(&e)

	// get the newly updated record
	err := DB.Preload("Creator").
		Preload("Attendees.User").
		Where(&models.Event{Uuid: uid}).
		First(&e).Error
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, e)
}

func DeleteEvent(c echo.Context) error {
	// get the event users details and check the permission
	i := c.Get("myEvent").(*models.EventUser)
	if !i.CanDeleteEvent() {
		return c.JSON(http.StatusNotFound, "Not found")
	}

	uid, _ := uuid.FromString(c.Param("uuid"))
	e := &models.Event{Uuid: uid}
	DB.Delete(&e)
	return c.NoContent(http.StatusNoContent)
}

// EventUser actions

func JoinEvent(c echo.Context) error {
	// get eventId and userId from context
	eId := c.Get("eventId").(int)
	u := c.Get("userId").(int)
	// Decode the request body and grab the role
	var memberRole memberRoleAndPerm
	err := json.NewDecoder(c.Request().Body).Decode(&memberRole)
	if err != nil {
		return err
	}

	eu := models.EventUser{
		EventID: eId,
		UserID:  u}

	if DB.Where(&eu).First(&eu).RecordNotFound() {
		fmt.Println("NEW RECORD")
		eu.MemberRole = memberRole.Role
		// verify the request has correct 'role'
		if memberRole.Role == "vendor" {
			eu.MemberPermission = eu.GrantMember()
		} else if memberRole.Role == "client" {
			// give a client read and write permissions by default
			eu.MemberPermission = eu.GrantOwner()
		} else {
			return c.JSON(http.StatusBadRequest, "Member role must be either 'vendor' or 'client'.")
		}
		// create new EventUser relationship row
		DB.Create(&eu)
	} else {
		return c.JSON(http.StatusBadRequest, "User already belongs to this event.")
	}

	// if the join is successful, return that event
	event := &models.Event{}
	dbErr := DB.Preload("Attendees.User").
		Preload("Creator").
		First(&event, eId).Error
	if dbErr != nil {
		return dbErr
	}
	return c.JSON(http.StatusOK, event)
}

func RemoveUserFromEvent(c echo.Context) error {
	eId := c.Get("eventId").(int)
	userUuid, _ := uuid.FromString(c.Param("userUuid"))
	u := &models.User{Uuid: userUuid}
	DB.Select("id").Where(&u).First(&u)

	eu := &models.EventUser{UserID: u.ID, EventID: eId}
	DB.Where(&eu).First(&eu)
	DB.Delete(&eu)

	return c.NoContent(http.StatusNoContent)
}

// Each member of an event has a "report" column
// spam accounts or vendors that were not actually present
// can be 'reported' by other users
// too many reports will remove that users from the event
// TODO: move this struct out to be used by any custom response
type userStatus struct {
	Message     string `json:"message"`
	UserRemoved bool   `json:"userRemoved"`
}

func ReportUser(c echo.Context) error {
	status := userStatus{}
	eId := c.Get("eventId").(int)
	userUuid, _ := uuid.FromString(c.Param("userUuid"))
	u := &models.User{Uuid: userUuid}
	self := c.Get("userId")
	DB.Select("id").Where(&u).First(&u)

	// Return bad request if the user tries to report themselves
	if self == u.ID {
		status.Message = "Incorrect user uuid"
		status.UserRemoved = false
		return c.JSON(http.StatusBadRequest, status)
	}

	eu := &models.EventUser{UserID: u.ID, EventID: eId}
	DB.Where(&eu).First(&eu)

	if eu.Reports >= 5 {
		status.Message = "User has been removed from event."
		status.UserRemoved = true
		DB.Delete(&eu)
		return c.JSON(http.StatusOK, status)
	} else {
		DB.Model(&eu).Update("reports", eu.Reports+1)
	}

	status.Message = "User has been reported."
	status.UserRemoved = false
	return c.JSON(http.StatusOK, status)
}
