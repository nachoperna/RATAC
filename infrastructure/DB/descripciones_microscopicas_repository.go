package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
)

type Descripcion_microscopicasRepository struct {
	queries *sqlc.Queries
}

func NewDescripcion_microscopicasRepository (queries *sqlc.Queries) *Descripcion_microscopicasRepository {
	return &Descripcion_microscopicasRepository {
		queries: queries,
	}
}

func (r *Descripcion_microscopicasRepository) CreateDescripcionMicroscopica(ctx context.Context, protocolo string, descripcion domain.Descripcion_microscopicas) error {
	_, err := r.queries.CreateDescripcionMicroscopica(ctx, sqlc.CreateDescripcionMicroscopicaParams{
		Descripcion: descripcion.Descripcion,
		// Envolvemos el string en sql.NullString
		Diagnostico:        sql.NullString{String: descripcion.Diagnostico.Descripcion, Valid: true},
		PacientesProtocolo: protocolo,
	})
	return err
}

func (r *Descripcion_microscopicasRepository) GetDescripcionMicroscopicaByProtocolo(ctx context.Context, protocolo string) ([]domain.Descripcion_microscopicas, error) {
	bd_descripciones, err := r.queries.GetDescripcionByPaciente(ctx, protocolo)
	if err != nil {
		return nil, err
	}

	var descripcionesDomain []domain.Descripcion_microscopicas
	for _, bd_desc := range bd_descripciones {
		descripcionesDomain = append(descripcionesDomain, domain.Descripcion_microscopicas{
			Descripcion: bd_desc.Descripcion,
			Diagnostico: domain.Diagnostico{
				// Extraemos el string del sql.NullString usando .String
				Descripcion: bd_desc.Diagnostico.String,
				Imagenes:    nil,
			},
		})
	}

	return descripcionesDomain, nil
}

func (r *Descripcion_microscopicasRepository) DeleteDescripcionMicroscopica(ctx context.Context, descripcion string, protocolo string) error {
	err := r.queries.DeleteDescripcionMicroscopica(ctx, sqlc.DeleteDescripcionMicroscopicaParams{
		Descripcion:        descripcion,
		PacientesProtocolo: protocolo,
	})
	return err
}

func (r *Descripcion_microscopicasRepository) CountDiagnosticos(ctx context.Context) (int64, error) {
	return r.queries.CountDiagnosticos(ctx)
}
