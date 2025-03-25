package adapters

import (
	application "Noisesubscribe/src/SoundSensor/Application"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type API2Adapter struct {
	apiURL string
}

// NewAPI2Adapter crea una nueva instancia del adaptador para enviar datos a la API 2
func NewAPI2Adapter(apiURL string) *API2Adapter {
	return &API2Adapter{
		apiURL: apiURL,
	}
}

// SendToAPI2 envía los datos del sensor a la API 2
func (adapter *API2Adapter) SendToAPI2(sensorData application.SoundSensorService) error {
	data, err := json.Marshal(sensorData)
	if err != nil {
		log.Println("Error al convertir los datos a JSON:", err)
		return err
	}

	resp, err := http.Post(adapter.apiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error al enviar los datos a la API 2:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error en la respuesta de la API 2: %d %s", resp.StatusCode, resp.Status)
		return err
	}

	log.Println("Datos enviados correctamente a la API 2.")
	return nil
}
