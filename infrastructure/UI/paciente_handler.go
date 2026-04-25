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
func (h *PacienteHandler) ShowHome(w http.ResponseWriter, r *http.Request) {
	err := views.ShowHome().Render(r.Context(), w)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
}
func (h *PacienteHandler) ListPacientes(w http.ResponseWriter, r *http.Request) {
	pacientes, err := h.pacienteService.ListUltimosPacientes(r.Context())
	if err != nil {
		// renderizar templ de error
	}
	views.ListPacientes(pacientes).Render(r.Context(), w)
}

func (h *PacienteHandler) APIPacientes(w http.ResponseWriter, r *http.Request) {
	pacientes, err := h.pacienteService.ListPacientes(r.Context())
	if err != nil {
		// renderizar templ de error
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pacientes[:5])
}
