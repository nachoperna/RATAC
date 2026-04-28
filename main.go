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
	
	var desc_microRepo domain.Descripcion_microscopicasRepository = dbrepo.NewDescripcion_microscopicasRepository(queries)
	desc_microServices := application.NewDescripcionMicroscopicaService(desc_microRepo)
	// desc_microHandler := ui.NewDescripcionMicroscopicaHandler(desc_microServices)

	var diagnosticoRepo domain.DiagnosticoRepository = dbrepo.NewDiagnosticoRepository(queries)
	diagnosticoServices := application.NewDiagnosticoService(diagnosticoRepo)
	// diagnosticoHandler := ui.NewDiagnosticoHandler(diagnosticoServices)

	homeHandler := ui.NewHomeHandler(pacienteServices, desc_microServices, diagnosticoServices)

	fs := http.FileServer(http.Dir("./infrastructure/UI/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", homeHandler.ShowHome)
	http.HandleFunc("/pacientes", pacienteHandler.ListPacientes)
	http.HandleFunc("/pacientes/", pacienteHandler.ListPacientesBy)
	http.HandleFunc("/apipacientes", pacienteHandler.APIPacientes)
	http.ListenAndServe(port, nil)
}
