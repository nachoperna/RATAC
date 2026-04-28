package ui

import (
	"RATAC/application"
	"RATAC/views"
	"encoding/json"
	"fmt"
	"net/http"
)

type PacienteHandler struct {
	pacienteService *application.PacienteService
}

func NewPacienteHandler(pacienteService *application.PacienteService) *PacienteHandler {
	return &PacienteHandler{pacienteService: pacienteService}
}

func (h *PacienteHandler) ListPacientes(w http.ResponseWriter, r *http.Request) {
	pacientes, err := h.pacienteService.ListUltimosPacientes(r.Context())
	if err != nil {
		// renderizar templ de error
	}
	views.ListPacientes(pacientes).Render(r.Context(), w)
}

func (h *PacienteHandler) ListPacientesBy(w http.ResponseWriter, r *http.Request) {
	paciente := r.URL.Query().Get("paciente")
	if paciente != ""{
		pacientes, err := h.pacienteService.GetPacienteByNombre(r.Context(), paciente)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
		views.ListPacientes(pacientes).Render(r.Context(), w)
	}
}

func (h *PacienteHandler) APIPacientes(w http.ResponseWriter, r *http.Request) {
	pacientes, err := h.pacienteService.ListPacientes(r.Context())
	if err != nil {
		// renderizar templ de error
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pacientes[:5])
}
