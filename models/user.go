package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	Base
	UserId    string `json:"user_id" gorm:"primaryKey;unique;not null;type:varchar(100)"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"-"`
}

// Base contains common columns for all tables
type Base struct {
	Counter   int       `json:"counter" gorm:"autoIncrement"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate will set User struct before every insert
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UserId = uuid.NewString()
	return nil
}
