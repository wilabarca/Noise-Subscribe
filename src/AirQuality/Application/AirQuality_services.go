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

	log.Println("[AirQualityService] 🚀 Iniciando conexión con el broker MQTT...")

	if err := service.mqttAdapter.Connect(); err != nil {
		log.Println("[AirQualityService] ❌ Error al conectar al broker MQTT:", err)
		return err
	}

	if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
		log.Println("[AirQualityService] ❌ Error al suscribirse al topic:", err)
		return err
	}

	log.Printf("[AirQualityService] ✅ Suscripción exitosa al topic: %s\n", topic)

	go service.consumeRabbitMQMessages("sensor_air")

	return nil
}

// ✅ Maneja los mensajes recibidos desde MQTT
func (service *AirQualityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("[AirQualityService] 📩 Mensaje MQTT recibido: %s\n", msg.Payload())

	var airData entities.AirQualitySensor
	if err := json.Unmarshal(msg.Payload(), &airData); err != nil {
		log.Println("[AirQualityService] ❌ Error al parsear el mensaje MQTT:", err)
		return
	}

	if airData.CO2PPM > 1200 {
		log.Printf("[AirQualityService] ⚠️ Nivel alto de CO₂ detectado: %.2d ppm\n", airData.CO2PPM)
		if err := service.repository.ProcessAndForward(airData); err != nil {
			log.Println("[AirQualityService] ❌ Error al reenviar los datos a la API:", err)
			return
		}
		log.Println("[AirQualityService] ✅ Datos enviados a la API Consumidora:", airData)
	} else {
		log.Printf("[AirQualityService] 🌿 Calidad del aire estable (%.2d ppm)\n", airData.CO2PPM)
	}
}

// ✅ Consume mensajes de la cola de RabbitMQ
func (service *AirQualityService) consumeRabbitMQMessages(queueName string) {
	log.Printf("[AirQualityService] 📦 Iniciando consumo de mensajes desde RabbitMQ en la cola '%s'...\n", queueName)

	messages, err := service.rabbitMQAdapter.Consume()
	if err != nil {
		log.Println("[AirQualityService] ❌ Error al consumir mensajes de RabbitMQ:", err)
		return
	}

	for msg := range messages {
		log.Printf("[AirQualityService] 📥 Mensaje recibido desde RabbitMQ: %s\n", msg.Body)

		var airData entities.AirQualitySensor
		if err := json.Unmarshal(msg.Body, &airData); err != nil {
			log.Println("[AirQualityService] ❌ Error al parsear el mensaje de RabbitMQ:", err)
			continue
		}

		if airData.CO2PPM > 1200 {
			log.Printf("[AirQualityService] ⚠️ Nivel alto de CO₂ desde RabbitMQ: %.2d ppm\n", airData.CO2PPM)
			if err := service.repository.ProcessAndForward(airData); err != nil {
				log.Println("[AirQualityService] ❌ Error al reenviar los datos a la API:", err)
				continue
			}
			log.Println("[AirQualityService] ✅ Datos reenviados a la API Consumidora desde RabbitMQ:", airData)
		} else {
			log.Printf("[AirQualityService] 🌿 Calidad de aire estable desde RabbitMQ (%.2d ppm), mensaje ignorado.\n", airData.CO2PPM)
		}
	}
}
