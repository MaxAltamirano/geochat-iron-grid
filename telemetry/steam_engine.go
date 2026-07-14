package telemetry

import (
	"math/rand"
	//"time"
)

// SteamEngineData representa la salida del motor a vapor/bio-digestor
type SteamEngineData struct {
	EficienciaTermica float64 `json:"eficiencia_termica"`
	PresionVapor      float64 `json:"presion_vapor"`
	EnergiaGenerada   float64 `json:"energia_generada"`
	Status            string  `json:"status"`
}

// SimularTelemetria genera datos en tiempo real de nuestra patente
func SimularTelemetria() SteamEngineData {
	return SteamEngineData{
		EficienciaTermica: 85.0 + rand.Float64()*10.0, // Alta eficiencia por el diseño
		PresionVapor:      2.5 + rand.Float64()*0.5,
		EnergiaGenerada:   1200.0 + rand.Float64()*200.0, // Watts generados por nodo
		Status:            "OPERATIVO (SÍNTESIS 432Hz)",
	}
}

