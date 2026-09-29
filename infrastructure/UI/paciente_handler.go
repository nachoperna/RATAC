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
		errorInline(w, r, MSJ_ERROR_BUSQUEDA, err, "ListPacientes: offset invalido")
		return
	}
	pacientes, resultados_total, err := h.pacienteService.ListPacientes(r.Context(), offset)
	if err != nil {
		errorInline(w, r, MSJ_ERROR_BUSQUEDA, err, "ListPacientes: consultar pacientes")
		return
	}
	// Con offset > 0 el usuario toco "Mostrar mas": el swap es beforeend, asi que
	// renderizar la lista vacia agregaria un cartel debajo de las tarjetas.
	if len(pacientes) == 0 && offset > 0 {
		errorInline(w, r, MSJ_SIN_MAS_RESULTADOS, nil, "")
		return
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
			errorInline(w, r, MSJ_ERROR_BUSQUEDA, err, "ListPacientesBy: buscar por nombre")
			return
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

	// El offset llega del cliente: se acota antes de convertirlo, porque un valor
	// mayor a 127 desborda el int8 y rompe la consulta.
	offsetPedido := req.Offset
	if offsetPedido < 0 {
		offsetPedido = 0
	}
	if offsetPedido > math.MaxInt8 {
		offsetPedido = math.MaxInt8
	}

	pacientes, resultados_total, err := h.pacienteService.GetPacienteByFiltro(r.Context(), req.Filtros, int8(offsetPedido))
	if err != nil {
		errorInline(w, r, MSJ_ERROR_BUSQUEDA, err, "ListPacientesByFiltro: consultar pacientes")
		return
	}
	if len(pacientes) == 0 {
		// Con offset > 0 el usuario toco "Mostrar mas": ya hay resultados en
		// pantalla y renderizar SinResultados los ensuciaria, porque el swap es
		// beforeend y el cartel se agregaria debajo de las tarjetas.
		if offsetPedido > 0 {
			errorInline(w, r, MSJ_SIN_MAS_RESULTADOS, nil, "")
			return
		}
		w.WriteHeader(http.StatusOK)
		render(w, r, views.SinResultados(), true)
	} else {
		w.WriteHeader(http.StatusOK)
		if offsetPedido == 0 {
			render(w, r, views.ShowResultados(pacientes, resultados_total, int8(offsetPedido), true), true)
		} else {
			render(w, r, views.ListPacientes(pacientes, resultados_total, int8(offsetPedido), true), true)
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
		errorAlUsuario(w, http.StatusInternalServerError, MSJ_ERROR_ELIMINAR, err, "BorrarPaciente: "+protocolo)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getOffset(offset string) (int8, error) {
	if offset == "" {
		return 0, nil
	}
	ioffset, err := strconv.Atoi(offset)
	if err != nil {
		return 0, err
	}
	// El offset viaja como int8 en toda la cadena: por encima de 127 desborda a
	// negativo y la conversion a uint64 del OFFSET de SQL da un numero enorme,
	// con lo que la consulta devuelve cualquier cosa. Se acota al maximo.
	if ioffset < 0 {
		ioffset = 0
	}
	if ioffset > math.MaxInt8 {
		ioffset = math.MaxInt8
	}
	return int8(ioffset), nil
}
