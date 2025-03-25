package entities

type Light struct {
    Intensidad int    `json:"intensidad"` // Valor de la intensidad de la luz
    Color      string `json:"color"`      // Color de la luz (por ejemplo, "rojo", "azul", etc.)
    Estado     bool   `json:"estado"`     // Estado de la luz (encendida o apagada)
}