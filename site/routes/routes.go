package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/site/controllers"
	"github.com/secretnamebasis/secret-app/site/middleware"
	"github.com/secretnamebasis/secret-app/site/views"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", views.Home)
	app.Get("/items/:id", controllers.DisplayItemByID)

	api := app.Group("/api", middleware.AuthReq(), middleware.LogRequests)
	api.Get("/info", controllers.APIInfo)
	api.Get("/items", controllers.AllItems)
	api.Post("/items", controllers.CreateItem)
	api.Get("/items/:id", controllers.ItemByID)
	api.Put("/items/:id", controllers.UpdateItem)
	api.Delete("/items/:id", controllers.DeleteItem)

	app.Use(views.NotFound)
}
