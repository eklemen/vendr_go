package main

import (
	"fmt"
	"github.com/eklemen/vendr/controllers"
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
	controllers.DB = db
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
	maxAge := 86400 * 30 // 30 days

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

	// User
	r := e.Group("/api")
	r.Use(middleware.JWT([]byte(os.Getenv("JWT_SECRET"))))
	r.Use(SetUserId)
	r.GET("/users", controllers.ListUsers)
	r.GET("/users/:uuid", controllers.GetUser)
	r.PUT("/users/:uuid", controllers.UpdateUser)
	r.DELETE("/users/:uuid", controllers.DeleteUser)

	// Event
	r.GET("/events", controllers.ListEvents)
	r.POST("/events", controllers.CreateEvent)
	r.GET("/events/:uuid", controllers.GetEvent)
	r.DELETE("/events/:uuid", controllers.DeleteEvent)

	// Start server
	e.Logger.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}
