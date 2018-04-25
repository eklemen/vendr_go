package controllers

import (
	"github.com/eklemen/vendr/models"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"net/http"
)

func ListEvents(c echo.Context) error {
	var e []models.Event
	r := DB.Find(&e)
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusOK, r.Value)
}

func CreateEvent(c echo.Context) error {
	e := new(models.Event)
	//n := models.User{Uuid: GetBearerUuid(c)}
	//f := DB.Where(&n).First(&n)
	if err := c.Bind(&e); err != nil {
		return err
	}
	e.Uuid = uuid.NewV4()
	//e.Creator = f.Value
	DB.Create(&e)
	return c.JSON(http.StatusCreated, &e)
}
