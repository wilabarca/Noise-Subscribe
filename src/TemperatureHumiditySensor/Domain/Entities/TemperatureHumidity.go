package entities

// TemperatureHumidity es la entidad que representa los datos de temperatura y humedad
type TemperatureHumidity struct {
	Temperature float64 `json:"temperature"` 
	Humidity    float64 `json:"humidity"`    
	Timestamp   string  `json:"timestamp"`   
}
