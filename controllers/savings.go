package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/khrees2412/evolvecredit/services"
	"github.com/khrees2412/evolvecredit/types"
	"github.com/khrees2412/evolvecredit/utils"
)

type ISavingsController interface {
	GetSavings(ctx *fiber.Ctx) error
	SaveFunds(ctx *fiber.Ctx) error
	RegisterRoutes(app *fiber.App)
}

type savingsController struct {
	accountService services.IAccountService
	savingsService services.ISavingsService
}

func NewSavingsController() ISavingsController {
	return &savingsController{
		accountService: services.NewAccountService(),
		savingsService: services.NewSavingsService(),
	}
}

func (ctl *savingsController) RegisterRoutes(app *fiber.App) {
	savings := app.Group("/v1/savings")
	savings.Get("/", utils.SecureAuth(), ctl.GetSavings)
	savings.Post("/", utils.SecureAuth(), ctl.SaveFunds)
}

func (ctl *savingsController) GetSavings(ctx *fiber.Ctx) error {
	_, err := utils.UserFromContext(ctx)
	if err != nil {
		return err
	}

	accountNumber := ctx.Params("account_number")

	savings, err := ctl.savingsService.GetSavings(accountNumber)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.GenericResponse{
		Success: true,
		Message: "Successfully retrieved Savings",
		Data: types.GetSavingsResponse{
			CurrentBalance: savings.CurrentBalance,
			LockedFunds:    savings.LockedAmount,
		},
	})
}

func (ctl *savingsController) SaveFunds(ctx *fiber.Ctx) error {
	userId, err := utils.UserFromContext(ctx)
	if err != nil {
		return err
	}

	var body *types.SavingsRequest
	if err = ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.GenericResponse{
			Success: false,
			Message: "Problem while parsing request body",
			Data:    err.Error(),
		})
	}

	errors := utils.ValidateStruct(body)
	if errors != nil {
		return ctx.JSON(errors)
	}

	err = ctl.savingsService.SaveFunds(userId, body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.GenericResponse{
		Success: true,
		Message: fmt.Sprintf("You have successfully locked  %d for savings, it will be unlocked in %d days", body.Amount, body.Duration),
	})
}
