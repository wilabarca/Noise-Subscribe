package controller

import (
	"encoding/json"
	"log"
	"net/http"

	application "Noisesubscribe/src/Light/Application"
	entities "Noisesubscribe/src/Light/Domain/Entities"
)

// LightController maneja las solicitudes relacionadas con los datos de luz
type LightController struct {
	service *application.LightService
}

// NewLightController crea una nueva instancia de LightController
func NewLightController(service *application.LightService) *LightController {
	return &LightController{service: service}
}

// GetLightStatus obtiene el estado actual de la luz (encendida o apagada)
func (controller *LightController) GetLightStatus(w http.ResponseWriter, r *http.Request) {
	// Obtener el estado de la luz a través del servicio
	lightStatus, err := controller.service.GetLightStatus()
	if err != nil {
		http.Error(w, "Error al obtener el estado de la luz", http.StatusInternalServerError)
		return
	}

	// Convertir el estado a JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lightStatus)
}

// TurnOnLight enciende la luz
func (controller *LightController) TurnOnLight(w http.ResponseWriter, r *http.Request) {
	// Llamar al servicio para encender la luz
	if err := controller.service.TurnOnLight(); err != nil {
		log.Println("❌ Error al encender la luz:", err)
		http.Error(w, "Error al encender la luz", http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ La luz ha sido encendida"))
}

// TurnOffLight apaga la luz
func (controller *LightController) TurnOffLight(w http.ResponseWriter, r *http.Request) {
	// Llamar al servicio para apagar la luz
	if err := controller.service.TurnOffLight(); err != nil {
		log.Println("❌ Error al apagar la luz:", err)
		http.Error(w, "Error al apagar la luz", http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ La luz ha sido apagada"))
}

// SetLightIntensity ajusta la intensidad de la luz
func (controller *LightController) SetLightIntensity(w http.ResponseWriter, r *http.Request) {
	var lightData entities.Light

	// Decodificar los datos de la solicitud JSON
	if err := json.NewDecoder(r.Body).Decode(&lightData); err != nil {
		log.Println("❌ Error al decodificar los datos:", err)
		http.Error(w, "Datos inválidos", http.StatusBadRequest)
		return
	}

	// Llamar al servicio para ajustar la intensidad de la luz
	if err := controller.service.SetLightIntensity(lightData); err != nil {
		log.Println("❌ Error al ajustar la intensidad de la luz:", err)
		http.Error(w, "Error al ajustar la intensidad", http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Intensidad de la luz ajustada"))
}
