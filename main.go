package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"go-image-upload/database"
	"go-image-upload/handlers"
	"go-image-upload/middleware"
	"go-image-upload/util"
)

func main() {
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	// Initialize database
	db := database.InitDatabase()

	// Generate AES key for JWT encryption
	key, err := util.GenerateRandomKey(32) // Generate a 32-byte key
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the app with dependencies
	myApp := handlers.NewApp(db, key)

	app.Use(middleware.LoggingMiddleware())

	// Routes
	app.Static("/", "./views")
	app.Static("/api/user/", "./views")
	app.Static("/api/image/", "./views")
	app.Static("/api/image/uploads/", "./uploads")

	api := app.Group("/api")
	user := api.Group("/user")
	image := api.Group("/image")

	app.Get("/", myApp.HandleIndex)
	user.Post("/login", middleware.LoginLimiter(), myApp.HandleLogin)
	user.Get("/dashboard", middleware.CheckJWT(myApp), middleware.CSRFMiddleware(), myApp.HandleDashboard)
	user.Get("/logout", myApp.HandleLogout)
	user.Post("/signup", myApp.HandleSignup)
	image.Post("/upload", middleware.CheckJWT(myApp), myApp.HandleUpload)
	image.Post("/delete-image/:id", middleware.CheckJWT(myApp), myApp.HandleDeleteImage)

	log.Fatal(app.Listen(":5000"))
}
