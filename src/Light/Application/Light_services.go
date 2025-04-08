package application

import (
	"encoding/json"
	"errors"
	"log"

	entities "Noisesubscribe/src/Light/Domain/Entities"
	repositories "Noisesubscribe/src/Light/Domain/Repositories"
	adapter "Noisesubscribe/src/Light/Infraestructure/Adapters" 

	"github.com/streadway/amqp"
)

type LightService struct {
	repository       repositories.LightRepository
	rabbitMQAdapter  *adapter.RabbitMQAdapter
	minLightLevel    float64
	maxLightLevel    float64
}

func NewLightService(repository repositories.LightRepository, rabbitMQAdapter *adapter.RabbitMQAdapter, minLightLevel, maxLightLevel float64) *LightService {
	if repository == nil {
		log.Fatal("El repositorio de luz no puede ser nulo")
	}
	if rabbitMQAdapter == nil {
		log.Fatal("El adaptador RabbitMQ no puede ser nulo")
	}

	return &LightService{
		repository:      repository,
		rabbitMQAdapter: rabbitMQAdapter,
		minLightLevel:   minLightLevel,
		maxLightLevel:   maxLightLevel,
	}
}

func (service *LightService) Start(queueName string, apiURL string) error {
	messages, err := service.rabbitMQAdapter.Consume()
	if err != nil {
		log.Println("❌ No se pudo consumir los mensajes:", err)
		return err
	}

	// Procesar los mensajes de la cola
	for msg := range messages {
		if err := service.processLightMessage(msg); err != nil {
			log.Println("❌ Error al procesar el mensaje de luz:", err)
		} else {
			log.Println("✅ Mensaje procesado correctamente")
		}
	}

	return nil
}

// processLightMessage procesa un mensaje recibido de la cola RabbitMQ
func (service *LightService) processLightMessage(msg amqp.Delivery) error {
	var lightData entities.Light

	if err := json.Unmarshal(msg.Body, &lightData); err != nil {
		return errors.New("Error al deserializar el mensaje: " + err.Error())
	}

	return service.checkLightLevel(lightData)
}


func (service *LightService) checkLightLevel(lightData entities.Light) error {
	if lightData.Nivel < service.minLightLevel {
		log.Println("Nivel de luz bajo, se requiere más luz.")
		return nil
	}
	if lightData.Nivel > service.maxLightLevel {
		log.Println("Nivel de luz alto, se requiere reducir la luz.")
		return nil
	}

	log.Println("Nivel de luz adecuado.")
	return nil
}
