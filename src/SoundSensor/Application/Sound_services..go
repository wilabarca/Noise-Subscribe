package application

import (
	entities "Noisesubscribe/src/SoundSensor/Domain/Entities"
	repositories "Noisesubscribe/src/SoundSensor/Domain/Repositories"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
)

// SoundSensorService es el servicio que maneja la lógica de negocio del sensor de sonido
type SoundSensorService struct {
	repository repositories.SoundSensorRepository
}

// NewSoundSensorService crea una nueva instancia del servicio de sensor de sonido
func NewSoundSensorService(repo repositories.SoundSensorRepository) *SoundSensorService {
	return &SoundSensorService{repository: repo}
}

// Handler para los mensajes recibidos del broker MQTT
func (s *SoundSensorService) messageHandler(client mqtt.Client, msg mqtt.Message) {
	// Deserializamos el mensaje recibido en la entidad SoundSensor
	var sensorData entities.SoundSensor
	if err := json.Unmarshal(msg.Payload(), &sensorData); err != nil {
		log.Printf("Error al deserializar el mensaje: %v\n", err)
		return
	}

	// Filtramos los datos: solo si el ruido es mayor o igual a 40 dB
	if sensorData.RuidoDB >= 40 {
		log.Printf("Datos relevantes recibidos: %+v\n", sensorData)
		// Llamamos al repositorio para reenviar los datos a la API 2
		err := s.repository.ProcessAndForward(sensorData)
		if err != nil {
			log.Printf("Error al reenviar los datos a la API 2: %v\n", err)
		}
	} else {
		// Si los datos no son relevantes, los ignoramos
		log.Printf("Datos ignorados para el sensor de sonido: %+v\n", sensorData)
	}
}

// Función que inicia la suscripción del servicio al broker MQTT
func (s *SoundSensorService) SubscribeToMQTT() {
	// Definimos las opciones del cliente MQTT
	opts := mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("soundSensorSubscriber")
	opts.SetDefaultPublishHandler(s.messageHandler)

	// Creamos el cliente MQTT
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Nos suscribimos al tópico de sensores de sonido
	if token := client.Subscribe("sensors/sound/#", 0, nil); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	// Mantiene la conexión abierta
	select {}
}

// Función que reenvía los datos de sonido a la API 2 usando HTTP
// SendToAPI2 reenvía los datos del sensor a la API 2
func SendToAPI2(sensorData entities.SoundSensor) error {
	// URL de la API 2
	url := "http://api2-url.com/endpoint" // URL de la API 2

	// Convertimos los datos de sensor a JSON
	payload, err := json.Marshal(sensorData)
	if err != nil {
		return err
	}

	// Enviamos los datos a la API 2
	resp, err := http.Post(url, "application/json", strings.NewReader(string(payload)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Verificamos la respuesta de la API 2
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error al reenviar los datos a la API 2: %v", resp.Status)
	}

	log.Printf("Datos de sonido reenviados correctamente a la API 2: %+v\n", sensorData)
	return nil
}
