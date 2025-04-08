package repositories

import "Noisesubscribe/src/Light/Domain/Entities"

type LightRepository interface {
	ProcessAndForward(lightData entities.Light) error
	
	GetLightData() (entities.Light, error) 
}
