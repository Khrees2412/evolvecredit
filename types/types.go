package types

import "github.com/golang-jwt/jwt"

type (
	// Claims represent the structure of the JWT token
	Claims struct {
		Email  string `json:"email"`
		UserId string `json:"user_id"`
		jwt.StandardClaims
	}
	GenericResponse struct {
		Success bool        `json:"success"`
		Message interface{} `json:"message"`
		Data    interface{} `json:"data,omitempty"`
	}
	RegisterRequest struct {
		FirstName string `json:"first_name" validate:"required"`
		LastName  string `json:"last_name" validate:"required"`
		Email     string `json:"email" validate:"required"`
		Password  string `json:"password" validate:"required"`
	}
	RegisterResponse struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}
	LoginRequest struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	LoginResponse struct {
		Token     *TokenResponse `json:"token"`
		FirstName string         `json:"first_name"`
		LastName  string         `json:"last_name"`
		Email     string         `json:"-"`
	}
	TokenResponse struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
		Issuer      string `json:"issuer"`
	}
	AccountResponse struct {
		AccountNumber   string `json:"account_number"`
		CurrentBalance  int64  `json:"current_balance"`
		PreviousBalance int64  `json:"previous_balance"`
	}
	SavingsResponse struct {
		AccountNumber   string `json:"account_number"`
		PreviousBalance int64  `json:"previous_balance"`
		CurrentBalance  int64  `json:"current_balance"`
		LockedAmount    int64  `json:"locked_amount"`
	}
	WithdrawalRequest struct {
		Amount        int64  `json:"amount" validate:"required"`
		Reason        string `json:"reason"`
		AccountNumber string `json:"account_number" validate:"required"`
	}
	DepositRequest struct {
		Amount        int64  `json:"amount" validate:"required"`
		Reason        string `json:"reason"`
		AccountNumber string `json:"account_number" validate:"required"`
	}
	SavingsRequest struct {
		AccountNumber string `json:"account_number"`
		Amount        int64  `json:"amount"`
		Duration      int    `json:"duration"`
	}
	GetSavingsResponse struct {
		CurrentBalance int64 `json:"current_balance"`
		LockedFunds    int64 `json:"locked_funds"`
	}
	Pagination struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}
)

type Currency string
type TransactionType string
type TransactionEntry string
type TransactionStatus string

const (
	NGN Currency = "NGN"

	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"
	Savings    TransactionType = "savings"

	Credit TransactionEntry = "credit"
	Debit  TransactionEntry = "debit"

	Success TransactionStatus = "successful"
	Failed  TransactionStatus = "failed"
	Pending TransactionStatus = "pending"
)
