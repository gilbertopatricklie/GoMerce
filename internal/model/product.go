package model

import (
	"time"
)

type Product struct {
	ID          int64     `json:"id"`
    Name        string    `json:"name"`
    Description *string   `json:"description,omitempty"`
    Price       int64     `json:"price"`
    Stock       int       `json:"stock"`
    IsActive    bool      `json:"is_active"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}