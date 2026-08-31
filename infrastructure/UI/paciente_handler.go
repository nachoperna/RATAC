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
	offset, err := getOffset(r.URL.Query().Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
		// renderizar templ de error
	}
	pacientes, resultados_total, err := h.pacienteService.ListPacientes(r.Context(), offset)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
		// renderizar templ de error
	}
	w.WriteHeader(http.StatusOK)
	if offset == 0 {
		render(w, r, views.ShowResultados(pacientes, resultados_total, offset, false), true)
	} else {
		render(w, r, views.ListPacientes(pacientes, resultados_total, offset, false), true)
	}
}

func (h *PacienteHandler) ListPacientesBy(w http.ResponseWriter, r *http.Request) {
	paciente := r.URL.Query().Get("paciente")
	offset, err := getOffset(r.URL.Query().Get("offset"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
		// renderizar templ de error
	}
	if paciente != "" {
		pacientes, resultados_total, err := h.pacienteService.GetPacienteByNombre(r.Context(), paciente, offset)
		if err != nil {
			fmt.Println("ERROR: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
			// renderizar templ de error
		}
		w.WriteHeader(http.StatusOK)
		if len(pacientes) == 0 {
			render(w, r, views.SinResultados(), true)
		} else {
			if offset == 0 {
				render(w, r, views.ShowResultados(pacientes, resultados_total, offset, false), true)
			} else {
				render(w, r, views.ListPacientes(pacientes, resultados_total, offset, false), true)
			}
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *PacienteHandler) ListPacientesByFiltro(w http.ResponseWriter, r *http.Request) {
	// Crear el slice que almacenará los filtros ordenados
	var req PayloadRequest

	// Decodificar el JSON del body en nuestro slice
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Error procesando el JSON", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pacientes, resultados_total, err := h.pacienteService.GetPacienteByFiltro(r.Context(), req.Filtros, int8(req.Offset))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
		// renderizar templ de error
	}
	w.WriteHeader(http.StatusOK)
	if len(pacientes) == 0 {
		render(w, r, views.SinResultados(), true)
	} else {
		if req.Offset == 0 {
			render(w, r, views.ShowResultados(pacientes, resultados_total, int8(req.Offset), true), true)
		} else {
			render(w, r, views.ListPacientes(pacientes, resultados_total, int8(req.Offset), true), true)
		}
	}
}

func (h *PacienteHandler) APIPacientes(w http.ResponseWriter, r *http.Request) {
	pacientes, _, err := h.pacienteService.ListPacientes(r.Context(), 0)
	if err != nil {
		// renderizar templ de error
		http.Error(w, "Error al obtener pacientes", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	render(w, r, views.ShowPaciente(*paciente), false)
}

func (h *PacienteHandler) BorrarPaciente(w http.ResponseWriter, r *http.Request) {
	protocolo := r.PathValue("protocolo")
	err := h.pacienteService.DeletePaciente(r.Context(), protocolo)
	if err != nil {
		// renderizar templ de error
		http.Error(w, "Error al eliminar pacientes", http.StatusInternalServerError)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getOffset(offset string) (int8, error) {
	var ioffset int

	if offset == "" {
		return 0, nil
	}
	ioffset, err := strconv.Atoi(offset)
	return int8(ioffset), err
}
