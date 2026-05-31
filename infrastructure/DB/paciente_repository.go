package db

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Masterminds/squirrel"
)

type PacienteRepository struct {
	queries *sqlc.Queries
	db *sql.DB
}

func NewPacienteRepository(queries *sqlc.Queries, db *sql.DB) *PacienteRepository {
	return &PacienteRepository{
		queries: queries,
		db: db,
	}
}

func (r *PacienteRepository) InsertarDiagnostico(ctx context.Context, paciente domain.Paciente) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return errors.New("Error al iniciar transacción")
	}
	defer tx.Rollback()
	qtx := r.queries.WithTx(tx)
	err = procesarDiagnostico(qtx, ctx, paciente)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return errors.New("Error al finalizar transacción")
	}
	return nil
}

func procesarDiagnostico(qtx *sqlc.Queries, ctx context.Context, paciente domain.Paciente) error{
	fecha_parseada, edad, err := TransformarDatos(paciente.Fecha, paciente.Edad)
	err = InsertarPaciente(paciente, fecha_parseada, edad, qtx, ctx)
	if err != nil{
		return err
	}
	for _, desc_micro := range paciente.Descripciones_microscopicas {
		err = InsertarDescMicro(desc_micro, paciente.Protocolo, qtx, ctx)
		if err != nil{
			return err
		}
	}
	return nil
}

func getValueOrNil(campo sql.NullString) *string{
	if campo.Valid{
		return &campo.String
	}
	return nil
}

func (r *PacienteRepository) ListPacientes(ctx context.Context) ([]domain.Paciente, error) {
	bd_pacientes, err := r.queries.ListPacientes(ctx)
	if err != nil {
		return nil, err
	}

	var pacientes []domain.Paciente
	for _, p := range bd_pacientes {
		paciente := domain.Paciente{
			Protocolo:                   p.Protocolo,
			Fecha:                       p.Fecha.GoString(),
			Solicitante:                 p.Solicitante,
			Tecnica:                     p.Tecnica,
			Familia:                     getValueOrNil(p.Familia),
			Especie:                     getValueOrNil(p.Especie),
			Raza:                        getValueOrNil(p.Raza),
			Edad:                        func() *string {if p.Edad.Valid {edad := strconv.Itoa(int(p.Edad.Int16)); return &edad}; return nil} (),
			NombrePaciente:              p.Paciente,
			Antecedentes:                getValueOrNil(p.Antecedentes),
			Descripciones_microscopicas: nil,
			DescripcionMacroscopica:     getValueOrNil(p.DescripcionMacroscopica),
			ReferenciasMastocitomas:     p.ReferenciasMastocitomas,
		}
		pacientes = append(pacientes, paciente)
	}
	return pacientes, nil
}

func (r *PacienteRepository) ListUltimosPacientes(ctx context.Context) ([]domain.Paciente, error) {
	bd_pacientes, err := r.queries.ListUltimosPacientes(ctx)
	if err != nil {
		return nil, err
	}

	var pacientes []domain.Paciente
	for _, p := range bd_pacientes {
		pacientes = append(pacientes, domain.Paciente{
			Protocolo:      p.Protocolo,
			Fecha:          p.Fecha.Format("02-01-2006"), // Formateo de fecha
			NombrePaciente: p.Paciente,
			Solicitante:    p.Solicitante,
			// ... (mapear el resto de campos como en ListPacientes)
		})
	}
	return pacientes, nil
}

func (r *PacienteRepository) CountPacientes(ctx context.Context) (int64, error) {
	return r.queries.CountPacientes(ctx)
}

func (r *PacienteRepository) GetPacienteByNombre(ctx context.Context, nombre string, offset int8) ([]domain.Paciente, int16, error) {
	bd_pacientes, err := r.queries.GetPacienteByNombre(ctx, sqlc.GetPacienteByNombreParams{
		Column1: nombre,
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	var pacientes []domain.Paciente
	var resultados_total int16
	if len(bd_pacientes) > 0{
		resultados_total = int16(bd_pacientes[0].Total)
	}
	for _, p := range bd_pacientes {
		pacientes = append(pacientes, domain.Paciente{
			Protocolo:      p.Protocolo,
			Fecha:          p.Fecha.Format("02-01-2006"), // Formateo de fecha
			NombrePaciente: p.Paciente,
			Solicitante:    p.Solicitante,
			Tecnica: p.Tecnica,
			Familia: getValueOrNil(p.Familia),
			Especie: getValueOrNil(p.Especie),
			Raza: getValueOrNil(p.Raza),
			Edad: func() *string {if p.Edad.Valid {edad := strconv.Itoa(int(p.Edad.Int16)); return &edad}; return nil} (),
		})
	}
	return pacientes, resultados_total, nil
}

func (r *PacienteRepository) GetPacienteByFiltro(ctx context.Context, filtros []domain.Filtro, offset int8) ([]domain.Paciente, int16, error){
	sqlQuery, args, _ := getQueryByFiltro(filtros, offset)

	// 5. Ejecutar la consulta en PostgreSQL
	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Error ejecutando query en Postgres: %v", err)
		return nil, 0, err
	}
	defer rows.Close()

	var pacientes []domain.Paciente
	var resultados_total int16 
	set_resultado := false

	// 6. Leer los resultados (Ejemplo genérico)
	for rows.Next() {
		var paciente domain.Paciente
		var err error
		dummy := 0
		if !set_resultado{
			err = rows.Scan(&paciente.Protocolo, &paciente.Fecha, &paciente.NombrePaciente, &paciente.Solicitante, &paciente.Tecnica, &paciente.Edad, &paciente.Familia, &paciente.Raza, &paciente.Especie, &resultados_total);
			set_resultado = true
		}else{
			err = rows.Scan(&paciente.Protocolo, &paciente.Fecha, &paciente.NombrePaciente, &paciente.Solicitante, &paciente.Tecnica, &paciente.Edad, &paciente.Familia, &paciente.Raza, &paciente.Especie, &dummy);
		}
		if  err != nil {
			log.Printf("Error leyendo fila: %v", err)
			continue
		}
		pacientes = append(pacientes, paciente)
	}

	return pacientes, resultados_total, nil
}

func getQueryByFiltro(filtros []domain.Filtro, offset int8) (string, []any, error){
	// 1. CRÍTICO PARA POSTGRESQL: Configurar el formato del dólar ($1, $2)
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	
	// Iniciamos la consulta base
	query_builder := psql.Select("protocolo", "fecha", "paciente", "solicitante", "tecnica", "edad", "familia", "raza", "especie", "COUNT(*) OVER() as total").From("pacientes")

	var cond squirrel.Sqlizer

	// 2. Iterar sobre los filtros enviados desde el Frontend
	for i, f := range filtros {
		if !domain.ColumnasPermitidas[f.Campo] {
			log.Printf("Advertencia: Columna ignorada por seguridad: %s\n", f.Campo)
			continue
		}

		var expresiones squirrel.Sqlizer

		for _, valor := range f.Valores {
			expresion, err := getOperadorSQL(f.Operador, f.Campo, valor)
			if err != nil {
				log.Println(err)
				continue
			}
			if expresiones == nil{
				expresiones = expresion
			}else{
				expresiones = squirrel.Or{expresiones, expresion}
			}
		}

		if f.Not {
			expresion_sql, args, err := expresiones.ToSql()
			if err != nil {
				log.Println("Error al intentar negar la expresión:", err)
				continue
			}
			negacion_sql := fmt.Sprintf("NOT (%s)", expresion_sql)
			expresiones = squirrel.Expr(negacion_sql, args...)
		}

		// Acoplar la lógica AND / OR
		if i == 0 {
			cond = expresiones
		} else {
			if f.Logica == "OR" {
				cond = squirrel.Or{cond, expresiones}
			} else {
				cond = squirrel.And{cond, expresiones} // AND por defecto
			}
		}
	}

	// 3. Aplicar las condiciones generadas
	if cond != nil {
		query_builder = query_builder.Where(cond).Limit(domain.LIMIT_RESULTADOS_PACIENTE).Offset(uint64(offset))
	}

	// 4. Generar el SQL final compatible con Postgres
	sqlQuery, args, err := query_builder.ToSql()
	if err != nil {
		log.Printf("Error construyendo query: %v", err)
		return "", nil, err
	}

	fmt.Printf("\n--- QUERY PARA POSTGRESQL ---\n")
	fmt.Printf("SQL: %s\n", sqlQuery)
	fmt.Printf("Argumentos: %v\n-----------------------------\n", args)
	return sqlQuery, args, nil
}

func getOperadorSQL(operador, campo, valor string) (squirrel.Sqlizer, error){
	switch operador {
	case "igual":
		if campo != "Edad" {
			return squirrel.ILike{campo: valor}, nil
		}
		return squirrel.Eq{campo: valor}, nil
	case "mayor":
		return squirrel.Gt{campo: valor}, nil
	case "menor":
		return squirrel.Lt{campo: valor}, nil
	case "menorigual":
		return squirrel.LtOrEq{campo: valor}, nil
	case "mayorigual":
		return squirrel.GtOrEq{campo: valor}, nil
	}
	return nil, errors.New("No se pudo obtener operador squirrel")
}

func (r *PacienteRepository) GetAllFromPaciente(ctx context.Context, protocolo string) (*domain.Paciente, error) {
	datos_paciente, err := r.queries.GetPaciente(ctx, protocolo)
	if err != nil{
		return nil, err
		// Generar error
	}

	paciente_repo := setDatosBasePaciente(datos_paciente)
	
	descripciones_micro, err := r.queries.GetDescripcionByPaciente(ctx, protocolo)
	if err != nil{
		return nil, err
		// Generar error
	}
	var descripciones_micro_repo []domain.Descripcion_microscopicas
	for _, desc := range descripciones_micro {
		imagenes, err := r.queries.GetImagenesByDescripcion(ctx, sqlc.GetImagenesByDescripcionParams{
			DescripcionesMicroscopicasPacientesProtocolo: protocolo,
			DescripcionesMicroscopicasDescripcion: desc.Descripcion,
		})
		if err != nil { // Aca se deberia retornar simplemente el repositorio sin las imagenes cargadas
			return nil, nil
		}
		var img_rutas []string
		for _, imagen := range imagenes {
			img_rutas = append(img_rutas, imagen.Ruta)
		}
		diagnostico := domain.Diagnostico{
			Descripcion: getValueOrNil(desc.Diagnostico),
			Imagenes: img_rutas,
		}
		descripcion_micro := domain.Descripcion_microscopicas{
			Descripcion: desc.Descripcion,
			Diagnostico: diagnostico,
			TablaGrado: nil,
		}
		descripciones_micro_repo = append(descripciones_micro_repo, descripcion_micro)
	}
	paciente_repo.Descripciones_microscopicas = descripciones_micro_repo
	return &paciente_repo, nil
}

func setDatosBasePaciente(datos sqlc.Paciente) domain.Paciente {
	return domain.Paciente{
		Protocolo: datos.Protocolo,
		Fecha: datos.Fecha.Format("02-01-2006"),
		Solicitante: datos.Solicitante,
		Tecnica: datos.Tecnica,
		Familia: getValueOrNil(datos.Familia),
		Especie: getValueOrNil(datos.Especie),
		Raza: getValueOrNil(datos.Raza),
		Edad: func() *string {if datos.Edad.Valid {edad := strconv.Itoa(int(datos.Edad.Int16)); return &edad}; return nil} (),
		NombrePaciente: datos.Paciente,
		ReferenciasMastocitomas: datos.ReferenciasMastocitomas,
		Antecedentes: getValueOrNil(datos.Antecedentes),
		DescripcionMacroscopica: getValueOrNil(datos.DescripcionMacroscopica),
		Descripciones_microscopicas: nil,
	}
}
