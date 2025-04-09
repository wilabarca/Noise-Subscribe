package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"Noisesubscribe/src/Light/Domain/Entities"
	"Noisesubscribe/src/Light/Domain/Repositories"
)

type LightAPIRepository struct {
	apiURL string
}

func NewLightAPIRepository(apiURL string) repositories.LightRepository {
	return &LightAPIRepository{apiURL: apiURL}
}

func (r *LightAPIRepository) ProcessAndForward(lightData entities.Light) error {
	jsonData, err := json.Marshal(lightData)
	if err != nil {
		return errors.New("error al serializar los datos de luz: " + err.Error())
	}

	resp, err := http.Post(r.apiURL+"/luz", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return errors.New("error al enviar datos de luz a la API: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("la API respondió con un error: " + resp.Status)
	}

	log.Println("✅ Datos de luz enviados correctamente")
	return nil
}

func (r *LightAPIRepository) GetLightData() (entities.Light, error) {
	return entities.Light{}, nil
}
