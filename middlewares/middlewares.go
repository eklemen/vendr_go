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
		// TODO: can remove type casting here
		id := int(claims["id"].(float64))
		c.Set("userId", id)

		uid := claims["uuid"].(string)
		u, _ := uuid.FromString(uid)
		c.Set("uuid", u)

		var user models.User
		DB.First(&user, id)
		c.Set("user", user)
		return next(c)
	}
}

func GetEventIDFromUUID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := c.Param("uuid")
		uid, _ := uuid.FromString(u)
		e := &models.Event{Uuid: uid}
		DB.Select([]string{"id"}).Where(&e).Find(&e)
		c.Set("eventId", e.ID)
		return next(c)
	}
}
