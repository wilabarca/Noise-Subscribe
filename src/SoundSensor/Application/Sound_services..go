package application

import (
	"encoding/json"
	"log"

	entities "Noisesubscribe/src/SoundSensor/Domain/Entities"
	repositories "Noisesubscribe/src/SoundSensor/Domain/Repositories"
	adapters "Noisesubscribe/src/SoundSensor/Infraestructure/Adapters"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type SoundSensorService struct {
	repository  repositories.SoundSensorRepository
	mqttAdapter *adapters.MQTTClientAdapter
	apiURL      string
}

func NewSoundSensorService(mqttAdapter *adapters.MQTTClientAdapter, apiURL string) *SoundSensorService {
	return &SoundSensorService{
		repository:  adapters.NewSoundSensorRepositoryAdapter(apiURL),
		mqttAdapter: mqttAdapter,
		apiURL:      apiURL,
	}
}

// Start inicia la suscripción con el topic y la URL de la API para el sensor de sonido
func (service *SoundSensorService) Start(topic string, apiURL string) error {
	// Establecer la URL de la API si no se pasó antes
	if apiURL != "" {
		service.apiURL = apiURL
	}

	if err := service.mqttAdapter.Connect(); err != nil {
		log.Println("Error al conectar al broker MQTT:", err)
		return err
	}

	if err := service.mqttAdapter.Subscribe(topic, 0, service.messageHandler); err != nil {
		log.Println("Error al suscribirse al topic:", err)
		return err
	}

	log.Println("Suscripción exitosa al topic:", topic)
	return nil
}

func (service *SoundSensorService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Mensaje recibido: %s\n", msg.Payload())

	var soundData entities.SoundSensor
	if err := json.Unmarshal(msg.Payload(), &soundData); err != nil {
		log.Println("Error al parsear el mensaje:", err)
		return
	}

	// Filtro: Nivel de sonido > 70 dB
	if soundData.RuidoDB > 70 {
		// No pasamos el apiURL aquí porque el adaptador ya lo maneja internamente
		if err := service.repository.ProcessAndForward(soundData); err != nil {
			log.Println("Error al reenviar los datos:", err)
			return
		}
		log.Println("Datos enviados a la API Consumidora:", soundData)
	} else {
		log.Println("Nivel de sonido normal, ignorando...")
	}
}
