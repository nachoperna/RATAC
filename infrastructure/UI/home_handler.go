package ui

import (
	"RATAC/application"
	"RATAC/views"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/a-h/templ"
)

type HomeHandler struct {
	pacienteService *application.PacienteService
	DescripcionMicroscopicaService *application.DescripcionMicroscopicaService
	DiagnosticoService *application.DiagnosticoService
}

func NewHomeHandler(
	pacienteService *application.PacienteService, 
	DescripcionMicroscopicaService *application.DescripcionMicroscopicaService,
	DiagnosticoService *application.DiagnosticoService) *HomeHandler {
	return &HomeHandler{
		pacienteService: pacienteService,
		DescripcionMicroscopicaService: DescripcionMicroscopicaService,
		DiagnosticoService: DiagnosticoService,
	}
}

func (h *HomeHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.ServeFile(w, r, "./infrastructure/UI/static/ruta_invalida.html")
		return
	}

	cant_docs, err := h.pacienteService.CountPacientes(r.Context())
	if err != nil {
		fmt.Println("Error obteniendo cantidad de pacientes", err)
	}
	cant_imgs, err := h.DiagnosticoService.CountImagenes(r.Context())
	if err != nil {
		fmt.Println("Error obteniendo cantidad de imagenes", err)
	}
	cant_diagnosticos, err := h.DescripcionMicroscopicaService.CountDiagnosticos(r.Context())
	if err != nil {
		fmt.Println("Error obteniendo cantidad de diagnosticos", err)
	}
	pacientes, tienen_diagnosticos, err := h.pacienteService.ListUltimosPacientes(r.Context())
	if err != nil {
		fmt.Println("Error obteniendo ultimos pacientes", err)
	}
	
	var casos templ.Component
	if pacientes != nil {
		casos = views.UltimosCasos(pacientes, tienen_diagnosticos)
	}
	tmp_aux := template.New("index.html").Funcs(template.FuncMap{ "render": renderTempl })
	tmp, err := tmp_aux.ParseFiles("./infrastructure/UI/static/index.html")
	if err != nil {
		fmt.Printf("Error al parsear el template: %v", err) // Esto saldrá en tu consola
		http.Error(w, "No se pudo cargar la página", http.StatusInternalServerError)
		return // Importante: salir de la función
	}

	datos := map[string]any{
		"cant_docs": cant_docs,
		"cant_imgs": cant_imgs,
		"cant_diagnosticos": cant_diagnosticos,
		"UltimosCasos": casos,
	}
	tmp.Execute(w, datos)
}

func renderTempl(c templ.Component) template.HTML {
	var buf strings.Builder
	c.Render(context.Background(), &buf)
	return template.HTML(buf.String())
}
