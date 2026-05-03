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

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://admin:password@localhost:5432/RATAC_DB?sslmode=disable"
	} else {
		dbURL = dbURL + "?sslmode=disable"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("❌ Error crítico al conectar con la Base de Datos: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)

	carpeta_jsons := "./JSONS/"
	archivos, err := os.ReadDir(carpeta_jsons)
	if err != nil {
		log.Fatalf("❌ Error al leer la carpeta %s: %v", carpeta_jsons, err)
	}

	archivosProcesados := 0

	for _, archivo := range archivos {
		if filepath.Ext(archivo.Name()) != ".json" {
			continue
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		func() {
			defer cancel()
			tx, err := db.BeginTx(ctx, nil)
			if err != nil {
				fmt.Printf("❌ [%s] Error al iniciar transaccion: %v\n", archivo.Name(), err)
				return
			}
			defer tx.Rollback()

			qtx := queries.WithTx(tx)

			err = procesarJson(archivo, carpeta_jsons, qtx, ctx)
			if err != nil {
				// ACÁ: Ahora el error va a burbujear con el mensaje literal de Postgres
				fmt.Printf("❌ [%s] Error al procesar: %v\n", archivo.Name(), err)
				return
			}

			err = tx.Commit()
			if err != nil {
				fmt.Printf("❌ [%s] Error crítico al guardar en BD (Commit falló): %v\n", archivo.Name(), err)
				return
			}

			fmt.Printf("✅ DB Insert: %s\n", archivo.Name())
			archivosProcesados++
		}()
	}

	if archivosProcesados == 0 {
		fmt.Println("🤷‍♂️ No se encontró ningún archivo JSON válido para insertar en ./JSONS/.")
	} else {
		fmt.Printf("🚀 ¡Proceso finalizado! Se insertaron %d pacientes en la base de datos.\n", archivosProcesados)
	}
}

func procesarJson(archivo os.DirEntry, carpeta string, qtx *sqlc.Queries, ctx context.Context) error {
	var paciente domain.Paciente
	ruta := filepath.Join(carpeta, archivo.Name())

	contenido, err := os.Open(ruta)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo: %v", err)
	}
	defer contenido.Close()

	err = json.NewDecoder(contenido).Decode(&paciente)
	if err != nil {
		return fmt.Errorf("no se pudo decodificar el JSON: %v", err)
	}

	fecha_parseada, edad, errTrans := transformarDatos(paciente.Fecha, paciente.Edad)
	if errTrans != nil {
		return fmt.Errorf("error transformando fecha/edad: %v", errTrans)
	}

	err = insertarPaciente(paciente, fecha_parseada, edad, qtx, ctx)
	if err != nil {
		return fmt.Errorf("fallo al insertar en tabla Paciente: %v", err)
	}

	for _, desc_micro := range paciente.Descripciones_microscopicas {
		err = insertarDescMicro(desc_micro, paciente.Protocolo, qtx, ctx)
		if err != nil {
			return err // Burbujeamos el error exacto que venga de abajo
		}
	}
	return nil
}

func transformarDatos(fecha, edad string) (time.Time, int64, error) {
	layout := "2-1-2006"
	fecha_parseada, err := time.Parse(layout, fecha)
	if err != nil {
		return time.Time{}, 0, fmt.Errorf("formato de fecha inválido (%s): %v", fecha, err)
	}

	var edad_ int64
	if edad != "" && edad != "null" {
		edad_, err = strconv.ParseInt(edad, 10, 16)
		if err != nil {
			return fecha_parseada, 0, fmt.Errorf("el campo edad contiene caracteres inválidos (%s): %v", edad, err)
		}
	}
	return fecha_parseada, edad_, nil
}

func insertarPaciente(paciente domain.Paciente, fecha_parseada time.Time, edad int64, qtx *sqlc.Queries, ctx context.Context) error {
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

func insertarDescMicro(desc domain.Descripcion_microscopicas, pk string, qtx *sqlc.Queries, ctx context.Context) error {
	_, err := qtx.CreateDescripcionMicroscopica(ctx, sqlc.CreateDescripcionMicroscopicaParams{
		Descripcion:        desc.Descripcion,
		Diagnostico:        sql.NullString{String: desc.Diagnostico.Descripcion, Valid: true},
		PacientesProtocolo: pk,
	})
	if err != nil {
		// Ahora sí, PostgreSQL nos va a decir la verdad
		return fmt.Errorf("DB Error en CreateDescripcionMicroscopica: %v", err)
	}

	err = insertarImagenes(desc.Diagnostico.Imagenes, desc.Descripcion, pk, qtx, ctx)
	if err != nil {
		return err
	}

	err = insertarTablaGrado(desc.TablaGrado, desc.Descripcion, pk, qtx, ctx)
	return err
}

func insertarImagenes(imagenes []string, pk1, pk2 string, qtx *sqlc.Queries, ctx context.Context) error {
	for _, imagen := range imagenes {
		_, err := qtx.CreateImagen(ctx, sqlc.CreateImagenParams{
			Ruta:                                  imagen,
			DescripcionesMicroscopicasDescripcion: pk1,
			DescripcionesMicroscopicasPacientesProtocolo: pk2,
		})
		if err != nil {
			return fmt.Errorf("DB Error en CreateImagen para la ruta '%s': %v", imagen, err)
		}
	}
	return nil
}

func insertarTablaGrado(tabla []domain.Grado_oncologico, pk1, pk2 string, qtx *sqlc.Queries, ctx context.Context) error {
	for _, fila := range tabla {
		puntaje, err := strconv.Atoi(fila.Puntaje)
		if err != nil {
			return fmt.Errorf("error parseando puntaje numérico '%s' de tabla de grado oncologico: %v", fila.Puntaje, err)
		}
		_, err = qtx.CreateGradoOncologico(ctx, sqlc.CreateGradoOncologicoParams{
			Caracteristica:                        fila.Caracteristica,
			MuestraAnalizada:                      sql.NullString{String: fila.Muestra_analizada, Valid: true},
			Puntaje:                               int16(puntaje),
			DescripcionesMicroscopicasDescripcion: pk1,
			DescripcionesMicroscopicasPacientesProtocolo: pk2,
		})
		if err != nil {
			return fmt.Errorf("DB Error en CreateGradoOncologico para '%s': %v", fila.Caracteristica, err)
		}
	}
	return nil
}
