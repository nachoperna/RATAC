package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
)

type DescripcionMicroscopicaRepository struct {
	queries *sqlc.Queries
}

func NewDescripcionMicroscopicaRepository(queries *sqlc.Queries) *DescripcionMicroscopicaRepository {
	return &DescripcionMicroscopicaRepository{
		queries: queries,
	}
}

func (r *DescripcionMicroscopicaRepository) CreateDescripcionMicroscopica(ctx context.Context, protocolo string, descripcion *domain.Descripcion_microscopicas) error {
	_, err := r.queries.CreateDescripcionMicroscopica(ctx, sqlc.CreateDescripcionMicroscopicaParams{
		Descripcion: descripcion.Descripcion,
		// Envolvemos el string en sql.NullString
		Diagnostico:        sql.NullString{String: descripcion.Diagnostico.Descripcion, Valid: true},
		PacientesProtocolo: protocolo,
	})
	return err
}

func (r *DescripcionMicroscopicaRepository) GetDescripcionMicroscopicaByProtocolo(ctx context.Context, protocolo string) ([]domain.Descripcion_microscopicas, error) {
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

func (r *DescripcionMicroscopicaRepository) DeleteDescripcionMicroscopica(ctx context.Context, descripcion string, protocolo string) error {
	err := r.queries.DeleteDescripcionMicroscopica(ctx, sqlc.DeleteDescripcionMicroscopicaParams{
		Descripcion:        descripcion,
		PacientesProtocolo: protocolo,
	})
	return err
}
