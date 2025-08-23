package main

import (
	"github.com/0xk4n3ki/OAuth-2.0-golang/config"
	"github.com/0xk4n3ki/OAuth-2.0-golang/controllers"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	config.GoogleConfig()

	app.Get("/google_login", controllers.GoogleLogin)
	app.Get("/google_callback", controllers.GoogleCallback)

	app.Listen(":9000")
}
