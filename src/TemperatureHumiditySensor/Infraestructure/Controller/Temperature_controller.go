package controller

import (
	application "Noisesubscribe/src/TemperatureHumiditySensor/Application"
	"log"
)

type TemperatureHumidityController struct {
	SensorService *application.TemperatureHumidityService
}

func NewTemperatureHumidityController(sensorService *application.TemperatureHumidityService) *TemperatureHumidityController {
	return &TemperatureHumidityController{SensorService: sensorService}
}

func (c *TemperatureHumidityController) Start(topic string) {
	if err := c.SensorService.Start(topic); err != nil {
		log.Println("Error al iniciar la suscripción de temperatura🌡️:", err)
	}
}
