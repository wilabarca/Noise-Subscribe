package adapters

import (
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClientAdapter manages the MQTT client connection and subscription.
type MQTTClientAdapter struct {
	client mqtt.Client
}

// NewMQTTClientAdapter creates a new MQTTClientAdapter instance with the broker URL.
func NewMQTTClientAdapter(brokerURL string) *MQTTClientAdapter {
	opts := mqtt.NewClientOptions().AddBroker(brokerURL)
	client := mqtt.NewClient(opts)
	return &MQTTClientAdapter{
		client: client,
	}
}

// Connect establishes a connection to the MQTT broker.
func (a *MQTTClientAdapter) Connect() error {
	if token := a.client.Connect(); token.Wait() && token.Error() != nil {
		log.Println(" Error al conectar con el broker MQTT:", token.Error())
		return token.Error()
	}
	log.Println("Conectado al broker MQTT.")
	return nil
}

// Subscribe subscribes to the specified topic and uses the provided message handler.
func (a *MQTTClientAdapter) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
	if token := a.client.Subscribe(topic, qos, handler); token.Wait() && token.Error() != nil {
		log.Println(" Error al suscribirse al topic:", token.Error())
		return token.Error()
	}
	log.Println("Suscripción exitosa al topic:", topic)
	return nil
}