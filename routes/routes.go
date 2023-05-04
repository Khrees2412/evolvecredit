package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khrees2412/evolvecredit/controllers"
)

func RegisterRoutes(router *fiber.App) {
	controllers.NewAuthController().RegisterRoutes(router)
}
