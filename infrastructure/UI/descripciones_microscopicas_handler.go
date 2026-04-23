package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"encoding/json"
	"net/http"
)

type DescripcionMicroscopicaHandler struct {
	service *application.DescripcionMicroscopicaService
}

func NewDescripcionMicroscopicaHandler(service *application.DescripcionMicroscopicaService) *DescripcionMicroscopicaHandler {
	return &DescripcionMicroscopicaHandler{
		service: service,
	}
}

// APICreateDescripcion lee un JSON del body y lo guarda en la base de datos
func (h *DescripcionMicroscopicaHandler) APICreateDescripcion(w http.ResponseWriter, r *http.Request) {
	// Asumimos que pasás el protocolo por query param: /api/descripciones?protocolo=123
	// Dependiendo de tu router (Chi, Gorilla Mux, o el estándar de Go 1.22), esto puede variar.
	protocolo := r.URL.Query().Get("protocolo")
	if protocolo == "" {
		http.Error(w, "Falta el parámetro protocolo", http.StatusBadRequest)
		return
	}

	var desc domain.Descripcion_microscopicas
	if err := json.NewDecoder(r.Body).Decode(&desc); err != nil {
		http.Error(w, "Error procesando el JSON", http.StatusBadRequest)
		return
	}

	err := h.service.CreateDescripcionMicroscopica(r.Context(), protocolo, &desc)
	if err != nil {
		http.Error(w, "Error guardando en la BD", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// APIGetDescripciones devuelve un JSON con todas las descripciones de un protocolo
func (h *DescripcionMicroscopicaHandler) APIGetDescripciones(w http.ResponseWriter, r *http.Request) {
	protocolo := r.URL.Query().Get("protocolo")
	if protocolo == "" {
		http.Error(w, "Falta el parámetro protocolo", http.StatusBadRequest)
		return
	}

	descripciones, err := h.service.GetDescripcionesByProtocolo(r.Context(), protocolo)
	if err != nil {
		http.Error(w, "Error consultando a la BD", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(descripciones)
}
