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

	if err := service.mqttAdapter.Connect(); err != nil {
		log.Println("Error al conectar al broker MQTT:", err)
		return err
	}

	if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
		log.Println("Error al suscribirse al topic:", err)
		return err
	}

	log.Println("Suscripción exitosa al topic:", topic)

	// Consumir mensajes de la cola en RabbitMQ
	go service.consumeRabbitMQMessages("SoundSensorQueue") // consume en paralelo
	return nil
}

// Lógica para procesar el mensaje de MQTT
func (service *SoundSensorService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensaje recibido: %s\n", msg.Payload())

	var soundData entities.SoundSensor
	if err := json.Unmarshal(msg.Payload(), &soundData); err != nil {
		log.Println("Error al parsear el mensaje:", err)
		return
	}

	if soundData.RuidoDB > 70 {
		if err := service.repository.ProcessAndForward(soundData); err != nil {
			log.Println("Error al reenviar los datos a la API:", err)
			return
		}
		log.Println("Datos enviados a la API Consumidora:", soundData)
	} else {
		log.Println("Nivel de sonido normal, ignorando...")
	}
}

func (service *SoundSensorService) consumeRabbitMQMessages(queueName string) {
	messages, err := service.rabbitMQAdapter.Consume()
	if err != nil {
		log.Println("❌ Error al consumir mensajes de RabbitMQ:", err)
		return
	}

	for msg := range messages {
		log.Printf("Mensaje recibido de RabbitMQ: %s\n", msg.Body)

		var soundData entities.SoundSensor
		if err := json.Unmarshal(msg.Body, &soundData); err != nil {
			log.Println("Error al parsear el mensaje de RabbitMQ:", err)
			continue
		}

		if soundData.RuidoDB > 70 {
			if err := service.repository.ProcessAndForward(soundData); err != nil {
				log.Println("Error al reenviar los datos a la API:", err)
				continue
			}
			log.Println("Datos reenviados a la API Consumidora desde RabbitMQ:", soundData)
		} else {
			log.Println("Nivel de sonido normal desde RabbitMQ, ignorando...")
		}
	}
}
