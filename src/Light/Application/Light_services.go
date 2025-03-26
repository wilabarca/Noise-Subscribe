package application

import (
	"errors"
	"log"

	entities "Noisesubscribe/src/Light/Domain/Entities"
	repositories "Noisesubscribe/src/Light/Domain/Repositories"
)

// LightService maneja la lógica de negocio para los datos de la luz
type LightService struct {
	repository repositories.LightRepository
}

// NewLightService crea una nueva instancia de LightService
func NewLightService(repository repositories.LightRepository) *LightService {
	if repository == nil {
		log.Fatal("❌ El repositorio de luz no puede ser nulo")
	}
	return &LightService{repository: repository}
}

// GetLightStatus obtiene el estado de la luz
func (service *LightService) GetLightStatus() (bool, error) {
	// Aquí simplemente se retorna el estado de la luz como ejemplo
	// En una implementación real, se consultarían los datos desde el repositorio o base de datos
	lightStatus := true // Asumimos que la luz está encendida

	return lightStatus, nil
}

// TurnOnLight enciende la luz
func (service *LightService) TurnOnLight() error {
	// Crear una entidad Light con estado encendido
	lightData := entities.Light{Estado: true}

	// Procesar los datos a través del repositorio
	if err := service.repository.ProcessAndForward(lightData); err != nil {
		return errors.New("❌ Error al encender la luz: " + err.Error())
	}

	return nil
}

// TurnOffLight apaga la luz
func (service *LightService) TurnOffLight() error {
	// Crear una entidad Light con estado apagado
	lightData := entities.Light{Estado: false}

	// Procesar los datos a través del repositorio
	if err := service.repository.ProcessAndForward(lightData); err != nil {
		return errors.New("❌ Error al apagar la luz: " + err.Error())
	}

	return nil
}

// SetLightIntensity ajusta la intensidad de la luz
func (service *LightService) SetLightIntensity(lightData entities.Light) error {
	// Procesar los datos a través del repositorio
	if err := service.repository.ProcessAndForward(lightData); err != nil {
		return errors.New("❌ Error al ajustar la intensidad de la luz: " + err.Error())
	}

	return nil
}
