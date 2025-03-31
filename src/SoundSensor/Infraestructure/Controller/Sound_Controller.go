package controller

import (
	"encoding/json"
	"log"
	"net/http"

	application "Noisesubscribe/src/SoundSensor/Application"
)

// SoundSensorController es el controlador que maneja las solicitudes HTTP relacionadas con el sensor de sonido
type SoundSensorController struct {
	service *application.SoundSensorService
}

// NewSoundSensorController crea una nueva instancia de SoundSensorController
func NewSoundSensorController(service *application.SoundSensorService) *SoundSensorController {
	return &SoundSensorController{service: service}
}

// StartSubscription maneja la solicitud para iniciar la suscripción al broker MQTT para el sensor de sonido
func (controller *SoundSensorController) StartSubscription(w http.ResponseWriter, r *http.Request) {
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
