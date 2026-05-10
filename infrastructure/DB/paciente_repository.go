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
	"time"

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

func (r *PacienteRepository) CreatePaciente(ctx context.Context, paciente domain.Paciente) error {
	layout := "02-01-2006"
	fecha_parseada, err := time.Parse(layout, paciente.Fecha)

	if err != nil {
		fmt.Printf("ERROR al parsear fecha: %s", err)
	}

	var edad int64
	if paciente.Edad != "" {
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
			Familia:                     p.Familia.String,
			Especie:                     p.Especie.String,
			Raza:                        p.Raza.String,
			Edad:                        strconv.Itoa(int(p.Edad.Int16)),
			NombrePaciente:              p.Paciente,
			Antecedentes:                p.Antecedentes.String,
			Descripciones_microscopicas: nil,
			DescripcionMacroscopica:     p.DescripcionMacroscopica.String,
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

func (r *PacienteRepository) GetPacienteByNombre(ctx context.Context, nombre string) ([]domain.Paciente, error) {
	bd_pacientes, err := r.queries.GetPacienteByNombre(ctx, nombre)
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
			Tecnica: p.Tecnica,
			Familia: p.Familia.String,
			Especie: p.Especie.String,
			Raza: p.Raza.String,
			Edad: strconv.FormatInt(int64(p.Edad.Int16),10),
		})
	}
	return pacientes, nil
}

func (r *PacienteRepository) GetPacienteByFiltro(ctx context.Context, filtros []domain.Filtro) ([]domain.Paciente, error){
	sqlQuery, args, _ := getQueryByFiltro(filtros)

	// 5. Ejecutar la consulta en PostgreSQL
	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		log.Printf("Error ejecutando query en Postgres: %v", err)
		return nil, err
	}
	defer rows.Close()

	var pacientes []domain.Paciente
	// 6. Leer los resultados (Ejemplo genérico)
	for rows.Next() {
		var paciente domain.Paciente
		
		if err := rows.Scan(&paciente.Protocolo, &paciente.Fecha, &paciente.NombrePaciente, &paciente.Solicitante, &paciente.Tecnica, &paciente.Edad, &paciente.NombrePaciente, &paciente.Raza, &paciente.Especie); err != nil {
			log.Printf("Error leyendo fila: %v", err)
			continue
		}
		pacientes = append(pacientes, paciente)
	}

	return pacientes, nil
}

func getQueryByFiltro(filtros []domain.Filtro) (string, []interface{}, error){
	// 1. CRÍTICO PARA POSTGRESQL: Configurar el formato del dólar ($1, $2)
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	
	// Iniciamos la consulta base
	query_builder := psql.Select("protocolo", "fecha", "paciente", "solicitante", "tecnica", "edad", "familia", "raza", "especie").From("pacientes")

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
		query_builder = query_builder.Where(cond)
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
