package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"encoding/json"
	"net/http"
)

type DiagnosticoHandler struct {
	service *application.DiagnosticoService
}

func NewDiagnosticoHandler(service *application.DiagnosticoService) *DiagnosticoHandler {
	return &DiagnosticoHandler{
		service: service,
	}
}

// APICreateDiagnostico recibe el JSON y lo asocia a un protocolo y una descripción microscópica
func (h *DiagnosticoHandler) APICreateDiagnostico(w http.ResponseWriter, r *http.Request) {
	protocolo := r.URL.Query().Get("protocolo")
	descripcionMicro := r.URL.Query().Get("descripcion")

	if protocolo == "" || descripcionMicro == "" {
		http.Error(w, "Faltan parámetros en la URL: protocolo y/o descripcion", http.StatusBadRequest)
		return
	}

	var diag domain.Diagnostico
	if err := json.NewDecoder(r.Body).Decode(&diag); err != nil {
		http.Error(w, "Error procesando el JSON", http.StatusBadRequest)
		return
	}

	err := h.service.CreateDiagnostico(r.Context(), protocolo, descripcionMicro, &diag)
	if err != nil {
		http.Error(w, "Error guardando el diagnóstico en la BD", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// APIGetDiagnostico devuelve el diagnóstico (con sus imágenes) en formato JSON
func (h *DiagnosticoHandler) APIGetDiagnostico(w http.ResponseWriter, r *http.Request) {
	protocolo := r.URL.Query().Get("protocolo")
	descripcionMicro := r.URL.Query().Get("descripcion")

	if protocolo == "" || descripcionMicro == "" {
		http.Error(w, "Faltan parámetros en la URL: protocolo y/o descripcion", http.StatusBadRequest)
		return
	}

	diagnostico, err := h.service.GetDiagnostico(r.Context(), protocolo, descripcionMicro)
	if err != nil {
		http.Error(w, "Error consultando el diagnóstico en la BD", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(diagnostico)
}

// APIDeleteImagen permite borrar una imagen suelta pasándole su ruta
func (h *DiagnosticoHandler) APIDeleteImagen(w http.ResponseWriter, r *http.Request) {
	ruta := r.URL.Query().Get("ruta")

	if ruta == "" {
		http.Error(w, "Falta el parámetro ruta en la URL", http.StatusBadRequest)
		return
	}

	err := h.service.DeleteImagen(r.Context(), ruta)
	if err != nil {
		http.Error(w, "Error eliminando la imagen de la BD", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
