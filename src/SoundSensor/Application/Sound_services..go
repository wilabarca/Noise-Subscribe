package application

import (
	entities "Noisesubscribe/src/SoundSensor/Domain/Entities"
	repositories "Noisesubscribe/src/SoundSensor/Domain/Repositories"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SoundSensorService struct {
	repository repositories.SoundSensorRepository
}

// NewSoundSensorService crea una nueva instancia del servicio de sensor de sonido
func NewSoundSensorService(repo repositories.SoundSensorRepository) *SoundSensorService {
	return &SoundSensorService{repository: repo}
}

// messageHandler procesa los mensajes recibidos del broker MQTT
func (s *SoundSensorService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	var sensorData entities.SoundSensor
	if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
		log.Printf("Error al deserializar el mensaje: %v\n", err)
		return
	}

	if sensorData.RuidoDB >= 40 {
		log.Printf("Datos relevantes recibidos: %+v\n", sensorData)
		if err := s.repository.ProcessAndForward(sensorData); err != nil {
			log.Printf("Error al reenviar los datos a la API 2: %v\n", err)
		}
	} else {
		log.Printf("Datos ignorados para el sensor de sonido: %+v\n", sensorData)
	}
}

// SubscribeToMQTT inicia la suscripción al broker MQTT
func (s *SoundSensorService) SubscribeToMQTT(topic string) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("soundSensorSubscriber")
	opts.SetDefaultPublishHandler(s.messageHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	select {}
}

// SendToAPI2 reenvía los datos del sensor a la API 2
func (s *SoundSensorService) SendToAPI2(sensorData entities.SoundSensor) error {
	url := "http://api2-url.com/endpoint"

	payload, err := json.Marshal(sensorData)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	log.Printf("Datos de sonido reenviados correctamente a la API 2: %+v\n", sensorData)
	return nil
}
