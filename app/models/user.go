package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Username  string    `json:"username" gorm:"unique;type:varchar(20)" validate:"required,min=6,max=20"`
	Password  string    `json:"password,omitempty" gorm:"type:varchar(255)" validate:"required,min=6"`
	FullName  string    `json:"full_name" gorm:"column:full_name;type:varchar(100)" validate:"required,min=6"`
}

func (u User) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

type UserSession struct {
	ID                  uint      `gorm:"primarykey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	UserID              uint      `json:"user_id" validate:"required"` // hilangkan unique di GORM
	Token               string    `json:"token" gorm:"unique;type:varchar(255)" validate:"required"`
	RefreshToken        string    `json:"refresh_token" gorm:"unique;type:varchar(255)" validate:"required"`
	TokenExpired        time.Time `json:"-" validate:"required"`
	RefreshTokenExpired time.Time `json:"-" validate:"required"`
}

func (s UserSession) Validate() error {
	v := validator.New()
	return v.Struct(s)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type LoginResponse struct {
	Username     string `json:"username"`
	FullName     string `json:"full_name"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
