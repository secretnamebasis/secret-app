package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/site/views"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", views.Home)

	app.Get("/entry", views.EntryPage)
	app.Get("/entries", views.EntriesPage)

	// journal
	app.Post("/submit/entry", views.SubmitEntryForm)

	// shop
	app.Post("/submit/order", views.SubmitOrderForm)
	app.Get("/order", views.OrderPage)
	app.Get("/pay", views.PayPage)
	app.Get("/success", views.SuccessPage)

	app.Use(views.NotFound)
}
