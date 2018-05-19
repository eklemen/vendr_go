package controllers

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/eklemen/vendr/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/markbates/goth"
	"github.com/satori/go.uuid"
	"net/http"
	"os"
	"time"
)

var DB *gorm.DB

type (
	createdUser struct {
		User  interface{} `json:"user"`
		Token string      `json:"token"`
	}
)

/////////////////////////////////
// User direct actions
/////////////////////////////////

func ListUsers(c echo.Context) error {
	var users []models.User
	r := DB.Find(&users)
	if r.Error != nil {
		return r.Error
	}
	return c.JSON(http.StatusOK, r.Value)
}

func GetUser(c echo.Context) error {
	u := models.NewUser()
	uid, _ := uuid.FromString(c.Param("uuid"))
	r := DB.Preload("CreatedEvents").
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

func CreateUser(c echo.Context, user goth.User) error {
	u := models.NewUser()

	// Search for existing user
	u.IgID = user.UserID
	u.IgUsername = user.NickName
	f := DB.Where(
		&models.User{
			IgID:       u.IgID,
			IgUsername: u.IgUsername,
		}).First(&u)

	// Create uuid to encode in JWT
	uid := uuid.NewV4()
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims (the DB id is encoded below)
	claims := token.Claims.(jwt.MapClaims)
	claims["uuid"] = uid

	cookie := new(http.Cookie)
	cookie.Name = "vendrToken"
	cookie.Expires = time.Now().Add(24 * time.Hour)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	cookie.Value = t
	c.SetCookie(cookie)

	if f.RecordNotFound() {
		u.Uuid = uid
		u.Email = user.Email
		u.IgID = user.UserID
		u.IgUsername = user.NickName
		u.IgFullName = user.FirstName + " " + user.LastName
		u.IgToken = user.AccessToken
		u.IgPic = user.AvatarURL
		DB.Create(&u)
		// Generate encoded token and send it as response.
		claims["id"] = u.ID
		if err != nil {
			return err
		}
		return c.Redirect(302, "http://localhost:3000/dashboard?token="+t)
	} else {
		claims["id"] = u.ID
		if err != nil {
			return err
		}
		return c.Redirect(302, "http://localhost:3000/dashboard?token="+t)
	}
}

func UpdateUser(c echo.Context) error {
	uid, _ := uuid.FromString(c.Param("uuid"))
	t := c.Get("uuid").(uuid.UUID)
	// reject if user tries to update another user
	if t != uid {
		return c.JSON(http.StatusUnauthorized, "You cannot update this user")
	}
	u := &models.User{Uuid: uid}
	if err := c.Bind(u); err != nil {
		return err
	}
	// update the user
	DB.Model(&u).Updates(&u)

	// return the full record after update
	DB.Where(&models.User{Uuid: uid}).
		First(&u)

	return c.JSON(http.StatusOK, u)
}

func DeleteUser(c echo.Context) error {
	uid, _ := uuid.FromString(c.Param("uuid"))
	u := &models.User{Uuid: uid}
	DB.Where(&u).First(&u)
	DB.Delete(&u)
	return c.NoContent(http.StatusNoContent)
}

/////////////////////////////////
// A User's events via EventUser
/////////////////////////////////
func GetSelfEventList(c echo.Context) error {
	userId := c.Get("userId").(int)
	var e []models.EventUser
	err := DB.Preload("Event").
		Where(&models.EventUser{UserID: userId}).
		Find(&e).Error
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, e)
}

func GetUsersEventList(c echo.Context) error {
	var (
		uid  uuid.UUID
		user *models.User
		e    []models.EventUser
	)
	uid, _ = uuid.FromString(c.Param("uuid"))
	user = &models.User{Uuid: uid}
	DB.Select("id").Where(&user).First(&user)

	DB.Preload("Event").
		Where(&models.EventUser{UserID: user.ID}).
		Find(&e)
	return c.JSON(http.StatusOK, e)
}

/////////////////////////////////
// Users contacts (address book)
/////////////////////////////////
func GetContactList(c echo.Context) error {
	var (
		uid uuid.UUID
		u   *models.User
	)
	uid, _ = uuid.FromString(c.Param("uuid"))
	u = &models.User{Uuid: uid, ContactList: []*models.User{}}
	r := DB.Preload("ContactList").
		Where(&models.User{Uuid: uid}).
		First(&u)
	if r.RecordNotFound() {
		return c.JSON(http.StatusNotFound, "Record not found")
	}
	if r.Error != nil {
		return r.Error
	}

	return c.JSON(http.StatusOK, u.ContactList)
}

func AddContact(c echo.Context) error {
	uid, _ := uuid.FromString(c.Param("uuid"))
	contact := models.User{Uuid: uid}
	err := DB.Where(&contact).First(&contact).Error
	if err != nil {
		return err
	}
	userId := c.Get("userId").(int)
	user := models.User{ID: userId}
	uerr := DB.First(&user).Error
	if uerr != nil {
		return uerr
	}

	DB.Model(&user).
		Association("ContactList").
		Append(&contact)

	return c.JSON(http.StatusOK, user.ContactList)
}

func RemoveContact(c echo.Context) error {
	userId := c.Get("userId").(int)
	user := &models.User{ID: userId, ContactList: []*models.User{}}
	uid, _ := uuid.FromString(c.Param("uuid"))
	contact := models.User{Uuid: uid}
	cerr := DB.First(&contact).Error
	if cerr != nil {
		return cerr
	}
	DB.Model(&user).Association("ContactList").Delete(&contact)

	count := DB.Model(&user).Association("ContactList").Count()
	if count == 0 {
		return c.JSON(http.StatusOK, []*models.User{})
	}
	err := DB.Preload("ContactList").First(&user).Error
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, user.ContactList)
}
