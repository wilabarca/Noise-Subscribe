package application

import (
	"encoding/json"
	"log"

	entities "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Entities"
	repositories "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Repositories"
	adapterRepo "Noisesubscribe/src/TemperatureHumiditySensor/Infraestructure/Adapters"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	rabbit "Noisesubscribe/src/TemperatureHumiditySensor/Infraestructure/Adapters"
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

func (service *TemperatureHumidityService) Start(topic string) error {
	if err := service.mqttAdapter.Connect(); err != nil {
		log.Println("Error al conectar al broker MQTT:", err)
		return err
	}

	if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
		log.Println("Error al suscribirse al topic:", err)
		return err
	}

	log.Println("Suscripción exitosa al topic:", topic)
	return nil
}

func (service *TemperatureHumidityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensaje recibido: %s\n", msg.Payload())

	var tempHumidityData entities.TemperatureHumidity
	if err := json.Unmarshal(msg.Payload(), &tempHumidityData); err != nil {
		log.Println("Error al parsear el mensaje:", err)
		return
	}

	// Filtro: Si la temperatura > 30°C
	if tempHumidityData.Temperature > 30 {
		if err := service.apiAdapter.SendToAPI(tempHumidityData); err != nil {
			log.Println("Error al enviar los datos a la API:", err)
			return
		}
		log.Println("Datos enviados a la API Consumidora:", tempHumidityData)

		// Publicar en RabbitMQ
		if err := service.rabbitMQAdapter.Publish("TemperatureHumidityQueue", msg.Payload()); err != nil {
			log.Println("Error al publicar en RabbitMQ:", err)
			return
		}
	} else {
		log.Println("Temperatura normal, ignorando...")
	}
}
