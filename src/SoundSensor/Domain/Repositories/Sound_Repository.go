package repositories

import entities "Noisesubscribe/src/SoundSensor/Domain/Entities"

type SoundSensorRepository interface {
	// filtramos y procesaos los datos del sensor
	ProcessAndForward(sensorData entities.SoundSensor) error
}