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
