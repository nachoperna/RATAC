package main

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

type Handler struct {
	queries *sqlc.Queries
	ctx     context.Context
}

func NewHandler(queries *sqlc.Queries, ctx context.Context) *Handler {
	return &Handler{
		queries: queries,
		ctx:     ctx,
	}
}

func main() {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=admin password=password dbname=RATAC_DB sslmode=disable")
	if err != nil {
		log.Fatalf("Error al conectar con la Base de Datos: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)
	ctx := context.Background()

	handler := NewHandler(queries, ctx)

	carpeta_jsons := "./JSONS/"
	archivos, _ := os.ReadDir(carpeta_jsons)
	for _, archivo := range archivos {
		procesarJson(archivo, carpeta_jsons, handler)
	}

}

func procesarJson(archivo os.DirEntry, carpeta_jsons string, h *Handler) {
	var paciente domain.Paciente
	ruta := filepath.Join(carpeta_jsons, archivo.Name())
	contenido, _ := os.Open(ruta)
	err := json.NewDecoder(contenido).Decode(&paciente)
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}

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

	_, err = h.queries.CreatePaciente(h.ctx, sqlc.CreatePacienteParams{
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

	if err != nil {
		fmt.Printf("\nERROR al crear paciente: %s", err)
	}
}
