
package entities

// Light representa los datos del sensor de luz
type Light struct {
	ID         int     `json:"id"`
	SensorID   string  `json:"sensor_id"`
	Nivel      float64 `json:"nivel"`      // Cambié Intensidad a Nivel (float para representar valores continuos)
	Timestamp  string  `json:"timestamp"`  // Momento en que se registró la medición
}
