package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/site/models"
	"go.etcd.io/bbolt"
)

var db *bbolt.DB
var bucket = []byte("items")

func APIInfo(c *fiber.Ctx) error {
	response := fiber.Map{
		"message": "Welcome to secret-swap API",
		"data":    "This is a sample API endpoint",
		"status":  "success",
	}
	return c.JSON(response)
}

func AllItems(c *fiber.Ctx) error {
	var items []models.Item
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var item models.Item
			if err := json.Unmarshal(v, &item); err != nil {
				return err
			}
			items = append(items, item)
			return nil
		})
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error retrieving items", "status": "error"})
	}

	return c.JSON(fiber.Map{"data": items, "status": "success"})
}

func CreateItem(c *fiber.Ctx) error {
	var newItem models.Item
	if err := c.BodyParser(&newItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body", "status": "error"})
	}

	newItem.ID, _ = NextID()
	newItem.CreatedAt = time.Now()

	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}

		itemJSON, err := json.Marshal(newItem)
		if err != nil {
			return err
		}

		return b.Put([]byte(strconv.Itoa(newItem.ID)), itemJSON)
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error creating item", "status": "error"})
	}

	return c.JSON(fiber.Map{"data": newItem, "status": "success"})
}

func ItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item models.Item

	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}

		itemJSON := b.Get([]byte(id))
		if itemJSON == nil {
			return fmt.Errorf("Item with ID %s not found", id)
		}

		return json.Unmarshal(itemJSON, &item)
	})

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error(), "status": "error"})
	}

	return c.JSON(fiber.Map{"data": item, "status": "success"})
}

func UpdateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedItem models.Item
	if err := c.BodyParser(&updatedItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body", "status": "error"})
	}

	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}

		itemJSON := b.Get([]byte(id))
		if itemJSON == nil {
			return fmt.Errorf("Item with ID %s not found", id)
		}

		var existingItem models.Item
		if err := json.Unmarshal(itemJSON, &existingItem); err != nil {
			return err
		}

		// Preserve the creation timestamp
		updatedItem.CreatedAt = existingItem.CreatedAt

		itemJSON, err := json.Marshal(updatedItem)
		if err != nil {
			return err
		}

		return b.Put([]byte(id), itemJSON)
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error updating item", "status": "error"})
	}

	return c.JSON(fiber.Map{"data": updatedItem, "status": "success"})
}

func DeleteItem(c *fiber.Ctx) error {
	id := c.Params("id")

	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}

		itemJSON := b.Get([]byte(id))
		if itemJSON == nil {
			return fmt.Errorf("Item with ID %s not found", id)
		}

		return b.Delete([]byte(id))
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Error deleting item", "status": "error"})
	}

	return c.JSON(fiber.Map{"message": "Item deleted successfully", "status": "success"})
}

func NextID() (int, error) {
	var id int
	err := db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return fmt.Errorf("Bucket %q not found!", bucket)
		}

		// Get the current sequence number
		seq, err := b.NextSequence()
		if err != nil {
			return err
		}

		id = int(seq)
		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}
