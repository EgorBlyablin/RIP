package ds

type TerrainType struct {
	Translation string
}

type SelectedTurbine struct {
	Turbine              Turbine
	AvgVelocity          float32
	Terrain              string
	CalculatedGeneration float32
}

type Request struct {
	ID                   uint
	Period               string
	SelectedTurbines     []SelectedTurbine
	CalculatedGeneration float32
}
