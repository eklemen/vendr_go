package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/markbates/goth/gothic"
	"github.com/satori/go.uuid"
	"net/http"
)

func GetBearerUuid(c echo.Context) uuid.UUID {
	usr := c.Get("user").(*jwt.Token)
	claims := usr.Claims.(jwt.MapClaims)
	uid := claims["uuid"].(string)
	id, _ := uuid.FromString(uid)
	return id
}
func GetBearerId(c echo.Context) int {
	usr := c.Get("user").(*jwt.Token)
	claims := usr.Claims.(jwt.MapClaims)
	id := claims["id"].(int)
	return id
}

func AuthInstagram(c echo.Context) error {
	res := c.Response().Writer
	req := c.Request()
	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		return c.JSON(http.StatusTemporaryRedirect, gothUser)
	} else {
		gothic.BeginAuthHandler(res, req)
		return err
	}
}

func AuthInstagramCB(c echo.Context) error {
	res := c.Response().Writer
	req := c.Request()
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		return err
	}
	u := CreateUser(c, user)
	return u
}
