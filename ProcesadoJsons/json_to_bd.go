package main

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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

const db_conection = "host=localhost port=5432 user=admin password=password dbname=RATAC_DB sslmode=disable"

func main() {
	db, err := sql.Open("postgres", db_conection)
	if err != nil {
		log.Fatalf("Error al conectar con la Base de Datos: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)

	carpeta_jsons := "./JSONS/"
	archivos, _ := os.ReadDir(carpeta_jsons)
	for _, archivo := range archivos {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		func ()  {
			defer cancel()
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				fmt.Println(errors.New("Error al iniciar transaccion: "), err)
			}
			qtx := queries.WithTx(tx)
			defer tx.Rollback()
			err = procesarJson(archivo, carpeta_jsons, qtx, ctx)
			if err != nil {
				fmt.Printf("\n[%s]: %v",archivo.Name(), err)
			}
			tx.Commit()
		}()
	}
}

func procesarJson(archivo os.DirEntry, carpeta_jsons string, qtx *sqlc.Queries, ctx context.Context) error{
	var paciente domain.Paciente
	ruta := filepath.Join(carpeta_jsons, archivo.Name())
	contenido, _ := os.Open(ruta)
	err := json.NewDecoder(contenido).Decode(&paciente)
	if err != nil {
		return errors.New("No se pudo decodificar JSON: ")
	}

	fecha_parseada, edad, err := transformarDatos(paciente.Fecha, paciente.Edad)
	err = insertarPaciente(paciente, fecha_parseada, edad, qtx, ctx)
	if err != nil{
		return err
	}
	for _, desc_micro := range paciente.Descripciones_microscopicas {
		err = insertarDescMicro(desc_micro, paciente.Protocolo, qtx, ctx)
		if err != nil{
			return err
		}
	}
	return nil
}

func transformarDatos(fecha,edad string) (time.Time, int64, error){
	layout := "2-1-2006"
	fecha_parseada, err := time.Parse(layout, fecha)

	if err != nil {
		err = errors.New("Error al parsear fecha")
	}

	var edad_ int64
	if edad != ""{
		edad_, err = strconv.ParseInt(edad, 10, 16)
		if err != nil {
			err = errors.New("Error al convertir edad")
		}
	}
	return fecha_parseada, edad_, err
}

func insertarPaciente(paciente domain.Paciente, fecha_parseada time.Time, edad int64, qtx *sqlc.Queries, ctx context.Context) error{
	_, err := qtx.CreatePaciente(ctx, sqlc.CreatePacienteParams{
		Protocolo:               paciente.Protocolo,
		Fecha:                   fecha_parseada,
		Solicitante:             paciente.Solicitante,
		Tecnica:                 paciente.Tecnica,
		Familia:                 sql.NullString{String: paciente.Familia, Valid: true},
		Especie:                 sql.NullString{String: paciente.Especie, Valid: true},
		Raza:                    sql.NullString{String: paciente.Raza, Valid: true},
		Edad:                    sql.NullInt16{Int16: int16(edad), Valid: true},
		Paciente:                paciente.NombrePaciente,
		Antecedentes:            sql.NullString{String: paciente.Antecedentes, Valid: true},
		DescripcionMacroscopica: sql.NullString{String: paciente.DescripcionMacroscopica, Valid: true},
		ReferenciasMastocitomas: paciente.ReferenciasMastocitomas,
	})
	
	return err
}

func insertarDescMicro(desc domain.Descripcion_microscopicas, pk string, qtx *sqlc.Queries, ctx context.Context) error{
	_, err := qtx.CreateDescripcionMicroscopica(ctx, sqlc.CreateDescripcionMicroscopicaParams{
		Descripcion: desc.Descripcion,
		Diagnostico: sql.NullString{String: desc.Diagnostico.Descripcion, Valid: true},
		PacientesProtocolo: pk,
	})
	if err != nil {
		err = errors.New("Error al insertar descripcion microscopica")
	}
	err = insertarImagenes(desc.Diagnostico.Imagenes, desc.Descripcion, pk, qtx, ctx)
	return err
}

func insertarImagenes(imagenes []string, pk1, pk2 string, qtx *sqlc.Queries, ctx context.Context) error{
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
