package application

import (
	entities "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Entities"
	repositories "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Repositories"
	"encoding/json"
	"log"

	"github.com/eclipse/paho.mqtt.golang"
)

// TemperatureHumidityService maneja la lógica de negocio para los datos de temperatura y humedad
type TemperatureHumidityService struct {
	repository repositories.TemperatureHumidityRepository
}

// NewTemperatureHumidityService crea una nueva instancia de TemperatureHumidityService
func NewTemperatureHumidityService(repository repositories.TemperatureHumidityRepository) *TemperatureHumidityService {
	return &TemperatureHumidityService{repository: repository}
}

// Start inicia la suscripción al broker MQTT y maneja los mensajes relacionados con la temperatura y humedad
func (service *TemperatureHumidityService) Start(mqttClient mqtt.Client, topic string) error {
	// Se suscribe al topic donde llegan los datos de temperatura y humedad
	if token := mqttClient.Subscribe(topic, 0, service.messageHandler); token.Wait() && token.Error() != nil {
		log.Println("❌ Error al suscribirse al topic:", token.Error())
		return token.Error()
	}

	log.Println("✅ Suscripción exitosa al topic:", topic)
	return nil
}

// messageHandler procesa los mensajes recibidos sobre la temperatura y la humedad
func (service *TemperatureHumidityService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	// Log para visualizar el mensaje recibido
	log.Printf("🔊 Mensaje recibido: %s\n", msg.Payload())

	// Deserializar los datos de temperatura y humedad
	var tempHumidityData entities.TemperatureHumidity
	if err := json.Unmarshal(msg.Payload(), &tempHumidityData); err != nil {
		log.Println("❌ Error al parsear el mensaje:", err)
		return
	}

	// Reenviar los datos de temperatura y humedad a la API o sistema correspondiente
	if err := service.repository.ProcessAndForward(tempHumidityData); err != nil {
		log.Println("❌ Error al reenviar los datos:", err)
		return
	}

	log.Println("✅ Datos de temperatura y humedad enviados a la API:", tempHumidityData)
}
