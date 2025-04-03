package repositories

import "Noisesubscribe/src/Light/Domain/Entities"

// LightRepository define las operaciones que el repositorio de Light debe implementar.
type LightRepository interface {
	// Procesa y reenvía los datos de la luz a otro servicio o base de datos
	ProcessAndForward(lightData entities.Light) error
	
	// Obtiene los datos del sensor de luz
	GetLightData() (entities.Light, error)  // Nuevo método para obtener los datos de luz
}
