package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

type eventId struct {
	ID int
}

func SetUserId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		usr := c.Get("user").(*jwt.Token)
		claims := usr.Claims.(jwt.MapClaims)
		// TODO: can remove type casting here
		id := int(claims["id"].(float64))
		c.Set("userId", id)

		uid := claims["uuid"].(string)
		u, _ := uuid.FromString(uid)
		c.Set("uuid", u)
		return next(c)
	}
}

func GetEventIDFromUUID(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		u := c.Param("uuid")
		if u != "" {
			uid, _ := uuid.FromString(u)
			var eventId eventId
			db.Raw("SELECT id FROM events WHERE uuid = ?", uid).Scan(&eventId)
			c.Set("eventId", eventId.ID)
			return next(c)
		}
		return next(c)
	}
}
