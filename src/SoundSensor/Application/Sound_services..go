package application

import (
	"encoding/json"
	"log"

	entities "Noisesubscribe/src/SoundSensor/Domain/Entities"
	repositories "Noisesubscribe/src/SoundSensor/Domain/Repositories"
	adapterRepo "Noisesubscribe/src/SoundSensor/Infraestructure/Adapters"
	rabbit "Noisesubscribe/src/SoundSensor/Infraestructure/Adapters"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SoundSensorService struct {
	repository      repositories.SoundSensorRepository
	mqttAdapter     *adapterRepo.MQTTClientAdapter
	apiURL          string
	rabbitMQAdapter *rabbit.RabbitMQAdapter
}

func NewSoundSensorService(mqttAdapter *adapterRepo.MQTTClientAdapter, apiURL string, rabbitMQAdapter *rabbit.RabbitMQAdapter) *SoundSensorService {
	return &SoundSensorService{
		repository:      adapterRepo.NewSoundSensorRepositoryAdapter(apiURL),
		mqttAdapter:     mqttAdapter,
		apiURL:          apiURL,
		rabbitMQAdapter: rabbitMQAdapter,
	}
}

func (service *SoundSensorService) Start(topic string, apiURL string) error {
	if apiURL != "" {
		service.apiURL = apiURL
	}

	log.Println("[SoundSensorService] 🔌 Iniciando conexión al broker MQTT...")
	if err := service.mqttAdapter.Connect(); err != nil {
		log.Println("[SoundSensorService] ❌ Error al conectar al broker MQTT:", err)
		return err
	}
	log.Println("[SoundSensorService] ✅ Conexión establecida con el broker MQTT.")

	log.Printf("[SoundSensorService] 📡 Intentando suscripción al topic: %s...\n", topic)
	if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
		log.Println("[SoundSensorService] ❌ Error al suscribirse al topic:", err)
		return err
	}
	log.Printf("[SoundSensorService] ✅ Suscripción exitosa al topic: %s\n", topic)

	log.Println("[SoundSensorService] 📥 Iniciando escucha de mensajes desde RabbitMQ para sonido...")
	go service.consumeRabbitMQMessages("SoundSensorQueue")

	return nil
}

func (service *SoundSensorService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Println("[SoundSensorService] 🔊 Mensaje recibido desde MQTT.")

	var soundData entities.SoundSensor
	if err := json.Unmarshal(msg.Payload(), &soundData); err != nil {
		log.Println("[SoundSensorService] ❌ Error al parsear el mensaje MQTT:", err)
		return
	}

	if soundData.RuidoDB > 70 {
		if err := service.repository.ProcessAndForward(soundData); err != nil {
			log.Println("[SoundSensorService] ❌ Error al reenviar los datos a la API:", err)
			return
		}
		log.Println("[SoundSensorService] 📤 Datos enviados a la API Consumidora:", soundData)
	} else {
		log.Println("[SoundSensorService] ✅ Nivel de sonido dentro de rango, sin acción requerida.")
	}
}

func (service *SoundSensorService) consumeRabbitMQMessages(queueName string) {
	messages, err := service.rabbitMQAdapter.Consume()
	if err != nil {
		log.Println("[SoundSensorService] ❌ Error al consumir mensajes de RabbitMQ:", err)
		return
	}

	for msg := range messages {
		log.Printf("[SoundSensorService] 📨 Mensaje recibido desde RabbitMQ: %s\n", msg.Body)

		var soundData entities.SoundSensor
		if err := json.Unmarshal(msg.Body, &soundData); err != nil {
			log.Println("[SoundSensorService] ❌ Error al parsear el mensaje de RabbitMQ:", err)
			continue
		}

		if soundData.RuidoDB > 70 {
			if err := service.repository.ProcessAndForward(soundData); err != nil {
				log.Println("[SoundSensorService] ❌ Error al reenviar los datos a la API:", err)
				continue
			}
			log.Println("[SoundSensorService] 📤 Datos reenviados a la API desde RabbitMQ:", soundData)
		} else {
			log.Println("[SoundSensorService] ✅ Nivel de sonido desde RabbitMQ dentro de rango, sin acción requerida.")
		}
	}
}
