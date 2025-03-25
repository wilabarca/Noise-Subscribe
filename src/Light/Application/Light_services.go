package application

import (
	entities "Noisesubscribe/src/Light/Domain/Entities"
	repositories "Noisesubscribe/src/Light/Domain/Repositories"

	"encoding/json"
	"log"

	"github.com/eclipse/paho.mqtt.golang"
)

// LightService maneja la lógica de negocio para los datos de la luz
type LightService struct {
    repository repositories.LightRepository
}

// NewLightService crea una nueva instancia de LightService
func NewLightService(repository repositories.LightRepository) *LightService {
    return &LightService{repository: repository}
}

// Start inicia la suscripción al broker MQTT y maneja los mensajes relacionados con la luz
func (service *LightService) Start(mqttClient mqtt.Client, topic string) error {
    // Se suscribe al topic donde llegan los datos sobre la luz
    if token := mqttClient.Subscribe(topic, 0, service.messageHandler); token.Wait() && token.Error() != nil {
        log.Println("❌ Error al suscribirse al topic:", token.Error())
        return token.Error()
    }

    log.Println("✅ Suscripción exitosa al topic:", topic)
    return nil
}

// messageHandler procesa los mensajes recibidos sobre la luz
func (service *LightService) messageHandler(client mqtt.Client, msg mqtt.Message) {
    // Log para visualizar el mensaje recibido
    log.Printf("🔊 Mensaje recibido: %s\n", msg.Payload())

    // Deserializar los datos sobre la luz
    var lightData entities.Light
    if err := json.Unmarshal(msg.Payload(), &lightData); err != nil {
        log.Println("❌ Error al parsear el mensaje:", err)
        return
    }

    // Reenviar los datos de la luz a la API o sistema correspondiente
    if err := service.repository.ProcessAndForward(lightData); err != nil {
        log.Println("❌ Error al reenviar los datos:", err)
        return
    }

    log.Println("✅ Datos de luz enviados a la API:", lightData)
}
