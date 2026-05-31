package main

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/domain"
	db "RATAC/infrastructure/DB"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

const db_conection = "host=db port=5432 user=admin password=password dbname=RATAC_DB sslmode=disable"
const RUTA_JSONS = "./JSONS/"
	
func main() {
	db, err := sql.Open("postgres", db_conection)
	if err != nil {
		log.Fatalf("Error al conectar con la Base de Datos: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)

	archivos, err := os.ReadDir(RUTA_JSONS)
	if err != nil {
		log.Fatalf("Error al leer la carpeta JSONS: %v", err)
	}

	for _, archivo := range archivos {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		func () {
			defer cancel()
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				fmt.Printf("\n[%s]: %v", archivo.Name(), err)
				return
			}
			defer tx.Rollback()
			qtx := queries.WithTx(tx)
			err = procesarJson(archivo, RUTA_JSONS, qtx, ctx)
			if err != nil {
				fmt.Printf("\n[%s]: %v",archivo.Name(), err)
				return
			}
			err = tx.Commit()
			if err != nil {
				fmt.Printf("\n[%s]: %v",archivo.Name(), errors.New("Error al hacer commit"))
				return
			}
		}()
	}
}

func procesarJson(archivo os.DirEntry, carpeta_jsons string, qtx *sqlc.Queries, ctx context.Context) error{
	var paciente domain.Paciente
	ruta := filepath.Join(carpeta_jsons, archivo.Name())
	contenido, err := os.Open(ruta)
	if err != nil {
		return errors.New("No se pudo abrir archivo")
	}
	defer contenido.Close()

	err = json.NewDecoder(contenido).Decode(&paciente)
	if err != nil {
		return errors.New("No se pudo decodificar JSON")
	}

	fecha_parseada, edad, err := db.TransformarDatos(paciente.Fecha, paciente.Edad)
	err = db.InsertarPaciente(paciente, fecha_parseada, edad, qtx, ctx)
	if err != nil{
		return err
	}
	for _, desc_micro := range paciente.Descripciones_microscopicas {
		err = db.InsertarDescMicro(desc_micro, paciente.Protocolo, qtx, ctx)
		if err != nil{
			return err
		}
	}
	return nil
}
