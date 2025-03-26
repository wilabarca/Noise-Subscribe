package adapters

import (
    mqtt "github.com/eclipse/paho.mqtt.golang"
    
)

type MQTTClientAdapter struct {
    client mqtt.Client
}

func NewMQTTClientAdapter(brokerURL string) *MQTTClientAdapter {
    opts := mqtt.NewClientOptions().AddBroker(brokerURL)
    return &MQTTClientAdapter{
        client: mqtt.NewClient(opts),
    }
}

func (a *MQTTClientAdapter) Connect() error {
    if token := a.client.Connect(); token.Wait() && token.Error() != nil {
        return token.Error()
    }
    return nil
}

func (a *MQTTClientAdapter) Subscribe(topic string, qos byte, handler mqtt.MessageHandler) error {
    if token := a.client.Subscribe(topic, qos, handler); token.Wait() && token.Error() != nil {
        return token.Error()
    }
    return nil
}