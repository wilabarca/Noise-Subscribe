package adapters

import (
	entities "Noisesubscribe/src/SoundSensor/Domain/Entities"
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

// SoundSensorRepositoryAdapter es un adaptador que maneja la interacción con una API externa para el sensor de sonido
type SoundSensorRepositoryAdapter struct {
	apiURL string
}

// NewSoundSensorRepositoryAdapter crea una nueva instancia del adaptador de API para el sensor de sonido
func NewSoundSensorRepositoryAdapter(apiURL string) *SoundSensorRepositoryAdapter {
	if apiURL == "" {
		apiURL = "http://localhost:8080/soundsensor/"
	}
	return &SoundSensorRepositoryAdapter{
		apiURL: apiURL,
	}
}

// ProcessAndForward procesa los datos del sensor de sonido y los envía a la API externa
func (adapter *SoundSensorRepositoryAdapter) ProcessAndForward(sensorData entities.SoundSensor) error {
	// No necesitamos el apiURL como parámetro, ya que lo estamos tomando de la propiedad del adaptador.
	data, err := json.Marshal(sensorData)
	if err != nil {
		log.Printf("Error al serializar datos: %v | Datos: %+v\n", err, sensorData)
		return err
	}

	resp, err := http.Post(adapter.apiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error en POST a %s: %v\n", adapter.apiURL, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Respuesta no exitosa: %d | URL: %s\n", resp.StatusCode, adapter.apiURL)
		return errors.New("código de estado: " + resp.Status)
	}

	log.Printf("Datos enviados correctamente a %s\n", adapter.apiURL)
	return nil
}
