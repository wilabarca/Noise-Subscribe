package controller

import (
	"encoding/json"
	"log"
	"net/http"

	application "Noisesubscribe/src/AirQuality/Application"
)

// AirQualityController es el controlador que maneja las solicitudes HTTP relacionadas con la calidad del aire
type AirQualityController struct {
	service *application.AirQualityService
}

// NewAirQualityController crea una nueva instancia de AirQualityController
func NewAirQualityController(service *application.AirQualityService) *AirQualityController {
	return &AirQualityController{service: service}
}

// StartSubscription maneja la solicitud para iniciar la suscripción al broker MQTT
func (controller *AirQualityController) StartSubscription(w http.ResponseWriter, r *http.Request) {
	// Verificar que la solicitud sea un POST
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Obtener el topic y la URL de la API desde el cuerpo de la solicitud
	var requestData struct {
		Topic  string `json:"topic"`
		ApiURL string `json:"apiURL"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Error al decodificar los datos", http.StatusBadRequest)
		return
	}

	// Depuración: imprimir los valores recibidos
	log.Println("Topic:", requestData.Topic)
	log.Println("ApiURL:", requestData.ApiURL)

	// Iniciar la suscripción al topic usando el servicio
	if err := controller.service.Start(requestData.Topic, requestData.ApiURL); err != nil {
		http.Error(w, "Error al iniciar la suscripción", http.StatusInternalServerError)
		log.Println("Error al iniciar la suscripción:", err)
		return
	}

	// Responder que la suscripción fue exitosa
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Suscripción iniciada correctamente",
	})
}
