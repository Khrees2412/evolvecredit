package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khrees2412/evolvecredit/services"
	"github.com/khrees2412/evolvecredit/utils"
)

type IAccountController interface {
	GetBalance(ctx *fiber.Ctx) error
	GetTotalSaved(ctx *fiber.Ctx) error
	CreateSavings(ctx *fiber.Ctx) error
	RegisterRoutes(app *fiber.App)
}

type accountController struct {
	accountService services.IAccountService
}

func NewAccountController() IAccountController {
	return &accountController{
		accountService: services.NewAccountService(),
	}
}

func (ctl *accountController) RegisterRoutes(app *fiber.App) {
	accounts := app.Group("/v1/accounts")
	accounts.Get("/balance", utils.SecureAuth(), ctl.GetBalance)
	accounts.Get("/savings", utils.SecureAuth(), ctl.GetTotalSaved)
	accounts.Post("/savings", utils.SecureAuth(), ctl.CreateSavings)

}

func (ctl *accountController) CreateSavings(ctx *fiber.Ctx) error {
	return nil
}

func (ctl *accountController) GetBalance(ctx *fiber.Ctx) error {
	return nil
}

func (ctl *accountController) GetTotalSaved(ctx *fiber.Ctx) error {
	return nil
}

