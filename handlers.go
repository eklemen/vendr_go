package main

//
//import (
//	"github.com/eklemen/vendr/models"
//	"github.com/labstack/echo"
//	"github.com/satori/go.uuid"
//	"net/http"
//	"strconv"
//)
//
//func CreateUser(c echo.Context) error {
//	//u := new(User) equivalent to line below
//	u := models.NewUser()
//	id := uuid.NewV4()
//	u.Uuid = id
//
//	// c.Bind maps the request to the given struct
//	if err := c.Bind(&u); err != nil {
//		return err
//	}
//	db.Create(&u)
//	return c.JSON(http.StatusCreated, &u)
//
//}
//
//func GetAllUsers(c echo.Context) error {
//	var users []models.User
//	res := db.Find(&users)
//	if res.Error != nil {
//		return res.Error
//	}
//	return c.JSON(http.StatusOK, res.Value)
//}
//
//func UpdateUser(c echo.Context) error {
//	//u := new(User)
//	id, _ := strconv.Atoi(c.Param("id"))
//	u := &models.User{ID: id}
//	if err := c.Bind(u); err != nil {
//		return err
//	}
//	db.Model(&u).Updates(&u)
//	return c.JSON(http.StatusOK, &u)
//}
//
//func DeleteUser(c echo.Context) error {
//	id, _ := strconv.Atoi(c.Param("id"))
//	u := &models.User{ID: id}
//	db.Delete(&u)
//	return c.NoContent(http.StatusNoContent)
//}
//
//func GetUser(c echo.Context) error {
//	var user models.User
//	res := db.Preload("CreatedEvents").First(&user, c.Param("id"))
//
//	if res.RecordNotFound() {
//		return c.JSON(http.StatusNotFound, "Record not found")
//	}
//	if res.Error != nil {
//		return res.Error
//	}
//	return c.JSON(http.StatusOK, res.Value)
//}
