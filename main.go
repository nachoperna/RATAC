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
	"os"

	_ "github.com/lib/pq"
)

const port = ":8080"

func getConnection() string {
	conn := os.Getenv("DATABASE_URL")
	if conn != "" {
		return conn
	}
	return "host=localhost port=5432 user=admin password=password dbname=RATAC_DB sslmode=disable"
}

func main() {
	db, err := sql.Open("postgres", getConnection())
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

	var usuarioRepo domain.UsuarioRepository = dbrepo.NewUsuarioRepository(queries)
	authService := application.NewAuthService(usuarioRepo)
	authHandler := ui.NewAuthHandler(authService)

	homeHandler := ui.NewHomeHandler(pacienteServices, desc_microServices, diagnosticoServices)
	adminHandler := ui.NewAdminHandler(application.NewAdminService(dbrepo.NewAdminRepository(queries)), pacienteServices)
	
	fs_static := http.FileServer(http.Dir("./infrastructure/UI/static"))
	fs_imagenes := http.FileServer(http.Dir("./IMAGENES/"))

	mux := http.NewServeMux()

	// --- Publico ---
	mux.Handle("/static/", http.StripPrefix("/static/", fs_static))
	mux.Handle("/IMAGENES/", http.StripPrefix("/IMAGENES/", fs_imagenes))
	mux.HandleFunc("/", homeHandler.ShowHome)
	mux.HandleFunc("/pacientes", pacienteHandler.ListPacientes)
	mux.HandleFunc("/pacientes/", pacienteHandler.ListPacientesByFiltro)
	mux.HandleFunc("/pacientes/nombre", pacienteHandler.ListPacientesBy)
	mux.HandleFunc("/paciente/protocolo/{protocolo}", pacienteHandler.ShowFullPaciente)
	mux.HandleFunc("/apipacientes", pacienteHandler.APIPacientes)

	// --- Sesion ---
	mux.HandleFunc("GET /login", authHandler.ShowLogin)
	mux.HandleFunc("POST /login", authHandler.Login)
	mux.HandleFunc("POST /logout", authHandler.Logout)

	// --- Protegido ---
	// Regla a mantener: todo lo que cuelgue de /admin/ o /diagnosticos/ se
	// registra con ui.Requiere. El eje rol -> permiso vive en domain/permisos.go.
	mux.HandleFunc("/admin/panel", ui.Requiere(domain.PermVerPanel, adminHandler.ShowAdminPanel))
	mux.HandleFunc("/diagnosticos", ui.Requiere(domain.PermVerDiagnosticos, adminHandler.DiagnosticosByUser))
	mux.HandleFunc("/diagnosticos/alta", ui.Requiere(domain.PermCargarDiag, adminHandler.ShowCargaDiagnostico))
	mux.HandleFunc("/diagnosticos/alta/procesado", ui.Requiere(domain.PermCargarDiag, adminHandler.ProcesarDocumento))
	mux.HandleFunc("/diagnosticos/alta/borrar_temporal", ui.Requiere(domain.PermCargarDiag, adminHandler.BorrarTemporal))
	mux.HandleFunc("/diagnosticos/alta/carga", ui.Requiere(domain.PermCargarDiag, adminHandler.AltaDiagnostico))
	mux.HandleFunc("/diagnosticos/baja/{protocolo}", ui.Requiere(domain.PermBorrarDiag, pacienteHandler.BorrarPaciente))

	err = http.ListenAndServe(port, ui.ConUsuario(authService, mux))
	if err != nil{
		log.Fatalf("Error al exponer puerto 8080: %v", err)
	}
}
