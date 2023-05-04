package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/google/uuid"
	_ "github.com/khrees2412/evolvecredit/config"
	"github.com/khrees2412/evolvecredit/database"
	"github.com/khrees2412/evolvecredit/routes"
	"log"
	"net/http"
	"os"
)

func main() {

	app := fiber.New()
	app.Use(requestid.New(requestid.Config{
		Header: "Evolve-Request-ID",
		Generator: func() string {
			return uuid.NewString()
		},
	}))
	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("You're home, yaay!!")
	})

	port := os.Getenv("PORT")

	if port == "" {
		port = "8020"
	}
	dbUrl := os.Getenv("DATABASE_URL")
	dbConnection, err := database.ConnectDB(dbUrl)

	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}

	err = database.MigrateAll(dbConnection)

	if err != nil {
		log.Fatalf("migration error: %v", err)
	}
	routes.RegisterRoutes(app)

	if err := app.Listen(":" + port); err != nil && err != http.ErrServerClosed {
		log.Fatal(fmt.Sprintf("listen: %s\n", err))
	}
}

func setupSystemRouteHandler(app *fiber.App) {
	// 404 Handler
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Sorry can't find that!")
	})
	// 405 Handler
	//router.NoMethod(handlers.Http405Handler())
}
