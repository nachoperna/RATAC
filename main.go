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

	authService := application.NewAuthService(dbrepo.NewAuthRepository(queries, db))
	adminService := application.NewAdminService(dbrepo.NewAdminRepository(queries))
	authHandler := ui.NewAuthHandler(authService, adminService)
	homeHandler := ui.NewHomeHandler(pacienteServices, desc_microServices, diagnosticoServices, authService)
	adminHandler := ui.NewAdminHandler(adminService, pacienteServices, authService)

	fs_static := http.FileServer(http.Dir("./infrastructure/UI/static"))
	fs_imagenes := http.FileServer(http.Dir("./IMAGENES/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs_static))
	http.Handle("/IMAGENES/", http.StripPrefix("/IMAGENES/", fs_imagenes))
	http.HandleFunc("/", homeHandler.ShowHome)
	http.HandleFunc("/pacientes", pacienteHandler.ListPacientes)
	http.HandleFunc("/pacientes/", pacienteHandler.ListPacientesByFiltro)
	http.HandleFunc("/pacientes/nombre", pacienteHandler.ListPacientesBy)
	http.HandleFunc("/paciente/protocolo/{protocolo}", pacienteHandler.ShowFullPaciente)
	http.HandleFunc("/apipacientes", pacienteHandler.APIPacientes)
	http.HandleFunc("/diagnosticos/alta", authHandler.LoggerChecker(func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./infrastructure/UI/static/carga_diagnostico.html")
	}))
	http.HandleFunc("/diagnosticos/alta/procesado", authHandler.LoggerChecker(adminHandler.ProcesarDocumento))
	http.HandleFunc("/diagnosticos/alta/borrar_temporal", authHandler.LoggerChecker(adminHandler.BorrarTemporal))
	http.HandleFunc("/diagnosticos/alta/carga", authHandler.LoggerChecker(adminHandler.AltaDiagnostico))
	http.HandleFunc("/diagnosticos/baja/{protocolo}", authHandler.LoggerChecker(pacienteHandler.BorrarPaciente))
	http.HandleFunc("/diagnosticos", authHandler.LoggerChecker(adminHandler.DiagnosticosByUser))
	http.HandleFunc("/admin/panel", authHandler.LoggerChecker(adminHandler.ShowAdminPanel))
	http.HandleFunc("/admin/solicitudes/aceptar/{email}", authHandler.LoggerChecker(adminHandler.SolicitudAprobada))

	// http.HandleFunc("/registrarse", authHandler.NoActionOk)
	// http.HandleFunc("/registrarse", authHandler.NoActionErr)
	http.HandleFunc("/registrarse", authHandler.Registrarse)
	http.HandleFunc("/login", authHandler.Login)
	// http.HandleFunc("/login", authHandler.NoActionOk)
	http.HandleFunc("/logout", authHandler.Logout)
	http.HandleFunc("/activo", authHandler.LoggerChecker(authHandler.SesionActiva))
	
	http.HandleFunc("/ingreso/formulario", func (w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./infrastructure/UI/static/panel_acceso.html")
	})

	err = http.ListenAndServe(port, nil)
	if err != nil{
		log.Fatalf("Error al exponer puerto 8080: %v", err)
	}
}
