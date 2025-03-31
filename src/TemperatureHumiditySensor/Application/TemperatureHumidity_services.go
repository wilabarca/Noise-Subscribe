package application

import (
	entities "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Entities"
	repositories "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Repositories"
	adapters "Noisesubscribe/src/TemperatureHumiditySensor/Infraestructure/Adapters"
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// TemperatureHumidityService maneja la lógica de negocio para los datos de temperatura y humedad
type TemperatureHumidityService struct {
	repository  repositories.TemperatureHumidityRepository
	mqttAdapter *adapters.MQTTClientAdapter
	apiAdapter  *adapters.APIAdapter // Aquí cambiamos a 'APIAdapter'
	apiURL      string
}

// NewTemperatureHumidityService crea una nueva instancia de TemperatureHumidityService
func NewTemperatureHumidityService(mqttAdapter *adapters.MQTTClientAdapter, apiURL string) *TemperatureHumidityService {
	return &TemperatureHumidityService{
		repository:  adapters.NewTemperatureHumidityRepositoryAdapter(apiURL),
		mqttAdapter: mqttAdapter,
		apiAdapter:  adapters.NewAPIAdapter(), // Aquí inicializamos el 'APIAdapter'
		apiURL:      apiURL,
	}
}

// Start comienza la suscripción al topic de MQTT y maneja la URL de la API.
func (service *TemperatureHumidityService) Start(topic string) error {
	// Conectar al broker MQTT
	if err := service.mqttAdapter.Connect(); err != nil {
		log.Println("Error al conectar al broker MQTT:", err)
		return err
	}

	// Suscribirse al topic
	if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
		log.Println("Error al suscribirse al topic:", err)
		return err
	}

	log.Println("Suscripción exitosa al topic:", topic)
	return nil
}

// messageHandler maneja los mensajes MQTT entrantes y los procesa.
func (service *TemperatureHumidityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensaje recibido: %s\n", msg.Payload())

	var tempHumidityData entities.TemperatureHumidity
	if err := json.Unmarshal(msg.Payload(), &tempHumidityData); err != nil {
		log.Println("Error al parsear el mensaje:", err)
		return
	}

	// Filtro: Si la temperatura > 30°C
	if tempHumidityData.Temperature > 30 {
		// Usar el adaptador de API para enviar los datos
		if err := service.apiAdapter.SendToAPI(tempHumidityData); err != nil {
			log.Println("Error al enviar los datos a la API:", err)
			return
		}
		log.Println("Datos de temperatura y humedad enviados a la API Consumidora:", tempHumidityData)
	} else {
		log.Println("Temperatura normal, ignorando...")
	}
}
