package adapters

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"Noisesubscribe/src/Light/Domain/Entities"
)

// APIAdapter es un adaptador que maneja la interacción con una API externa
type APIAdapter struct {
	apiURL string
}

// NewAPIAdapter crea una nueva instancia del adaptador de API
func NewAPIAdapter(apiURL string) *APIAdapter {
	return &APIAdapter{apiURL: apiURL}
}

// SendLightDataEnvio procesa los datos de luz y los envía a la API externa
func (adapter *APIAdapter) SendLightData(lightData entities.Light) error {
	// Convertir la entidad Light a JSON
	lightDataJSON, err := json.Marshal(lightData)
	if err != nil {
		log.Println("Error al convertir los datos de luz a JSON:", err)
		return err
	}

	// Realizar una solicitud HTTP POST a la API externa
	resp, err := http.Post(adapter.apiURL, "application/json", bytes.NewBuffer(lightDataJSON))
	if err != nil {
		log.Println("Error al enviar los datos a la API:", err)
		return err
	}
	defer resp.Body.Close()

	// Verificar la respuesta de la API
	if resp.StatusCode != http.StatusOK {
		log.Println("Error en la respuesta de la API, código de estado:", resp.StatusCode)
		return err
	}

	log.Println("Datos de luz enviados correctamente a la API.")
	return nil
}
