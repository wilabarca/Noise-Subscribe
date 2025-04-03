package application

import (
	"encoding/json"
	"log"

	entities "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Entities"
	repositories "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Repositories"
	adapterRepo "Noisesubscribe/src/TemperatureHumiditySensor/Infraestructure/Adapters"
	rabbit "Noisesubscribe/src/TemperatureHumiditySensor/Infraestructure/Adapters"
	"github.com/eclipse/paho.mqtt.golang"
)

type TemperatureHumidityService struct {
	repository      repositories.TemperatureHumidityRepository
	mqttAdapter     *adapterRepo.MQTTClientAdapter
	apiAdapter      *adapterRepo.APIAdapter
	apiURL          string
	rabbitMQAdapter *rabbit.RabbitMQAdapter
}

func NewTemperatureHumidityService(mqttAdapter *adapterRepo.MQTTClientAdapter, apiURL string, rabbitMQAdapter *rabbit.RabbitMQAdapter) *TemperatureHumidityService {
	return &TemperatureHumidityService{
		repository:      adapterRepo.NewTemperatureHumidityRepositoryAdapter(apiURL),
		mqttAdapter:     mqttAdapter,
		apiAdapter:      adapterRepo.NewAPIAdapter(), // Inicialización del APIAdapter
		apiURL:          apiURL,
		rabbitMQAdapter: rabbitMQAdapter,
	}
}

// Método para iniciar el consumo de mensajes de RabbitMQ
func (service *TemperatureHumidityService) StartConsuming(queueName string) error {
	messages, err := service.rabbitMQAdapter.Consume(queueName)
	if err != nil {
		log.Println("❌ Error al consumir los mensajes de RabbitMQ:", err)
		return err
	}

	// Procesar los mensajes consumidos
	for msg := range messages {
		log.Printf("Mensaje recibido de RabbitMQ: %s\n", msg.Body)

		var tempHumidityData entities.TemperatureHumidity
		if err := json.Unmarshal(msg.Body, &tempHumidityData); err != nil {
			log.Println("Error al parsear el mensaje:", err)
			continue
		}

		// Filtro: Si la temperatura > 30°C
		if tempHumidityData.Temperature > 30 {
			if err := service.apiAdapter.SendToAPI(tempHumidityData); err != nil {
				log.Println("Error al enviar los datos a la API:", err)
				continue
			}
			log.Println("Datos enviados a la API Consumidora:", tempHumidityData)
		} else {
			log.Println("Temperatura normal, ignorando...")
		}
	}

	return nil
}

// Método para iniciar la suscripción MQTT
func (service *TemperatureHumidityService) Start(topic string) error {
	if err := service.mqttAdapter.Connect(); err != nil {
		log.Println("❌ Error al conectar al broker MQTT:", err)
		return err
	}

	if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
		log.Println("❌ Error al suscribirse al topic:", err)
		return err
	}

	log.Println("✅ Suscripción exitosa al topic:", topic)
	return nil
}

// Handler para procesar los mensajes MQTT
func (service *TemperatureHumidityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensaje recibido de MQTT: %s\n", msg.Payload())

	var tempHumidityData entities.TemperatureHumidity
	if err := json.Unmarshal(msg.Payload(), &tempHumidityData); err != nil {
		log.Println("❌ Error al parsear el mensaje de MQTT:", err)
		return
	}

	// Filtro: Si la temperatura > 30°C
	if tempHumidityData.Temperature > 30 {
		if err := service.apiAdapter.SendToAPI(tempHumidityData); err != nil {
			log.Println("❌ Error al enviar los datos a la API:", err)
			return
		}
		log.Println("✅ Datos enviados a la API Consumidora:", tempHumidityData)
	} else {
		log.Println("Temperatura normal, ignorando...")
	}
}
