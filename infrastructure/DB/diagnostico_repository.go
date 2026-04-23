package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
)

type DiagnosticoRepository struct {
	queries *sqlc.Queries
}

func NewDiagnosticoRepository(queries *sqlc.Queries) *DiagnosticoRepository {
	return &DiagnosticoRepository{
		queries: queries,
	}
}

// CreateDiagnostico actualiza el texto en la descripción microscópica y guarda las imágenes asociadas
func (r *DiagnosticoRepository) CreateDiagnostico(ctx context.Context, protocolo string, descripcionMicro string, diagnostico *domain.Diagnostico) error {
	// 1. Actualizamos el texto del diagnóstico en la tabla Descripciones_microscopicas
	_, err := r.queries.UpdateDescripcionMicroscopica(ctx, sqlc.UpdateDescripcionMicroscopicaParams{
		Descripcion:        descripcionMicro,
		PacientesProtocolo: protocolo,
		Diagnostico:        sql.NullString{String: diagnostico.Descripcion, Valid: true},
	})
	if err != nil {
		return err
	}

	// 2. Insertamos cada imagen en la tabla Imagenes
	for _, imgRuta := range diagnostico.Imagenes {
		_, err := r.queries.CreateImagen(ctx, sqlc.CreateImagenParams{
			Ruta:                                  imgRuta,
			DescripcionesMicroscopicasDescripcion: descripcionMicro,
			DescripcionesMicroscopicasPacientesProtocolo: protocolo,
		})
		if err != nil {
			return err // Si falla una, retornamos el error
		}
	}

	return nil
}

// GetDiagnostico obtiene el texto y las imágenes de un diagnóstico específico
func (r *DiagnosticoRepository) GetDiagnostico(ctx context.Context, protocolo string, descripcionMicro string) (*domain.Diagnostico, error) {
	// Obtenemos las imágenes
	bdImagenes, err := r.queries.GetImagenesByDescripcion(ctx, sqlc.GetImagenesByDescripcionParams{
		DescripcionesMicroscopicasDescripcion:        descripcionMicro,
		DescripcionesMicroscopicasPacientesProtocolo: protocolo,
	})
	if err != nil {
		return nil, err
	}

	var rutas []string
	for _, img := range bdImagenes {
		rutas = append(rutas, img.Ruta)
	}

	// Nota: El texto del diagnóstico ya suele venir cuando hacés el Get de la Descripción Microscópica,
	// pero acá armamos la entidad de dominio con las imágenes correspondientes.
	return &domain.Diagnostico{
		// Si necesitas el texto acá, tendrías que hacer una consulta extra a la tabla Descripciones_microscopicas
		Imagenes: rutas,
	}, nil
}

// DeleteImagen borra una imagen específica por su ruta
func (r *DiagnosticoRepository) DeleteImagen(ctx context.Context, ruta string) error {
	return r.queries.DeleteImagen(ctx, ruta)
}
