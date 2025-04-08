package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	aq_app "Noisesubscribe/src/AirQuality/Application"
	aq_adapters "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"
	light_app "Noisesubscribe/src/Light/Application"
	light_adapters "Noisesubscribe/src/Light/Infraestructure/Adapters"
	ss_app "Noisesubscribe/src/SoundSensor/Application"
	ss_adapters "Noisesubscribe/src/SoundSensor/Infraestructure/Adapters"
	th_app "Noisesubscribe/src/TemperatureHumiditySensor/Application"
	th_adapters "Noisesubscribe/src/TemperatureHumiditySensor/Infraestructure/Adapters"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando archivo .env: %v", err)
	}

	brokerURL := os.Getenv("BROKER_URL")
	apiURL := os.Getenv("API_URL")
	amqpURL := os.Getenv("AMQP_URL")

	temperatureTopic := "sensor.temperature"
	soundTopic := "sensor.sound"
	airQualityTopic := "sensor.air"
	lightTopic := "sensor.light"

	thRabbitAdapter, err := th_adapters.NewRabbitMQAdapter(amqpURL)
	if err != nil {
		log.Fatalf("Error al conectar RabbitMQ para Temperatura/Humedad: %v", err)
	}
	defer thRabbitAdapter.Close()

	ssRabbitAdapter, err := ss_adapters.NewRabbitMQAdapter(amqpURL)
	if err != nil {
		log.Fatalf("Error al conectar RabbitMQ para Sonido: %v", err)
	}
	defer ssRabbitAdapter.Close()

	aqRabbitAdapter, err := aq_adapters.NewRabbitMQAdapter(amqpURL)
	if err != nil {
		log.Fatalf("Error al conectar RabbitMQ para Calidad del Aire: %v", err)
	}
	defer aqRabbitAdapter.Close()

	lightRabbitAdapter, err := light_adapters.NewRabbitMQAdapter(amqpURL)
	if err != nil {
		log.Fatalf("Error al conectar RabbitMQ para Luz: %v", err)
	}
	defer lightRabbitAdapter.Close()

	thMQTT := th_adapters.NewMQTTClientAdapter(brokerURL)
	thService := th_app.NewTemperatureHumidityService(thMQTT, apiURL, thRabbitAdapter)
	ssMQTT := ss_adapters.NewMQTTClientAdapter(brokerURL)
	ssService := ss_app.NewSoundSensorService(ssMQTT, apiURL, ssRabbitAdapter)

	aqMQTT := aq_adapters.NewMQTTClientAdapter(brokerURL)
	aqService := aq_app.NewAirQualityService(aqMQTT, apiURL, aqRabbitAdapter)

	lightRepo := light_adapters.NewLightAPIRepository(apiURL)

	lightService := light_app.NewLightService(
		lightRepo,
		lightRabbitAdapter,
		100.0,
		800.0,
	)

	go func() {
		if err := thService.Start(temperatureTopic, ""); err != nil {
			log.Fatalf("Error iniciando sensor Temperatura/Humedad🌡️: %v", err)
		}
	}()

	go func() {
		if err := ssService.Start(soundTopic, ""); err != nil {
			log.Fatalf("Error iniciando sensor de Sonido: %v", err)
		}
	}()

	go func() {
		if err := aqService.Start(airQualityTopic, ""); err != nil {
			log.Fatalf("Error iniciando sensor de Calidad del Aire: %v", err)
		}
	}()

	go func() {
		if err := lightService.Start(lightTopic, ""); err != nil {
			log.Fatalf("Error iniciando sensor de Luz: %v", err)
		}
	}()

	log.Println("Todos los sensores están operativos y suscritos a sus topics")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Apagando sensores...")
}
