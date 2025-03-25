package adapters

import (
	"encoding/json"
	"log"
	application "Noisesubscribe/src/SoundSensor/Application"
	entities "Noisesubscribe/src/SoundSensor/Domain/Entities"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTAdapter struct {
	client  mqtt.Client
	service *application.SoundSensorService
}

// NewMQTTAdapter crea una nueva instancia del adaptador MQTT
func NewMQTTAdapter(brokerURL string, service *application.SoundSensorService) *MQTTAdapter {
	opts := mqtt.NewClientOptions().AddBroker(brokerURL).SetClientID("sound-sensor-subscriber")
	client := mqtt.NewClient(opts)

	return &MQTTAdapter{
		client:  client,
		service: service,
	}
}

// SubscribeToMQTT se suscribe al broker MQTT y maneja los mensajes recibidos
func (adapter *MQTTAdapter) SubscribeToMQTT(topic string) error {
	if token := adapter.client.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Error al conectar con el broker MQTT:", token.Error())
		return token.Error()
	}

	// Suscribirse al topic donde llegan los datos de los sensores
	if token := adapter.client.Subscribe(topic, 0, adapter.messageHandler); token.Wait() && token.Error() != nil {
		log.Println("Error al suscribirse al topic:", token.Error())
		return token.Error()
	}

	log.Println("Suscripción exitosa al topic:", topic)
	return nil
}

// messageHandler maneja los mensajes recibidos del broker MQTT
func (adapter *MQTTAdapter) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensaje recibido: %s\n", msg.Payload())

	// Convertir el payload a la estructura esperada
	var sensorData entities.SoundSensor
	if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
		log.Println("Error al parsear el mensaje:", err)
		return
	}

	// Filtrar los datos si el ruido es menor de 40 dB
	if sensorData.RuidoDB < 40 {
		log.Println("Ruido demasiado bajo, ignorado.")
		return
	}

	// Reenviar los datos relevantes a la API 2
	if err := adapter.service.SendToAPI2(sensorData); err != nil {
		log.Println("Error al reenviar datos a la API 2:", err)
		return
	}

	log.Println("Datos enviados a la API 2:", sensorData)
}
