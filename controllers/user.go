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

	// Create uuid to encode in JWT
	uuid := uuid.NewV4()
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = uuid
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	if f.RecordNotFound() {
		fmt.Println("NOT FOUND")

		u.Uuid = uuid
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
	res := DB.Preload("CreatedEvents").Find(&users)
	if res.Error != nil {
		return res.Error
	}
	return c.JSON(http.StatusOK, res.Value)
}

func UpdateUser(c echo.Context) error {
	uuid, _ := uuid.FromString(c.Param("uuid"))
	u := &models.User{Uuid: uuid}
	if err := c.Bind(u); err != nil {
		return err
	}
	DB.Model(&u).Updates(&u)
	r := DB.Preload("CreatedEvents").
		Where(&models.User{Uuid: uuid}).
		First(&u)
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusOK, r.Value)
}

func DeleteUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	u := &models.User{ID: id}
	DB.Delete(&u)
	return c.NoContent(http.StatusNoContent)
}

func GetUser(c echo.Context) error {
	//usr := c.Get("user").(*jwt.Token)
	//claims := usr.Claims.(jwt.MapClaims)
	//uid := claims["uuid"]
	//fmt.Println("uid", uid)
	//fmt.Println("usr", usr)

	var u models.User
	uid, _ := uuid.FromString(c.Param("uuid"))
	r := DB.Preload("CreatedEvents").
		// Gives error type string and type uuid
		Where(&models.User{Uuid: uid}).
		First(&u)

	if r.RecordNotFound() {
		return c.JSON(http.StatusNotFound, "Record not found")
	}
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusOK, r.Value)
}
