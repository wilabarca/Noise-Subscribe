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

func (service *TemperatureHumidityService) Start(topic string, apiURL string) error {
	if apiURL != "" {
		service.apiURL = apiURL
	}

	log.Println("[TemperatureHumidityService] 🔌 Iniciando conexión al broker MQTT...")
	if err := service.mqttAdapter.Connect(); err != nil {
		log.Println("[TemperatureHumidityService] ❌ Error al conectar al broker MQTT:", err)
		return err
	}
	log.Println("[TemperatureHumidityService] ✅ Conexión establecida con el broker MQTT.")

	log.Printf("[TemperatureHumidityService] 📡 Intentando suscripción al topic: %s...\n", topic)
	if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
		log.Println("[TemperatureHumidityService] ❌ Error al suscribirse al topic:", err)
		return err
	}
	log.Printf("[TemperatureHumidityService] ✅ Suscripción exitosa al topic: %s\n", topic)

	log.Println("[TemperatureHumidityService] 📥 Iniciando escucha de mensajes desde RabbitMQ para temperatura...")
	go service.consumeRabbitMQMessages("TemperatureSensorQueue")

	return nil
}


func (service *TemperatureHumidityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Println("[TemperatureHumidityService] 📥 Mensaje MQTT recibido para sensor de Temperatura y Humedad.")

	var tempHumidityData entities.TemperatureHumidity
	if err := json.Unmarshal(msg.Payload(), &tempHumidityData); err != nil {
		log.Println("❌ Error al parsear mensaje MQTT:", err)
		return
	}

	log.Printf("🌡️  Temperatura: %.2f°C | 💧 Humedad: %.2f%%\n", tempHumidityData.Temperature, tempHumidityData.Humidity)

	if tempHumidityData.Temperature > 30 {
		log.Println("🚨 Temperatura elevada detectada, enviando a la API...")
		if err := service.apiAdapter.SendToAPI(tempHumidityData); err != nil {
			log.Println("❌ Error al enviar datos a la API:", err)
			return
		}
		log.Println("✅ Datos enviados exitosamente a la API.")
	} 
}


func (service *TemperatureHumidityService) consumeRabbitMQMessages(queueName string) {
	log.Printf("🐇 Iniciando consumo de mensajes desde la cola RabbitMQ: %s\n", queueName)

	messages, err := service.rabbitMQAdapter.Consume()
	if err != nil {
		log.Printf("❌ Error al consumir mensajes del sensor de temperatura: %v\n", err)
		return
	}

	for msg := range messages {
		var TemData entities.TemperatureHumidity
		if err := json.Unmarshal(msg.Body, &TemData); err != nil {
			log.Printf("❌ Error al parsear mensaje RabbitMQ: %v\n", err)
			continue
		}

		log.Printf("🌡️  Temperatura: %.2f°C | 💧 Humedad: %.2f%%\n", TemData.Temperature, TemData.Humidity)

		if TemData.Temperature > 26 {
			log.Println("🚨 Temperatura elevada, reenviando a la API...")
			if err := service.repository.ProcessAndForward(TemData); err != nil {
				log.Printf("❌ Error al enviar los datos a la API: %v\n", err)
				continue
			}
			log.Println("✅ Datos reenviados exitosamente a la API.")
		} else {
			log.Println("🌤️  Temperatura agradable, no se requiere acción.")
		}
	}
}