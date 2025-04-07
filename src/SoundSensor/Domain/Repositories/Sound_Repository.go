package repositories

import entities "Noisesubscribe/src/SoundSensor/Domain/Entities"

type SoundSensorRepository interface {
	ProcessAndForward(sensorData entities.SoundSensor) error
}