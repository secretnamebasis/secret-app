package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/site/middleware"
	"github.com/secretnamebasis/secret-app/site/views"
)

func SetupRoutes(app *fiber.App) {
	app.Static("/", "./site/assets")
	app.Get("/", views.Home)
	app.Get("/about", views.About)
	app.Get("/entry", views.EntryPage)
	app.Get("/entries", views.EntriesPage)

	// Apply AuthReq middleware for /entry and /entries routes
	authRoutes := app.Group("/", middleware.AuthReq())
	authRoutes.Post("/submit/entry", views.SubmitEntryForm)

	// shop
	authRoutes.Post("/submit/order", views.SubmitOrderForm)
	authRoutes.Get("/order", views.OrderPage)
	authRoutes.Get("/pay", views.PayPage)
	authRoutes.Get("/success", views.SuccessPage)

	app.Use(views.NotFound)
}
