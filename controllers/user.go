package controllers

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/eklemen/vendr/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/markbates/goth"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"strconv"
)

var DB *gorm.DB

type (
	createdUser struct {
		User  interface{} `json:"user"`
		Token string      `json:"token"`
	}
)

func CreateUser(c echo.Context, user goth.User) error {
	//u := new(User) equivalent to line below
	u := models.NewUser()

	// Search for existing user
	u.IgID = user.UserID
	u.IgUsername = user.NickName
	// Can this be shortened to
	// DB.First(&u) ?
	f := DB.Where(
		&models.User{
			IgID:       u.IgID,
			IgUsername: u.IgUsername,
		}).First(&u)

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	if f.RecordNotFound() {
		fmt.Println("NOT FOUND")
		id := uuid.NewV4()
		u.Uuid = id
		// TODO: make an interface for this?
		u.Email = user.Email
		u.IgID = user.UserID
		u.IgUsername = user.NickName
		u.IgFullName = user.FirstName + " " + user.LastName
		u.IgToken = user.AccessToken
		u.IgPic = user.AvatarURL
		DB.Create(&u)
		nu := &createdUser{
			Token: t,
			User:  u,
		}
		return c.JSON(http.StatusCreated, &nu)
	} else {
		nu := &createdUser{
			Token: t,
			User:  f.Value,
		}
		return c.JSON(http.StatusOK, nu)
	}
	// c.Bind maps the request to the given struct
	//if err := c.Bind(&u); err != nil {
	//	return err
	//}
}

func GetAllUsers(c echo.Context) error {
	var users []models.User
	res := DB.Find(&users)
	if res.Error != nil {
		return res.Error
	}
	return c.JSON(http.StatusOK, res.Value)
}

func UpdateUser(c echo.Context) error {
	//u := new(User)
	id, _ := strconv.Atoi(c.Param("id"))
	u := &models.User{ID: id}
	if err := c.Bind(u); err != nil {
		return err
	}
	DB.Model(&u).Updates(&u)
	return c.JSON(http.StatusOK, &u)
}

func DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	u := &models.User{ID: id}
	DB.Delete(&u)
	return c.NoContent(http.StatusNoContent)
}

func GetUser(c echo.Context) error {
	var u models.User
	uuid := c.Param("uuid")
	fmt.Println("ID", uuid)
	r := DB.Preload("CreatedEvents").
		// Gives error type string and type uuid
		//Where(&models.User{Uuid: uuid}).
		Where("uuid = ?", uuid).
		First(&u, c.Param("id"))

	if r.RecordNotFound() {
		return c.JSON(http.StatusNotFound, "Record not found")
	}
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusOK, r.Value)
}
