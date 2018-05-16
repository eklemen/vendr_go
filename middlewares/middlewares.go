package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/eklemen/vendr/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

var DB *gorm.DB

func LoadUserIntoContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		usr := c.Get("user").(*jwt.Token)
		claims := usr.Claims.(jwt.MapClaims)

		// set the user id into Context
		id := int(claims["id"].(float64))
		c.Set("userId", id)

		// set the user uuid into Context
		uid, _ := uuid.FromString(claims["uuid"].(string))
		c.Set("uuid", uid)

		//user := models.User{}
		//err := DB.First(&user, id).Error
		//if err != nil {
		//	c.JSON(http.StatusNotFound, "Not found")
		//}
		//c.Set("user", user)
		return next(c)
	}
}

func GetEventIDFromUUID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// find event by uuid and set event id into context
		euid, _ := uuid.FromString(c.Param("uuid"))
		e := &models.Event{Uuid: euid}
		DB.Select("id").Where(&e).Find(&e)
		c.Set("eventId", e.ID)

		// Get the EventUser row for the current User in the selected Event
		userId := c.Get("userId").(int)
		eu := &models.EventUser{EventID: e.ID, UserID: userId}
		DB.Where(&eu).First(&eu)
		c.Set("myEvent", eu)

		return next(c)
	}
}
