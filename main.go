package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	th_app "Noisesubscribe/src/TemperatureHumiditySensor/Application"
	th_adapters "Noisesubscribe/src/TemperatureHumiditySensor/Infraestructure/Adapters"
	ss_app "Noisesubscribe/src/SoundSensor/Application"
	ss_adapters "Noisesubscribe/src/SoundSensor/Infraestructure/Adapters"
)

func main() {
	// Configuración común
	brokerURL := "tcp://broker.emqx.io:1883" // Broker MQTT público de prueba
	

	// Configuración para el sensor de Temperatura y Humedad
	thTopic := "casa/sensor/temperatura_humedad"
	

	// Configuración para el sensor de Sonido
	ssTopic := "casa/sensor/sonido"
	

	// Inicializar sensor de Temperatura y Humedad
	thMQTT := th_adapters.NewMQTTClientAdapter(brokerURL)
	thService := th_app.NewTemperatureHumidityService(thMQTT)

	// Inicializar sensor de Sonido
	ssMQTT := ss_adapters.NewMQTTClientAdapter(brokerURL)
	ssService := ss_app.NewSoundSensorService(ssMQTT)

	// Iniciar ambos sensores
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

	log.Println("Ambos sensores están operativos y suscritos a sus topics")

	// Mantener la aplicación corriendo
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Apagando sensores...")
}