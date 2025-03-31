package entities


type Light struct {
	ID         int    `json:"id"`
	SensorID   string `json:"sensor_id"`
	Intensidad int    `json:"intensidad"`
	Color      string `json:"color"`
	Estado     bool   `json:"estado"`
	Timestamp  string `json:"timestamp"`
}