package application

import (
	"encoding/json"
	"log"

	entities "Noisesubscribe/src/AirQuality/Domain/Entities"
	repositories "Noisesubscribe/src/AirQuality/Domain/Repositories"
	adapterRepo "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	rabbit "Noisesubscribe/src/AirQuality/Infraestructure/Adapters" 
)

type AirQualityService struct {
	repository      repositories.AirQualityRepository
	mqttAdapter     *adapterRepo.MQTTClientAdapter
	apiURL          string
	rabbitMQAdapter *rabbit.RabbitMQAdapter
}

func NewAirQualityService(mqttAdapter *adapterRepo.MQTTClientAdapter, apiURL string, rabbitMQAdapter *rabbit.RabbitMQAdapter) *AirQualityService {
	return &AirQualityService{
		repository:      adapterRepo.NewAirQualityRepositoryAdapter(apiURL),
		mqttAdapter:     mqttAdapter,
		apiURL:          apiURL,
		rabbitMQAdapter: rabbitMQAdapter,
	}
}

func (service *AirQualityService) Start(topic string, apiURL string) error {
	if apiURL != "" {
		service.apiURL = apiURL
	}

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

func (service *AirQualityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensaje recibido: %s\n", msg.Payload())

	var airData entities.AirQualitySensor
	if err := json.Unmarshal(msg.Payload(), &airData); err != nil {
		log.Println("Error al parsear el mensaje:", err)
		return
	}

	// Filtro: Temperatura > 30°C
	if airData.Temperatura > 30 {
		if err := service.repository.ProcessAndForward(airData); err != nil {
			log.Println("Error al reenviar los datos a la API:", err)
			return
		}
		log.Println("Datos enviados a la API Consumidora:", airData)

		// Publicar en RabbitMQ
		if err := service.rabbitMQAdapter.Publish("AirQualityQueue", msg.Payload()); err != nil {
			log.Println("Error al publicar en RabbitMQ:", err)
			return
		}
	} else {
		log.Println("Temperatura normal, ignorando...")
	}
}
