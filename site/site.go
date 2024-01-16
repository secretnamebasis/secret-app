// site/site.go

package site

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/site/config"
	"github.com/secretnamebasis/secret-app/site/db"
	"github.com/secretnamebasis/secret-app/site/middleware"
	"github.com/secretnamebasis/secret-app/site/routes"
)

func setupMiddleware(app *fiber.App) {
	app.Use(middleware.LogRequests)
}

func MakeWebsite(config config.Config) *fiber.App {
	app := fiber.New()

	setupMiddleware(app)

	// Initialize the database
	if err := db.InitDB(); err != nil {
		log.Fatal(err)
	}

	routes.SetupRoutes(app)
	return app
}

func StartServer(app *fiber.App, port int) error {
	return app.Listen(fmt.Sprintf(":%d", port))
}