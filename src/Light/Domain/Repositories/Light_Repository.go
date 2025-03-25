package repositories

import "Noisesubscribe/src/Light/Domain/Entities"

// LightRepository define las operaciones que el repositorio de Light debe implementar.
type LightRepository interface {
    // Procesa y reenvía los datos de la luz a otro servicio o base de datos
    ProcessAndForward(lightData entities.Light) error
}
