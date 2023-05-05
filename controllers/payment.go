package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khrees2412/evolvecredit/services"
	"github.com/khrees2412/evolvecredit/types"
	"github.com/khrees2412/evolvecredit/utils"
)

type IPaymentController interface {
	Deposit(ctx *fiber.Ctx) error
	Withdraw(ctx *fiber.Ctx) error
	RegisterRoutes(app *fiber.App)
}

type paymentController struct {
	paymentService services.IPaymentService
}

func NewPaymentController() IPaymentController {
	return &paymentController{
		paymentService: services.NewPaymentService(),
	}
}

func (ctl *paymentController) RegisterRoutes(app *fiber.App) {
	payments := app.Group("/v1/payments")
	payments.Post("/deposit", utils.SecureAuth(), ctl.Deposit)
	payments.Post("/withdraw", utils.SecureAuth(), ctl.Withdraw)
}

func (ctl *paymentController) Deposit(ctx *fiber.Ctx) error {
	userId, err := utils.UserFromContext(ctx)
	if err != nil {
		return err
	}

	var body *types.DepositRequest
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

	err = ctl.paymentService.Deposit(userId, body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.GenericResponse{
		Success: true,
		Message: "You have successfully made a deposit",
	})
}

func (ctl *paymentController) Withdraw(ctx *fiber.Ctx) error {
	userId, err := utils.UserFromContext(ctx)
	if err != nil {
		return err
	}

	var body *types.WithdrawalRequest
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

	err = ctl.paymentService.Withdrawal(userId, body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(types.GenericResponse{
		Success: true,
		Message: "You have successfully withdrawn funds",
	})
}
