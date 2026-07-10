package ui

import (
	"RATAC/application"
	"RATAC/domain"
	"RATAC/views"
	"fmt"
	"net/http"
)

type AdminHandler struct {
	adminService *application.AdminService
}

func NewAdminHandler(adminService *application.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}

func (h *AdminHandler) ProcesarDocumento(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	archivos := r.MultipartForm.File["archivos"]
	var pacientes []domain.Paciente

	for _, archivo := range archivos {
		contenido, err := archivo.Open() // abrimos el archivo
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer contenido.Close() // cerramos el archivo luego de usarlo
		paciente, err := h.adminService.ConvertirDocumento(contenido, archivo.Filename)
		if err != nil {
			// w.WriteHeader(http.StatusBadRequest)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
			// Renderizar templ de error
		} else {
			pacientes = append(pacientes, *paciente)
		}
	}

	fmt.Println("RETORNA PACIENTE CON PROT: ", pacientes[0].Protocolo)
	w.WriteHeader(http.StatusOK)
	views.InformacionExtraida(pacientes).Render(r.Context(), w)
}
