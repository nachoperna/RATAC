package main

import (
	sqlc "RATAC/DB/sqlc"
	"RATAC/application"
	"RATAC/domain"
	dbrepo "RATAC/infrastructure/DB"
	ui "RATAC/infrastructure/UI"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

const db_conection = "host=localhost port=5432 user=admin password=password dbname=RATAC_DB sslmode=disable"
const port = ":8080"

func main() {
	db, err := sql.Open("postgres", db_conection)
	if err != nil {
		log.Fatalf("Error al conectar con la Base de Datos: %v", err)
	}
	defer db.Close()

	queries := sqlc.New(db)

	var pacienteRepo domain.PacienteRepository = dbrepo.NewPacienteRepository(queries)
	pacienteServices := application.NewPacienteService(pacienteRepo)
	pacienteHandler := ui.NewPacienteHandler(pacienteServices)

	http.HandleFunc("/", pacienteHandler.ShowHome)
	http.HandleFunc("/pacientes", pacienteHandler.ListPacientes)
	http.HandleFunc("/apipacientes", pacienteHandler.APIPacientes)
	http.ListenAndServe(port, nil)
}
