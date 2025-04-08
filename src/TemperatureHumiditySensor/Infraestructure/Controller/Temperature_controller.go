package controller

import (
	"encoding/json"
	"log"
	"net/http"

	application "Noisesubscribe/src/TemperatureHumiditySensor/Application"
)

type TemperatureHumidityController struct {
	SensorService *application.TemperatureHumidityService
}

// NewTemperatureHumidityController crea una nueva instancia del controlador
func NewTemperatureHumidityController(sensorService *application.TemperatureHumidityService) *TemperatureHumidityController {
	return &TemperatureHumidityController{SensorService: sensorService}
}

func (c *TemperatureHumidityController) StartSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		Topic  string `json:"topic"`
		ApiURL string `json:"apiURL"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Error al decodificar los datos", http.StatusBadRequest)
		return
	}

	// Logs de depuración
	log.Println("🌡️ Topic recibido:", requestData.Topic)

	// Iniciar la suscripción usando el servicio
	if err := c.SensorService.Start(requestData.Topic, requestData.ApiURL); err != nil {
		http.Error(w, "Error al iniciar la suscripción", http.StatusInternalServerError)
		log.Println("Error al iniciar la suscripción de temperatura y humedad:", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "✅ Suscripción iniciada correctamente",
	})
}
