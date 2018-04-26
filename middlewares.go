package main

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
)

func SetUserId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		usr := c.Get("user").(*jwt.Token)
		claims := usr.Claims.(jwt.MapClaims)
		id := int(claims["id"].(float64))
		c.Set("userId", id)

		uid := claims["uuid"].(string)
		u, _ := uuid.FromString(uid)
		c.Set("uuid", u)
		return next(c)
	}
}
