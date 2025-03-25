package adapters

import (
    "github.com/eclipse/paho.mqtt.golang"
    "log"
)

// MQTTClientAdapter es un adaptador para manejar la conexión y suscripción a MQTT
type MQTTClientAdapter struct {
    client mqtt.Client
}

// NewMQTTClientAdapter crea una nueva instancia del adaptador MQTT
func NewMQTTClientAdapter(brokerURL string) *MQTTClientAdapter {
    opts := mqtt.NewClientOptions().AddBroker(brokerURL)
    client := mqtt.NewClient(opts)

    return &MQTTClientAdapter{
        client: client,
    }
}

// Connect establece la conexión con el broker MQTT
func (adapter *MQTTClientAdapter) Connect() error {
    if token := adapter.client.Connect(); token.Wait() && token.Error() != nil {
        log.Println("❌ Error al conectar al broker MQTT:", token.Error())
        return token.Error()
    }
    log.Println("✅ Conectado al broker MQTT exitosamente.")
    return nil
}

// Subscribe se suscribe a un topic de MQTT
func (adapter *MQTTClientAdapter) Subscribe(topic string, qos byte, messageHandler mqtt.MessageHandler) error {
    if token := adapter.client.Subscribe(topic, qos, messageHandler); token.Wait() && token.Error() != nil {
        log.Println("❌ Error al suscribirse al topic:", token.Error())
        return token.Error()
    }
    log.Println("✅ Suscripción exitosa al topic:", topic)
    return nil
}
