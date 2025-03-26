package repositories

import entities "Noisesubscribe/src/TemperatureHumiditySensor/Domain/Entities"

// TemperatureHumidityRepository define la interfaz para procesar y reenviar los datos de temperatura y humedad
type TemperatureHumidityRepository interface {
    ProcessAndForward(tempHumidityData entities.TemperatureHumidity) error
}
