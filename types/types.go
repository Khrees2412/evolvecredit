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
		FirstName   string `json:"first_name" validate:"required"`
		LastName    string `json:"last_name" validate:"required"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email" validate:"required"`
		Password    string `json:"password" validate:"required"`
		Code        string `json:"code" validate:"required"`
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
)

type Currency string
type TransactionType string
type TransactionEntry string
type TransactionStatus string

const (
	NGN Currency = "NGN"

	Deposit    TransactionType = "deposit"
	Withdrawal TransactionType = "withdrawal"

	Credit TransactionEntry = "credit"
	Debit  TransactionEntry = "debit"

	Success TransactionStatus = "success"
	Failed  TransactionStatus = "failed"
	Pending TransactionStatus = "pending"
)
