package ui

import (
	"RATAC/application"
	"RATAC/domain"
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
type PayloadRequest struct {
	Filtros []domain.Filtro `json:"filtros"`
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
			// renderizar templ de error
		}
		if len(pacientes) == 0{
			views.SinResultados().Render(r.Context(), w)
		}else{
			views.ListPacientes(pacientes).Render(r.Context(), w)
		}
	}
}

func (h *PacienteHandler) ListPacientesByFiltro(w http.ResponseWriter, r *http.Request) {
	// Crear el slice que almacenará los filtros ordenados
	var req PayloadRequest

	// Decodificar el JSON del body en nuestro slice
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error procesando el JSON", http.StatusBadRequest)
		return
	}
	pacientes, err := h.pacienteService.GetPacienteByFiltro(r.Context(), req.Filtros)
	if err != nil {
		// renderizar templ de error
	}
	if len(pacientes) == 0{
		views.SinResultados().Render(r.Context(), w)
	}else{
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
