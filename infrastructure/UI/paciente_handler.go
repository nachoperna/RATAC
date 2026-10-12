package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/views"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

type PacienteHandler struct {
	pacienteService *application.PacienteService
	pageHandler *PageHandler
}

func NewPacienteHandler(pacienteService *application.PacienteService, pageHandler *PageHandler) *PacienteHandler {
	return &PacienteHandler{
		pacienteService: pacienteService,
		pageHandler: pageHandler,
	}
}

type PayloadRequest struct {
	Filtros []domain.Filtro `json:"filtros"`
	Offset  int             `json:"offset"`
}

func (h *PacienteHandler) ListPacientes(w http.ResponseWriter, r *http.Request) {
	offset, err := getOffset(r.URL.Query().Get("offset"))
	if err != nil {
		errorInline(w, r, MSJ_ERROR_BUSQUEDA, err, "ListPacientes: offset invalido")
		return
	}
	pacientes, resultados_total, err := h.pacienteService.ListPacientes(r.Context(), offset)
	if err != nil {
		errorInline(w, r, MSJ_ERROR_BUSQUEDA, err, "ListPacientes: consultar pacientes")
		return
	}
	w.WriteHeader(http.StatusOK)
	if offset == 0 {
		h.pageHandler.Render(w, r, views.ShowResultados(pacientes, resultados_total, offset, false), true)
	} else {
		h.pageHandler.Render(w, r, views.ListPacientes(pacientes, resultados_total, offset, false), true)
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
			errorInline(w, r, MSJ_ERROR_BUSQUEDA, err, "ListPacientesBy: buscar por nombre")
			return
		}
		w.WriteHeader(http.StatusOK)
		if len(pacientes) == 0 {
			h.pageHandler.Render(w, r, views.SinResultados(), true)
		} else {
			if offset == 0 {
				h.pageHandler.Render(w, r, views.ShowResultados(pacientes, resultados_total, offset, false), true)
			} else {
				h.pageHandler.Render(w, r, views.ListPacientes(pacientes, resultados_total, offset, false), true)
			}
		}
	} else {
		errorInline(w, r, "Escriba el nombre de un paciente para buscar.", nil, "")
		return
	}
}

func (h *PacienteHandler) ListPacientesByFiltro(w http.ResponseWriter, r *http.Request) {
	// Crear el slice que almacenará los filtros ordenados
	var req PayloadRequest

	// Decodificar el JSON del body en nuestro slice
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorInline(w, r, MSJ_ERROR_BUSQUEDA, err, "ListPacientesByFiltro: JSON invalido")
		return
	}

	offsetPedido := req.Offset
	if offsetPedido < 0 {
		offsetPedido = 0
	}
	if offsetPedido > math.MaxInt16 {
		offsetPedido = math.MaxInt16
	}

	pacientes, resultados_total, err := h.pacienteService.GetPacienteByFiltro(r.Context(), req.Filtros, int16(offsetPedido))
	if err != nil {
		errorInline(w, r, MSJ_ERROR_BUSQUEDA, err, "ListPacientesByFiltro: consultar pacientes")
		return
	}
	if len(pacientes) == 0 {
		w.WriteHeader(http.StatusOK)
		h.pageHandler.Render(w, r, views.SinResultados(), true)
	} else {
		w.WriteHeader(http.StatusOK)
		if offsetPedido == 0 {
			h.pageHandler.Render(w, r, views.ShowResultados(pacientes, resultados_total, int16(offsetPedido), true), true)
		} else {
			h.pageHandler.Render(w, r, views.ListPacientes(pacientes, resultados_total, int16(offsetPedido), true), true)
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
	token, err := r.Cookie(NOMBRE_TOKEN)
	if err != nil && err != http.ErrNoCookie {
		// renderizar templ de error
		http.Error(w, "Error al cargar paciente", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	valor_token := ""
	if token != nil {
		valor_token = token.Value
	}

	paciente, es_dueño, err := h.pacienteService.GetAllFromPaciente(r.Context(), protocolo, valor_token)
	if err != nil || paciente == nil {
		// renderizar templ de error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	h.pageHandler.Render(w, r, views.ShowPaciente(*paciente, es_dueño), false)
}

func (h *PacienteHandler) BorrarPaciente(w http.ResponseWriter, r *http.Request) {
	protocolo := r.PathValue("protocolo")
	err := h.pacienteService.DeletePacienteYArchivos(r.Context(), protocolo)
	if err != nil {
		errorAlUsuario(w, http.StatusInternalServerError, MSJ_ERROR_ELIMINAR, err, "BorrarPaciente: "+protocolo)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getOffset(offset string) (int16, error) {
	if offset == "" {
		return 0, nil
	}
	ioffset, err := strconv.Atoi(offset)
	if err != nil {
		return 0, err
	}
	if ioffset < 0 {
		ioffset = 0
	}
	if ioffset > math.MaxInt16 {
		ioffset = math.MaxInt16
	}
	return int16(ioffset), nil
}
