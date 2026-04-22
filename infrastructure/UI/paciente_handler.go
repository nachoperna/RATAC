package ui

import (
	"RATAC/application"
	"net/http"
)

type PacienteHandler struct{
	pacienteService *application.PacienteService
}

func NewPacienteHandler(pacienteService *application.PacienteService) *PacienteHandler {
	return &PacienteHandler{pacienteService: pacienteService}
}

func (h *PacienteHandler) ListPacientes(w http.ResponseWriter, r *http.Request) {
	pacientes, err := h.pacienteService.ListPacientes(r.Context())
	if err != nil{
		// renderizar templ de error
	}
	print(pacientes)
	// renderizar templ de muestra de pacientes
}
