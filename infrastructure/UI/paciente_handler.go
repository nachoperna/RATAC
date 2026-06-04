package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/views"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type PacienteHandler struct {
	pacienteService *application.PacienteService
}

func NewPacienteHandler(pacienteService *application.PacienteService) *PacienteHandler {
	return &PacienteHandler{pacienteService: pacienteService}
}

type PayloadRequest struct {
	Filtros []domain.Filtro `json:"filtros"`
	Offset  int             `json:"offset"`
}

func (h *PacienteHandler) ListPacientes(w http.ResponseWriter, r *http.Request) {
	pacientes, err := h.pacienteService.ListPacientes(r.Context())
	if err != nil {
		// renderizar templ de error
	}
	views.ShowResultados(pacientes, 0, 0, false).Render(r.Context(), w)
}

func (h *PacienteHandler) ListPacientesBy(w http.ResponseWriter, r *http.Request) {
	paciente := r.URL.Query().Get("paciente")
	offset, _ := getOffset(r.URL.Query().Get("offset"))

	if paciente != "" {
		pacientes, resultados_total, err := h.pacienteService.GetPacienteByNombre(r.Context(), paciente, offset)
		if err != nil {
			fmt.Println("ERROR: ", err)
			// renderizar templ de error
		}
		if len(pacientes) == 0 {
			views.SinResultados().Render(r.Context(), w)
		} else {
			if offset == 0 {
				views.ShowResultados(pacientes, resultados_total, offset, false).Render(r.Context(), w)
			} else {
				views.ListPacientes(pacientes, resultados_total, offset, false).Render(r.Context(), w)
			}
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

	pacientes, resultados_total, err := h.pacienteService.GetPacienteByFiltro(r.Context(), req.Filtros, int8(req.Offset))
	if err != nil {
		// renderizar templ de error
	}
	if len(pacientes) == 0 {
		views.SinResultados().Render(r.Context(), w)
	} else {
		if req.Offset == 0 {
			views.ShowResultados(pacientes, resultados_total, int8(req.Offset), true).Render(r.Context(), w)
		} else {
			views.ListPacientes(pacientes, resultados_total, int8(req.Offset), true).Render(r.Context(), w)
		}
	}
}

func (h *PacienteHandler) APIPacientes(w http.ResponseWriter, r *http.Request) {
	pacientes, err := h.pacienteService.ListPacientes(r.Context())
	if err != nil {
		// renderizar templ de error
		http.Error(w, "Error al obtener pacientes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	// Lógica segura para evitar el panic de "out of range"
	limite := min(5, len(pacientes))

	json.NewEncoder(w).Encode(pacientes[:limite])
}

func (h *PacienteHandler) ShowFullPaciente(w http.ResponseWriter, r *http.Request) {
	protocolo := r.PathValue("protocolo")
	paciente, err := h.pacienteService.GetAllFromPaciente(r.Context(), protocolo)
	if err != nil || paciente == nil {
		// renderizar templ de error
		fmt.Println("Algo salio mal")
	}
	views.ShowPaciente(*paciente).Render(r.Context(), w)
}

func getOffset(offset string) (int8, error) {
	var ioffset int

	if offset == "" {
		return 0, nil
	}
	ioffset, err := strconv.Atoi(offset)
	return int8(ioffset), err
}
