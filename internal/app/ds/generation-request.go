package ds

import (
	"time"
)

type GenerationRequest struct {
	ID     uint   `json:"id"`
	Status string `gorm:"type:varchar(10);check:status IN ('draft','sent','completed','rejected','deleted');default:'draft';not null" json:"status"`

	CreatedByID uint       `gorm:"not null" json:"created_by_id"`
	CreatedBy   User       `json:"-"`
	CreatedAt   time.Time  `gorm:"not null" json:"created_at"`
	FormedAt    *time.Time `json:"formed_at"`

	ClosedByID *uint      `json:"closed_by_id"`
	ClosedBy   *User      `json:"-"`
	ClosedAt   *time.Time `json:"closed_at"`
	DeletedAt  *time.Time `json:"deleted_at"`

	PeriodDays *uint `gorm:"check:(period_days > 0)" json:"period_days"`

	TurbineGenerationRequests *[]TurbineGenerationRequest `json:"turbine_generation_requests,omitempty"`
}

type UpdateGenerationRequest struct {
	PeriodDays *uint `json:"period_days" binding:"omitnil,gt=0,lt=36500"`
}
