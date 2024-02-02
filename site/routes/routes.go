package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/site/views"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", views.Home)
	app.Post("/submit", views.SubmitForm)
	app.Get("/pay", views.PayPage)
	app.Get("/success", views.SuccessPage)

	app.Use(views.NotFound)
}
