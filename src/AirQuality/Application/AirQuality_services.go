package application

import (
	"encoding/json"
	"log"

	entities "Noisesubscribe/src/AirQuality/Domain/Entities"
	repositories "Noisesubscribe/src/AirQuality/Domain/Repositories"
	adapterRepo "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"
	rabbit "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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

	log.Println("✅ Suscripción exitosa al topic:", topic)

	// Lanzamos el consumo de la cola en segundo plano
	go service.consumeRabbitMQMessages("sensor_air")

	return nil
}

// ✅ Maneja los mensajes de MQTT
func (service *AirQualityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensaje recibido: %s\n", msg.Payload())

	var airData entities.AirQualitySensor
	if err := json.Unmarshal(msg.Payload(), &airData); err != nil {
		log.Println("Error al parsear el mensaje:", err)
		return
	}

	if airData.CO2PPM > 1200 {
		if err := service.repository.ProcessAndForward(airData); err != nil {
			log.Println("Error al reenviar los datos a la API:", err)
			return
		}
		log.Println("Datos enviados a la API Consumidora:", airData)
	} else {
		log.Println("Calidad de aire estable, ignorando...")
	}
}

// ✅ Maneja los mensajes de RabbitMQ con loop como el sensor de sonido
func (service *AirQualityService) consumeRabbitMQMessages(queueName string) {
	messages, err := service.rabbitMQAdapter.Consume()
	if err != nil {
		log.Println("❌ Error al consumir mensajes de RabbitMQ:", err)
		return
	}

	for msg := range messages {
		log.Printf("Mensaje recibido desde la cola '%s': %s\n", queueName, msg.Body)

		var airData entities.AirQualitySensor
		if err := json.Unmarshal(msg.Body, &airData); err != nil {
			log.Println("Error al parsear el mensaje de RabbitMQ:", err)
			continue
		}

		if airData.CO2PPM > 1200 {
			if err := service.repository.ProcessAndForward(airData); err != nil {
				log.Println("Error al reenviar los datos a la API:", err)
				continue
			}
			log.Println("Datos reenviados a la API Consumidora desde RabbitMQ:", airData)
		} else {
			log.Println("Calidad de aire estable desde RabbitMQ, ignorando...")
		}
	}
}
