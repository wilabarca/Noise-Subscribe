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
	repository      repositories.LightRepository
	rabbitMQAdapter *adapter.RabbitMQAdapter
	minLightLevel   float64
	maxLightLevel   float64
}

func NewLightService(repository repositories.LightRepository, rabbitMQAdapter *adapter.RabbitMQAdapter, minLightLevel, maxLightLevel float64) *LightService {
	if repository == nil {
		log.Fatal("[LightService] ❌ El repositorio de luz no puede ser nulo")
	}
	if rabbitMQAdapter == nil {
		log.Fatal("[LightService] ❌ El adaptador RabbitMQ no puede ser nulo")
	}

	return &LightService{
		repository:      repository,
		rabbitMQAdapter: rabbitMQAdapter,
		minLightLevel:   minLightLevel,
		maxLightLevel:   maxLightLevel,
	}
}

func (service *LightService) Start(queueName string, apiURL string) error {
	log.Printf("[LightService] 🚦 Iniciando consumo de mensajes en la cola: %s...\n", queueName)

	messages, err := service.rabbitMQAdapter.Consume()
	if err != nil {
		log.Println("[LightService] ❌ No se pudo consumir los mensajes desde RabbitMQ:", err)
		return err
	}

	for msg := range messages {
		log.Println("[LightService] 📥 Mensaje recibido desde RabbitMQ.")
		if err := service.processLightMessage(msg); err != nil {
			log.Println("[LightService] ❌ Error al procesar el mensaje de luz:", err)
		} else {
			log.Println("[LightService] ✅ Mensaje procesado correctamente.")
		}
	}

	return nil
}

func (service *LightService) processLightMessage(msg amqp.Delivery) error {
	var lightData entities.Light

	if err := json.Unmarshal(msg.Body, &lightData); err != nil {
		return errors.New("[LightService] ❌ Error al deserializar el mensaje: " + err.Error())
	}

	log.Printf("[LightService] 🔎 Datos recibidos: %+v\n", lightData)
	return service.checkLightLevel(lightData)
}

func (service *LightService) checkLightLevel(lightData entities.Light) error {
	if lightData.Nivel < service.minLightLevel {
		log.Printf("[LightService] 🔅 Nivel de luz bajo (%.2f). Acción recomendada: Aumentar iluminación.\n", lightData.Nivel)
		return nil
	}
	if lightData.Nivel > service.maxLightLevel {
		log.Printf("[LightService] 🔆 Nivel de luz alto (%.2f). Acción recomendada: Disminuir iluminación.\n", lightData.Nivel)
		return nil
	}

	log.Printf("[LightService] ✅ Nivel de luz adecuado (%.2f).\n", lightData.Nivel)
	return nil
}
