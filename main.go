package main

import (
	"fmt"
	"github.com/eklemen/vendr/controllers"
	"github.com/eklemen/vendr/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/markbates/goth"
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

	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Authentication strategies
	goth.UseProviders(
		instagram.New(
			"0c6f9953a60e4ceaa28ef77f53f8a1d1",
			"8aa0026c837245189f8f94ee9a05a08e",
			"http://localhost:8080/auth/instagram/callback"),
	)

	// Routes

	// Auth
	e.GET("/auth/{provider}", controllers.AuthInstagram)
	e.GET("/auth/{provider}/callback", controllers.AuthInstagramCB)

	// User
	e.GET("/users", controllers.GetAllUsers)
	e.POST("/users", controllers.CreateUser)
	e.GET("/users/:id", controllers.GetUser)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	// Start server
	e.Logger.Fatal(e.Start(os.Getenv("SERVER_PORT")))
}

//func wrapHandler(h http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		h.ServeHTTP(w, r)
//	})
//}
