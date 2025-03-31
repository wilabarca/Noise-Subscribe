package entities

type AirQualitySensor struct {
	ID         int     `json:"id"`
	SensorID   string  `json:"sensor_id"`
	CO2PPM     int     `json:"co2_ppm"`
	PM25       int     `json:"pm25"`
	Temperatura float64 `json:"temperatura"`
	Timestamp  string  `json:"timestamp"`
}
