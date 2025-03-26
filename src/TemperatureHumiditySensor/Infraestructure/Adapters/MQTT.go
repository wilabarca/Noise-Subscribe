package adapters

import (
	"encoding/json"
	"log"

	application "Noisesubscribe/src/TemperatureHumiditySensor/Application"
	entities "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Entities"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTAdapter maneja la comunicación con el broker MQTT
type MQTTAdapter struct {
	client  mqtt.Client
	service *application.TemperatureHumidityService
}

// NewMQTTAdapter crea una nueva instancia del adaptador MQTT
func NewMQTTAdapter(brokerURL string, service *application.TemperatureHumidityService) *MQTTAdapter {
	opts := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID("temperature-humidity-subscriber")
	client := mqtt.NewClient(opts)

	return &MQTTAdapter{
		client:  client,
		service: service,
	}
}

// SubscribeToMQTT se suscribe al broker MQTT y maneja los mensajes recibidos
func (adapter *MQTTAdapter) SubscribeToMQTT(topic string) error {
	if token := adapter.client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("❌ Error al conectar con el broker MQTT:", token.Error())
		return token.Error()
	}

	// Suscribirse al topic donde llegan los datos de temperatura y humedad
	if token := adapter.client.Subscribe(topic, 0, adapter.messageHandler); token.Wait() && token.Error() != nil {
		log.Println("❌ Error al suscribirse al topic:", token.Error())
		return token.Error()
	}

	log.Println("✅ Suscripción exitosa al topic:", topic)
	return nil
}

// messageHandler maneja los mensajes recibidos del broker MQTT
func (adapter *MQTTAdapter) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("🔊 Mensaje recibido: %s\n", msg.Payload())

	// Convertir el payload a la estructura esperada
	var sensorData entities.TemperatureHumidity
	if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
		log.Println("❌ Error al parsear el mensaje:", err)
		return
	}

	// Filtrar los datos si la temperatura o la humedad son valores extremos
	if sensorData.Temperature < -50 || sensorData.Temperature > 60 {
		log.Println("⚠️ Temperatura fuera de rango, ignorando.")
		return
	}
	if sensorData.Humidity < 0 || sensorData.Humidity > 100 {
		log.Println("⚠️ Humedad fuera de rango, ignorando.")
		return
	}

	// Reenviar los datos relevantes a la API
	if err := adapter.service.Repository.ProcessAndForward(sensorData); err != nil {
		log.Println("❌ Error al reenviar datos a la API:", err)
		return
	}

	log.Println("✅ Datos de temperatura y humedad enviados a la API:", sensorData)
}
