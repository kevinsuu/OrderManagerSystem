package model

import "time"

// Address 地址模型
type Address struct {
	ID         string    `json:"id" db:"id"`
	UserID     string    `json:"user_id" db:"user_id"`
	Name       string    `json:"name" db:"name"`
	Phone      string    `json:"phone" db:"phone"`
	Street     string    `json:"street" db:"street"`
	City       string    `json:"city" db:"city"`
	District   string    `json:"district" db:"district"`
	PostalCode string    `json:"postal_code" db:"postal_code"`
	IsDefault  bool      `json:"is_default" db:"is_default"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// AddressRequest 地址請求
type AddressRequest struct {
	Name       string `json:"name" binding:"required"`
	Phone      string `json:"phone" binding:"required"`
	Street     string `json:"street" binding:"required"`
	City       string `json:"city" binding:"required"`
	District   string `json:"district" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
	IsDefault  bool   `json:"is_default"`
}
