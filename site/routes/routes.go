package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/site/controllers"
	"github.com/secretnamebasis/secret-app/site/views"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", views.Home)
	app.Get("/api/info", controllers.APIInfo)
	app.Get("/api/items", controllers.AllItems)
	app.Post("/api/items", controllers.CreateItem)
	app.Get("/api/items/:id", controllers.ItemByID)
	app.Put("/api/items/:id", controllers.UpdateItem)
	app.Delete("/api/items/:id", controllers.DeleteItem)
	app.Use(views.NotFound)
}
