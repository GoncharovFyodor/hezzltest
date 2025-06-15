package domain

import (
	"github.com/GoncharovFyodor/hezzltest/internal/models"
	"time"
)

type Good struct {
	ID          int       `json:"id"`
	ProjectID   int       `json:"projectId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Priority    int       `json:"priority"`
	Removed     bool      `json:"removed"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CreateGoodRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateGoodRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type Meta struct {
	Total   int `json:"total"`
	Removed int `json:"removed"`
	Limit   int `json:"limit"`
	Offset  int `json:"offset"`
}

type GoodList struct {
	Meta  `json:"meta"`
	Goods []models.Good `json:"goods"`
}

type ReprioritizeRequest struct {
	NewPriority int `json:"newPriority" binding:"required"`
}
