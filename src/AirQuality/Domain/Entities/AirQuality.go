package entities

// AirQuality representa los datos de calidad del aire de un sensor
type AirQuality struct {
    SensorID    string  `json:"sensor_id"`
    CO2PPM      int     `json:"co2_ppm"`
    PM25        int     `json:"pm25"`
    Temperatura float64 `json:"temperatura"`
    Timestamp   string  `json:"timestamp"`
}