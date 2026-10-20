package domain

import "context"

type Grado_oncologico struct {
	Caracteristica	string	`json:"Caracteristica"`
	Muestra_analizada	string	`json:"Muestra analizada"`
	Puntaje		string	`json:"Puntaje"`
}

type Grado_oncologicoRepository interface {
	CreateTablaGrado(ctx context.Context, tabla []Grado_oncologico) error
}

func GetGrado(puntaje uint8) uint8 {
	switch {
		case puntaje <= 5:
			return 1
		case puntaje <= 7:
			return 2
		default:
			return 3
	}
}
