package controller

import (
	application "Noisesubscribe/src/SoundSensor/Application"
	"log"
)

type SoundSensorController struct {
	SensorService *application.SoundSensorService
}

// NewSoundSensorController crea una nueva instancia del controlador de SoundSensor
func NewSoundSensorController(sensorService *application.SoundSensorService) *SoundSensorController {
	return &SoundSensorController{SensorService: sensorService}
}

// Start inicia la suscripción al broker MQTT y la escucha de los mensajes
func (c *SoundSensorController) Start(topic string) {
	log.Println("Iniciando la suscripción al broker MQTT...")
	c.SensorService.SubscribeToMQTT(topic)
}
