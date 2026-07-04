package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"
)

func TransformarDatos(fecha string, edad *string) (time.Time, *int64, error){
	layout := "2-1-2006"
	fecha_parseada, err := time.Parse(layout, fecha)

	if err != nil {
		return fecha_parseada, nil, errors.New("Error al parsear fecha")
	}

	var edad_ *int64
	if edad != nil{
		var aux int64
		aux, err = strconv.ParseInt(*edad, 10, 16)
		if err != nil {
			err = errors.New("Error al convertir edad")
		}else{
			edad_ = &aux
		}
	}
	return fecha_parseada, edad_, err
}

func InsertarPaciente(paciente domain.Paciente, fecha_parseada time.Time, edad *int64, qtx *sqlc.Queries, ctx context.Context) error{
	_, err := qtx.CreatePaciente(ctx, sqlc.CreatePacienteParams{
		Protocolo:               paciente.Protocolo,
		Fecha:                   fecha_parseada,
		Solicitante:             paciente.Solicitante,
		Tecnica:                 paciente.Tecnica,
		Familia:                 SetvalueOrNull(paciente.Familia),
		Especie:                 SetvalueOrNull(paciente.Especie),
		Raza:                    SetvalueOrNull(paciente.Raza),
		Edad:                    func() sql.NullInt16 {if edad != nil {return sql.NullInt16{Int16: int16(*edad), Valid: true}} else {return sql.NullInt16{Valid: false}}}() ,
		Paciente:                paciente.NombrePaciente,
		Antecedentes:            SetvalueOrNull(paciente.Antecedentes),
		DescripcionMacroscopica: SetvalueOrNull(paciente.DescripcionMacroscopica),
		ReferenciasMastocitomas: paciente.ReferenciasMastocitomas,
	})
	if err != nil {
		return errors.New("Error al insertar paciente")
	}
	return nil
}

func SetvalueOrNull(campo *string) sql.NullString{
	if campo != nil{
		return sql.NullString{String: *campo, Valid: true}
	}
	return sql.NullString{Valid: false}
}

func InsertarDescMicro(desc domain.Descripcion_microscopicas, pk string, qtx *sqlc.Queries, ctx context.Context) error{
	_, err := qtx.CreateDescripcionMicroscopica(ctx, sqlc.CreateDescripcionMicroscopicaParams{
		Descripcion: desc.Descripcion,
		Diagnostico: SetvalueOrNull(desc.Diagnostico.Descripcion),
		PacientesProtocolo: pk,
	})
	if err != nil {
		return errors.New("Error al insertar descripcion microscopica")
	}
	err = InsertarImagenes(desc.Diagnostico.Imagenes, desc.Descripcion, pk, qtx, ctx)
	if err != nil {
		return err
	}
	err = InsertarTablaGrado(desc.TablaGrado, desc.Descripcion, pk, qtx, ctx)
	return err
}

func InsertarImagenes(imagenes []string, pk1, pk2 string, qtx *sqlc.Queries, ctx context.Context) error{
	for _, imagen := range imagenes {
		_, err := qtx.CreateImagen(ctx, sqlc.CreateImagenParams{
			Ruta: imagen,
			DescripcionesMicroscopicasDescripcion: pk1,
			DescripcionesMicroscopicasPacientesProtocolo: pk2,
		})
		if err != nil {
			return errors.New("Error al insertar ruta de imagen")
		}
	}
	return nil
}

func InsertarTablaGrado(tabla []domain.Grado_oncologico, pk1, pk2 string, qtx *sqlc.Queries, ctx context.Context) error{
	for _, fila := range tabla {
		puntaje, err := strconv.Atoi(fila.Puntaje)
		if err != nil {
			return errors.New("Error parseando puntaje de tabla de grado oncologico")
		}
		_, err = qtx.CreateGradoOncologico(ctx, sqlc.CreateGradoOncologicoParams{
			Caracteristica: fila.Caracteristica,
			MuestraAnalizada: sql.NullString{String: fila.Muestra_analizada, Valid: true},
			Puntaje: int16(puntaje),
			DescripcionesMicroscopicasDescripcion: pk1,
			DescripcionesMicroscopicasPacientesProtocolo: pk2,
		})
		if err != nil {
			return errors.New("Error al insertar fila de grado oncologico")
		}
	}
	return nil
}
