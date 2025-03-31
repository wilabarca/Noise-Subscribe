package application

import (
	"encoding/json"
	"log"

	entities "Noisesubscribe/src/AirQuality/Domain/Entities"
	repositories "Noisesubscribe/src/AirQuality/Domain/Repositories"
	adapters "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type AirQualityService struct {
	repository  repositories.AirQualityRepository
	mqttAdapter *adapters.MQTTClientAdapter
	apiURL      string
}

func NewAirQualityService(mqttAdapter *adapters.MQTTClientAdapter, apiURL string) *AirQualityService {
	return &AirQualityService{
		repository:  adapters.NewAirQualityRepositoryAdapter(apiURL),
		mqttAdapter: mqttAdapter,
		apiURL:      apiURL,
	}
}

// Start inicia la suscripción con el topic y la URL de la API
func (service *AirQualityService) Start(topic string, apiURL string) error {
	// Establecer la URL de la API si no se pasó antes
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
		// No pasamos el apiURL aquí porque el adaptador ya lo maneja internamente
		if err := service.repository.ProcessAndForward(airData); err != nil {
			log.Println("Error al reenviar los datos:", err)
			return
		}
		log.Println("Datos enviados a la API Consumidora:", airData)
	} else {
		log.Println("Temperatura normal, ignorando...")
	}
}
