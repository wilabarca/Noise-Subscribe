package adapters

import (
    "bytes"
    "encoding/json"
    "errors"
    "log"
    "net/http"

    entities "Noisesubscribe/src/AirQuality/Domain/Entities"
)

type AirQualityRepositoryAdapter struct {
    ApiURL string
}

func NewAirQualityRepositoryAdapter(apiURL string) *AirQualityRepositoryAdapter {
    return &AirQualityRepositoryAdapter{
        ApiURL: apiURL,
    }
}

func (adapter *AirQualityRepositoryAdapter) ProcessAndForward(airData entities.AirQuality) error {
    data, err := json.Marshal(airData)
    if err != nil {
        log.Println("❌ Error al convertir los datos a JSON:", err)
        return err
    }

    resp, err := http.Post(adapter.ApiURL, "application/json", bytes.NewBuffer(data))
    if err != nil {
        log.Println("❌ Error al enviar los datos a la API 2:", err)
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return errors.New("respuesta no exitosa: " + resp.Status)
    }

    return nil
}