package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"Noisesubscribe/src/Light/Domain/Entities"
)

// APIAdapter es una estructura que maneja la interacción con una API externa
type APIAdapter struct{}

// NewAPIAdapter crea una nueva instancia del adaptador para la API del sensor de luz.
func NewAPIAdapter() *APIAdapter {
	return &APIAdapter{}
}

// LightRepositoryAdapter es el adaptador que maneja la interacción con el repositorio de datos del sensor de luz.
type LightRepositoryAdapter struct {
	apiURL string
}

// NewLightRepositoryAdapter crea una nueva instancia del adaptador para el repositorio de luz.
func NewLightRepositoryAdapter(apiURL string) *LightRepositoryAdapter {
	// Aseguramos que la URL predeterminada se utiliza si apiURL está vacío
	if apiURL == "" {
		apiURL = "http://localhost:8080/lightsensor/" // URL predeterminada
	}
	return &LightRepositoryAdapter{
		apiURL: apiURL,
	}
}

// ProcessAndForward procesa los datos del sensor de luz y los envía a la API.
func (adapter *LightRepositoryAdapter) ProcessAndForward(lightData entities.Light) error {
	// Convertir la entidad Light a JSON
	data, err := json.Marshal(lightData)
	if err != nil {
		log.Printf("Error al serializar datos: %v | Datos: %+v\n", err, lightData)
		return err
	}

	// Realizar una solicitud HTTP POST a la API externa
	resp, err := http.Post(adapter.apiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error en POST a %s: %v\n", adapter.apiURL, err)
		return err
	}
	defer resp.Body.Close()

	// Verificar el código de respuesta de la API
	if resp.StatusCode != http.StatusOK {
		log.Printf("Respuesta no exitosa: %d | URL: %s\n", resp.StatusCode, adapter.apiURL)
		return errors.New("código de estado: " + resp.Status)
	}

	log.Printf("Datos de luz enviados correctamente a %s\n", adapter.apiURL)
	return nil
}
