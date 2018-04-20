package main

import (
	"fmt"
	"github.com/eklemen/vendr/controllers"
	"github.com/eklemen/vendr/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/instagram"
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

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Authentication strategies
	gomniauth.SetSecurityKey("REPLACE_ME")
	gomniauth.WithProviders(
		instagram.New(
			os.Getenv("INSTAGRAM_CLIENT_ID"),
			"INSTAGRAM_CLIENT_SECRET",
			"callback",
		),
	)

	// Routes
	e.GET("/users", controllers.GetAllUsers)
	e.POST("/users", controllers.CreateUser)
	e.GET("/users/:id", controllers.GetUser)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	// Start server
	e.Logger.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}
