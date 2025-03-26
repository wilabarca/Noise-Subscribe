package controller

import (
	application "Noisesubscribe/src/TemperatureHumiditySensor/Application"
	"log"
	"github.com/eclipse/paho.mqtt.golang"
)

// TemperatureHumidityController maneja la lógica para la suscripción y procesamiento de datos del sensor
type TemperatureHumidityController struct {
	SensorService *application.TemperatureHumidityService
}

// NewTemperatureHumidityController crea una nueva instancia del controlador
func NewTemperatureHumidityController(sensorService *application.TemperatureHumidityService) *TemperatureHumidityController {
	return &TemperatureHumidityController{SensorService: sensorService}
}

// Start inicia la suscripción al broker MQTT y la escucha de los mensajes
func (c *TemperatureHumidityController) Start(mqttClient mqtt.Client, topic string) {
	log.Println("🚀 Iniciando la suscripción al broker MQTT...")
	if err := c.SensorService.Start(mqttClient, topic); err != nil {
		log.Println("❌ Error al iniciar la suscripción:", err)
	}
}
