package adapters

import (
	application "Noisesubscribe/src/TemperatureHumiditySensor/Application"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

// APIAdapter maneja la comunicación con la API para enviar datos de temperatura y humedad
type APIAdapter struct {
	apiURL string
}

// NewAPIAdapter crea una nueva instancia del adaptador para la API
func NewAPIAdapter(apiURL string) *APIAdapter {
	return &APIAdapter{
		apiURL: apiURL,
	}
}

// SendToAPI envía los datos de temperatura y humedad a la API
func (adapter *APIAdapter) SendToAPI(sensorData application.TemperatureHumidityService) error {
	data, err := json.Marshal(sensorData)
	if err != nil {
		log.Println("❌ Error al convertir los datos a JSON:", err)
		return err
	}

	resp, err := http.Post(adapter.apiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("❌ Error al enviar los datos a la API:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("❌ Error en la respuesta de la API: %d %s", resp.StatusCode, resp.Status)
		return err
	}

	log.Println("✅ Datos de temperatura y humedad enviados correctamente a la API.")
	return nil
}
