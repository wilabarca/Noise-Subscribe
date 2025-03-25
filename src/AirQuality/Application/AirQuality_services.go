package application

import (
	"Noisesubscribe/src/AirQuality/Domain/Entities"
	"Noisesubscribe/src/AirQuality/Domain/Repositories"
	adapters "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"
	"encoding/json"
	"log"

	"github.com/eclipse/paho.mqtt.golang"
)

// AirQualityService es el servicio que maneja los datos de calidad del aire
type AirQualityService struct {
    repository   repositories.AirQualityRepository
    mqttAdapter  *adapters.MQTTClientAdapter  // Adaptador MQTT
}

// NewAirQualityService crea una nueva instancia de AirQualityService
func NewAirQualityService(repository repositories.AirQualityRepository, mqttAdapter *adapters.MQTTClientAdapter) *AirQualityService {
    return &AirQualityService{
        repository:  repository,
        mqttAdapter: mqttAdapter,
    }
}

// Start inicia la suscripción al broker MQTT y maneja los mensajes de calidad del aire
func (service *AirQualityService) Start(topic string) error {
    // Conectar al broker MQTT
    if err := service.mqttAdapter.Connect(); err != nil {
        log.Println("❌ Error al conectar al broker MQTT:", err)
        return err
    }

    // Se suscribe al topic donde llegan los datos de calidad del aire
    if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
        log.Println("❌ Error al suscribirse al topic:", err)
        return err
    }

    log.Println("✅ Suscripción exitosa al topic:", topic)
    return nil
}

// messageHandler procesa los mensajes recibidos del broker MQTT
func (service *AirQualityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
    // Log para visualizar el mensaje recibido
    log.Printf("🔊 Mensaje recibido: %s\n", msg.Payload())

    // Deserializar los datos de calidad del aire
    var airData entities.AirQuality
    if err := json.Unmarshal(msg.Payload(), &airData); err != nil {
        log.Println("❌ Error al parsear el mensaje:", err)
        return
    }

    // Filtrar los datos antes de reenviarlos
    if airData.Temperatura < 35 {
        log.Println("🌡️ Temperatura demasiado baja, ignorado.")
        return
    }

    // Reenviar los datos relevantes a la API 2
    if err := service.repository.ProcessAndForward(airData); err != nil {
        log.Println("❌ Error al reenviar los datos:", err)
        return
    }

    log.Println("✅ Datos de calidad del aire enviados a la API 2:", airData)
}
