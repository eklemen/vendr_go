package controllers

import (
	"github.com/labstack/echo"
	"github.com/markbates/goth/gothic"
	"net/http"
)

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
