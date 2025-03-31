package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	entities "Noisesubscribe/src/AirQuality/Domain/Entities"
)

type AirQualityRepositoryAdapter struct {
	apiURL string
}

func NewAirQualityRepositoryAdapter(apiURL string) *AirQualityRepositoryAdapter {
	// Aseguramos que la URL predeterminada se utiliza si apiURL está vacío
	if apiURL == "" {
		apiURL = "http://localhost:8080/airqualitysensor/"
	}
	return &AirQualityRepositoryAdapter{
		apiURL: apiURL,
	}
}

func (adapter *AirQualityRepositoryAdapter) ProcessAndForward(airData entities.AirQualitySensor) error {
	// No necesitamos el apiURL como parámetro, ya que lo estamos tomando de la propiedad del adaptador.
	data, err := json.Marshal(airData)
	if err != nil {
		log.Printf("Error al serializar datos: %v | Datos: %+v\n", err, airData)
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
