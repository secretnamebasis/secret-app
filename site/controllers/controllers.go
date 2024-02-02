package controllers

import (
	"github.com/gofiber/fiber/v2"
	"go.etcd.io/bbolt"
)

var db *bbolt.DB
var bucket = []byte("order")

type OrderData struct {
	Name    string
	Address string
}

func SetDB(database *bbolt.DB) {
	db = database
}

func APIInfo(c *fiber.Ctx) error {
	response := fiber.Map{
		"message": "Welcome to secret-swap API",
		"data":    "pong",
		"status":  "success",
	}
	return c.JSON(response)
}

func Home(c *fiber.Ctx) error {
	// Your existing home controller logic
	// ...
	return nil
}
