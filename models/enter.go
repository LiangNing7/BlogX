package models

import (
	"time"
)

type Model struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type IDRequest struct {
	ID uint `json:"id" form:"id" uri:"id"`
}
