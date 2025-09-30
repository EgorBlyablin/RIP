package ds

type SelectedTurbine struct {
	Turbine              Turbine
	AvgVelocity          float32
	Alpha                float32
	CalculatedGeneration float32
}

type GenerationCalculationRequest struct {
	ID                   uint
	Period               string
	SelectedTurbines     []SelectedTurbine
	CalculatedGeneration float32
}
