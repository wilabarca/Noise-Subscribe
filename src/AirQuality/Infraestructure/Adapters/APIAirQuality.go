package adapters



import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	entities "Noisesubscribe/src/AirQuality/Domain/Entities"
)

// AirQualityRepositoryAdapter es el adapter que implementa AirQualityRepository
type AirQualityRepositoryAdapter struct {
	ApiURL string
}

// NewAirQualityRepositoryAdapter crea una nueva instancia del adapter
func NewAirQualityRepositoryAdapter(apiURL string) *AirQualityRepositoryAdapter {
	return &AirQualityRepositoryAdapter{
		ApiURL: apiURL,
	}
}

// ProcessAndForward procesa los datos de calidad del aire y los envía a la API 2
func (adapter *AirQualityRepositoryAdapter) ProcessAndForward(airData entities.AirQuality) error {
	// Convertir los datos a JSON
	data, err := json.Marshal(airData)
	if err != nil {
		log.Println("❌ Error al convertir los datos a JSON:", err)
		return err
	}

	// Realizar la solicitud HTTP a la API 2
	resp, err := http.Post(adapter.ApiURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Println("❌ Error al enviar los datos a la API 2:", err)
		return err
	}
	defer resp.Body.Close()

	// Verificar que la respuesta sea exitosa
	if resp.StatusCode != http.StatusOK {
		return errors.New("❌ Error al reenviar los datos a la API 2: " + resp.Status)
	}

	log.Println("✅ Datos enviados a la API 2 exitosamente:", airData)
	return nil
}
