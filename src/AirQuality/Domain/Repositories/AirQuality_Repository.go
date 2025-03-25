package repositories

import entities "Noisesubscribe/src/AirQuality/Domain/Entities"

type AirQualityRepository interface {
	ProcessAndForward(airData entities.AirQuality)error
}