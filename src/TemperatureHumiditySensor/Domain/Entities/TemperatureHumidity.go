package entities

// TemperatureHumidity es la entidad que representa los datos de temperatura y humedad
type TemperatureHumidity struct {
	Temperature float64 `json:"temperature"` // Temperatura en grados Celsius
	Humidity    float64 `json:"humidity"`    // Humedad en porcentaje
	Timestamp   string  `json:"timestamp"`   // Marca de tiempo del registro
}
