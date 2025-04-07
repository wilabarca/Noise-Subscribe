package entities

type SoundSensor struct {
	ID          int     `json:"id"`           
	SensorID    string  `json:"sensor_id"`    
	RuidoDB     float64 `json:"ruido_dB"`    
	Timestamp   string  `json:"timestamp"`   
}
