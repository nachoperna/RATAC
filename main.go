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

	var pacienteRepo domain.PacienteRepository = dbrepo.NewPacienteRepository(queries, db)
	pacienteServices := application.NewPacienteService(pacienteRepo)
	pacienteHandler := ui.NewPacienteHandler(pacienteServices)
	
	var desc_microRepo domain.Descripcion_microscopicasRepository = dbrepo.NewDescripcion_microscopicasRepository(queries)
	desc_microServices := application.NewDescripcionMicroscopicaService(desc_microRepo)
	// desc_microHandler := ui.NewDescripcionMicroscopicaHandler(desc_microServices)

	var diagnosticoRepo domain.DiagnosticoRepository = dbrepo.NewDiagnosticoRepository(queries)
	diagnosticoServices := application.NewDiagnosticoService(diagnosticoRepo)
	// diagnosticoHandler := ui.NewDiagnosticoHandler(diagnosticoServices)

	homeHandler := ui.NewHomeHandler(pacienteServices, desc_microServices, diagnosticoServices)
	adminHandler := ui.NewAdminHandler(application.NewAdminService(&dbrepo.AdminRepository{}))
	
	fs_static := http.FileServer(http.Dir("./infrastructure/UI/static"))
	fs_imagenes := http.FileServer(http.Dir("./IMAGENES/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs_static))
	http.Handle("/imagenes/", http.StripPrefix("/imagenes/", fs_imagenes))
	http.HandleFunc("/", homeHandler.ShowHome)
	http.HandleFunc("/pacientes", pacienteHandler.ListPacientes)
	http.HandleFunc("/pacientes/", pacienteHandler.ListPacientesByFiltro)
	http.HandleFunc("/pacientes/nombre", pacienteHandler.ListPacientesBy)
	http.HandleFunc("/paciente/protocolo/{protocolo}", pacienteHandler.ShowFullPaciente)
	http.HandleFunc("/apipacientes", pacienteHandler.APIPacientes)
	http.HandleFunc("/diagnosticos/alta", adminHandler.ProcesarDocumento)

	err = http.ListenAndServe(port, nil)
	if err != nil{
		log.Fatalf("Error al exponer puerto 8080: %v", err)
	}
}
