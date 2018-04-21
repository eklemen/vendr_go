package controllers

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/markbates/goth/gothic"
	"net/http"
)

func AuthInstagram(c echo.Context) error {
	res := c.Response().Writer
	req := c.Request()
	fmt.Println("RES", res)
	fmt.Println("REQ", req)
	if gothUser, err := gothic.CompleteUserAuth(res, req); err == nil {
		fmt.Println("Woo", gothUser)
		return c.JSON(http.StatusTemporaryRedirect, gothUser)
	} else {
		fmt.Println("foo")
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

	//return c.Redirect(http.StatusFound, authUrl)
	return c.JSON(http.StatusOK, user)
}
