package db

import (
	"RATAC/domain"
	"bytes"
	"encoding/json"
	"errors"
)

type AdminRepository struct {}

func (r *AdminRepository) MapeoDocumento(contenido bytes.Buffer) (*domain.Paciente, error) {
	var paciente domain.Paciente
	err := json.NewDecoder(&contenido).Decode(&paciente)
	if err != nil {
		return nil, errors.New("No se pudo decodificar JSON")
	}
	return &paciente, nil
}
