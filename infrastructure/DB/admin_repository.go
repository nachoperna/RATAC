package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"bytes"
	"context"
	"encoding/json"
	"errors"
)

type AdminRepository struct {
	queries *sqlc.Queries
}

func NewAdminRepository(queries *sqlc.Queries) *AdminRepository {
	return &AdminRepository { queries: queries }
}

func (r *AdminRepository) MapeoDocumento(contenido bytes.Buffer) (*domain.Paciente, error) {
	var paciente domain.Paciente
	err := json.NewDecoder(&contenido).Decode(&paciente)
	if err != nil {
		return nil, errors.New("No se pudo decodificar JSON")
	}
	return &paciente, nil
}

func (r *AdminRepository) PacienteYaRegistrado(ctx context.Context, protocolo string) bool {
	if protocolo == "" { // sin protocolo no hay nada que buscar; evita falsos positivos
		return false
	}
	_, err := r.queries.GetPaciente(ctx, protocolo)
	if err != nil { // ErrNoRows o cualquier otro error => no podemos afirmar que ya existe
		return false
	}
	return true
}
