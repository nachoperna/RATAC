package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

type PacienteRepository struct {
	queries *sqlc.Queries
}

func NewPacienteRepository(queries *sqlc.Queries) *PacienteRepository {
	return &PacienteRepository{
		queries: queries,
	}
}

func (r *PacienteRepository) CreatePaciente(ctx context.Context, paciente domain.Paciente) error  {
	layout := "02-01-2006"
	fecha_parseada, err := time.Parse(layout, paciente.Fecha)

	if err != nil {
		fmt.Printf("ERROR al parsear fecha: %s", err)
	}

	var edad int64
	if paciente.Edad != ""{
		edad, err = strconv.ParseInt(paciente.Edad, 10, 16)
		if err != nil {
			fmt.Printf("ERROR al parsear edad: %s", err)
		}
	}

	_, err = r.queries.CreatePaciente(ctx, sqlc.CreatePacienteParams{
		Protocolo:               paciente.Protocolo,
		Fecha:                   fecha_parseada,
		Solicitante:             paciente.Solicitante,
		Tecnica:                 paciente.Tecnica,
		Familia:                 sql.NullString{String: paciente.Familia, Valid: true},
		Especie:                 sql.NullString{String: paciente.Especie, Valid: true},
		Raza:                    sql.NullString{String: paciente.Raza, Valid: true}, //Linea 74
		Edad:                    sql.NullInt16{Int16: int16(edad), Valid: true},
		Paciente:                paciente.NombrePaciente,
		Antecedentes:            sql.NullString{String: paciente.Antecedentes, Valid: true},
		DescripcionMacroscopica: sql.NullString{String: paciente.DescripcionMacroscopica, Valid: true},
		ReferenciasMastocitomas: paciente.ReferenciasMastocitomas,
	})
	return err
}

func (r *PacienteRepository) ListPacientes(ctx context.Context) ([]domain.Paciente, error)  {
	bd_pacientes, err := r.queries.ListPacientes(ctx)
	if err != nil{
		return nil, err
	}
	
	var pacientes []domain.Paciente
	for _, p := range bd_pacientes {
		paciente := domain.Paciente{
			Protocolo: p.Protocolo,
			Fecha: p.Fecha.GoString(),
			Solicitante: p.Solicitante,
			Tecnica: p.Tecnica,
			Familia: p.Familia.String,
			Especie: p.Especie.String,
			Raza: p.Raza.String,
			Edad: strconv.Itoa(int(p.Edad.Int16)),
			NombrePaciente: p.Paciente,
			Antecedentes: p.Antecedentes.String,
			Descripciones_microscopicas: nil,
			DescripcionMacroscopica: p.DescripcionMacroscopica.String,
			ReferenciasMastocitomas: p.ReferenciasMastocitomas,
		}
		pacientes = append(pacientes, paciente)
	}
	return pacientes, nil
}
