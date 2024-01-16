package code

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/functions/wallet/dero"
	"go.etcd.io/bbolt"
)

// Item represents a sample data structure for demonstration
type Item struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// Config struct to hold configuration parameters
type Config struct {
	Port int
}

var (
	db     *bbolt.DB
	bucket = []byte("items")
)

func initDB() error {
	var err error
	db, err = bbolt.Open("items.db", 0600, nil)
	if err != nil {
		return err
	}

	// Create a bucket if it doesn't exist
	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucket)
		return err
	})

	return err
}

func setupRoutes(app *fiber.App) {
	app.Get("/", handleHome)
	app.Get("/api/info", handleAPIInfo)
	app.Get("/api/items", getAllItems)
	app.Post("/api/items", createItem)
	app.Get("/api/items/:id", getItemByID)
	app.Put("/api/items/:id", updateItem)
	app.Delete("/api/items/:id", deleteItem)
	app.Use(handleNotFound)
}

func setupMiddleware(app *fiber.App) {
	app.Use(logRequests)
}

func logRequests(c *fiber.Ctx) error {
	log.Printf("Request: %s %s", c.Method(), c.OriginalURL())
	return c.Next()
}

func handleHome(c *fiber.Ctx) error {
	message := fmt.Sprintf("Welcome to secret-swap!\n%s", dero.CreateServiceAddress(dero.Address()))
	return c.SendString(message)
}

func handleAPIInfo(c *fiber.Ctx) error {
	response := fiber.Map{
		"message": "Welcome to secret-swap API",
		"data":    "This is a sample API endpoint",
		"status":  "success",
	}
	return c.JSON(response)
}

func getAllItems(c *fiber.Ctx) error {
	var items []Item
	err := db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return nil
		}
		return b.ForEach(func(k, v []byte) error {
			var item Item
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

func createItem(c *fiber.Ctx) error {
	var newItem Item
	if err := c.BodyParser(&newItem); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid request body", "status": "error"})
	}

	newItem.ID, _ = getNextID()
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

func getItemByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var item Item

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

func updateItem(c *fiber.Ctx) error {
	id := c.Params("id")
	var updatedItem Item
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

		var existingItem Item
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

func deleteItem(c *fiber.Ctx) error {
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

func getNextID() (int, error) {
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

func handleNotFound(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotFound).SendString("404 Not Found")
}

func makeWebsite(config Config) *fiber.App {
	app := fiber.New()
	setupMiddleware(app)

	// Initialize the database
	if err := initDB(); err != nil {
		log.Fatal(err)
	}

	setupRoutes(app)
	return app
}

func startServer(app *fiber.App, port int) error {
	return app.Listen(fmt.Sprintf(":%d", port))
}
