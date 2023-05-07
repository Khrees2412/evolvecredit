package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestAccountController_GetAccount(t *testing.T) {
	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "This gets an account",
			route:        "/accounts",
			expectedCode: 200,
		},
	}

	// Define Fiber app.
	app := fiber.New()

	// Create route with GET method for test
	app.Get("/hello", func(c *fiber.Ctx) error {
		// Return simple string as response
		return c.SendString("Hello, World!")
	})

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest("GET", "/v1"+test.route, nil)

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, 1)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
