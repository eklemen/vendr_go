package controllers

import (
	"github.com/eklemen/vendr/models"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"net/http"
)

func ListEvents(c echo.Context) error {
	var e []models.Event
	r := DB.Preload("Creator").Find(&e)
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusOK, r.Value)
}

func GetEvent(c echo.Context) error {
	e := new(models.Event)
	uid, _ := uuid.FromString(c.Param("uuid"))
	r := DB.Preload("Creator").
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
	if err := c.Bind(&e); err != nil {
		return err
	}
	e.CreatorID = c.Get("userId").(int)
	e.Uuid = uuid.NewV4()
	DB.Create(&e)
	r := DB.Preload("Creator").First(&e, e.ID)
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusCreated, r.Value)
}

func DeleteEvent(c echo.Context) error {
	uid, _ := uuid.FromString(c.Param("uuid"))
	e := &models.Event{Uuid: uid}
	DB.Delete(&e)
	return c.NoContent(http.StatusNoContent)
}
