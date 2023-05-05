package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khrees2412/evolvecredit/services"
	"github.com/khrees2412/evolvecredit/types"
	"github.com/khrees2412/evolvecredit/utils"
)

type IAccountController interface {
	GetAccount(ctx *fiber.Ctx) error
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
	accounts.Get("/balance", utils.SecureAuth(), ctl.GetAccount)

}

func (ctl *accountController) GetAccount(ctx *fiber.Ctx) error {
	userId, err := utils.UserFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	account, err := ctl.accountService.GetAccount(userId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.GenericResponse{
		Success: true,
		Message: "Successfully retrieved account",
		Data:    account,
	})

}
