package model

import (
	"time"
)

type CreateCategoryRequest struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *string `json:"parentId"`
	Level    int     `json:"level" binding:"gte=0"`
	Sort     int     `json:"sort" binding:"required"`
}

type UpdateCategoryRequest struct {
	Name     *string `json:"name,omitempty"`
	ParentID *string `json:"parentId,omitempty"`
	Level    *int    `json:"level,omitempty"`
	Sort     *int    `json:"sort,omitempty"`
}

type Category struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	Name      string     `json:"name"`
	ParentID  *string    `json:"parentId" gorm:"index"`
	Level     int        `json:"level"`
	Sort      int        `json:"sort"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt,omitempty" gorm:"index"`
}
