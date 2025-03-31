package controller

import (
	application "Noisesubscribe/src/TemperatureHumiditySensor/Application"
	"log"
)

// TemperatureHumidityController gestiona la suscripción y lógica de procesamiento de los datos de temperatura y humedad
type TemperatureHumidityController struct {
	SensorService *application.TemperatureHumidityService
}

// NewTemperatureHumidityController crea una nueva instancia de TemperatureHumidityController
func NewTemperatureHumidityController(sensorService *application.TemperatureHumidityService) *TemperatureHumidityController {
	return &TemperatureHumidityController{SensorService: sensorService}
}

// Start inicia la suscripción al broker MQTT y maneja los mensajes entrantes.
func (c *TemperatureHumidityController) Start(topic string) {
	log.Println("Iniciando la suscripción al broker MQTT...")
	if err := c.SensorService.Start(topic); err != nil {
		log.Println("Error al iniciar la suscripción:", err)
	}
}
