package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khrees2412/evolvecredit/services"
	"github.com/khrees2412/evolvecredit/utils"
)

type ISavingsController interface {
	GetSavings(ctx *fiber.Ctx) error
	SaveFunds(ctx *fiber.Ctx) error
	RegisterRoutes(app *fiber.App)
}

type savingsController struct {
	accountService services.IAccountService
}

func NewSavingsController() ISavingsController {
	return &savingsController{
		accountService: services.NewAccountService(),
	}
}

func (ctl *savingsController) RegisterRoutes(app *fiber.App) {
	savings := app.Group("/v1/savings")
	savings.Get("/", utils.SecureAuth(), ctl.GetSavings)
	savings.Post("/", utils.SecureAuth(), ctl.SaveFunds)
}

func (ctl *savingsController) GetSavings(ctx *fiber.Ctx) error {
	return nil
}

func (ctl *savingsController) SaveFunds(ctx *fiber.Ctx) error {
	return nil
}
