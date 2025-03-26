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
    repository    repositories.AirQualityRepository
    mqttAdapter   *adapters.MQTTClientAdapter
}

func NewAirQualityService(
    repository repositories.AirQualityRepository,
    mqttAdapter *adapters.MQTTClientAdapter,
) *AirQualityService {
    return &AirQualityService{
        repository:    repository,
        mqttAdapter:   mqttAdapter,
    }
}

func (service *AirQualityService) Start(topic string) error {
    if err := service.mqttAdapter.Connect(); err != nil {
        log.Println("❌ Error al conectar al broker MQTT:", err)
        return err
    }

    if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
        log.Println("❌ Error al suscribirse al topic:", err)
        return err
    }

    log.Println("✅ Suscripción exitosa al topic:", topic)
    return nil
}

func (service *AirQualityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
    log.Printf("🔊 Mensaje recibido: %s\n", msg.Payload())

    var airData entities.AirQuality
    if err := json.Unmarshal(msg.Payload(), &airData); err != nil {
        log.Println("❌ Error al parsear el mensaje:", err)
        return
    }

    // Filtro ajustado: Temperatura > 30°C
    if airData.Temperatura > 30 {
        if err := service.repository.ProcessAndForward(airData); err != nil {
            log.Println("❌ Error al reenviar los datos:", err)
            return
        }
        log.Println("✅ Datos enviados a la API 2:", airData)
    } else {
        log.Println("🌡️ Temperatura normal, ignorando...")
    }
}