package ds

import "fmt"

type Turbine struct {
	ID          uint
	Title       string  `gorm:"type:varchar(100);not null"`
	Description string  `gorm:"type:text;not null"`
	Status      string  `gorm:"type:varchar(8);default:'active';check:status IN ('active', 'inactive');not null"`
	Image       *string `gorm:"type:varchar(100)"`

	Power  uint32 `gorm:"check:power > 0;not null"`
	Height uint16 `gorm:"check:height > 0;not null"`
}

func (turbine *Turbine) ImageSrc() string {
	if turbine.Image != nil {
		return fmt.Sprintf("http://localhost:9001/api/v1/buckets/images/objects/download?preview=True&prefix=%s", *turbine.Image)
	}
	return ""
}
