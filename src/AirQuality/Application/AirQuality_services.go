package application

import (
	"encoding/json"
	"log"

	entities "Noisesubscribe/src/AirQuality/Domain/Entities"
	repositories "Noisesubscribe/src/AirQuality/Domain/Repositories"
	adapterRepo "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"
	rabbit "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/streadway/amqp"
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
		return err.(error)
	}

	if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
		log.Println("Error al suscribirse al topic:", err)
		return err.(error)
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

		// Ahora no publicamos en RabbitMQ, en lugar de eso vamos a consumir la cola
		// Esto solo se consume y no se publica nada.
		if err := service.StartQueueConsumer("AirQualityQueue"); err != nil {
			log.Println("Error al consumir la cola de RabbitMQ:", err)
			return
		}
	} else {
		log.Println("Temperatura normal, ignorando...")
	}
}

// Consume los mensajes de RabbitMQ
func (service *AirQualityService) StartQueueConsumer(queueName string) error {
	// Conectar al RabbitMQ
	if err := service.rabbitMQAdapter.Connect(); err != nil {
		log.Println("Error al conectar a RabbitMQ:", err)
		return err.(error)
	}

	// Consumir mensajes de la cola y pasar el handler
	if err := service.rabbitMQAdapter.Consume(queueName, service.messageHandlerFromQueue); err != nil {
		log.Println("Error al consumir mensajes de la cola:", err)
		return err
	}

	log.Println("Consumiendo mensajes de la cola:", queueName)

	return nil
}

// Manejar el mensaje recibido de la cola RabbitMQ
func (service *AirQualityService) messageHandlerFromQueue(msg amqp.Delivery) {
	// Log del mensaje recibido
	log.Printf("Mensaje recibido desde RabbitMQ: %s\n", msg.Body)

	// Deserializar los datos del mensaje (Body es un []byte)
	var airData entities.AirQualitySensor
	if err := json.Unmarshal(msg.Body, &airData); err != nil {
		log.Println("Error al parsear el mensaje:", err)
		msg.Ack(false)  // Acknowledge el mensaje aún si hubo error en el parseo
		return
	}

	// Filtro: Temperatura > 30°C
	if airData.Temperatura > 30 {
		if err := service.repository.ProcessAndForward(airData); err != nil {
			log.Println("Error al reenviar los datos a la API:", err)
			msg.Ack(false)  // Acknowledge el mensaje aún si hubo error en el reenvío
			return
		}
		log.Println("Datos enviados a la API Consumidora:", airData)
	} else {
		log.Println("Temperatura normal, ignorando...")
	}

	// Acknowledge el mensaje después de procesarlo correctamente
	msg.Ack(false)
}
