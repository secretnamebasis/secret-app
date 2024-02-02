package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/secretnamebasis/secret-app/site/controllers"
	"github.com/stretchr/testify/assert"
)

func TestAPIInfo(t *testing.T) {
	t.Run("Given APIInfo handler", func(t *testing.T) {
		app := fiber.New()
		app.Get("/api/info", controllers.APIInfo)

		t.Run("When requesting APIInfo", func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
			resp, err := app.Test(req)

			t.Run("Then return APIInfo successfully", func(t *testing.T) {
				assert.NoError(t, err)
				assert.Equal(t, http.StatusOK, resp.StatusCode)

				var responseBody map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&responseBody)
				assert.NoError(t, err)

				assert.Equal(t, "Welcome to secret-swap API", responseBody["message"])
				assert.Equal(t, "pong", responseBody["data"])
				assert.Equal(t, "success", responseBody["status"])
			})
		})
	})
}
