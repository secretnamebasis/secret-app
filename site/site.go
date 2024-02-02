// site/site.go

package site

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/site/middleware"
	"github.com/secretnamebasis/secret-app/site/routes"
)

func setupMiddleware(app *fiber.App) {
	app.Use(middleware.LogRequests)
}

func MakeWebsite() *fiber.App {
	app := fiber.New()

	setupMiddleware(app)

	routes.SetupRoutes(app)
	return app
}

func StartServer(app *fiber.App, port int) error {
	return app.Listen(fmt.Sprintf(":%d", port))
}
