package ds

import (
	"time"
)

type GenerationRequest struct {
	ID     uint
	Status string `gorm:"type:varchar(10);check:status IN ('draft','sent','completed','rejected','deleted');default:'draft';not null"`

	CreatedByID uint `gorm:"not null"`
	CreatedBy   User
	CreatedAt   time.Time `gorm:"not null"`
	SentAt      *time.Time

	ClosedByID *uint
	ClosedBy   *User
	ClosedAt   *time.Time
	DeletedAt  *time.Time

	PeriodDays *uint16 `gorm:"check:period_days > 0"`

	TurbineGenerationRequests []TurbineGenerationRequest
}
