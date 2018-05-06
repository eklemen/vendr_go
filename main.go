package main

import (
	"fmt"
	"github.com/eklemen/vendr/controllers"
	"github.com/eklemen/vendr/middlewares"
	"github.com/eklemen/vendr/models"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/instagram"
	"github.com/subosito/gotenv"
	"os"
)

var db *gorm.DB

func main() {
	gotenv.Load()
	var err error
	// Dont forget to add postgres adapter to imports
	// _ "github.com/jinzhu/gorm/dialects/postgres"
	db, err = gorm.Open(
		"postgres",
		"host="+os.Getenv("DB_HOST")+" user="+os.Getenv("DB_USERNAME")+
			" dbname="+os.Getenv("DB_DATABASE")+" sslmode=disable password="+
			os.Getenv("DB_PASSWORD"))

	if err != nil {
		panic("failed to connect database")
	} else {
		fmt.Println("DB Connected...")
	}
	defer db.Close()
	// TODO: create a struct for these
	controllers.DB = db
	middlewares.DB = db
	db.LogMode(true)
	// Migrate the schema
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Event{})
	db.AutoMigrate(&models.EventUser{})

	e := echo.New()
	e.Debug = true
	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Authentication strategies
	key := os.Getenv("GOTH_SESSION_SECRET")
	maxAge := 86400 * 90 // 90 days

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true // HttpOnly should always be enabled
	store.Options.Secure = false

	gothic.Store = store
	goth.UseProviders(
		instagram.New(
			os.Getenv("INSTAGRAM_CLIENT_ID"),
			os.Getenv("INSTAGRAM_CLIENT_SECRET"),
			"http://localhost:8080/auth/instagram/callback?provider=instagram"),
	)

	// Routes

	// Auth
	e.GET("/auth/:provider", controllers.AuthInstagram)
	e.GET("/auth/:provider/callback", controllers.AuthInstagramCB)
	e.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))

	// User
	u := e.Group("/api")
	u.Use(middlewares.LoadUserIntoContext)
	u.GET("/users", controllers.ListUsers)
	u.GET("/users/:uuid", controllers.GetUser)
	u.PUT("/users/:uuid", controllers.UpdateUser)
	u.DELETE("/users/:uuid", controllers.DeleteUser)
	u.GET("/users/self/events", controllers.GetSelfEventList)
	u.GET("/users/:uuid/events", controllers.GetUsersEventList)

	// Event
	event := u.Group("/events")
	event.Use(middlewares.GetEventIDFromUUID)
	event.GET("", controllers.ListEvents)
	event.POST("", controllers.CreateEvent)
	event.GET("/:uuid", controllers.GetEvent)
	event.PUT("/:uuid", controllers.UpdateEvent)
	event.POST("/:uuid/join", controllers.JoinEvent)
	event.DELETE("/:uuid", controllers.DeleteEvent)

	// Start server
	e.Logger.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}
