package types

import "github.com/golang-jwt/jwt"

type GenericResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Claims represent the structure of the JWT token
type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.StandardClaims
}
