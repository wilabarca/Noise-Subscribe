package entities

type SoundSensor struct {
	ID        string  `json:"id"`
	RuidoDB   float64 `json:"ruido_dB"`
	Timestamp string  `json:"timestamp"`
}