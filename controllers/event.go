package controllers

import (
	"encoding/json"
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
	r := DB.Preload("Creator").
		Preload("Attendees.User").
		Find(&e)
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusOK, r.Value)
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
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusOK, r.Value)
}

func CreateEvent(c echo.Context) error {
	e := new(models.Event)

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
			MemberPermission: "edit",
			MemberRole:       "vendor",
		},
	}

	// create the event
	DB.Set("gorm:association_autoupdate", false).Save(&e)

	r := DB.Preload("Attendees.User").
		Preload("Creator").
		First(&e, e.ID)
	if r.Error != nil {
		return r.Error
	}

	return c.JSON(http.StatusCreated, r.Value)
}

func UpdateEvent(c echo.Context) error {
	uid, _ := uuid.FromString(c.Param("uuid"))
	e := &models.Event{Uuid: uid}
	if err := c.Bind(e); err != nil {
		return err
	}
	DB.Model(&e).Updates(&e)
	r := DB.Preload("Creator").
		Preload("Attendees.User").
		Where(&models.Event{Uuid: uid}).
		First(&e)
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusOK, r.Value)
}

func DeleteEvent(c echo.Context) error {
	uid, _ := uuid.FromString(c.Param("uuid"))
	e := &models.Event{Uuid: uid}
	DB.Delete(&e)
	return c.NoContent(http.StatusNoContent)
}

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

	// for permissions
	// (user.priv & 2 = 2)  <- if that is true then the user has that permission
	// check the permission in the event middleware
	eu := models.EventUser{
		EventID:          eId,
		UserID:           u,
		MemberRole:       memberRole.Role,
		MemberPermission: "read"}
	DB.FirstOrCreate(&eu, eu)

	r := DB.Preload("Attendees.User").
		Preload("Creator").
		First(&models.Event{}, eId)
	if r.Error != nil {
		return r.Error
	}

	return c.JSON(http.StatusOK, r.Value)
}
