package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/khrees2412/evolvecredit/services"
	"github.com/khrees2412/evolvecredit/types"
	"github.com/khrees2412/evolvecredit/utils"
)

type ITransactionController interface {
	GetTransaction(ctx *fiber.Ctx) error
	GetAllTransactions(ctx *fiber.Ctx) error
	RegisterRoutes(app *fiber.App)
}

type transactionController struct {
	transactionService services.ITransactionService
}

func NewTransactionController() ITransactionController {
	return &transactionController{
		transactionService: services.NewTransactionService(),
	}
}

func (ctl *transactionController) RegisterRoutes(app *fiber.App) {
	transactions := app.Group("/v1/transactions")
	transactions.Get("/", utils.SecureAuth(), ctl.GetAllTransactions)
	transactions.Get("/:id", utils.SecureAuth(), ctl.GetTransaction)
}

func (ctl *transactionController) GetTransaction(ctx *fiber.Ctx) error {
	userId, err := utils.UserFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	transactionId := ctx.Params("id")
	transaction, err := ctl.transactionService.GetTransaction(userId, transactionId)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	res := types.GenericResponse{
		Success: true,
		Message: "Successfully retrieved transaction",
		Data:    transaction,
	}
	return ctx.Status(fiber.StatusOK).JSON(res)

}

func (ctl *transactionController) GetAllTransactions(ctx *fiber.Ctx) error {
	userId, err := utils.UserFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	page := ctx.QueryInt("page")
	pageSize := ctx.QueryInt("page_size")

	status := ctx.Query("status")
	entry := ctx.Query("entry")

	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}
	pagination := types.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
	transactions, err := ctl.transactionService.GetAllTransactions(userId, types.TransactionEntry(entry), types.TransactionStatus(status), pagination)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(types.GenericResponse{
			Success: false,
			Message: err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(types.GenericResponse{
		Success: true,
		Message: "Successfully retrieved transactions",
		Data:    transactions,
	})
}
