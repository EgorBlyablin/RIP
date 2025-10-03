package ds

type TurbineGenerationRequest struct {
	TurbineID uint `gorm:"not null;uniqueIndex:idx_turbine_request"`
	Turbine   Turbine

	GenerationRequestID uint `gorm:"not null;uniqueIndex:idx_turbine_request"`
	GenerationRequest   GenerationRequest

	AvgVelocity          *float32 `gorm:"precision:2;scale:1;check: avg_velocity > 0"`
	Alpha                *float32 `gorm:"precision:1;scale:2;check: alpha > 0"`
	CalculatedGeneration *uint64
}
