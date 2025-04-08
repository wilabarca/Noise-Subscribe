package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	entities "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Entities"
)

// APIAdapter es una estructura para manejar interacciones con la API
type APIAdapter struct{}

func (a *APIAdapter) SendToAPI(tempHumidityData entities.TemperatureHumidity) any {
	panic("unimplemented")
}

// NewAPIAdapter crea una nueva instancia de APIAdapter
func NewAPIAdapter() *APIAdapter {
	return &APIAdapter{}
}

type TemperatureHumidityRepositoryAdapter struct {
	apiURL string
}

// NewTemperatureHumidityRepositoryAdapter crea una nueva instancia del adaptador para el repositorio de temperatura y humedad.
func NewTemperatureHumidityRepositoryAdapter(apiURL string) *TemperatureHumidityRepositoryAdapter {
	// Aseguramos que la URL predeterminada se utiliza si apiURL está vacío
	if apiURL == "" {
		apiURL = "http://localhost:8080/temperaturehumiditysensor/" // URL predeterminada
	}
	return &TemperatureHumidityRepositoryAdapter{
		apiURL: apiURL,
	}
}

// ProcessAndForward procesa los datos de temperatura y humedad y los reenvía a la API.
func (adapter *TemperatureHumidityRepositoryAdapter) ProcessAndForward(tempHumidityData entities.TemperatureHumidity) error {
	// No necesitamos el apiURL como parámetro, ya que lo estamos tomando de la propiedad del adaptador.
	data, err := json.Marshal(tempHumidityData)
	if err != nil {
		log.Printf("Error al serializar datos: %v | Datos: %+v\n", err, tempHumidityData)
		return err
	}

	resp, err := http.Post(adapter.apiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Printf("Error en POST a %s: %v\n", adapter.apiURL, err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf(" Respuesta no exitosa: %d | URL: %s\n", resp.StatusCode, adapter.apiURL)
		return errors.New("código de estado: " + resp.Status)
	}

	log.Printf("Datos de temperatura y humedad enviados correctamente a %s\n", adapter.apiURL)
	return nil
}