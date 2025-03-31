package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	th_app "Noisesubscribe/src/TemperatureHumiditySensor/Application"
	th_adapters "Noisesubscribe/src/TemperatureHumiditySensor/Infraestructure/Adapters"
	th_rabbit "Noisesubscribe/src/TemperatureHumiditySensor/Infraestructure/Adapters"

	ss_app "Noisesubscribe/src/SoundSensor/Application"
	ss_adapters "Noisesubscribe/src/SoundSensor/Infraestructure/Adapters"
	ss_rabbit "Noisesubscribe/src/SoundSensor/Infraestructure/Adapters"

	aq_app "Noisesubscribe/src/AirQuality/Application"
	aq_adapters "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"
	aq_rabbit "Noisesubscribe/src/AirQuality/Infraestructure/Adapters"
)

func main() {
	// Cargar variables de entorno desde .env
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando archivo .env: %v", err)
	}

	// Configuración obtenida desde .env
	brokerURL := os.Getenv("BROKER_URL")
	apiURL := os.Getenv("API_URL")
	amqpURL := os.Getenv("AMQP_URL")

	// Configuración de topics (dejados en el código)
	thTopic := "casa/sensor/temperatura_humedad"
	ssTopic := "casa/sensor/sonido"
	aqTopic := "casa/sensor/airquality"

	// Inicializar adaptadores de RabbitMQ por cada sensor
	thRabbitAdapter, err := th_rabbit.NewRabbitMQAdapter(amqpURL)
	if err != nil {
		log.Fatalf("Error al conectar RabbitMQ para Temperatura/Humedad: %v", err)
	}
	defer thRabbitAdapter.Close()

	ssRabbitAdapter, err := ss_rabbit.NewRabbitMQAdapter(amqpURL)
	if err != nil {
		log.Fatalf("Error al conectar RabbitMQ para Sonido: %v", err)
	}
	defer ssRabbitAdapter.Close()

	aqRabbitAdapter, err := aq_rabbit.NewRabbitMQAdapter(amqpURL)
	if err != nil {
		log.Fatalf("Error al conectar RabbitMQ para Calidad del Aire: %v", err)
	}
	defer aqRabbitAdapter.Close()

	// Inicializar sensor de Temperatura y Humedad
	thMQTT := th_adapters.NewMQTTClientAdapter(brokerURL)
	thService := th_app.NewTemperatureHumidityService(thMQTT, apiURL, thRabbitAdapter)

	// Inicializar sensor de Sonido
	ssMQTT := ss_adapters.NewMQTTClientAdapter(brokerURL)
	ssService := ss_app.NewSoundSensorService(ssMQTT, apiURL, ssRabbitAdapter)

	// Inicializar sensor de Calidad del Aire
	aqMQTT := aq_adapters.NewMQTTClientAdapter(brokerURL)
	aqService := aq_app.NewAirQualityService(aqMQTT, apiURL, aqRabbitAdapter)

	// Iniciar cada sensor en goroutines
	go func() {
		if err := thService.Start(thTopic); err != nil {
			log.Fatalf("Error iniciando sensor Temperatura/Humedad: %v", err)
		}
	}()

	go func() {
		if err := ssService.Start(ssTopic, ""); err != nil {
			log.Fatalf("Error iniciando sensor de Sonido: %v", err)
		}
	}()

	go func() {
		if err := aqService.Start(aqTopic, ""); err != nil {
			log.Fatalf("Error iniciando sensor de Calidad del Aire: %v", err)
		}
	}()

	log.Println("Todos los sensores están operativos y suscritos a sus topics")

	// Mantener la aplicación corriendo hasta recibir una señal de terminación
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Apagando sensores...")
}
