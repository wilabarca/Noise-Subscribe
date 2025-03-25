package controllers

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
	// Aquí, dependiendo de tu implementación, se puede obtener el estado de la luz desde algún repositorio o servicio.
	// Por ejemplo, si es un estado guardado en base de datos, se consulta y se devuelve.
	lightStatus := map[string]bool{"estado": true} // Estado de la luz (true = encendida, false = apagada)

	// Convertir el estado a JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lightStatus)
}

// TurnOnLight enciende la luz
func (controller *LightController) TurnOnLight(w http.ResponseWriter, r *http.Request) {
	// Crear una entidad Light con estado encendido
	lightData := entities.Light{Estado: true}

	// Procesar y reenviar los datos
	if err := controller.service.repository.ProcessAndForward(lightData); err != nil {
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
	// Crear una entidad Light con estado apagado
	lightData := entities.Light{Estado: false}

	// Procesar y reenviar los datos
	if err := controller.service.repository.ProcessAndForward(lightData); err != nil {
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

	// Procesar y reenviar los datos
	if err := controller.service.repositories.ProcessAndForward(lightData); err != nil {
		log.Println("❌ Error al ajustar la intensidad de la luz:", err)
		http.Error(w, "Error al ajustar la intensidad", http.StatusInternalServerError)
		return
	}

	// Responder con éxito
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("✅ Intensidad de la luz ajustada"))
}
