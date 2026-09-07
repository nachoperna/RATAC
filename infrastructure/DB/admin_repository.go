package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"bytes"
	"context"
	"database/sql"
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
	_, err := r.queries.GetPaciente(ctx, protocolo)
	if err == sql.ErrNoRows{
		return false
	}
	return true
}

func (r  *AdminRepository) GetUltimosDiagnosticosCargados(ctx context.Context, offset int8, token string) ([]domain.Paciente, int16, error) {
	pacientes_rows, err := r.queries.ListPacientesByLab(ctx, sqlc.ListPacientesByLabParams{
		Offset: int32(offset),
		Token: token,
	})
	if err != nil {
		return nil, 0, err
	}
	var pacientes []domain.Paciente
	var total int16 = 0
	if len(pacientes_rows) > 0 {
		total = int16(pacientes_rows[0].Total)
	}
	for _, p := range pacientes_rows {
		paciente := pacienteByLabORM(p)
		pacientes = append(pacientes, paciente)
	}
	return pacientes, total, nil
}
